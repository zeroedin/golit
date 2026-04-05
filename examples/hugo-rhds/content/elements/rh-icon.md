---
title: "Icon"
imports:
  - rh-button
  - rh-icon
---

<p>3 demos for <code>&lt;rh-icon&gt;</code></p>


### accessibility

{{< raw >}}
<rh-icon set="ui" icon="arrow-down" accessible-label="Page down"></rh-icon>
<rh-icon icon="info" aria-labelledby="info" role="img"></rh-icon>
<span id="info">Information</span>
{{< /raw >}}


### index

{{< raw >}}
<rh-icon icon="hat"></rh-icon>
{{< /raw >}}


### test remove icon dynamic

{{< raw >}}
<rh-icon icon="hat"></rh-icon>
<rh-button id="remove">Remove Icon</rh-button>
<rh-icon icon="">test</rh-icon>
{{< /raw >}}

