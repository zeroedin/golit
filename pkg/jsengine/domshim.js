/**
 * DOM shim for running Lit components in QuickJS.
 * Event/EventTarget modeled on @lit-labs/ssr-dom-shim/lib/events.js
 * for full compatibility with @lit/context, controllers, and Lit lifecycle.
 *
 * Based on @lit-labs/ssr-dom-shim but adapted for QuickJS.
 */

// --- Event phases ---
const _NONE = 0;
const _CAPTURING_PHASE = 1;
const _AT_TARGET = 2;
const _BUBBLING_PHASE = 3;

const _isCaptureOption = (options) =>
  typeof options === 'boolean' ? options : !!(options?.capture);

// --- Event ---
globalThis.Event = globalThis.Event || class Event {
  constructor(type, options) {
    if (arguments.length === 0) throw new Error('The type argument must be specified');
    const opts = (typeof options === 'object' && options) ? options : {};
    this._type = String(type);
    this._bubbles = !!opts.bubbles;
    this._composed = !!opts.composed;
    this._cancelable = !!opts.cancelable;
    this._defaultPrevented = false;
    this._propagationStopped = false;
    this._immediatePropagationStopped = false;
    this._target = null;
    this._currentTarget = null;
    this._eventPhase = _NONE;
    this._timestamp = Date.now();
    this._isBeingDispatched = false;

    this.NONE = _NONE;
    this.CAPTURING_PHASE = _CAPTURING_PHASE;
    this.AT_TARGET = _AT_TARGET;
    this.BUBBLING_PHASE = _BUBBLING_PHASE;
  }

  get type() { return this._type; }
  get bubbles() { return this._bubbles; }
  get composed() { return this._composed; }
  get cancelable() { return this._cancelable; }
  get defaultPrevented() { return this._cancelable && this._defaultPrevented; }
  get timeStamp() { return this._timestamp; }
  get target() { return this._target; }
  get currentTarget() { return this._currentTarget; }
  get srcElement() { return this._target; }
  get isTrusted() { return false; }
  get returnValue() { return !this._cancelable || !this._defaultPrevented; }

  get eventPhase() {
    return this._isBeingDispatched ? this._eventPhase : _NONE;
  }

  get cancelBubble() { return this._propagationStopped; }
  set cancelBubble(value) { if (value) this._propagationStopped = true; }

  composedPath() {
    return this._isBeingDispatched ? [this._target] : [];
  }

  stopPropagation() {
    this._propagationStopped = true;
  }

  stopImmediatePropagation() {
    this._propagationStopped = true;
    this._immediatePropagationStopped = true;
  }

  preventDefault() {
    this._defaultPrevented = true;
  }
};
globalThis.Event.NONE = _NONE;
globalThis.Event.CAPTURING_PHASE = _CAPTURING_PHASE;
globalThis.Event.AT_TARGET = _AT_TARGET;
globalThis.Event.BUBBLING_PHASE = _BUBBLING_PHASE;

// --- CustomEvent ---
globalThis.CustomEvent = globalThis.CustomEvent || class CustomEvent extends Event {
  constructor(type, options) {
    super(type, options);
    this._detail = (typeof options === 'object' && options) ? (options.detail ?? null) : null;
  }
  get detail() { return this._detail; }
};

// --- EventTarget ---
class EventTarget {
  constructor() {
    this.__eventListeners = new Map();
    this.__captureEventListeners = new Map();
  }

  addEventListener(type, callback, options) {
    if (callback == null) return;
    const listenerMap = _isCaptureOption(options)
      ? this.__captureEventListeners
      : this.__eventListeners;
    let listeners = listenerMap.get(type);
    if (listeners === undefined) {
      listeners = new Map();
      listenerMap.set(type, listeners);
    } else if (listeners.has(callback)) {
      return;
    }
    const normalizedOpts = (typeof options === 'object' && options) ? options : {};
    normalizedOpts.signal?.addEventListener?.('abort', () =>
      this.removeEventListener(type, callback, options));
    listeners.set(callback, normalizedOpts);
  }

  removeEventListener(type, callback, options) {
    if (callback == null) return;
    const listenerMap = _isCaptureOption(options)
      ? this.__captureEventListeners
      : this.__eventListeners;
    const listeners = listenerMap.get(type);
    if (listeners !== undefined) {
      listeners.delete(callback);
      if (!listeners.size) listenerMap.delete(type);
    }
  }

