---
title: "Menu"
imports:
  - rh-button
  - rh-menu
---

<p>6 demos for <code>&lt;rh-menu&gt;</code></p>


### color context

<rh-context-demo>
  <rh-menu>
    <a href="#">Link1</a>
    <a href="#">Link2</a>
    <a href="#">Link3</a>
    <a href="#">Link4</a>
  </rh-menu>
</rh-context-demo>


### index

<rh-menu id="rh-buttons">
  <rh-button data-item="1" variant="link">Menuitem1</rh-button>
  <rh-button data-item="2" variant="link">Menuitem2</rh-button>
  <rh-button data-item="3" variant="link">Menuitem3</rh-button>
  <rh-button data-item="4" variant="link">Menuitem4</rh-button>
</rh-menu>


### menu item

<rh-menu>
  <rh-menu-item>Menuitem1</rh-menu-item>
  <rh-menu-item>Menuitem2</rh-menu-item>
  <rh-menu-item>Menuitem3</rh-menu-item>
</rh-menu>


### position left

<rh-menu position="left" persistent>
  <rh-button variant="link">Menuitem1</rh-button>
  <rh-button variant="link">Menuitem2</rh-button>
  <rh-button variant="link">Menuitem3</rh-button>
  <rh-button variant="link">Menuitem4</rh-button>
</rh-menu>


### position right

<rh-menu position="right">
  <a href="#">Link1</a>
  <a href="#">Link2</a>
  <a href="#">Link3</a>
  <a href="#">Link4</a>
</rh-menu>


### position top

<rh-menu position="top">
  <rh-button variant="link">Menuitem1</rh-button>
  <rh-button variant="link">Menuitem2</rh-button>
  <rh-button variant="link">Menuitem3</rh-button>
  <rh-button variant="link">Menuitem4</rh-button>
</rh-menu>

