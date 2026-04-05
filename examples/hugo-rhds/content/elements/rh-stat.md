---
title: "Statistic"
imports:
  - rh-icon
  - rh-stat
---

<p>7 demos for <code>&lt;rh-stat&gt;</code></p>


### color context

{{< raw >}}
<rh-context-demo>
  <rh-stat>
    <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 36 36">
      <path d="M17.37 8v10a.63.63 0 0 0 1.25 0V8a.63.63 0 0 0-1.25 0Zm7 0v7a.63.63 0 0 0 1.25 0V8a.63.63 0 0 0-1.25 0Zm-14 0v12a.63.63 0 0 0 1.25 0V8a.63.63 0 0 0-1.25 0ZM31 17.89a.63.63 0 0 0-.63.62v11.87H5.62v-2.93a.63.63 0 0 0-1.25 0V31a.63.63 0 0 0 .63.62h26a.62.62 0 0 0 .62-.62V18.51a.62.62 0 0 0-.62-.62Z" />
      <path d="M5 21a.63.63 0 0 0 .62-.63V5.62h24.75V9a.63.63 0 0 0 1.25 0V5a.62.62 0 0 0-.62-.62H5a.63.63 0 0 0-.63.62v15.36A.63.63 0 0 0 5 21Zm27.35-9.24a.62.62 0 0 0-.87.17C28.73 16 21.5 22.93 4 23.27a.63.63 0 0 0 0 1.25c18.07-.34 25.64-7.61 28.52-11.9a.62.62 0 0 0-.17-.86Z" />
    </svg>
    <span slot="statistic">Statistic Placeholder</span>
    <span>Description Placeholder</span>
  </rh-stat>
</rh-context-demo>
{{< /raw >}}


### icon slot

{{< raw >}}
<rh-stat>
  <rh-icon icon="code" slot="icon"></rh-icon>
  <span slot="statistic">Statistic Placeholder</span>
  <span>Description Placeholder</span>
</rh-stat>
{{< /raw >}}


### icon svg

{{< raw >}}
<rh-stat>
  <rh-icon slot="icon" icon="experimental"></rh-icon>
  <span slot="statistic">Statistic Placeholder</span>
  <span>Description Placeholder</span>
</rh-stat>
{{< /raw >}}


### icon

{{< raw >}}
<rh-stat icon="code">
  <span slot="statistic">Statistic Placeholder</span>
  <span>Description Placeholder</span>
</rh-stat>
{{< /raw >}}


### index

{{< raw >}}
<rh-stat>
  <span slot="statistic">Statistic Placeholder</span>
  <span>Description Placeholder</span>
</rh-stat>
{{< /raw >}}


### large

{{< raw >}}
<rh-stat size="large" icon="code">
  <span slot="statistic">Statistic Placeholder</span>
  <span>Description Placeholder</span>
</rh-stat>
{{< /raw >}}


### slotted content

{{< raw >}}
<rh-stat top="statistic">
  <rh-icon slot="icon" icon="code" set="ui"></rh-icon>
  <span slot="title">Overwrite Title</span>
  <p>Stat body that includes two lines and a footnote.</p>
  <span slot="statistic">Overwrite Statistic</span>
</rh-stat>
{{< /raw >}}