  dispatchEvent(event) {
    const composedPath = [this];
    let parent = this.__eventTargetParent;
    if (event.composed) {
      while (parent) {
        composedPath.push(parent);
        parent = parent.__eventTargetParent;
      }
    } else {
      while (parent && parent !== this.__host) {
        composedPath.push(parent);
        parent = parent.__eventTargetParent;
      }
    }

    let stopProp = false;
    let stopImmediate = false;
    let eventPhase = _NONE;
    let target = null;
    let tmpTarget = null;
    let currentTarget = null;

    const origStop = event.stopPropagation.bind(event);
    const origImmediate = event.stopImmediatePropagation.bind(event);

    Object.defineProperties(event, {
      target: { get() { return target ?? tmpTarget; }, configurable: true, enumerable: true },
      srcElement: { get() { return target ?? tmpTarget; }, configurable: true, enumerable: true },
      currentTarget: { get() { return currentTarget; }, configurable: true, enumerable: true },
      eventPhase: { get() { return eventPhase; }, configurable: true, enumerable: true },
      composedPath: { value: () => composedPath, configurable: true, enumerable: true },
      stopPropagation: {
        value: () => { stopProp = true; origStop(); },
        configurable: true, enumerable: true,
      },
      stopImmediatePropagation: {
        value: () => { stopImmediate = true; origImmediate(); },
        configurable: true, enumerable: true,
      },
    });

    event._isBeingDispatched = true;

    const invoke = (listener, opts, listenersMap) => {
      if (typeof listener === 'function') {
        listener(event);
      } else if (typeof listener?.handleEvent === 'function') {
        listener.handleEvent(event);
      }
      if (opts.once) listenersMap.delete(listener);
    };

    const finish = () => {
      currentTarget = null;
      eventPhase = _NONE;
      event._isBeingDispatched = false;
      return !event.defaultPrevented;
    };

    // Retarget event.target across shadow boundaries.
    target = (!this.__host || !event.composed) ? this : null;
    const retarget = (eventTargets) => {
      tmpTarget = this;
      while (tmpTarget.__host && eventTargets.includes(tmpTarget.__host)) {
        tmpTarget = tmpTarget.__host;
      }
    };

    // Capture phase (root -> target)
    const capturePath = composedPath.slice().reverse();
    for (const et of capturePath) {
      if (!target && (!tmpTarget || tmpTarget === et.__host)) {
        retarget(capturePath.slice(capturePath.indexOf(et)));
      }
      currentTarget = et;
      eventPhase = (et === (target ?? tmpTarget)) ? _AT_TARGET : _CAPTURING_PHASE;
      const listeners = et.__captureEventListeners?.get(event.type);
      if (listeners) {
        for (const [listener, opts] of listeners) {
          invoke(listener, opts, listeners);
          if (stopImmediate) return finish();
        }
      }
      if (stopProp) return finish();
    }

    // Bubble phase (target -> root), or just [this] if non-bubbling.
    const bubblePath = event.bubbles ? composedPath : [this];
    tmpTarget = null;
    for (const et of bubblePath) {
      if (!target && (!tmpTarget || et === tmpTarget.__host)) {
        retarget(bubblePath.slice(0, bubblePath.indexOf(et) + 1));
      }
      currentTarget = et;
      eventPhase = (et === (target ?? tmpTarget)) ? _AT_TARGET : _BUBBLING_PHASE;
      const listeners = et.__eventListeners?.get(event.type);
      if (listeners) {
        for (const [listener, opts] of listeners) {
          invoke(listener, opts, listeners);
          if (stopImmediate) return finish();
        }
      }
      if (stopProp) return finish();
    }

    return finish();
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
    if (this.__definitions.has(name)) return;
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
// does not replace classes or the registry (and wipe previously registered elements).
globalThis.EventTarget ??= EventTarget;
globalThis.Element ??= Element;
globalThis.HTMLElement ??= HTMLElement;
globalThis.CustomElementRegistry ??= CustomElementRegistry;
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
  constructor(type, init) { super(type, init); this.message = init?.message || ''; this.error = init?.error || null; }
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
// Must be a real EventTarget so ContextProvider/ContextRoot can dispatch on it.
globalThis.litServerRoot = globalThis.litServerRoot || (() => {
  const root = new HTMLElement();
  Object.defineProperty(root, 'localName', { get() { return 'lit-server-root'; } });
  return root;
})();

// Dynamic import() shim for preloaded modules.
// When golit pre-loads a module (e.g. prism-esm) into QJS as a script,
// the module's exports become properties on __preloadedModules[name].
// This shim makes import('module-name') resolve to the preloaded exports.
globalThis.__preloadedModules = {};
globalThis.__registerPreload = function(name, exports) {
  globalThis.__preloadedModules[name] = exports;
};

// --- SSR URL / location (align with @lit/ssr getWindow: stable base URL) ---
// matchMedia is intentionally not shimmed (viewport-specific; no server semantics).

(function () {
  const base = typeof globalThis.__golitLocationHref === 'string'
    ? globalThis.__golitLocationHref
    : 'http://localhost/';
  try {
    globalThis.location ??= new URL(base);
  } catch (_) {
    globalThis.location ??= {
      href: base,
      origin: '',
      protocol: 'http:',
      host: 'localhost',
      hostname: 'localhost',
      port: '',
      pathname: '/',
      search: '',
      hash: '',
      assign() {},
      replace() {},
      reload() {},
      toString() { return this.href; },
    };
  }
})();

globalThis.URLSearchParams ??= class URLSearchParams {
  constructor(init) {
    this._map = new Map();
    if (init == null || init === '') return;
    if (typeof init === 'string') {
      let s = init.startsWith('?') ? init.slice(1) : init;
      for (const pair of s.split('&')) {
        if (!pair) continue;
        const i = pair.indexOf('=');
        const k = i >= 0 ? pair.slice(0, i) : pair;
        const v = i >= 0 ? pair.slice(i + 1) : '';
        try {
          this._map.set(decodeURIComponent(k), decodeURIComponent(v));
        } catch (_) {
          this._map.set(k, v);
        }
      }
    } else if (typeof init === 'object') {
      for (const key of Object.keys(init)) {
        this._map.set(key, String(init[key]));
      }
    }
  }
  get(k) { return this._map.has(k) ? this._map.get(k) : null; }
  set(k, v) { this._map.set(String(k), String(v)); }
  has(k) { return this._map.has(k); }
  append(k, v) {
    const key = String(k);
    const cur = this._map.get(key);
    this._map.set(key, cur == null ? String(v) : cur + ',' + String(v));
  }
  toString() {
    const parts = [];
    for (const [k, v] of this._map) {
      parts.push(encodeURIComponent(k) + '=' + encodeURIComponent(v));
    }
    return parts.join('&');
  }
};

if (typeof globalThis.URL === 'undefined') {
  globalThis.URL = class URL {
    constructor(url, base) {
      let s = String(url);
      if (base != null && base !== '') {
        const b = String(base).replace(/\/+$/, '');
        s = s.startsWith('/') ? b + s : b + '/' + s.replace(/^\/+/, '');
      }
      this.href = s;
    }
    toString() { return this.href; }
    get pathname() {
      const q = this.href.indexOf('?');
      const h = this.href.indexOf('#');
      const end = q >= 0 && h >= 0 ? Math.min(q, h) : (q >= 0 ? q : (h >= 0 ? h : this.href.length));
      const start = this.href.indexOf('/', this.href.indexOf('//') + 2);
      if (start < 0) return '/';
      const p = this.href.slice(start, end);
      return p || '/';
    }
    get search() {
      const q = this.href.indexOf('?');
      if (q < 0) return '';
      const h = this.href.indexOf('#', q);
      return h >= 0 ? this.href.slice(q, h) : this.href.slice(q);
    }
    get hash() {
      const h = this.href.indexOf('#');
      return h >= 0 ? this.href.slice(h) : '';
    }
    get host() {
      const m = this.href.match(/^https?:\/\/([^/?#]+)/i);
      return m ? m[1] : '';
    }
    get hostname() { return this.host.split(':')[0]; }
    get port() {
      const h = this.host;
      const i = h.lastIndexOf(':');
      return i > 0 && /^\d+$/.test(h.slice(i + 1)) ? h.slice(i + 1) : '';
    }
    get protocol() {
      const m = this.href.match(/^([a-z]+:)/i);
      return m ? m[1].toLowerCase() : 'http:';
    }
    get origin() {
      const m = this.href.match(/^(https?:\/\/[^/?#]+)/i);
      return m ? m[1] : '';
    }
  };
}

(function () {
  if (typeof globalThis.__golitFetch !== 'function') return;
  globalThis.fetch ??= function (input, init) {
    const url = typeof input === 'string'
      ? input
      : (input && typeof input.url === 'string')
        ? input.url
        : String(input);
    const o = init && typeof init === 'object' ? init : {};
    const headersObj = o.headers;
    const headers = {};
    if (headersObj && typeof headersObj === 'object') {
      if (typeof headersObj.forEach === 'function') {
        headersObj.forEach((v, k) => { headers[String(k)] = String(v); });
      } else {
        for (const k of Object.keys(headersObj)) {
          headers[k] = String(headersObj[k]);
        }
      }
    }
    let bodyStr;
    if (o.body != null) {
      bodyStr = typeof o.body === 'string' ? o.body : String(o.body);
    }
    const initJson = JSON.stringify({
      method: o.method,
      headers,
      body: bodyStr,
    });
    let raw;
    try {
      raw = globalThis.__golitFetch(url, initJson);
    } catch (e) {
      return Promise.reject(e);
    }
    let d;
    try {
      d = JSON.parse(raw);
    } catch (e) {
      return Promise.reject(new Error('__golitFetch returned invalid JSON'));
    }
    if (d.error) {
      return Promise.reject(new Error(d.error));
    }
    const textBody = d.body != null ? String(d.body) : '';
    return Promise.resolve({
      ok: !!d.ok,
      status: d.status | 0,
      statusText: d.statusText || '',
      text() { return Promise.resolve(textBody); },
      json() {
        try {
          return Promise.resolve(JSON.parse(textBody));
        } catch (e) {
          return Promise.reject(e);
        }
      },
      headers: {
        get() { return null; },
      },
    });
  };
})();

// Lit's isServer check -- set via esbuild's define option.
// No assignment here to avoid esbuild warnings.
