---
title: "Call to Action"
imports:
  - rh-cta
  - rh-surface
lightdom:
  - rh-cta-lightdom-shim.css
---

<p>11 demos for <code>&lt;rh-cta&gt;</code></p>


### analytics

{{< raw >}}
<p>In this demo, analytics events involving <code>&lt;rh-cta></code> elements are parsed by a
  document-level analytics event listener. Unlike <code>&lt;pfe-cta></code>, which implemented
  support for analytics internally, <code>&lt;rh-cta></code> users must implement their own
  analytics code, taking this demo as an example.</p>
<section data-analytics-category="simple">
  <h2>Simple Case: Light DOM</h2>
  <rh-cta data-analytics-linktype="cta"
          data-analytics-text="Default">
    <a href="#default">Default</a>
  </rh-cta>
  <rh-cta variant="primary"
          data-analytics-linktype="cta"
          data-analytics-text="Primary">
    <a href="#primary">Primary</a>
  </rh-cta>
  <rh-cta variant="secondary"
          data-analytics-linktype="cta"
          data-analytics-text="Secondary">
    <a href="#secondary">Secondary</a>
  </rh-cta>
  <rh-cta variant="brick"
          data-analytics-linktype="cta"
          data-analytics-text="Brick">
    <a href="#brick">Brick</a>
  </rh-cta>
  <section data-analytics-category="href-attr">
    <h3>Using <code>href</code> attribute</h3>
    <rh-cta href="#default"
            data-analytics-linktype="cta"
            data-analytics-text="Default">Default</rh-cta>
    <rh-cta href="#primary"
            variant="primary"
            data-analytics-linktype="cta"
            data-analytics-text="Primary">Primary</rh-cta>
    <rh-cta href="#secondary"
            variant="secondary"
            data-analytics-linktype="cta"
            data-analytics-text="Secondary">Secondary</rh-cta>
    <rh-cta href="#brick"
            variant="brick"
            data-analytics-linktype="cta"
            data-analytics-text="Brick">Brick</rh-cta>
  </section>
</section>
<section data-analytics-category="shadow">
  <h2>Advanced Case: Deep Shadow</h2>
  <shadow-root>
    <template shadowrootmode="open">
      <rh-cta data-analytics-linktype="cta"
              data-analytics-text="Default">
        <a href="#default">Default</a>
      </rh-cta>
      <rh-cta variant="primary"
              data-analytics-linktype="cta"
              data-analytics-text="Primary">
        <a href="#primary">Primary</a>
      </rh-cta>
      <rh-cta variant="secondary"
              data-analytics-linktype="cta"
              data-analytics-text="Secondary">
        <a href="#secondary">Secondary</a>
      </rh-cta>
      <rh-cta variant="brick"
              style="width:auto;"
              data-analytics-linktype="cta"
              data-analytics-text="Brick">
        <a href="#brick">Brick</a>
      </rh-cta>
    </template>
  </shadow-root>
  <h3>Using <code>href</code> attribute</h3>
  <shadow-root data-analytics-category="href-attr">
    <template shadowrootmode="open">
      <rh-cta href="#default"
              data-analytics-linktype="cta"
              data-analytics-text="Default">Default</rh-cta>
      <rh-cta href="#primary"
              data-analytics-linktype="cta"
              data-analytics-text="Primary"
              variant="primary">Primary</rh-cta>
      <rh-cta href="#secondary"
              data-analytics-linktype="cta"
              data-analytics-text="Secondary"
              variant="secondary">Secondary</rh-cta>
      <rh-cta href="#brick"
              style="width:auto;"
              variant="brick"
              data-analytics-linktype="cta"
              data-analytics-text="Brick">Brick</rh-cta>
    </template>
  </shadow-root>
</section>
<section data-analytics-category="slotted">
  <h2>Complex Case: Slotted Link, Deep CTA</h2>
  <slotted-link>
    <a slot="default"
       data-analytics-linktype="cta"
       data-analytics-text="Default"
       href="#default">Default</a>
    <a slot="primary"
       data-analytics-linktype="cta"
       data-analytics-text="Primary"
       href="#primary">Primary</a>
    <a slot="secondary"
       data-analytics-linktype="cta"
       data-analytics-text="Secondary"
       href="#secondary">Secondary</a>
    <a slot="brick"
       data-analytics-linktype="cta"
       data-analytics-text="Brick"
       href="#brick">Brick</a>
  </slotted-link>
