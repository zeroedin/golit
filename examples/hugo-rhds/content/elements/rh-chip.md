---
title: "Chip"
imports:
  - rh-button
  - rh-chip
---

<p>9 demos for <code>&lt;rh-chip&gt;</code></p>


### attributes and states

{{< raw >}}
<rh-chip-group accessible-label="Various attributes and states:">
  <rh-chip value="checked" checked>I am checked</rh-chip>
  <rh-chip value="i-am-disabled" disabled>I am disabled</rh-chip>
  <rh-chip value="devops">Custom value unchecked</rh-chip>
  <rh-chip value="i-am-disabled-checked" disabled checked>I am disabled &amp; checked</rh-chip>
</rh-chip-group>
{{< /raw >}}


### chip group

{{< raw >}}
<rh-chip-group>
  <rh-chip>Edge</rh-chip>
  <rh-chip checked>AI/ML</rh-chip>
  <rh-chip>DevOps</rh-chip>
</rh-chip-group>
{{< /raw >}}


### clear all

{{< raw >}}
<rh-chip-group accessible-label="Custom filter text">
  <rh-chip checked>Bluetooth</rh-chip>
  <rh-chip>Wi-Fi</rh-chip>
  <rh-chip>RFID</rh-chip>
  <rh-chip checked disabled>Chip and pin</rh-chip>
  <span slot="clear-all">Clear technologies</span>
</rh-chip-group>
{{< /raw >}}


### color context

{{< raw >}}
<rh-context-demo>
  <rh-chip-group>
    <rh-chip>Edge</rh-chip>
    <rh-chip checked>AI/ML</rh-chip>
    <rh-chip>DevOps</rh-chip>
    <rh-chip disabled>Disabled</rh-chip>
  </rh-chip-group>
</rh-context-demo>
{{< /raw >}}


### custom label

{{< raw >}}
<rh-chip-group accessible-label="Filter technologies:">
  <rh-chip>Edge</rh-chip>
  <rh-chip checked>AI/ML</rh-chip>
  <rh-chip>OpenShift</rh-chip>
  <rh-chip disabled>COBOL</rh-chip>
</rh-chip-group>
{{< /raw >}}


### events

{{< raw >}}
<rh-chip-group>
  <rh-chip value="custom-value">Check me</rh-chip>
</rh-chip-group>
{{< /raw >}}


### form control

{{< raw >}}
<form id="chip-form">
  <rh-chip-group>
    <rh-chip name="edge">Edge</rh-chip>
    <rh-chip name="ai" checked>AI/ML</rh-chip>
    <rh-chip name="devops">DevOps</rh-chip>
  </rh-chip-group>
  <hr>
  <rh-button type="submit">Submit</rh-button>
  <output name="output"></output>
</form>
{{< /raw >}}


### index

{{< raw >}}
<rh-chip>Chip</rh-chip>
{{< /raw >}}


### size

{{< raw >}}
<rh-chip-group size="sm">
  <span slot="accessible-label">Filter by (small):</span>
  <rh-chip>Automation</rh-chip>
  <rh-chip checked>Security</rh-chip>
  <rh-chip>Containers</rh-chip>
</rh-chip-group>
<rh-chip-group>
  <rh-chip>Open Source</rh-chip>
  <rh-chip checked>ARO</rh-chip>
  <rh-chip>RHEL</rh-chip>
</rh-chip-group>
{{< /raw >}}

