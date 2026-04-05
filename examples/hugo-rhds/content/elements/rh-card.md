---
title: "Card"
imports:
  - rh-blockquote
  - rh-button
  - rh-card
  - rh-cta
lightdom:
  - rh-cta-lightdom-shim.css
---

<p>11 demos for <code>&lt;rh-card&gt;</code></p>


### color context

{{< raw >}}
<rh-context-demo>
  <rh-card>
    <h2 slot="header">Default card, slotted content and footer</h2>
    <p>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est egestas, a
      sollicitudin mauris tincidunt. Pellentesque vel dapibus risus. Nullam aliquam felis orci, eget cursus mi
      lacinia quis. Vivamus at felis sem.
    </p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card id="custom">
    <h2 slot="header">Slotted title, content, and footer</h2>
    <label for="picker">Change this card's color-palette</label>
    <rh-context-picker id="picker" target="custom"></rh-context-picker>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="full">
    <svg slot="header" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 120">
      <title>A placeholder image in a card header slot</title>
      <rect x="0" y="0" width="320" height="120" fill="light-dark(var(--rh-color-surface-light, #e0e0e0), var(--rh-color-surface-dark, #383838))" />
      <line x1="0" y1="0" x2="320" y2="120" stroke="var(--rh-color-border-subtle)" />
      <line x1="320" y1="0" x2="0" y2="120" stroke="var(--rh-color-border-subtle)" />
    </svg>
    <h2 slot="header">Card with slotted image header. Full width image.</h2>
    <p>
      lorem ipsum dolor sit amet, consectetur adipiscing elit. nullam eleifend elit sed est egestas, a
      sollicitudin mauris
      tincidunt. pellentesque vel dapibus risus. nullam aliquam felis orci, eget cursus mi lacinia quis.
      vivamus
      at felis
      sem.
    </p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="bar">
    <h2 slot="header">Custom header</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="bar">
    <h2 slot="header">Custom header</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
</rh-context-demo>
{{< /raw >}}


### grid

{{< raw >}}
<div id="card-grid">
  <rh-card>
    <h2 slot="header">Grid Card</h2>
    <p>In a grid, cards will fill all the available vertical space.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card>
    <h2 slot="header">Longer Content</h2>
    <p>
      In such a case, all of the grid card's footers should be
      vertically aligned. Meaning, they should always rest in
      the bottom of their card. Even when one card has much more
      content than its neighbour, and thus fill more vertical space
      in it's body, the footers should still be aligned.
    </p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card>
    <h2 slot="header">Grid Card</h2>
    <p>
      These kinds of situations should best be avoided. See the
      guidelines for more information.
    </p>
    <rh-cta variant="primary" slot="footer">
      <a href="https://ux.redhat.com/elements/card/guidelines/#vertical-height">Read the Guidelines</a>
    </rh-cta>
  </rh-card>
  <rh-card>
    <p>This card has no header and short body content.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card>
    <h2 slot="header">No footer, More Content</h2>
    <p>
      In such a case, all of the grid card's footers should be
      vertically aligned. Meaning, they should always rest in
      the bottom of their card. Even when one card has much more
      content than its neighbour, and thus fill more vertical space
      in it's body, the footers should still be aligned.
    </p>
  </rh-card>
  <rh-card>
    <h2 slot="header">No body</h2>
    <rh-cta variant="primary" slot="footer">
      <a href="https://ux.redhat.com/elements/card/guidelines/#vertical-height">Read the Guidelines</a>
    </rh-cta>
  </rh-card>
</div>
{{< /raw >}}


### index

{{< raw >}}
<rh-card>
  <h2 slot="header">Card</h2>
  <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
     Nullam eleifend elit sed est egestas, a sollicitudin mauris
     tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
     felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
  <rh-cta slot="footer" priority="primary">
    <a href="#">Call to action</a>
  </rh-cta>
</rh-card>
{{< /raw >}}


### promo full width