</section>
<h2>Last CTA Analytics Event</h2>
<json-viewer>{}</json-viewer>
{{< /raw >}}


### brick

{{< raw >}}
<div id="grid">
  <rh-cta variant="brick"><a href="#">Link #1</a></rh-cta>
  <rh-cta variant="brick"><a href="#">Link #2</a></rh-cta>
  <rh-cta variant="brick"><a href="#">Link #3</a></rh-cta>
  <rh-cta variant="brick">
    <a href="#default">Supercalifragilisticexpialidocious</a>
  </rh-cta>
</div>
{{< /raw >}}


### button

{{< raw >}}
<rh-cta>
  <button>Button</button>
</rh-cta>
<rh-cta variant="primary">
  <button>Button</button>
</rh-cta>
<rh-cta variant="secondary">
  <button>Button</button>
</rh-cta>
{{< /raw >}}


### color context with lightdom css

{{< raw >}}
<rh-context-demo>
  <rh-cta><a href="#default">Default</a></rh-cta>
  <rh-cta icon="play-circle"><a href="#default-video">Default Video</a></rh-cta>
  <rh-cta variant="primary"><a href="#primary">Primary</a></rh-cta>
  <rh-cta variant="primary" icon="play-circle"><a href="#primary-video">Video</a></rh-cta>
  <rh-cta variant="secondary"><a href="#secondary">Secondary</a></rh-cta>
  <rh-cta variant="brick"><a href="#brick">Brick</a></rh-cta>
  <rh-cta variant="brick" icon="play-circle"><a href="#brick-video">Brick Video</a></rh-cta>
</rh-context-demo>
{{< /raw >}}


### color context

{{< raw >}}
<rh-context-demo>
  <rh-cta><a href="#default">Default</a></rh-cta>
  <rh-cta icon="play-circle"><a href="#default-video">Default Video</a></rh-cta>
  <rh-cta variant="primary"><a href="#primary">Primary</a></rh-cta>
  <rh-cta variant="primary" icon="play-circle"><a href="#primary-video">Video</a></rh-cta>
  <rh-cta variant="secondary"><a href="#secondary">Secondary</a></rh-cta>
  <rh-cta variant="brick"><a href="#brick">Brick</a></rh-cta>
  <rh-cta variant="brick" icon="play-circle"><a href="#brick-video">Brick Video</a></rh-cta>
</rh-context-demo>
{{< /raw >}}


### href attribute

{{< raw >}}
<section>
  <rh-cta href="#default">Default</rh-cta>
  <rh-cta href="#default-video" icon="play-circle">Default Video</rh-cta>
  <rh-cta href="#primary" variant="primary">Primary</rh-cta>
  <rh-cta href="#primary-video" variant="primary" icon="play-circle">Video</rh-cta>
  <rh-cta href="#secondary" variant="secondary">Secondary</rh-cta>
  <div id="brick">
    <rh-cta href="#brick" variant="brick">Brick</rh-cta>
    <rh-cta href="#brick-icon" variant="brick" icon="user">Brick Icon</rh-cta>
  </div>
</section>
{{< /raw >}}


### index

{{< raw >}}
<rh-cta href="#">Call to Action</rh-cta>
{{< /raw >}}


### no cta javascript

{{< raw >}}
<p>&lt;rh-cta&gt; where JavaScript is not loaded, while loading a `lightdom-shim.css`.</p>
<section id="variants">
  <h2>Variants</h2>
  <rh-cta><a href="#default">Default</a></rh-cta>
  <rh-cta icon="play-circle"><a href="#default-video">Default Video</a></rh-cta>
  <rh-cta variant="primary"><a href="#primary">Primary</a></rh-cta>
  <rh-cta variant="primary" icon="play-circle"><a href="#primary-video">Video</a></rh-cta>
  <rh-cta variant="secondary"><a href="#secondary">Secondary</a></rh-cta>
  <rh-cta variant="brick"><a href="#brick">Brick</a></rh-cta>
  <rh-cta variant="brick" icon="play-circle"><a href="#brick-video">Brick Video</a></rh-cta>
