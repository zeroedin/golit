---
title: "Tile"
imports:
  - rh-button
  - rh-cta
  - rh-icon
  - rh-tile
lightdom:
  - rh-cta-lightdom-shim.css
  - rh-tile-lightdom.css
---

<p>26 demos for <code>&lt;rh-tile&gt;</code></p>


### accented tiles

{{< raw >}}
<rh-context-demo>
  <rh-tile-group>
    <rh-tile class="accented-tile">
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    </rh-tile>
    <rh-tile class="accented-tile">
      <h2 slot="headline"><a href="#top">Link 2</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    </rh-tile>
    <rh-tile class="accented-tile">
      <h2 slot="headline"><a href="#top">Link 2</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    </rh-tile>
  </rh-tile-group>
</rh-context-demo>
{{< /raw >}}


### checkable

{{< raw >}}
<rh-tile checkable>
  <h2 slot="headline">Headline</h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile checkable checked>
  <h2 slot="headline">Headline</h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile checkable bleed disabled>
  <h2 slot="headline">Headline</h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### color context

{{< raw >}}
<rh-context-demo>
  <div class="wrap layout">
    <h2>Basic</h2>
    <rh-tile>
      <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 296 50">
        <title>A placeholder image in a tile image slot</title>
        <rect x="0"
              y="0"
              width="100%"
              height="100%"
              fill="var(--rh-color-interactive-primary-default)"
              fill-opacity="0.2"
              stroke="var(--rh-color-interactive-primary-hover)"
              stroke-width="1"
              stroke-dasharray="1 1"
        />
      </svg>
      <div slot="title">Title</div>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
    <rh-tile>
      <div slot="title">Title</div>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
  </div>
  <div class="wrap layout">
    <h2>Full-width images</h2>
    <rh-tile bleed>
      <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 200">
        <title>A placeholder image in a tile image slot</title>
        <rect x="0"
              y="0"
              width="100%"
              height="100%"
              fill="var(--rh-color-interactive-primary-default)"
              fill-opacity="0.2"
              stroke="var(--rh-color-interactive-primary-hover)"
              stroke-width="1"
              stroke-dasharray="1 1"
        />
      </svg>
      <div slot="title">Title</div>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
    <rh-tile bleed>
      <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 170">
        <title>A placeholder image in a tile image slot</title>
        <rect x="0"
              y="0"
              width="100%"
              height="100%"
              fill="var(--rh-color-interactive-primary-default)"
              fill-opacity="0.2"
              stroke="var(--rh-color-interactive-primary-hover)"
              stroke-width="1"
              stroke-dasharray="1 1"
        />
      </svg>
      <rh-icon slot="icon" set="standard" icon="cloud-automation"></rh-icon>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
  </div>
  <div class="wrap layout">
    <h2>Desaturated heading</h2>
    <rh-tile desaturated>
      <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
        <path d="M40.3 20.3c.1.3.1.6.1.9 0 3.7-4.5 4.4-7.7 4.4-12.2 0-21.2-7.6-21.2-9.9 0-.2 0-.3.1-.5l-1 2.3c-.2.6-.4 1.2-.4 1.8 0 4.5 10.2 11.4 21.9 11.4 5.2 0 9.1-1.9 9.1-5.4-.1-.9-.2-1.7-.4-2.5l-.5-2.5z"/>
        <path d="M32.7 25.5c3.1 0 7.7-.6 7.7-4.4 0-.3 0-.6-.1-.9l-1.9-8.1c-.4-1.8-.8-2.6-3.9-4.2-2.4-1.2-7.7-3.3-9.3-3.3-1.5 0-1.9 1.9-3.6 1.9s-2.9-1.4-4.5-1.4c-1.5 0-2.5 1-3.2 3.1 0 0-2.1 5.9-2.4 6.8 0 .2-.1.3-.1.5 0 2.5 9.1 10 21.3 10m8.1-2.8c.2.8.4 1.7.4 2.5 0 3.5-3.9 5.4-9.1 5.4-11.7 0-21.9-6.8-21.9-11.4 0-.6.1-1.3.4-1.8-4.2.2-9.6 1-9.6 5.8 0 7.9 18.6 17.6 33.4 17.6 11.3 0 14.2-5.1 14.2-9.2-.1-3.1-2.8-6.8-7.8-8.9" fill="#e00"/>
      </svg>
      <div slot="title">Title</div>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
  </div>
  <div class="wrap layout">
    <h2>Disabled</h2>
    <rh-tile disabled>
      <div slot="title">Title</div>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
  </div>
  <div class="wrap layout">
    <h2 id="compact">Compact</h2>
    <rh-tile compact>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
    <rh-tile compact>
      <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
        <path d="M40.3 20.3c.1.3.1.6.1.9 0 3.7-4.5 4.4-7.7 4.4-12.2 0-21.2-7.6-21.2-9.9 0-.2 0-.3.1-.5l-1 2.3c-.2.6-.4 1.2-.4 1.8 0 4.5 10.2 11.4 21.9 11.4 5.2 0 9.1-1.9 9.1-5.4-.1-.9-.2-1.7-.4-2.5l-.5-2.5z"/>
        <path d="M32.7 25.5c3.1 0 7.7-.6 7.7-4.4 0-.3 0-.6-.1-.9l-1.9-8.1c-.4-1.8-.8-2.6-3.9-4.2-2.4-1.2-7.7-3.3-9.3-3.3-1.5 0-1.9 1.9-3.6 1.9s-2.9-1.4-4.5-1.4c-1.5 0-2.5 1-3.2 3.1 0 0-2.1 5.9-2.4 6.8 0 .2-.1.3-.1.5 0 2.5 9.1 10 21.3 10m8.1-2.8c.2.8.4 1.7.4 2.5 0 3.5-3.9 5.4-9.1 5.4-11.7 0-21.9-6.8-21.9-11.4 0-.6.1-1.3.4-1.8-4.2.2-9.6 1-9.6 5.8 0 7.9 18.6 17.6 33.4 17.6 11.3 0 14.2-5.1 14.2-9.2-.1-3.1-2.8-6.8-7.8-8.9" fill="#e00"/>
      </svg>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
    <rh-tile compact>
      <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
        <path d="M40.3 20.3c.1.3.1.6.1.9 0 3.7-4.5 4.4-7.7 4.4-12.2 0-21.2-7.6-21.2-9.9 0-.2 0-.3.1-.5l-1 2.3c-.2.6-.4 1.2-.4 1.8 0 4.5 10.2 11.4 21.9 11.4 5.2 0 9.1-1.9 9.1-5.4-.1-.9-.2-1.7-.4-2.5l-.5-2.5z"/>
        <path d="M32.7 25.5c3.1 0 7.7-.6 7.7-4.4 0-.3 0-.6-.1-.9l-1.9-8.1c-.4-1.8-.8-2.6-3.9-4.2-2.4-1.2-7.7-3.3-9.3-3.3-1.5 0-1.9 1.9-3.6 1.9s-2.9-1.4-4.5-1.4c-1.5 0-2.5 1-3.2 3.1 0 0-2.1 5.9-2.4 6.8 0 .2-.1.3-.1.5 0 2.5 9.1 10 21.3 10m8.1-2.8c.2.8.4 1.7.4 2.5 0 3.5-3.9 5.4-9.1 5.4-11.7 0-21.9-6.8-21.9-11.4 0-.6.1-1.3.4-1.8-4.2.2-9.6 1-9.6 5.8 0 7.9 18.6 17.6 33.4 17.6 11.3 0 14.2-5.1 14.2-9.2-.1-3.1-2.8-6.8-7.8-8.9" fill="#e00"/>
      </svg>
      <h2 slot="headline"><a href="#top">Link</a></h2>
    </rh-tile>
    <rh-tile compact bleed>
      <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 50">
        <title>A placeholder image in a tile image slot</title>
        <rect x="0"
              y="0"
              width="100%"
              height="100%"
              fill="var(--rh-color-interactive-primary-default)"
              fill-opacity="0.2"
              stroke="var(--rh-color-interactive-primary-hover)"
              stroke-width="1"
              stroke-dasharray="1 1"
        />
      </svg>
      <rh-icon slot="icon" set="standard" icon="cloud-automation"></rh-icon>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
    <rh-tile compact>
      <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 50">
        <title>A placeholder image in a tile image slot</title>
        <rect x="0"
              y="0"
              width="100%"
              height="100%"
              fill="var(--rh-color-interactive-primary-default)"
              fill-opacity="0.2"
              stroke="var(--rh-color-interactive-primary-hover)"
              stroke-width="1"
              stroke-dasharray="1 1"
        />
      </svg>
      <h2 slot="headline"><a href="#top">Link</a></h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
    <rh-tile compact bleed>
      <a href="#top" slot="image">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 50">
          <title>A placeholder image in a tile image slot</title>
          <rect x="0"
                y="0"
                width="100%"
                height="100%"
                fill="var(--rh-color-interactive-primary-default)"
                fill-opacity="0.2"
                stroke="var(--rh-color-interactive-primary-hover)"
                stroke-width="1"
                stroke-dasharray="1 1"
          />
        </svg>
      </a>
    </rh-tile>
  </div>
  <div class="wrap layout">
    <h2>Checkboxes</h2>
    <rh-tile checkable>
      <h2 slot="headline">Headline</h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
    <rh-tile checkable checked>
      <h2 slot="headline">Headline</h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
    <rh-tile checkable bleed disabled>
      <h2 slot="headline">Headline</h2>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit.
      <div slot="footer">Suspendisse eu turpis elementum</div>
    </rh-tile>
  </div>
  <div class="wrap">
    <h2>Tile Group</h2>
    <rh-tile-group>
      <rh-tile checked>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
      <rh-tile>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
      <rh-tile>
        <div slot="title">Title</div>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
    </rh-tile-group>
  </div>
  <div class="wrap">
    <h2>Tile Group, Disabled</h2>
    <rh-tile-group disabled>
      <rh-tile checked>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
      <rh-tile>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
      <rh-tile>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
    </rh-tile-group>
  </div>
  <div class="wrap">
    <h2>Tile Group, Radio</h2>
    <rh-tile-group radio>
      <rh-tile checked>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
      <rh-tile>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
      <rh-tile>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
    </rh-tile-group>
  </div>
  <div class="wrap">
    <h2>Tile Group, Radio, Disabled</h2>
    <rh-tile-group radio disabled>
      <rh-tile checked>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
      <rh-tile>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
      <rh-tile>
        <h2 slot="headline">Headline</h2>
        Lorem ipsum dolor sit amet, consectetur adipiscing elit.
        <div slot="footer">Suspendisse eu turpis elementum</div>
      </rh-tile>
    </rh-tile-group>
  </div>
