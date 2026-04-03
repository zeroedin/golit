/**
 * SSR rendering helpers loaded once at engine initialization.
 * Used by RenderElement to map HTML attributes to Lit reactive properties.
 */

function attributeToProperty(Ctor, attrName) {
  if (Ctor.__attributeToPropertyMap) {
    return Ctor.__attributeToPropertyMap.get(attrName);
  }
  if (Ctor.elementProperties) {
    for (const [propName, config] of Ctor.elementProperties) {
      const mappedAttr = (config && config.attribute !== undefined)
        ? (config.attribute === false ? null : config.attribute)
        : propName.toLowerCase();
      if (mappedAttr === attrName) return propName;
    }
  }
  return attrName;
}

function getPropertyConfig(Ctor, propName) {
  if (Ctor.elementProperties) {
    return Ctor.elementProperties.get(propName) || {};
  }
  return {};
}

function coerceValue(value, config) {
  const type = config && config.type;
  if (type === Number) return Number(value);
  if (type === Boolean) return value !== 'false';
  return value;
}

function extractStyles(styles) {
  if (!styles) return '';
  if (typeof styles === 'string') return styles;
  if (Array.isArray(styles)) {
    return styles.map(s => extractStyles(s)).filter(Boolean).join('\n');
  }
  if (styles.cssText !== undefined) return styles.cssText;
  if (styles._$cssResult$) return styles.cssText || '';
  return '';
}