{{< raw >}}
<div id="full-width">
  <rh-card variant="promo" full-width>
    <svg slot="image"
        role="img"
        aria-label="sample image"
        width="1920" height="158">
      <rect fill="var(--rh-color-border-interactive)"
            fill-opacity="0.1"
            stroke="var(--rh-color-border-interactive)"
            stroke-width="2px"
            width="100%"
            height="100%"
            stroke-dasharray="4 4">
    </svg>
    <h2 slot="header">Full Width Promo</h2>
    <p>Promos which span the entire viewport do not have a border.</p>
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
  </rh-card>
  <rh-card variant="promo" full-width color-palette="darkest">
    <svg slot="image"
        role="img"
        aria-label="sample image"
        width="1920" height="158">
      <rect fill="var(--rh-color-border-interactive)"
            fill-opacity="0.1"
            stroke="var(--rh-color-border-interactive)"
            stroke-width="2px"
            width="100%"
            height="100%"
            stroke-dasharray="4 4">
    </svg>
    <h2 slot="header">Full Width Promo Dark</h2>
    <p>Promos can have a dark color palette too. You can add the <code>color-palette</code> attribute or wrap them with <code>rh-surface</code>.</p>
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
  </rh-card>
  <rh-card variant="promo" full-width class="reverse">
    <svg slot="image"
        role="img"
        aria-label="sample image"
        width="1920" height="158">
      <rect fill="var(--rh-color-border-interactive)"
            fill-opacity="0.1"
            stroke="var(--rh-color-border-interactive)"
            stroke-width="2px"
            width="100%"
            height="100%"
            stroke-dasharray="4 4">
    </svg>
    <h2 slot="header">Full Width Promo - reverse</h2>
    <p>Using the <code>container</code> and <code>image</code> parts to reverse and style the Promo.</p>
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
  </rh-card>
  <rh-card variant="promo" full-width>
    <h2 slot="header">Full Width Promo - No Image</h2>
    <p>This Promo is just like the last version, except this version doesn't have an
       image in the <code>image</code> slot. Without an image, this turns into a single column component.</p>
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
  </rh-card>
</div>
{{< /raw >}}


### promo narrow

{{< raw >}}
<section id="narrow">
  <rh-card variant="promo">
    <svg slot="image"
        role="img"
        aria-label="sample image"
        width="360" height="158">
      <rect fill="var(--rh-color-border-interactive)"
            fill-opacity="0.1"
            stroke="var(--rh-color-border-interactive)"
            stroke-width="2px"
            width="100%"
            height="100%"
            stroke-dasharray="4 4">
    </svg>
    <h2 slot="header">Narrow Promo</h2>
    <p>Promos narrower than 296px prioritize text by moving the image to the end.</p>
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
  </rh-card>
  <rh-card variant="promo">
    <h2 slot="header">Narrow Promo</h2>
    <p>A promo can optionally omit the image.</p>
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
  </rh-card>
  <rh-card variant="promo">
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
    <p>The &rdquo;standard&ldquo; Promo. It has a paragraph and a CTA in the &rdquo;image&ldquo; slot.</p>
  </rh-card>
  <rh-card variant="promo" color-palette="darkest">
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
    <p>A &rdquo;standard&ldquo; Promo with a <code>color-palette</code> attribute added to the <code>&lt;rh-card&gt;</code> tag.</p>
  </rh-card>
</section>
{{< /raw >}}


### promo standard

{{< raw >}}
<section id="standard">
  <rh-card variant="promo" color-palette="lighter">
    <p>The &rdquo;standard&ldquo; Promo. It has a paragraph in the default slot and a CTA in the footer slot.</p>
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
  </rh-card>
  <rh-card variant="promo" color-palette="darker" class="reverse">
    <p>In this &rdquo;standard&ldquo; Promo, the slots can be reversed using a class and the <code>::part()</code> selector.</p>
    <rh-cta slot="footer" href="#">Call to action</rh-cta>
  </rh-card>
</section>
{{< /raw >}}


### promo

{{< raw >}}
<rh-card variant="promo">
  <svg slot="image"
       width="1920"
       height="250"
       role="img"
       aria-label="sample image">
    <rect fill="var(--rh-color-border-interactive)"
      fill-opacity="0.1"
      stroke="var(--rh-color-border-interactive)"
      stroke-width="2px"
      width="100%"
      height="100%"
      stroke-dasharray="4 4">
  </svg>
  <h2 slot="header">Promo</h2>
  <p>Featured promo card has an image in the <code>image</code> slot,
    a heading in the <code>header</code> slot, and a <abbr title="call to action">CTA</abbr>
    in the <code>footer</code> slot, and body content.</p>
  <rh-cta slot="footer" href="#">Call to action</rh-cta>