</rh-context-demo>
{{< /raw >}}


### color palettes

{{< raw >}}
<rh-tile color-palette="darkest">
  <h2 slot="headline">Darkest</h2>
  <rh-cta href="#">Call to Action</rh-cta>
</rh-tile>
<rh-tile color-palette="darker">
  <h2 slot="headline">Darker</h2>
  <rh-cta href="#">Call to Action</rh-cta>
</rh-tile>
<rh-tile color-palette="dark">
  <h2 slot="headline">Dark</h2>
  <rh-cta href="#">Call to Action</rh-cta>
</rh-tile>
<rh-tile color-palette="light">
  <h2 slot="headline">Light</h2>
  <rh-cta href="#">Call to Action</rh-cta>
</rh-tile>
<rh-tile color-palette="lighter">
  <h2 slot="headline">Lighter</h2>
  <rh-cta href="#">Call to Action</rh-cta>
</rh-tile>
<rh-tile color-palette="lightest">
  <h2 slot="headline">Lightest</h2>
  <rh-cta href="#">Call to Action</rh-cta>
</rh-tile>
{{< /raw >}}


### compact link with fullwidth image and icon

{{< raw >}}
<rh-tile compact bleed>
  <rh-icon slot="icon" set="standard" icon="mug"></rh-icon>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 120">
    <title>A placeholder image in a tile image slot</title>
    <rect x="0"
          y="0"
          width="320"
          height="120"
          fill="var(--rh-color-interactive-primary-default)"
          fill-opacity="0.2"
          stroke="var(--rh-color-interactive-primary-hover)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <h2 slot="headline"><a href="#top">Compact link tile</a></h2>
  <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit set est egestat, a sollicitudn mauris tincidunt.</p>
