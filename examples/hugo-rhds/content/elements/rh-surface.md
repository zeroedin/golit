---
title: "Surface"
imports:
  - rh-blockquote
  - rh-card
  - rh-cta
  - rh-spinner
  - rh-surface
  - rh-tag
lightdom:
  - rh-cta-lightdom-shim.css
---

<p>3 demos for <code>&lt;rh-surface&gt;</code></p>


### color palettes

{{< raw >}}
<rh-surface id="surface" color-palette="darkest">
  <h2>Darkest</h2>
  <rh-cta><a href="#">Call to Action</a></rh-cta>
</rh-surface>
<rh-surface color-palette="darker">
  <h2>Darker</h2>
  <rh-cta><a href="#">Call to Action</a></rh-cta>
</rh-surface>
<rh-surface color-palette="dark">
  <h2>Dark</h2>
  <rh-cta><a href="#">Call to Action</a></rh-cta>
</rh-surface>
<rh-surface color-palette="light">
  <h2>Light</h2>
  <rh-cta><a href="#">Call to Action</a></rh-cta>
</rh-surface>
<rh-surface color-palette="lighter">
  <h2>Lighter</h2>
  <rh-cta><a href="#">Call to Action</a></rh-cta>
</rh-surface>
<rh-surface color-palette="lightest">
  <h2>Lightest</h2>
  <rh-cta><a href="#">Call to Action</a></rh-cta>
</rh-surface>
{{< /raw >}}


### index

{{< raw >}}
<rh-surface id="surface" color-palette="darkest">
  <rh-blockquote align="inline-start" size="default">
    <p>Surface is used to provide a theme to children</p>
  </rh-blockquote>
  <rh-spinner size="lg">Loading...</rh-spinner>
  <rh-tag color="green">Sold</rh-tag>
</rh-surface>
{{< /raw >}}


### nested combination elements

{{< raw >}}
<rh-surface color-palette="darkest">
  <rh-card>
    <p>The card has no color-palette. It's nested CTA should therefore inherit
       context from the grandparent, rh-surface.</p>
    <rh-surface color-palette="light">
      <p>The nested surface should have lighter color <a href="#">even for links</a></p>
      <rh-card color-palette="dark">
        <p>and the nested card should likewise set it's own <a href="#">scheme</a></p>
      </rh-card>
      <rh-cta href="#">light</rh-cta>
    </rh-surface>
    <rh-cta href="#">Should be on dark</rh-cta>
  </rh-card>
</rh-surface>
{{< /raw >}}

