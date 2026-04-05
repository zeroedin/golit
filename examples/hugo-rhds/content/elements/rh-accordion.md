---
title: "Accordion"
imports:
  - rh-accordion
  - rh-cta
  - rh-tag
lightdom:
  - rh-cta-lightdom-shim.css
---

<p>9 demos for <code>&lt;rh-accordion&gt;</code></p>


### accents bottom

{{< raw >}}
<rh-accordion accents="bottom">
  <rh-accordion-header>
    Item One
    <rh-tag slot="accents" color="green" variant="filled" icon="information">Green</rh-tag>
    <rh-tag slot="accents" color="red" variant="filled" icon="information">Red</rh-tag>
    <rh-tag slot="accents" color="orange" variant="filled" icon="notification-fill">Orange</rh-tag>
    <rh-tag slot="accents" color="purple" variant="filled" icon="ban">Purple</rh-tag>
  </rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <rh-accordion-header>
    Item Two
    <rh-tag slot="accents" color="green" variant="filled" icon="information">Green</rh-tag>
  </rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <rh-accordion-header>
    Item Three
    <rh-tag slot="accents" color="red" variant="filled" icon="information">Red</rh-tag>
  </rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
</rh-accordion>
{{< /raw >}}


### accents

{{< raw >}}
<rh-accordion>
  <rh-accordion-header>
    Item One
    <rh-tag slot="accents" color="green" variant="filled" icon="information">Green</rh-tag>
    <rh-tag slot="accents" color="red" variant="filled" icon="information">Red</rh-tag>
    <rh-tag slot="accents" color="orange" variant="filled" icon="notification-fill">Orange</rh-tag>
    <rh-tag slot="accents" color="purple" variant="filled" icon="ban">Purple</rh-tag>
  </rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <rh-accordion-header>
    Item Two
    <rh-tag slot="accents" color="green" variant="filled" icon="information">Green</rh-tag>
  </rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <rh-accordion-header>
    Item Three
    <rh-tag slot="accents" color="red" variant="filled" icon="information">Red</rh-tag>
  </rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
</rh-accordion>
{{< /raw >}}


### color context

{{< raw >}}
<rh-context-demo>
  <rh-accordion>
    <rh-accordion-header expanded>Item One</rh-accordion-header>
    <rh-accordion-panel>
      <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
      <rh-cta href="#">Call To Action</rh-cta>
      <rh-accordion id="nested">
        <rh-accordion-header expanded>Forced Palette</rh-accordion-header>
        <rh-accordion-panel>
<label> No matter the parent context, this panel should always be
  <rh-context-picker target="nested" value="lightest"></rh-context-picker>
</label>
<rh-cta href="#">Call To Action</rh-cta>
        </rh-accordion-panel>
        <rh-accordion-header>Item Two</rh-accordion-header>
        <rh-accordion-panel>
<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
<rh-cta href="#">Call To Action</rh-cta>
        </rh-accordion-panel>
      </rh-accordion>
    </rh-accordion-panel>
    <rh-accordion-header>Item Two</rh-accordion-header>
    <rh-accordion-panel>
      <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
    </rh-accordion-panel>
    <rh-accordion-header>Item Three</rh-accordion-header>
    <rh-accordion-panel>
      <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
    </rh-accordion-panel>
  </rh-accordion>
</rh-context-demo>
{{< /raw >}}


### expanded attribute

{{< raw >}}
<rh-accordion>
  <h2><rh-accordion-header expanded>Item One</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <h2><rh-accordion-header expanded>Item Two</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <h2><rh-accordion-header>Item Three</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
</rh-accordion>
{{< /raw >}}


### expanded index

{{< raw >}}
<rh-accordion expanded-index="0, 2">
  <h2><rh-accordion-header>Item One</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <h2><rh-accordion-header>Item Two</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <h2><rh-accordion-header>Item Three</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
</rh-accordion>
{{< /raw >}}


### index

{{< raw >}}
<rh-accordion>
  <!-- H2 tags will be removed on upgrade, rh-accordion-header will set the correct heading level internally using the header tag that wraps it -->
  <h2><rh-accordion-header>Item One</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <h2><rh-accordion-header>Item Two</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <h2><rh-accordion-header>Item Three</rh-accordion-header></h2>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
</rh-accordion>
{{< /raw >}}