</rh-tile>
{{< /raw >}}


### compact link with icon

{{< raw >}}
<rh-tile compact desaturated>
  <h2 slot="headline"><a href="#top">Compact link tile</a></h2>
  <rh-icon slot="icon" set="standard" icon="mug"></rh-icon>
  <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit set est egestat, a sollicitudn mauris tincidunt.</p>
</rh-tile>
{{< /raw >}}


### compact link with image and icon

{{< raw >}}
<rh-tile compact>
  <rh-icon slot="icon" set="standard" icon="mug"></rh-icon>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 120">
    <title>A placeholder image in a card header slot</title>
    <rect x="0"
          y="0"
          width="320"
          height="120"
          fill="var(--rh-color-interactive-primary-default)"
          fill-opacity="0.2"
          stroke="var(--rh-color-interactive-primary-hover)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <h2 slot="headline"><a href="#top">Compact link tile</a></h2>
  <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit set est egestat, a sollicitudn mauris tincidunt.</p>
</rh-tile>
{{< /raw >}}


### compact link with image

{{< raw >}}
<rh-tile compact desaturated>
  <h2 slot="headline"><a href="#top">Compact link tile</a></h2>
  <rh-icon slot="icon" set="standard" icon="mug"></rh-icon>
  <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit set est egestat, a sollicitudn mauris tincidunt.</p>