</rh-card>
<rh-card variant="promo" class="reverse">
  <svg slot="image"
       width="1920"
       height="250"
       role="img"
       aria-label="sample image">
    <rect fill="var(--rh-color-border-interactive)"
      fill-opacity="0.1"
      stroke="var(--rh-color-border-interactive)"
      stroke-width="2px"
      width="100%"
      height="100%"
      stroke-dasharray="4 4">
  </svg>
  <h2 slot="header">Promo reversed</h2>
  <p>By selecting the <code>image</code>and <code>body</code> CSS Shadow Parts, you can reverse the <code>image</code> and <code>body</code> shadow parts and adjust column widths.</p>
  <rh-cta slot="footer" href="#">Call to action</rh-cta>
</rh-card>
<rh-card variant="promo">
  <h2 slot="header">Featured: no image is okay too</h2>
  <p>Sometimes, you may not have an image to go with your card. By omitting the slotted image,
    the promo reverts to a single column.</p>
  <rh-cta slot="footer" href="#">Call to action</rh-cta>
</rh-card>
<rh-card variant="promo" color-palette="darkest">
  <svg slot="image"
       width="1920"
       height="250"
       role="img"
       aria-label="sample image">
    <rect fill="var(--rh-color-border-interactive)"
      fill-opacity="0.1"
      stroke="var(--rh-color-border-interactive)"
      stroke-width="2px"
      width="100%"
      height="100%"
      stroke-dasharray="4 4">
  </svg>
  <h2 slot="header">Promo Darkest</h2>
  <p>Promos can have a darkest color palette too. You can add the <code>color-palette</code> attribute or wrap them with <code>rh-surface</code>.</p>
  <rh-cta slot="footer" href="#">Call to action</rh-cta>
</rh-card>
{{< /raw >}}


### sticky pattern

{{< raw >}}
<section id="card-pattern-sticky">
  <rh-card>
    <rh-button slot="header" variant="close">Close</rh-button>
    <h2 slot="header" class="title">Title, lg</h2>
    <h2 slot="header">Heading, xs</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
       Nullam eleifend elit sed est egestas, a sollicitudin mauris tincidunt.
       Pellentesque vel dapibus risus.</p>
    <rh-cta slot="footer"><a href="#">Call to action</a></rh-cta>
  </rh-card>
  <rh-card color-palette="lighter">
    <rh-button slot="header" variant="close">Close</rh-button>
    <h2 slot="header" class="title">Title, lg</h2>
    <h2 slot="header">Heading, xs</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
       Nullam eleifend elit sed est egestas, a sollicitudin mauris tincidunt.
       Pellentesque vel dapibus risus.</p>
    <rh-cta slot="footer"><a href="#">Call to action</a></rh-cta>
  </rh-card>
</section>
{{< /raw >}}


### title pattern

{{< raw >}}
<section id="card-pattern-title">
  <rh-card>
    <h2 slot="header">Title, lg</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
       Nullam eleifend elit sed est egestas, a sollicitudin mauris tincidunt.
       Pellentesque vel dapibus risus.</p>
    <rh-cta><a href="#">Call to action</a></rh-cta>
  </rh-card>
  <rh-card>
    <h2 slot="header">Title, lg</h2>
    <svg class="sample" width="160" height="80"><rect/></svg>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
       Nullam eleifend elit sed est egestas, a sollicitudin mauris tincidunt.
       Pellentesque vel dapibus risus.</p>
    <rh-cta><a href="#">Call to action</a></rh-cta>
  </rh-card>
  <rh-card color-palette="lighter">
    <h2 slot="header">Title, lg</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
       Nullam eleifend elit sed est egestas, a sollicitudin mauris tincidunt.
       Pellentesque vel dapibus risus.</p>
    <rh-cta><a href="#">Call to action</a></rh-cta>
  </rh-card>
  <rh-card color-palette="lighter">
    <h2 slot="header">Title, lg</h2>
    <svg class="sample" width="160" height="80"><rect/></svg>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
       Nullam eleifend elit sed est egestas, a sollicitudin mauris tincidunt.
       Pellentesque vel dapibus risus.</p>
    <rh-cta><a href="#">Call to action</a></rh-cta>
  </rh-card>
</section>
{{< /raw >}}


### variants

