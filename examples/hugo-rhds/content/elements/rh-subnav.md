---
title: "Subnavigation"
imports:
  - rh-icon
  - rh-navigation-link
  - rh-subnav
lightdom:
  - rh-subnav-lightdom.css
---

<p>7 demos for <code>&lt;rh-subnav&gt;</code></p>


### color context

{{< raw >}}
<rh-context-demo>
  <rh-subnav>
    <a href="#">Users</a>
    <a href="#">Containers</a>
    <a href="#">Databases</a>
    <a href="#" aria-current="page">Servers</a>
    <a href="#">System</a>
    <a href="#"><rh-icon slot="icon" icon="data" set="ui" size="sm"></rh-icon> Network</a>
    <a href="#">Cloud</a>
  </rh-subnav>
</rh-context-demo>
{{< /raw >}}


### dynamic

{{< raw >}}
<rh-subnav id="demo" accessible-label="dynamic">
  <a href="#">Users</a>
  <a href="#">Containers</a>
  <a href="#">Databases</a>
  <a href="#" aria-current="page">Servers</a>
  <a href="#">System</a>
  <a href="#">Network</a>
  <a href="#">Cloud</a>
</rh-subnav>
{{< /raw >}}


### index

{{< raw >}}
<rh-subnav>
  <a href="#">Users</a>
  <a href="#">Containers</a>
  <a href="#">Databases</a>
  <a href="#" aria-current="page">Servers</a>
  <a href="#">System</a>
  <a href="#">Network</a>
  <a href="#">Cloud</a>
</rh-subnav>
{{< /raw >}}


### overflow scroll

{{< raw >}}
<div class="resizable">
  <rh-subnav label-scroll-left="Scroll back" label-scroll-right="Scroll forward">
    <a href="#">Users</a>
    <a href="#">Containers</a>
    <a href="#">Databases</a>
    <a href="#" aria-current="page">Servers</a>
    <a href="#">System</a>
    <a href="#">Network</a>
    <a href="#">Cloud</a>
  </rh-subnav>
</div>
{{< /raw >}}


### padded

{{< raw >}}
<h2>In a container with padding</h2>
<div class="padded">
  <rh-subnav>
    <a href="#">Users</a>
    <a href="#">Containers</a>
    <a href="#">Databases</a>
    <a href="#" aria-current="page">Servers</a>
    <a href="#">System</a>
    <a href="#"><rh-icon slot="icon" icon="data" set="ui" size="sm"></rh-icon> Network</a>
    <a href="#">Cloud</a>
  </rh-subnav>
</div>
{{< /raw >}}


### right to left

{{< raw >}}
<div dir="rtl">
  <rh-subnav>
    <a href="#">משתמשים</a>
    <a href="#">מיכלים</a>
    <a href="#">מאגרי מידע</a>
    <a href="#" aria-current="page">שרתים</a>
    <a href="#">מַעֲרֶכֶת</a>
    <a href="#">רֶשֶׁת</a>
    <a href="#">עָנָן</a>
  </rh-subnav>
</div>
{{< /raw >}}


### with navigation link

{{< raw >}}
<rh-subnav>
  <rh-navigation-link href="#">Users</rh-navigation-link>
  <rh-navigation-link href="#">Containers</rh-navigation-link>
  <rh-navigation-link href="#">Databases</rh-navigation-link>
  <rh-navigation-link href="#" current-page>Servers</rh-navigation-link>
  <rh-navigation-link href="#">System</rh-navigation-link>
  <rh-navigation-link href="#">Network</rh-navigation-link>
  <rh-navigation-link href="#">Cloud</rh-navigation-link>
</rh-subnav>
{{< /raw >}}