</rh-tile>
{{< /raw >}}


### compact

{{< raw >}}
<rh-tile compact>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile compact>
  <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
    <title>Red Hat</title>
    <path d="M40.3 20.3c.1.3.1.6.1.9 0 3.7-4.5 4.4-7.7 4.4-12.2 0-21.2-7.6-21.2-9.9 0-.2 0-.3.1-.5l-1 2.3c-.2.6-.4 1.2-.4 1.8 0 4.5 10.2 11.4 21.9 11.4 5.2 0 9.1-1.9 9.1-5.4-.1-.9-.2-1.7-.4-2.5l-.5-2.5z"/>
    <path d="M32.7 25.5c3.1 0 7.7-.6 7.7-4.4 0-.3 0-.6-.1-.9l-1.9-8.1c-.4-1.8-.8-2.6-3.9-4.2-2.4-1.2-7.7-3.3-9.3-3.3-1.5 0-1.9 1.9-3.6 1.9s-2.9-1.4-4.5-1.4c-1.5 0-2.5 1-3.2 3.1 0 0-2.1 5.9-2.4 6.8 0 .2-.1.3-.1.5 0 2.5 9.1 10 21.3 10m8.1-2.8c.2.8.4 1.7.4 2.5 0 3.5-3.9 5.4-9.1 5.4-11.7 0-21.9-6.8-21.9-11.4 0-.6.1-1.3.4-1.8-4.2.2-9.6 1-9.6 5.8 0 7.9 18.6 17.6 33.4 17.6 11.3 0 14.2-5.1 14.2-9.2-.1-3.1-2.8-6.8-7.8-8.9" fill="#e00"/>
  </svg>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile compact>
  <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
    <title>Red Hat</title>
    <path d="M40.3 20.3c.1.3.1.6.1.9 0 3.7-4.5 4.4-7.7 4.4-12.2 0-21.2-7.6-21.2-9.9 0-.2 0-.3.1-.5l-1 2.3c-.2.6-.4 1.2-.4 1.8 0 4.5 10.2 11.4 21.9 11.4 5.2 0 9.1-1.9 9.1-5.4-.1-.9-.2-1.7-.4-2.5l-.5-2.5z"/>
    <path d="M32.7 25.5c3.1 0 7.7-.6 7.7-4.4 0-.3 0-.6-.1-.9l-1.9-8.1c-.4-1.8-.8-2.6-3.9-4.2-2.4-1.2-7.7-3.3-9.3-3.3-1.5 0-1.9 1.9-3.6 1.9s-2.9-1.4-4.5-1.4c-1.5 0-2.5 1-3.2 3.1 0 0-2.1 5.9-2.4 6.8 0 .2-.1.3-.1.5 0 2.5 9.1 10 21.3 10m8.1-2.8c.2.8.4 1.7.4 2.5 0 3.5-3.9 5.4-9.1 5.4-11.7 0-21.9-6.8-21.9-11.4 0-.6.1-1.3.4-1.8-4.2.2-9.6 1-9.6 5.8 0 7.9 18.6 17.6 33.4 17.6 11.3 0 14.2-5.1 14.2-9.2-.1-3.1-2.8-6.8-7.8-8.9" fill="#e00"/>
  </svg>
  <h2 slot="headline"><a href="#top">Link</a></h2>