{{< raw >}}
<section id="card-patterns">
  <rh-card>
    <h2 slot="header">Default card, slotted content and footer</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card>
    <h2 slot="header">Slotted title, content, and footer</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card color-palette="lighter">
    <h2 slot="header">Lighter color palette</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card color-palette="lighter" class="bar">
    <h2 slot="header">Lighter color palette and custom header</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="bar">
    <h2 slot="header">Custom header</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="center">
    <h2 slot="header">Center aligned content, footer</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="end">
    <h2 slot="header">End aligned content, footer</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="image">
    <svg slot="header" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 120">
        <title>A placeholder image in a card header slot</title>
        <rect x="0" y="0" width="320" height="120" fill="light-dark(var(--rh-color-surface-light, #e0e0e0), var(--rh-color-surface-dark, #383838))" />
        <line x1="0" y1="0" x2="320" y2="120" stroke="var(--rh-color-border-subtle)" />
        <line x1="320" y1="0" x2="0" y2="120" stroke="var(--rh-color-border-subtle)" />
      </svg>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="image">
    <svg slot="header" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 120">
        <title>A placeholder image in a card header slot</title>
        <rect x="0" y="0" width="320" height="120" fill="light-dark(var(--rh-color-surface-light, #e0e0e0), var(--rh-color-surface-dark, #383838))" />
        <line x1="0" y1="0" x2="320" y2="120" stroke="var(--rh-color-border-subtle)" />
        <line x1="320" y1="0" x2="0" y2="120" stroke="var(--rh-color-border-subtle)" />
      </svg>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
  <rh-card class="image heading">
    <svg slot="header" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 120">
      <title>A placeholder image in a card header slot</title>
      <rect x="0" y="0" width="320" height="120" fill="light-dark(var(--rh-color-surface-light, #e0e0e0), var(--rh-color-surface-dark, #383838))" />
      <line x1="0" y1="0" x2="320" y2="120" stroke="var(--rh-color-border-subtle)" />
      <line x1="320" y1="0" x2="0" y2="120" stroke="var(--rh-color-border-subtle)" />
    </svg>
    <h2>Card with slotted image header. Full width image.</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      Nullam eleifend elit sed est egestas, a sollicitudin mauris
      tincidunt. Pellentesque vel dapibus risus. Nullam aliquam
      felis orci, eget cursus mi lacinia quis. Vivamus at felis sem.</p>
    <rh-cta variant="primary" slot="footer">
      <a href="#">Call to action</a>
    </rh-cta>
  </rh-card>
</section>
{{< /raw >}}


### video pattern

{{< raw >}}
<section id="card-pattern-quote">
  <rh-card>
    <rh-video slot="image"></rh-video>
    <a href="#">Lorem ipsum dolor sit amet consectetur adipisicing elit.</a>
  </rh-card>
  <rh-card>
    <rh-video slot="image"></rh-video>
    <p>Lorem ipsum dolor sit amet consectetur adipisicing elit.
       Quaerat quae accusantium velit sed amet, praesentium maiores illum ad.
       Odio vero molestiae sint animi. Vero in ad fugit sit, possimus explicabo.</p>
    <rh-cta slot="footer"><a href="#">Call to action</a></rh-cta>
  </rh-card>
  <rh-card>
    <rh-video slot="image"></rh-video>
    <rh-blockquote>
      <p>Lorem ipsum dolor sit amet consectetur adipisicing elit.
         Quaerat quae accusantium velit sed amet, praesentium maiores illum ad.
         Odio vero molestiae sint animi. Vero in ad fugit sit, possimus explicabo.</p>
      <span slot="author">First name Last name</span>
      <span slot="title">Job title, Company name</span>
    </rh-blockquote>
  </rh-card>
  <rh-card color-palette="lighter">
    <rh-video slot="image"></rh-video>
    <a href="#">Lorem ipsum dolor sit amet consectetur adipisicing elit.</a>
  </rh-card>
  <rh-card color-palette="lighter">
    <rh-video slot="image"></rh-video>
    <p>Lorem ipsum dolor sit amet consectetur adipisicing elit.
       Quaerat quae accusantium velit sed amet, praesentium maiores illum ad.
       Odio vero molestiae sint animi. Vero in ad fugit sit, possimus explicabo.</p>
    <rh-cta slot="footer"><a href="#">Call to action</a></rh-cta>
  </rh-card>
  <rh-card color-palette="lighter">
    <rh-video slot="image"></rh-video>
    <rh-blockquote>
      <p>Lorem ipsum dolor sit amet consectetur adipisicing elit.
         Quaerat quae accusantium velit sed amet, praesentium maiores illum ad.
         Odio vero molestiae sint animi. Vero in ad fugit sit, possimus explicabo.</p>
      <span slot="author">First name Last name</span>
      <span slot="title">Job title, Company name</span>
    </rh-blockquote>
  </rh-card>
</section>
{{< /raw >}}

