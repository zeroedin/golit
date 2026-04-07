/**
 * Template Result Collector with Lit Hydration Markers
 *
 * Converts Lit's TemplateResult objects into HTML strings with
 * Lit-compatible hydration comment markers:
 *   <!--lit-part DIGEST--> / <!--/lit-part-->  (TemplateResult boundaries)
 *   <!--lit-part--> / <!--/lit-part-->          (primitive child values)
 *   <!--lit-node N-->                           (before elements with bindings)
 *
 * Uses a proper HTML state machine (modeled after Lit's getTemplateHtml)
 * to classify bindings, ensuring hydration compatibility.
 */

// PartType constants matching Lit's internal values
const ATTRIBUTE_PART = 1;
const CHILD_PART = 2;
const PROPERTY_PART = 3;
const BOOLEAN_ATTRIBUTE_PART = 4;
const EVENT_PART = 5;
const ELEMENT_PART = 6;

// ============================================================
// Digest computation (DJB2, matching Lit's digestForTemplateResult)
// ============================================================

function computeDigest(strings) {
  const hashes = new Uint32Array(2).fill(5381);
  for (const s of strings) {
    for (let i = 0; i < s.length; i++) {
      hashes[i % 2] = (hashes[i % 2] * 33) ^ s.charCodeAt(i);
    }
  }
  const bytes = new Uint8Array(hashes.buffer);
  let binary = '';
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  if (typeof btoa === 'function') return btoa(binary);
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
  let result = '';
  for (let i = 0; i < binary.length; i += 3) {
    const a = binary.charCodeAt(i);
    const b = i + 1 < binary.length ? binary.charCodeAt(i + 1) : 0;
    const c = i + 2 < binary.length ? binary.charCodeAt(i + 2) : 0;
    result += chars[a >> 2];
    result += chars[((a & 3) << 4) | (b >> 4)];
    result += i + 1 < binary.length ? chars[((b & 15) << 2) | (c >> 6)] : '=';
    result += i + 2 < binary.length ? chars[c & 63] : '=';
  }
  return result;
}

// ============================================================
// Binding classifier — HTML state machine
// ============================================================

// Regexes ported from Lit's getTemplateHtml (lit-html.ts).
// These track the five parsing states in the HTML scanner.
const SPACE_CHAR = `[ \\t\\n\\f\\r]`;
const NAME_CHAR = `[^\\s"'>=/]`;
const ATTR_VALUE_CHAR = `[^ \\t\\n\\f\\r"'\`<>=]`;

const textEndRegex = /<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g;
const COMMENT_START = 1;
const TAG_NAME = 2;

const commentEndRegex = /-->/g;
const comment2EndRegex = />/g;

const tagEndRegex = new RegExp(
  `>|${SPACE_CHAR}(?:(${NAME_CHAR}+)(${SPACE_CHAR}*=${SPACE_CHAR}*(?:${ATTR_VALUE_CHAR}|("|')|))|$)`, 'g'
);
const ENTIRE_MATCH = 0;
const ATTRIBUTE_NAME = 1;
const SPACES_AND_EQUALS = 2;
const QUOTE_CHAR = 3;

const singleQuoteAttrEndRegex = /'/g;
const doubleQuoteAttrEndRegex = /"/g;
const rawTextElement = /^(?:script|style|textarea|title)$/i;

/**
 * Classify each binding in a TemplateResult by walking the template strings
 * with an HTML tokenizer. Returns one descriptor per expression.
 *
 * Each descriptor: { type, name, partType, tagName, strings }
 *   type:     'child' | 'attribute' | 'boolean' | 'property' | 'event' | 'element'
 *   partType: CHILD_PART | ATTRIBUTE_PART | BOOLEAN_ATTRIBUTE_PART | PROPERTY_PART | EVENT_PART | ELEMENT_PART
 *   name:     attribute name without prefix (null for child/element)
 *   strings:  interleaved strings for the attribute (for multi-expr attrs)
 *   tagName:  tag name of the element this binding is on
 */