</rh-tile>
<rh-tile compact bleed>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 50">
    <title>A placeholder image in a tile image slot</title>
    <rect x="0"
          y="0"
          width="300"
          height="50"
          fill="var(--rh-color-interactive-primary-default)"
          fill-opacity="0.2"
          stroke="var(--rh-color-interactive-primary-hover)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <rh-icon slot="icon" set="standard" icon="cloud-automation"></rh-icon>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile compact>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 296 50">
    <title>A placeholder image in a tile image slot</title>
    <rect x="0"
          y="0"
          width="296"
          height="50"
          fill="var(--rh-color-interactive-primary-default)"
          fill-opacity="0.2"
          stroke="var(--rh-color-interactive-primary-hover)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile compact bleed>
  <a href="#top" slot="image">
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 50">
      <title>A placeholder image in a tile image slot</title>
      <rect x="0"
            y="0"
            width="300"
            height="50"
            fill="var(--rh-color-interactive-primary-default)"
            fill-opacity="0.2"
            stroke="var(--rh-color-interactive-primary-hover)"
            stroke-width="1"
            stroke-dasharray="1 1"
      />
    </svg>
  </a>
</rh-tile>
{{< /raw >}}


### custom props

{{< raw >}}
<rh-context-demo color-palette="light">
  <rh-tile>
    <h2 slot="headline"><a href="#top">Link</a></h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  </rh-tile>
  <rh-tile>
    <h2 slot="headline"><a href="#top">Link 2</a></h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  </rh-tile>
  <rh-tile>
    <h2 slot="headline"><a href="#top">Link 2</a></h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  </rh-tile>
</rh-context-demo>
{{< /raw >}}


### desaturated heading slotted icon

{{< raw >}}
<rh-tile desaturated>
  <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48">
    <path d="M40.3 20.3c.1.3.1.6.1.9 0 3.7-4.5 4.4-7.7 4.4-12.2 0-21.2-7.6-21.2-9.9 0-.2 0-.3.1-.5l-1 2.3c-.2.6-.4 1.2-.4 1.8 0 4.5 10.2 11.4 21.9 11.4 5.2 0 9.1-1.9 9.1-5.4-.1-.9-.2-1.7-.4-2.5l-.5-2.5z"/>
    <path d="M32.7 25.5c3.1 0 7.7-.6 7.7-4.4 0-.3 0-.6-.1-.9l-1.9-8.1c-.4-1.8-.8-2.6-3.9-4.2-2.4-1.2-7.7-3.3-9.3-3.3-1.5 0-1.9 1.9-3.6 1.9s-2.9-1.4-4.5-1.4c-1.5 0-2.5 1-3.2 3.1 0 0-2.1 5.9-2.4 6.8 0 .2-.1.3-.1.5 0 2.5 9.1 10 21.3 10m8.1-2.8c.2.8.4 1.7.4 2.5 0 3.5-3.9 5.4-9.1 5.4-11.7 0-21.9-6.8-21.9-11.4 0-.6.1-1.3.4-1.8-4.2.2-9.6 1-9.6 5.8 0 7.9 18.6 17.6 33.4 17.6 11.3 0 14.2-5.1 14.2-9.2-.1-3.1-2.8-6.8-7.8-8.9" fill="#e00"/>
  </svg>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### desaturated heading

{{< raw >}}
<rh-tile desaturated icon="cloud-automation">
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### disabled

{{< raw >}}
<rh-tile disabled>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### form control

{{< raw >}}
<form id="form">
  <rh-tile-group radio>
    <rh-tile name="radio" value="1"><span slot="headline">Radio 1</span></rh-tile>
    <rh-tile name="radio" value="2"><span slot="headline">Radio 2</span></rh-tile>
  </rh-tile-group>
  <rh-tile-group checkable>
    <rh-tile name="check" value="1"><span slot="headline">Check 1</span></rh-tile>
    <rh-tile name="check" value="2"><span slot="headline">Check 2</span></rh-tile>
  </rh-tile-group>
  <rh-button type="submit">Submit</rh-button>
  <output name="output"></output>
</form>
{{< /raw >}}


### full width images

{{< raw >}}
<rh-tile bleed>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 200" width="100%">
    <title>300 X 200 placeholder image</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile bleed>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 170" width="100%">
    <title>300 X 170 placeholder image</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <rh-icon slot="headline" set="standard" icon="cloud-automation"></rh-icon>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### icon with full width image

