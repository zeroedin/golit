---
title: "Announcement"
imports:
  - rh-announcement
  - rh-cta
lightdom:
  - rh-announcement-lightdom-shim.css
  - rh-cta-lightdom-shim.css
---

<p>5 demos for <code>&lt;rh-announcement&gt;</code></p>


### color context


<rh-context-demo>
  <rh-announcement>
    <svg slot="image" width="80" height="48" role="img" aria-label="Sample image">
      <rect fill="light-dark(var(--rh-color-surface-dark, #383838), var(--rh-color-surface-light, #e0e0e0))" stroke="var(--rh-color-border-subtle)" stroke-width="2px" width="100%" height="100%" stroke-dasharray="4 4">
      </rect>
      <text x="17" y="30" style="font-family: var(--rh-font-family-code, RedHatMono, 'Red Hat Mono', 'Courier New', Courier, monospace); font-size: var(--rh-font-size-body-text-md, 1rem);" fill="light-dark(var(--rh-color-text-primary-on-dark, #ffffff), var(--rh-color-text-primary-on-light, #151515))">
              Image
      </text>
    </svg>
    <p>Lorem ipsum dolor sit amet consectetur adipisicing elit adipisicing elit adipisicing elit.</p>
    <rh-cta slot="cta" href="#">Learn More</rh-cta>
  </rh-announcement>
</rh-context-demo>



### dismissable


<rh-announcement dismissable>
    <svg slot="image" width="80" height="48" role="img" aria-label="Sample image">
      <rect fill="light-dark(var(--rh-color-surface-dark, #383838), var(--rh-color-surface-light, #e0e0e0))" stroke="var(--rh-color-border-subtle)" stroke-width="2px" width="100%" height="100%" stroke-dasharray="4 4">
      </rect>
      <text x="17" y="30" style="font-family: var(--rh-font-family-code, RedHatMono, 'Red Hat Mono', 'Courier New', Courier, monospace); font-size: var(--rh-font-size-body-text-md, 1rem);" fill="light-dark(var(--rh-color-text-primary-on-dark, #ffffff), var(--rh-color-text-primary-on-light, #151515))">
              Image
      </text>
    </svg>
  <p>Lorem ipsum dolor sit amet consectetur adipisicing elit ipsum dolor sit.</p>
  <rh-cta slot="cta" href="#">Learn More</rh-cta>
</rh-announcement>
<rh-announcement dismissable color-palette="dark" image-position="inline-start">
    <svg slot="image" width="80" height="48" role="img" aria-label="Sample image">
      <rect fill="light-dark(var(--rh-color-surface-dark, #383838), var(--rh-color-surface-light, #e0e0e0))" stroke="var(--rh-color-border-subtle)" stroke-width="2px" width="100%" height="100%" stroke-dasharray="4 4">
      </rect>
      <text x="17" y="30" style="font-family: var(--rh-font-family-code, RedHatMono, 'Red Hat Mono', 'Courier New', Courier, monospace); font-size: var(--rh-font-size-body-text-md, 1rem);" fill="light-dark(var(--rh-color-text-primary-on-dark, #ffffff), var(--rh-color-text-primary-on-light, #151515))">
              Image
      </text>
    </svg>
  <p>An announcement with <code>inline-start</code> for the value of <code>image position</code>.</p>
  <rh-cta slot="cta" href="#">Learn More</rh-cta>
</rh-announcement>
<rh-announcement dismissable image-position="block-start">
  <svg slot="image" width="80" height="48" role="img" aria-label="Sample image">
    <rect fill="light-dark(var(--rh-color-surface-dark, #383838), var(--rh-color-surface-light, #e0e0e0))" stroke="var(--rh-color-border-subtle)" stroke-width="2px" width="100%" height="100%" stroke-dasharray="4 4">
    </rect>
    <text x="17" y="30" style="font-family: var(--rh-font-family-code, RedHatMono, 'Red Hat Mono', 'Courier New', Courier, monospace); font-size: var(--rh-font-size-body-text-md, 1rem);" fill="light-dark(var(--rh-color-text-primary-on-dark, #ffffff), var(--rh-color-text-primary-on-light, #151515))">
            Image
    </text>
  </svg>
  <p>An announcement with <code>block-start</code> for the value of <code>image position</code>.</p>
  <rh-cta slot="cta" href="#">Learn More</rh-cta>