</section>
<rh-surface color-palette="darkest">
  <h2>Dark Color Context</h2>
  <rh-cta><a href="#default">Default</a></rh-cta>
  <rh-cta icon="play-circle"><a href="#default-video">Default Video</a></rh-cta>
  <rh-cta variant="primary"><a href="#primary">Primary</a></rh-cta>
  <rh-cta variant="primary" icon="play-circle"><a href="#primary-video">Video</a></rh-cta>
  <rh-cta variant="secondary"><a href="#secondary">Secondary</a></rh-cta>
  <rh-cta variant="brick"><a href="#brick">Brick</a></rh-cta>
  <rh-cta variant="brick" icon="play-circle"><a href="#brick-video">Brick Video</a></rh-cta>
</rh-surface>
<section dir="rtl" lang="he">
  <header lang="en" dir="ltr">
    <h2>Right-to-Left Languages</h2>
  </header>
  <div>
    <rh-cta>
      <a href="#default">ברירת מחדל</a>
    </rh-cta>
    <rh-cta icon="play-circle">
      <a href="#default-video">ברירת מחדל - וידאו</a>
    </rh-cta>
    <rh-cta variant="primary">
      <a href="#primary">ראשי</a>
    </rh-cta>
    <rh-cta variant="primary" icon="play-circle">
      <a href="#primary-video">ראשי - וידאו</a>
    </rh-cta>
    <rh-cta variant="secondary">
      <a href="#secondary">משני</a>
    </rh-cta>
    <rh-cta variant="brick">
      <a href="#brick">לבנה</a>
    </rh-cta>
    <rh-cta variant="brick" icon="user">
      <a href="#brick-icon">לבנה עם אייקון</a>
    </rh-cta>
  </div>
</section>
<section dir="rtl" lang="he">
  <header dir="ltr" lang="en">
    <h2>Deep Shadow RTL</h2>
    <p>
      When the CTA is found within the shadow root of an element which is itself within a RTL context,
      it's own contents should also be displayed right-to-left, including the various box models, the placement
      of the icon, and the direction of the arrow. Authors should not need to specify the direction with the
      <code>dir="rtl"</code> attribute if it is added to a containing element.
      In this demo, the shadow root's host element has a light blue background color, in order distinguish it from the document content.
    </p>
  </header>
  <shadow-root>
    <template shadowrootmode="open">
      <rh-cta id="deep">
        <a href="#default">ברירת מחדל</a>
      </rh-cta>
      <rh-cta icon="play-circle">
        <a href="#default-video">ברירת מחדל - וידאו</a>
      </rh-cta>
      <rh-cta variant="primary">
        <a href="#primary">ראשי</a>
      </rh-cta>
      <rh-cta variant="primary" icon="play-circle">
        <a href="#primary-video">ראשי - וידאו</a>
      </rh-cta>
      <rh-cta variant="secondary">
        <a href="#secondary">משני</a>
      </rh-cta>
      <rh-cta variant="brick">
        <a href="#brick">לבנה</a>
      </rh-cta>
      <rh-cta variant="brick" icon="user">
        <a href="#brick-icon">לבנה עם אייקון</a>
      </rh-cta>
    </template>
  </shadow-root>
</section>
{{< /raw >}}


### resizing