### large

{{< raw >}}
<rh-accordion large>
  <rh-accordion-header>Item One</rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <rh-accordion-header>Item Two</rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <rh-accordion-header>Item Three</rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
</rh-accordion>
{{< /raw >}}


### nested

{{< raw >}}
<rh-accordion large>
  <rh-accordion-header>Item One</rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
    <rh-cta href="#">Call To Action</rh-cta>
    <rh-accordion>
      <rh-accordion-header>Item One</rh-accordion-header>
      <rh-accordion-panel>
<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
<rh-cta href="#">Call To Action</rh-cta>
      </rh-accordion-panel>
      <rh-accordion-header>Item Two</rh-accordion-header>
      <rh-accordion-panel>
<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
<rh-cta href="#">Call To Action</rh-cta>
      </rh-accordion-panel>
    </rh-accordion>
  </rh-accordion-panel>
  <rh-accordion-header>Item Two</rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
  <rh-accordion-header>Item Three</rh-accordion-header>
  <rh-accordion-panel>
    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </rh-accordion-panel>
</rh-accordion>
{{< /raw >}}


### right to left

{{< raw >}}
<section dir="rtl" lang="he" id="rtl">
  <rh-accordion>
    <rh-accordion-header expanded>תוכנה חופשית</rh-accordion-header>
    <rh-accordion-panel>
      <p> ”תוכנה חופשית“ זה עניין של חירות, לא של מחיר. כדי להבין את העקרון, צריך לחשוב על ”חופש“ כמו ב”חופש הביטו“ ולא כמו ב”בירה חופשי“. </p>
      <rh-cta href="#">קראו עוד</rh-cta>
      <rh-accordion>
        <rh-accordion-header>תוכנה חופשית</rh-accordion-header>
        <rh-accordion-panel>
<p>תוכנה היא תוכנה חופשית אם למשתמשים יש את כל החירויות הללו. לפיכך אתם צריכים להיות חופשיים להפיץ עותקים בהפצת-המשך, עם או בלי שינויים, חינם או בעבור תשלום, לכל אחד בכל מקום. החירות לעשות את הדברים האלו פירושו (בין שאר הדברים) שאינכם חייבים לבקש רשות ו/או לשלם בשבילה. </p>
<rh-cta href="#">קראו עוד</rh-cta>
        </rh-accordion-panel>
        <rh-accordion-header>תוכנה חופשית</rh-accordion-header>
        <rh-accordion-panel>
<p> החופש להריץ את התוכנה פירושו החופש לכל אדם או ארגון להשתמש בתוכנה על כל סוג של מערכת מחשב, לכל מטרה שהיא, ומבלי להדרש ליצור קשר כתוצאה מכך עם המפתח או כל ישות מסוימת אחרת. </p>
<rh-cta href="#">קראו עוד</rh-cta>
        </rh-accordion-panel>
      </rh-accordion>
    </rh-accordion-panel>
    <rh-accordion-header>תוכנה חופשית</rh-accordion-header>
    <rh-accordion-panel>
      <p>" החופש להפיץ עותקים בהפצת-המשך חייב לכלול צורות בינאריות או ניתנות-להרצה של התוכנה, כמו גם את קוד-המקור, לגרסאות שעברו שינוי כמו גם לגרסאות שלא שונו. (הפצת תוכנות בצורה ניתנת-להרצה היא חיונית למערכות הפעלה חופשיות נוחות להתקנה.) זה בסדר אם אין דרך להפיק צורה בינארית או ניתנת-להרצה של תוכנה מסוימת (מאחר ומספר שפות לא תומכות בתכונה הזו), אך חייב להיות לכם החופש להפיץ צורות כאלה בהפצת-המשך במידה ומצאתם או פיתחתם דרך לעשות זאת. </p>
    </rh-accordion-panel>
    <rh-accordion-header>תוכנה חופשית</rh-accordion-header>
    <rh-accordion-panel>
      <p> כדי שהחופש לשנות, והחופש לפרסם גרסאות שעברו שינוי יהיו בעלי משמעות, חייבת להיות לכם גישה לקוד-המקור של התוכנה. לכן נגישות של קוד-המקור היא תנאי הכרחי לתוכנה חופשית. </p>
    </rh-accordion-panel>
  </rh-accordion>
</section>
{{< /raw >}}