</rh-announcement>



### events


<rh-announcement dismissable>
  <svg slot="image" width="80" height="48" role="img" aria-label="Sample image">
    <rect fill="light-dark(var(--rh-color-surface-dark, #383838), var(--rh-color-surface-light, #e0e0e0))" stroke="var(--rh-color-border-subtle)" stroke-width="2px" width="100%" height="100%" stroke-dasharray="4 4">
    </rect>
    <text x="17" y="30" style="font-family: var(--rh-font-family-code, RedHatMono, 'Red Hat Mono', 'Courier New', Courier, monospace); font-size: var(--rh-font-size-body-text-md, 1rem);" fill="light-dark(var(--rh-color-text-primary-on-dark, #ffffff), var(--rh-color-text-primary-on-light, #151515))">
            Image
    </text>
  </svg>
  <p>Click the close button to fire the <code>close</code> event.</p>
  <rh-cta slot="cta" href="#">Learn More</rh-cta>
</rh-announcement>



### image position


<rh-announcement image-position="inline-start">
    <svg slot="image" width="80" height="48" role="img" aria-label="Sample image">
      <rect fill="light-dark(var(--rh-color-surface-dark, #383838), var(--rh-color-surface-light, #e0e0e0))" stroke="var(--rh-color-border-subtle)" stroke-width="2px" width="100%" height="100%" stroke-dasharray="4 4">
      </rect>
      <text x="17" y="30" style="font-family: var(--rh-font-family-code, RedHatMono, 'Red Hat Mono', 'Courier New', Courier, monospace); font-size: var(--rh-font-size-body-text-md, 1rem);" fill="light-dark(var(--rh-color-text-primary-on-dark, #ffffff), var(--rh-color-text-primary-on-light, #151515))">
              Image
      </text>
    </svg>
  <p>On mobile viewports, the image will stay to the left (inline-start) of this main body content.</p>
  <rh-cta slot="cta" href="#">Learn More</rh-cta>
</rh-announcement>
<rh-announcement image-position="block-start" color-palette="dark">
    <svg slot="image" width="80" height="48" role="img" aria-label="Sample image">
      <rect fill="light-dark(var(--rh-color-surface-dark, #383838), var(--rh-color-surface-light, #e0e0e0))" stroke="var(--rh-color-border-subtle)" stroke-width="2px" width="100%" height="100%" stroke-dasharray="4 4">
      </rect>
      <text x="17" y="30" style="font-family: var(--rh-font-family-code, RedHatMono, 'Red Hat Mono', 'Courier New', Courier, monospace); font-size: var(--rh-font-size-body-text-md, 1rem);" fill="light-dark(var(--rh-color-text-primary-on-dark, #ffffff), var(--rh-color-text-primary-on-light, #151515))">
              Image
      </text>
    </svg>
  <p>In this <code>block-start</code> version, the image stays on top of the body content on mobile viewports.</p>
  <rh-cta slot="cta" href="#">Learn More</rh-cta>
</rh-announcement>



### index


<rh-announcement>
    <svg slot="image" width="80" height="48" role="img" aria-label="Sample image">
      <rect fill="light-dark(var(--rh-color-surface-dark, #383838), var(--rh-color-surface-light, #e0e0e0))" stroke="var(--rh-color-border-subtle)" stroke-width="2px" width="100%" height="100%" stroke-dasharray="4 4">
      </rect>
      <text x="17" y="30" style="font-family: var(--rh-font-family-code, RedHatMono, 'Red Hat Mono', 'Courier New', Courier, monospace); font-size: var(--rh-font-size-body-text-md, 1rem);" fill="light-dark(var(--rh-color-text-primary-on-dark, #ffffff), var(--rh-color-text-primary-on-light, #151515))">
              Image
      </text>
    </svg>
  <p>Lorem ipsum dolor sit amet consectetur adipisicing elit adipisicing elit adipisicing elit.</p>
  <rh-cta slot="cta" href="#">Learn More</rh-cta>
</rh-announcement>