function classifyBindings(strings) {
  const l = strings.length - 1;
  const bindings = [];

  let rawTextEndRegex2 = undefined;
  let regex = textEndRegex;

  for (let i = 0; i < l; i++) {
    const s = strings[i];
    let attrNameEndIndex = -1;
    let attrName = undefined;
    let lastIndex = 0;
    let match;
    let currentTagName = '';

    while (lastIndex < s.length) {
      regex.lastIndex = lastIndex;
      match = regex.exec(s);
      if (match === null) break;
      if (regex.lastIndex === lastIndex) {
        // Zero-length match — advance to avoid infinite loop
        lastIndex++;
        continue;
      }
      lastIndex = regex.lastIndex;

      if (regex === textEndRegex) {
        if (match[COMMENT_START] === '!--') {
          regex = commentEndRegex;
        } else if (match[COMMENT_START] !== undefined) {
          regex = comment2EndRegex;
        } else if (match[TAG_NAME] !== undefined) {
          if (rawTextElement.test(match[TAG_NAME])) {
            rawTextEndRegex2 = new RegExp(`</${match[TAG_NAME]}`, 'g');
          }
          currentTagName = match[TAG_NAME];
          regex = tagEndRegex;
        }
      } else if (regex === tagEndRegex) {
        if (match[ENTIRE_MATCH] === '>') {
          regex = rawTextEndRegex2 || textEndRegex;
          attrNameEndIndex = -1;
        } else if (match[ATTRIBUTE_NAME] === undefined) {
          // Attribute name position (no name matched = element binding position)
          attrNameEndIndex = -2;
        } else {
          attrNameEndIndex = regex.lastIndex - match[SPACES_AND_EQUALS].length;
          attrName = match[ATTRIBUTE_NAME];
          regex =
            match[QUOTE_CHAR] === undefined
              ? tagEndRegex
              : match[QUOTE_CHAR] === '"'
                ? doubleQuoteAttrEndRegex
                : singleQuoteAttrEndRegex;
        }
      } else if (regex === doubleQuoteAttrEndRegex || regex === singleQuoteAttrEndRegex) {
        regex = tagEndRegex;
      } else if (regex === commentEndRegex || regex === comment2EndRegex) {
        regex = textEndRegex;
      } else {
        // Raw text end regex
        regex = tagEndRegex;
        rawTextEndRegex2 = undefined;
      }
    }

    // Classify this binding based on current parser state
    if (regex === textEndRegex) {
      // Child position (between tags)
      bindings.push({
        type: 'child',
        partType: CHILD_PART,
        name: null,
        tagName: null,
        strings: null,
      });
    } else if (attrNameEndIndex >= 0) {
      // First/only expression in an attribute
      const prefix = attrName[0];
      let type, partType, cleanName;
      if (prefix === '.') {
        type = 'property'; partType = PROPERTY_PART; cleanName = attrName.slice(1);
      } else if (prefix === '?') {
        type = 'boolean'; partType = BOOLEAN_ATTRIBUTE_PART; cleanName = attrName.slice(1);
      } else if (prefix === '@') {
        type = 'event'; partType = EVENT_PART; cleanName = attrName.slice(1);
      } else {
        type = 'attribute'; partType = ATTRIBUTE_PART; cleanName = attrName;
      }
      bindings.push({
        type: type,
        partType: partType,
        name: cleanName,
        tagName: currentTagName,
        strings: null, // filled in below for multi-expr
      });
    } else if (attrNameEndIndex === -2) {
      // Element binding position (after tag name)
      bindings.push({
        type: 'element',
        partType: ELEMENT_PART,
        name: null,
        tagName: currentTagName,
        strings: null,
      });
    } else {
      // Inside a quoted or unquoted attribute value — multi-expression attribute.
      // This is a continuation of the previous attribute binding.
      const prev = bindings.length > 0 ? bindings[bindings.length - 1] : null;
      bindings.push({
        type: prev ? prev.type : 'attribute',
        partType: prev ? prev.partType : ATTRIBUTE_PART,
        name: prev ? prev.name : attrName,
        tagName: prev ? prev.tagName : currentTagName,
        strings: null,
        _isMultiContinuation: true,
      });
    }
  }

  // Build attribute strings arrays for multi-expression attributes.
  // Group consecutive bindings for the same attribute.
  for (let i = 0; i < bindings.length; i++) {
    const b = bindings[i];
    if (b.type === 'child' || b.type === 'element') continue;
    if (b._isMultiContinuation) continue;

    // Find how many consecutive bindings share this attribute
    let end = i + 1;
    while (end < bindings.length && bindings[end]._isMultiContinuation) end++;
    const count = end - i;

    // Build the strings array: [before_first_expr, between_1_2, between_2_3, ..., after_last_expr]
    // For a single-expression attribute: strings = ['', ''] (the value between the quotes)
    // We extract these from the raw template strings
    const attrStrings = [];
    for (let j = i; j <= end && j < strings.length; j++) {
      if (j === i) {
        // Extract the part after the last '=' and optional quote
        const s = strings[j];
        const eqIdx = s.lastIndexOf('=');
        if (eqIdx >= 0) {
          let start = eqIdx + 1;
          while (start < s.length && (s[start] === '"' || s[start] === "'" || s[start] === ' ')) start++;
          attrStrings.push(s.slice(start));
        } else {
          attrStrings.push('');
        }
      } else if (j === end) {
        // Extract the part before the closing quote
        const s = strings[j];
        const quoteIdx = s.search(/["'>\s]/);
        attrStrings.push(quoteIdx >= 0 ? s.slice(0, quoteIdx) : s);
      } else {
        // Middle string (between expressions in same attribute)
        attrStrings.push(strings[j]);
      }
    }

    // Assign to all bindings in this group
    for (let j = i; j < end; j++) {
      bindings[j].strings = attrStrings;
      bindings[j]._groupStart = i;
      bindings[j]._groupEnd = end;
      bindings[j]._indexInGroup = j - i;
    }
  }

  return bindings;
}

// ============================================================
// Main entry point
// ============================================================

globalThis.__collectTemplateResult = function collectTemplateResult(value, isRoot) {
  if (value === null || value === undefined) {
    return isRoot ? '<!--lit-part--><!--/lit-part-->' : '';
  }
  if (typeof value === 'symbol') {
    return isRoot ? '<!--lit-part--><!--/lit-part-->' : '';
  }
  if (typeof value === 'string') return escapeHTML(value);
  if (typeof value === 'number') return String(value);
  if (typeof value === 'boolean') return value ? 'true' : '';

  // TemplateResult
  if (value && value['_$litType$'] !== undefined) {
    const digest = computeDigest(value.strings);
    return `<!--lit-part ${digest}-->` + renderTemplateResult(value) + '<!--/lit-part-->';
  }

  // DirectiveResult
  if (value && value['_$litDirective$'] !== undefined) {
    return renderDirective(value, { type: CHILD_PART, name: null, strings: null });
  }

  // Iterable
  if (Array.isArray(value)) {
    return value.map(v => collectTemplateResult(v)).join('');
  }
  if (value && typeof value[Symbol.iterator] === 'function') {
    let result = '';
    for (const item of value) result += collectTemplateResult(item);
    return result;
  }

  return escapeHTML(String(value));
};

// ============================================================
// Template rendering with classified bindings
// ============================================================

function renderTemplateResult(result) {
  const strings = result.strings;
  const values = result.values;
  const bindings = classifyBindings(strings);

  let html = '';
  let nodeIndex = 0;

  // Track which elements need <!--lit-node--> markers.
  // Collect element indices that have any binding.
  const boundElementNodes = new Set();
  let tempNodeIdx = 0;
  for (let i = 0; i < strings.length; i++) {
    const elemCount = countNodes(strings[i]);
    if (i < bindings.length) {
      const b = bindings[i];
      if (b.type !== 'child') {
        // This binding is on the element at (tempNodeIdx + elemCount - 1) if the element
        // is in this string, or on the current element context.
        // We use the last element in this static string as the bound element.
        if (elemCount > 0) {
          boundElementNodes.add(tempNodeIdx + elemCount - 1);
        } else {
          // Element opened in a previous string, still open
          // The bound element is the last one we saw
          boundElementNodes.add(tempNodeIdx - 1 >= 0 ? tempNodeIdx - 1 : 0);
        }
      }
    }
    tempNodeIdx += elemCount;
    if (i < bindings.length && bindings[i].type === 'child') {
      tempNodeIdx++; // child comment marker is a node
    }
  }

  for (let i = 0; i < strings.length; i++) {
    let staticPart = strings[i];

    // Inject <!--lit-node N--> before elements that have bindings
    if (i < bindings.length && bindings[i].type !== 'child') {
      // Only inject on the first binding for this element
      const b = bindings[i];
      const isFirstForElement = (i === 0 || bindings[i - 1].type === 'child' ||
        (b._groupStart !== undefined && b._indexInGroup === 0 && (i === 0 || bindings[i-1].type === 'child' || bindings[i-1]._groupEnd <= i)));

      if (isFirstForElement || (b._groupStart === undefined && !b._isMultiContinuation)) {
        // The bound element is the LAST element open tag in this static string.
        // nodeIndex is the count at the START of the string, so the bound element's
        // index is nodeIndex + (nodes in this string) - 1.
        const boundNodeIndex = nodeIndex + countNodes(staticPart) - 1;
        staticPart = injectLitNodeMarker(staticPart, boundNodeIndex);
      }
    }

    // Strip binding syntax from static part for non-child bindings
    if (i < bindings.length) {
      const b = bindings[i];
      if (b.type === 'attribute' && !b._isMultiContinuation) {
        // First expression in a regular attribute: strip trailing `attrname="` or `attrname=`
        staticPart = stripTrailingAttr(staticPart);
      } else if (b.type === 'attribute' && b._isMultiContinuation) {
        // Middle/continuation of multi-expr attr: don't emit the interstitial string
        // (it's handled when we emit the full attribute value)
        staticPart = '';
      } else if (b.type === 'boolean') {
        staticPart = stripTrailingAttr(staticPart);
      } else if (b.type === 'property' || b.type === 'event') {
        staticPart = stripTrailingAttr(staticPart);
      } else if (b.type === 'element') {
        // Element binding: no attribute to strip, just emit static part
      }
    }

    // For the string AFTER a non-child binding, strip the leading closing quote
    if (i > 0 && i <= bindings.length) {
      const prevB = bindings[i - 1];
      if (prevB.type !== 'child') {
        // Check if this is the last expression in an attribute group
        const isLastInGroup = !prevB._isMultiContinuation ||
          (i < bindings.length && !bindings[i]._isMultiContinuation) ||
          i === bindings.length;
        if (isLastInGroup && !bindings[i]?._isMultiContinuation) {
          staticPart = staticPart.replace(/^['"]/, '');
        }
      }
    }

    html += staticPart;
    nodeIndex += countNodes(strings[i]);

    // Emit binding value
    if (i < values.length && i < bindings.length) {
      const b = bindings[i];

      switch (b.type) {
        case 'child':
          html += resolveChildValue(values[i]);
          nodeIndex++; // child part comment is a node
          break;

        case 'attribute': {
          // For multi-expression attributes, only emit at the last expression
          const groupEnd = b._groupEnd !== undefined ? b._groupEnd : i + 1;
          const isLast = (i === groupEnd - 1);
          if (isLast) {
            // Collect all values for this attribute
            const groupStart = b._groupStart !== undefined ? b._groupStart : i;
            const attrValues = [];
            for (let j = groupStart; j < groupEnd; j++) {
              attrValues.push(resolveAttrValue(values[j], bindings[j]));
            }
            // Interpolate with the attribute's interleaved strings
            const attrStrings = b.strings || ['', ''];
            let attrValue = '';
            for (let j = 0; j < attrStrings.length; j++) {
              attrValue += attrStrings[j];
              if (j < attrValues.length) attrValue += attrValues[j];
            }
            html += ` ${b.name}="${escapeAttr(attrValue)}"`;
          }
          break;
        }

        case 'boolean':
          if (values[i]) {
            html += ` ${b.name}`;
          }
          break;

        case 'property':
        case 'event':
        case 'element':
          // No output for these binding types — hydration handles them client-side
          break;
      }
    }
  }

  return html;
}

// ============================================================
// Value resolution
// ============================================================

let _depth = 0;

function resolveChildValue(val) {
  if (val === null || val === undefined) return '<!--lit-part--><!--/lit-part-->';
  if (typeof val === 'symbol') return '<!--lit-part--><!--/lit-part-->';
  if (_depth > 10) return '<!--lit-part--><!--/lit-part-->';

  _depth++;
  try {

  // TemplateResult
  if (val && val['_$litType$'] !== undefined) {
    const digest = computeDigest(val.strings);
    return `<!--lit-part ${digest}-->` + renderTemplateResult(val) + '<!--/lit-part-->';
  }

  // DirectiveResult
  if (val && val['_$litDirective$'] !== undefined) {
    const rendered = renderDirective(val, { type: CHILD_PART, name: null, strings: null });
    if (rendered && typeof rendered === 'string' && rendered.startsWith('<!--lit-part')) return rendered;
    return '<!--lit-part-->' + (rendered || '') + '<!--/lit-part-->';
  }

  // Primitive (check BEFORE iterable — strings are iterable in JS!)
  if (typeof val === 'string') return '<!--lit-part-->' + escapeHTML(val) + '<!--/lit-part-->';
  if (typeof val === 'number' || typeof val === 'boolean') return '<!--lit-part-->' + String(val) + '<!--/lit-part-->';

  // Iterable (arrays and other iterables, but NOT strings)
  if (Array.isArray(val)) {
    let result = '<!--lit-part-->';
    for (const item of val) result += resolveChildValue(item);
    return result + '<!--/lit-part-->';
  }
  if (val && typeof val[Symbol.iterator] === 'function') {
    let result = '<!--lit-part-->';
    for (const item of val) result += resolveChildValue(item);
    return result + '<!--/lit-part-->';
  }

  // Fallback
  return '<!--lit-part-->' + escapeHTML(String(val)) + '<!--/lit-part-->';

  } finally { _depth--; }
}

function resolveAttrValue(val, binding) {
  if (val === null || val === undefined) return '';
  if (typeof val === 'symbol') return '';

  // DirectiveResult (e.g., classMap, styleMap, ifDefined)
  if (val && val['_$litDirective$'] !== undefined) {
    return renderDirective(val, {
      type: binding.partType,
      name: binding.name,
      strings: binding.strings,
    });
  }

  if (typeof val === 'string') return val;
  if (typeof val === 'number') return String(val);
  if (typeof val === 'boolean') return val ? '' : null;
  return String(val);
}

// ============================================================
// Directive resolution with proper PartInfo
// ============================================================

function renderDirective(directiveResult, partInfo) {
  const DirectiveClass = directiveResult['_$litDirective$'];
  const values = directiveResult.values || [];

  try {
    // Construct PartInfo matching Lit's interface
    const info = {
      type: partInfo.type || CHILD_PART,
      name: partInfo.name || undefined,
      strings: partInfo.strings || undefined,
      tagName: partInfo.tagName || undefined,
    };

    const instance = new DirectiveClass(info);
    const result = instance.render(...values);

    // Handle directive render results
    if (result === null || result === undefined || typeof result === 'symbol') return '';
    if (typeof result === 'string') return result;
    if (typeof result === 'number') return String(result);

    // Nested TemplateResult from directive
    if (result && result['_$litType$'] !== undefined) {
      const digest = computeDigest(result.strings);
      return `<!--lit-part ${digest}-->` + renderTemplateResult(result) + '<!--/lit-part-->';
    }

    // Nested DirectiveResult
    if (result && result['_$litDirective$'] !== undefined) {
      return renderDirective(result, partInfo);
    }

    return String(result);
  } catch(e) {
    // Fallback for directives that fail: try to extract a sensible value
    if (values.length === 1) {
      const val = values[0];
      if (typeof val === 'string') return partInfo.type === CHILD_PART ? escapeHTML(val) : val;
      if (typeof val === 'number') return String(val);
      if (typeof val === 'object' && val !== null && !Array.isArray(val)) {
        // Likely classMap or styleMap — detect by attribute name
        if (partInfo.name === 'class') {
          return ' ' + Object.entries(val).filter(([_, v]) => v).map(([k]) => k).join(' ') + ' ';
        }
        if (partInfo.name === 'style') {
          return Object.entries(val).filter(([_, v]) => v != null && v !== '').map(([k, v]) => `${k}:${v}`).join(';');
        }
      }
    }
    return '';
  }
}

// ============================================================
// Helper functions
// ============================================================

/**
 * Strip the trailing attribute declaration from a static string.
 * Handles: ` class="`, ` @click="`, ` .value="`, ` ?hidden="`,
 * unquoted variants like ` class=`, and attributes with static value
 * prefixes like ` class="language-` (where "language-" is a static
 * prefix before the first binding expression in the attribute value).
 */
function stripTrailingAttr(str) {
  // Match: whitespace + optional prefix + attribute name + = + optional whitespace
  //        + optional opening quote + optional static value content
  // The key improvement: after the opening quote, allow non-quote characters
  // to handle static prefixes like class="language-${expr}"
  const match = str.match(/^([\s\S]*?)\s+[@.?]?[^\s"'>=/]+=\s*(?:"[^"]*|'[^']*|[^\s"'>]*)$/);
  if (match) return match[1];
  // Also handle attribute at the very start of a string (e.g., after a tag name)
  const match2 = str.match(/^([\s\S]*?)[@.?]?[^\s"'>=/]+=\s*(?:"[^"]*|'[^']*|[^\s"'>]*)$/);
  if (match2 && match2[1].length < str.length) return match2[1];
  return str;
}

// Count nodes visible to Lit's TreeWalker (SHOW_ELEMENT | SHOW_COMMENT).
// This includes element open tags AND HTML comments (<!-- ... -->).
// IMPORTANT: We must skip over comment *contents* because a TreeWalker
// visits a comment as a single node — it never descends into comment text.
// Without this, element-like sequences inside comments (e.g. `<pre>`)
// would be falsely counted as element nodes.
function countNodes(html) {
  let count = 0;
  let i = 0;
  while (i < html.length) {
    if (html[i] === '<') {
      if (i + 3 < html.length && html[i + 1] === '!' && html[i + 2] === '-' && html[i + 3] === '-') {
        count++; // HTML comment: count it as one node
        // Skip past the closing --> so we don't scan comment contents
        const closeIdx = html.indexOf('-->', i + 4);
        if (closeIdx >= 0) {
          i = closeIdx + 3;
        } else {
          i = html.length; // unclosed comment — skip to end
        }
        continue;
      } else if (i + 1 < html.length && html[i + 1] !== '/' && html[i + 1] !== '!') {
        count++; // Element open tag: <div, <span, etc.
      }
    }
    i++;
  }
  return count;
}

function injectLitNodeMarker(html, nodeIndex) {
  let lastTagStart = -1;
  for (let i = html.length - 1; i >= 0; i--) {
    if (html[i] === '<' && i + 1 < html.length && html[i + 1] !== '/' && html[i + 1] !== '!') {
      lastTagStart = i;
      break;
    }
  }
  if (lastTagStart < 0) return html;
  return html.substring(0, lastTagStart) + `<!--lit-node ${nodeIndex}-->` + html.substring(lastTagStart);
}

function escapeHTML(str) {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function escapeAttr(str) {
  if (str === null || str === undefined) return '';
  return String(str).replace(/&/g, '&amp;').replace(/"/g, '&quot;');
}