{{< raw >}}
<section id="resize">
  <div>
    <rh-cta>
      <a href="#default">Get product details</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta icon="play-circle">
      <a href="#default-video">Get product details</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="primary">
      <a href="#primary">Get product details</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="primary" icon="play-circle">
      <a href="#primary-video">Get product details</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="secondary">
      <a href="#secondary">Get product details</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="brick">
      <a href="#brick">Get product details</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="brick" icon="user">
      <a href="#brick-icon">Get product details</a>
    </rh-cta>
  </div>
  <p>No matter the container size, the arrow trailing the CTA message should never appear on a line by itself.</p>
  <div>
    <rh-cta><a href="#default">Default link cta with longer text</a></rh-cta>
  </div>
  <div dir="rtl">
    <rh-cta>
      <a href="#default">קריאה לפעולה בררית מחדל עם טקסט ארוך</a>
    </rh-cta>
  </div>
  <p>Long words should break in the middle</p>
  <div>
    <rh-cta>
      <a href="#default">Supercalifragilisticexpialidocious</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta icon="play-circle">
      <a href="#default-video">Supercalifragilisticexpialidocious</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="primary">
      <a href="#primary">Supercalifragilisticexpialidocious</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="primary" icon="play-circle">
      <a href="#primary-video">Supercalifragilisticexpialidocious</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="secondary">
      <a href="#secondary">Supercalifragilisticexpialidocious</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="brick">
      <a href="#brick">Supercalifragilisticexpialidocious</a>
    </rh-cta>
  </div>
  <div>
    <rh-cta variant="brick" icon="user">
      <a href="#brick-icon">Supercalifragilisticexpialidocious</a>
    </rh-cta>
  </div>
</section>
{{< /raw >}}


### right to left

{{< raw >}}
<section dir="rtl"
         lang="he">
  <header lang="en"
          dir="ltr">
    <h2>Right-to-Left Languages</h2>
  </header>
  <div>
    <rh-cta> href="#default">ברירת מחדל</rh-cta>
    <rh-cta icon="play-circle"
            href="#default-video">ברירת מחדל - וידאו</rh-cta>
    <rh-cta variant="primary"
            href="#primary">ראשי</rh-cta>
    <rh-cta variant="primary"
            icon="play-circle"
            href="#primary-video">ראשי - וידאו</rh-cta>
    <rh-cta variant="secondary"
            href="#secondary">משני</rh-cta>
    <rh-cta variant="brick"
            href="#brick">לבנה</rh-cta>
    <rh-cta variant="brick"
            icon="user"
            href="#brick-icon">לבנה עם אייקון</rh-cta>
  </div>
</section>
<section dir="rtl"
         lang="he">
  <header dir="ltr"
          lang="en">
    <h2>Deep Shadow RTL</h2>
    <p>
      When the CTA is found within the shadow root of an element which is itself within a RTL context,
      it's own contents should also be displayed right-to-left, including the various box models, the placement
      of the icon, and the direction of the arrow. Authors should not need to specify the direction with the
      <code>dir="rtl"</code> attribute if it is added to a containing element.
      In this demo, the shadow root's host element has a light blue background color, in order distinguish it from the document content.
    </p>
  </header>
  <shadow-root>
    <template shadowrootmode="open">
      <rh-cta id="deep"
              href="#default">ברירת מחדל</rh-cta>
      <rh-cta icon="play-circle"
              href="#default-video">ברירת מחדל - וידאו</rh-cta>
      <rh-cta variant="primary"
              href="#primary">ראשי</rh-cta>
      <rh-cta variant="primary"
              icon="play-circle"
              href="#primary-video">ראשי - וידאו</rh-cta>
      <rh-cta variant="secondary"
              href="#secondary">משני</rh-cta>
      <rh-cta variant="brick"
              href="#brick">לבנה</rh-cta>
      <rh-cta variant="brick"
              icon="user"
              href="#brick-icon">לבנה עם אייקון</rh-cta>
    </template>
  </shadow-root>
</section>
{{< /raw >}}


### variants

{{< raw >}}
<section id="cta-variants">
  <rh-cta href="#default">Default</rh-cta>
  <rh-cta icon="play-circle" href="#default-video">Default Video</rh-cta>
  <rh-cta variant="primary" href="#primary">Primary</rh-cta>
  <rh-cta variant="primary" icon="play-circle" href="#primary-video">Video</rh-cta>
  <rh-cta variant="secondary" href="#secondary">Secondary</rh-cta>
  <div id="brick">
    <rh-cta variant="brick" href="#brick">Brick</rh-cta>
    <rh-cta variant="brick" icon="users" href="#brick-icon">Brick Icon</rh-cta>
  </div>
</section>
{{< /raw >}}

