/**
 * Minimal DOM shim for running Lit components in QuickJS.
 * Provides just enough of the DOM API for Lit to:
 * 1. Register custom elements
 * 2. Instantiate components
 * 3. Set attributes
 * 4. Call render()
 *
 * Based on @lit-labs/ssr-dom-shim but stripped to essentials.
 */

// --- EventTarget ---
globalThis.Event = globalThis.Event || class Event {
  constructor(type, options) {
    this.type = type;
    this.bubbles = options?.bubbles || false;
    this.composed = options?.composed || false;
    this.cancelable = options?.cancelable || false;
  }
};

globalThis.CustomEvent = globalThis.CustomEvent || class CustomEvent extends Event {
  constructor(type, options) {
    super(type, options);
    this.detail = options?.detail || null;
  }
};

class EventTarget {
  constructor() {
    this.__listeners = {};
  }
  addEventListener(type, listener) {
    (this.__listeners[type] = this.__listeners[type] || []).push(listener);
  }
  removeEventListener(type, listener) {
    const listeners = this.__listeners[type];
    if (listeners) {
      this.__listeners[type] = listeners.filter(l => l !== listener);
    }
  }
  dispatchEvent(event) {
    const listeners = this.__listeners[event.type] || [];
    for (const l of listeners) l.call(this, event);
    return true;
  }
}

// --- Element ---
class Element extends EventTarget {
  constructor() {
    super();
    this.__attributes = new Map();
    this.__shadowRoot = null;
    this.__shadowRootMode = null;
  }

  get attributes() {
    return Array.from(this.__attributes).map(([name, value]) => ({ name, value }));
  }

  get shadowRoot() {
    if (this.__shadowRootMode === 'closed') return null;
    return this.__shadowRoot;
  }

  get localName() {
    return this.constructor.__localName || '';
  }

  get tagName() {
    return this.localName.toUpperCase();
  }

  setAttribute(name, value) {
    this.__attributes.set(name, String(value));
  }

  getAttribute(name) {
    return this.__attributes.has(name) ? this.__attributes.get(name) : null;
  }

  hasAttribute(name) {
    return this.__attributes.has(name);
  }

  removeAttribute(name) {
    this.__attributes.delete(name);
  }

  toggleAttribute(name, force) {
    if (this.hasAttribute(name)) {
      if (force === undefined || !force) {
        this.removeAttribute(name);
        return false;
      }
    } else {
      if (force === undefined || force) {
        this.setAttribute(name, '');
        return true;
      } else {
        return false;
      }
    }
    return true;
  }

  attachShadow(init) {
    const shadowRoot = { host: this };
    this.__shadowRootMode = init.mode;
    if (init.mode === 'open') {
      this.__shadowRoot = shadowRoot;
    }
    return shadowRoot;
  }

  attachInternals() {
    return {
      role: null,
      ariaLabel: null,
      ariaLabelledByElements: null,
      states: new Set(),
    };
  }
}

// --- HTMLElement ---
class HTMLElement extends Element {
  constructor() {
    super();
  }
}

// --- CustomElementRegistry ---
class CustomElementRegistry {
  constructor() {
    this.__definitions = new Map();
  }

  define(name, ctor) {
    ctor.__localName = name;
    this.__definitions.set(name, {
      ctor,
      observedAttributes: ctor.observedAttributes || [],
    });
  }

  get(name) {
    const def = this.__definitions.get(name);
    return def?.ctor;
  }

  getName(ctor) {
    for (const [name, def] of this.__definitions) {
      if (def.ctor === ctor) return name;
    }
    return null;
  }

  whenDefined(name) {
    const def = this.__definitions.get(name);
    if (def) return Promise.resolve(def.ctor);
    return new Promise(() => {}); // Never resolves for undefined elements
  }
}