{{< raw >}}
<rh-tile bleed>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 200">
    <title>300 X 200 placeholder image</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
    <title>A placeholder image in a tile icon slot</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile bleed>
  <rh-icon slot="icon" icon="cloud-automation"></rh-icon>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 200">
    <title>300 X 200 placeholder image</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile bleed icon="cloud-automation">
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 200">
    <title>300 X 200 placeholder image</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### icon with image

{{< raw >}}
<rh-tile>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" width="300" height="200">
    <title>300 X 200 placeholder image</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
    <title>A placeholder image in a tile icon slot</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile icon="cloud-automation">
  <rh-icon slot="icon" icon="cloud-automation"></rh-icon>
  <svg slot="image" xmlns="http://www.w3.org/2000/svg" width="300" height="200">
    <title>300 X 200 placeholder image</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### icon

{{< raw >}}
<rh-tile>
  <rh-icon slot="icon" icon="cloud-automation"></rh-icon>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile icon="cloud-automation">
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile>
  <svg slot="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
    <title>A placeholder image in a tile icon slot</title>
    <rect x="0" y="0" width="99%" height="99%"
          fill="none"
          stroke="var(--rh-color-icon-primary)"
          stroke-width="1"
          stroke-dasharray="1 1"
    />
  </svg>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### in flex container

{{< raw >}}
<rh-tile>
  <img slot="image" src="https://fakeimg.pl/296x50" alt="296 X 50 placeholder image">
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile>
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### index

{{< raw >}}
<rh-tile>
  <img slot="image" src="https://fakeimg.pl/296x50" alt="296 X 50 placeholder image">
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#top">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### link types

{{< raw >}}
<rh-tile link="private"
         icon="bank-safe"
         icon-set="standard">
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#super-secret-section">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
<rh-tile link="external"
         icon="globe-abstract"
         icon-set="standard">
  <div slot="title">Title</div>
  <h2 slot="headline"><a href="#super-public-section">Link</a></h2>
  Lorem ipsum dolor sit amet, consectetur adipiscing elit.
  <div slot="footer">Suspendisse eu turpis elementum</div>
</rh-tile>
{{< /raw >}}


### tile group disabled

{{< raw >}}
<rh-tile-group disabled>
  <rh-tile checked>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
  <rh-tile>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
  <rh-tile>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
</rh-tile-group>
{{< /raw >}}


### tile group radio disabled

{{< raw >}}
<rh-tile-group radio disabled>
  <rh-tile checked>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
  <rh-tile>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
  <rh-tile>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
</rh-tile-group>
{{< /raw >}}


### tile group radio

{{< raw >}}
<rh-tile-group radio>
  <rh-tile checked>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
  <rh-tile>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
  <rh-tile>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
</rh-tile-group>
{{< /raw >}}


### tile group

{{< raw >}}
<rh-tile-group>
  <rh-tile checked>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
  <rh-tile>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
  <rh-tile>
    <div slot="title">Title</div>
    <h2 slot="headline">Headline</h2>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    <div slot="footer">Suspendisse eu turpis elementum</div>
  </rh-tile>
</rh-tile-group>
{{< /raw >}}


### without footer content

{{< raw >}}
<section class="grid">
  <rh-tile>
    <h2 slot="headline"><a href="#top">Tile w/ no footer</a></h2>
    Lorem ipsum dolor sit, amet consectetur adipisicing elit. Atque tenetur quae omnis beatae, distinctio natus quo hic est iste tempora vitae nesciunt, minus totam itaque modi! Laborum obcaecati
    soluta sapiente?
  </rh-tile>
  <rh-tile>
    <h2 slot="headline"><a href="#top">Tile w/ no footer</a></h2>
    Lorem ipsum dolor sit, amet consectetur adipisicing elit. Atque tenetur quae omnis beatae, distinctio natus quo hic est iste tempora vitae nesciunt, minus totam itaque modi! Laborum obcaecati
    soluta sapiente?
  </rh-tile>
  <rh-tile>
    <h2 slot="headline"><a href="#top">Tile w/ only headline</a></h2>
  </rh-tile>
  <rh-tile>
    <h2 slot="headline"><a href="#top">Tile w/ only headline</a></h2>
  </rh-tile>
</section>
{{< /raw >}}