// --- Globals ---
// Use ??= so that loading multiple bundles into the same QJS engine
// does not replace the registry (and wipe previously registered elements).
globalThis.EventTarget = EventTarget;
globalThis.Element = Element;
globalThis.HTMLElement = HTMLElement;
globalThis.CustomElementRegistry = CustomElementRegistry;
globalThis.customElements ??= new CustomElementRegistry();
globalThis.document = globalThis.document || {
  nodeType: 9,
  createComment: (data) => ({ data, nodeType: 8, textContent: data }),
  createTextNode: (data) => ({ data, nodeType: 3, textContent: data }),
  createElement: (tag) => {
    const el = new HTMLElement();
    el.__localName = tag;
    return el;
  },
  createDocumentFragment: () => ({
    nodeType: 11,
    childNodes: [],
    firstChild: null,
    appendChild(n) {
      this.childNodes.push(n);
      this.firstChild = this.firstChild || n;
      return n;
    },
  }),
  createTreeWalker: (root, whatToShow, filter) => {
    // Minimal TreeWalker that Lit uses for template parsing.
    // Returns a walker that traverses a simple tree structure.
    const nodes = [];
    function collect(node) {
      if (!node) return;
      const type = node.nodeType || 1;
      const show = (whatToShow & (1 << (type - 1)));
      if (show) nodes.push(node);
      const children = node.childNodes || [];
      for (const child of children) collect(child);
    }
    collect(root);
    let index = -1;
    return {
      currentNode: root,
      nextNode() {
        index++;
        if (index < nodes.length) {
          this.currentNode = nodes[index];
          return this.currentNode;
        }
        return null;
      },
    };
  },
  importNode: (node, deep) => node,
  adoptedStyleSheets: [],
  querySelector: () => null,
  querySelectorAll: () => [],
  addEventListener: () => {},
  removeEventListener: () => {},
  dispatchEvent: () => true,
  head: { appendChild() {} },
  body: { appendChild() {} },
};
globalThis.window = globalThis;
// Ensure window-level event methods exist
if (!globalThis.addEventListener) globalThis.addEventListener = () => {};
if (!globalThis.removeEventListener) globalThis.removeEventListener = () => {};
if (!globalThis.dispatchEvent) globalThis.dispatchEvent = () => true;
try { globalThis.navigator = globalThis.navigator || { userAgent: 'golit-qjs' }; } catch(e) {}
globalThis.Document = globalThis.Document || class Document {};
globalThis.HTMLDocument = globalThis.HTMLDocument || class HTMLDocument {};
globalThis.DOMParser = globalThis.DOMParser || class DOMParser {
  parseFromString() { return globalThis.document; }
};
globalThis.requestAnimationFrame = globalThis.requestAnimationFrame || ((cb) => setTimeout(cb, 0));
globalThis.cancelAnimationFrame = globalThis.cancelAnimationFrame || (() => {});
// setTimeout/setInterval stubs for QJS
if (typeof globalThis.setTimeout === 'undefined') {
  globalThis.setTimeout = (fn) => { fn(); return 0; };
  globalThis.clearTimeout = () => {};
  globalThis.setInterval = () => 0;
  globalThis.clearInterval = () => {};
}
globalThis.ErrorEvent = globalThis.ErrorEvent || class ErrorEvent extends Event {
  constructor(type, init) { super(type); this.message = init?.message || ''; this.error = init?.error || null; }
};
globalThis.MutationObserver = globalThis.MutationObserver || class MutationObserver {
  constructor() {}
  observe() {}
  disconnect() {}
};
globalThis.IntersectionObserver = globalThis.IntersectionObserver || class IntersectionObserver {
  constructor() {}
  observe() {}
  unobserve() {}
  disconnect() {}
};
globalThis.ResizeObserver = globalThis.ResizeObserver || class ResizeObserver {
  constructor() {}
  observe() {}
  disconnect() {}
};
globalThis.CSSStyleSheet = globalThis.CSSStyleSheet || class CSSStyleSheet {
  constructor() { this.cssRules = []; }
  replaceSync() {}
  replace() { return Promise.resolve(this); }
};
globalThis.ShadowRoot = globalThis.ShadowRoot || class ShadowRoot {};
globalThis.DocumentFragment = globalThis.DocumentFragment || class DocumentFragment {};
globalThis.Node = globalThis.Node || class Node {};
globalThis.NodeFilter = globalThis.NodeFilter || { SHOW_COMMENT: 128 };

// Lit checks for adoptedStyleSheets support
try {
  if (globalThis.ShadowRoot && globalThis.ShadowRoot.prototype &&
      !('adoptedStyleSheets' in globalThis.ShadowRoot.prototype)) {
    Object.defineProperty(globalThis.ShadowRoot.prototype, 'adoptedStyleSheets', {
      get() { return []; },
      set() {},
    });
  }
} catch(e) {}

// Lit context root for SSR -- @lit/context checks globalThis.litServerRoot
// when running in server mode and attaches event listeners to it.
globalThis.litServerRoot = globalThis.litServerRoot || {
  addEventListener: () => {},
  removeEventListener: () => {},
  dispatchEvent: () => true,
};

// Dynamic import() shim for preloaded modules.
// When golit pre-loads a module (e.g. prism-esm) into QJS as a script,
// the module's exports become properties on __preloadedModules[name].
// This shim makes import('module-name') resolve to the preloaded exports.
globalThis.__preloadedModules = {};
globalThis.__registerPreload = function(name, exports) {
  globalThis.__preloadedModules[name] = exports;
};

// Lit's isServer check -- set via esbuild's define option.
// No assignment here to avoid esbuild warnings.
