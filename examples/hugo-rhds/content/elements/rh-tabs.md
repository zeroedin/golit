---
title: "Tabs"
imports:
  - rh-cta
  - rh-icon
  - rh-tabs
lightdom:
  - rh-cta-lightdom-shim.css
---

<p>11 demos for <code>&lt;rh-tabs&gt;</code></p>


### box inset

{{< raw >}}
<rh-tabs box="inset"
         id="inset">
  <rh-tab id="users"
          slot="tab"
          active>Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Containers</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
{{< /raw >}}


### centered

{{< raw >}}
<rh-tabs centered>
  <rh-tab id="users"
          slot="tab"
          active>Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Containers</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
{{< /raw >}}


### color context

{{< raw >}}
<rh-context-demo>
  <form>
    <fieldset>
      <legend>Box Layout</legend>
      <label>
        <input type="radio"
               name="variant"
               value=""> None
      </label>
      <label>
        <input type="radio"
               name="variant"
               value="inset"
               checked> Inset
      </label>
      <label>
        <input type="radio"
               name="variant"
               value="box"> Box
      </label>
    </fieldset>
    <fieldset>
      <legend>Layout</legend>
      <label>
        <input type="checkbox"
               name="centered"> Centered
      </label>
      <label>
        <input type="checkbox"
               name="vertical"> Vertical
      </label>
    </fieldset>
  </form>
  <rh-tabs>
    <rh-tab id="users"
            slot="tab"
            active>Users</rh-tab>
    <rh-tab-panel>
      <p>The control of <a href="https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/system_administrators_guide/ch-managing_users_and_groups">users and groups</a> is a core
        element of Red Hat Enterprise Linux system administration. This chapter explains how to add, manage, and delete users and groups in the graphical user
        interface and on the command line, and covers advanced topics, such as creating group directories. </p>
      <rh-cta>
        <a href="https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/system_administrators_guide/ch-managing_users_and_groups">Read about user management</a>
      </rh-cta>
    </rh-tab-panel>
    <rh-tab slot="tab">Containers</rh-tab>
    <rh-tab-panel>
      <p>Containers are technologies that allow the packaging and isolation of applications with their entire runtime environment—all of the files necessary to run. This makes it easy to move the
        contained application between environments (dev, test, production, etc.) while retaining full functionality. Containers are also an important part of <a href="/en/topics/security">IT
          security</a>. By <a href="/en/topics/security/container-security">building security into the container pipeline</a>&nbsp;and defending infrastructure, containers stay reliable, scalable, and
        trusted. You can also easily move the containerized application between public, private and hybrid cloud environments and data centers (or on-premises) with consistent behavior and
        functionality.</p>
      <rh-cta>
        <a href="https://www.redhat.com/en/topics/containers/whats-a-linux-container">Continue reading</a>
      </rh-cta>
    </rh-tab-panel>
    <rh-tab slot="tab">Database</rh-tab>
    <rh-tab-panel>
      <p>Modern application development teams must use multiple cloud databases to access various features and capabilities. As database usage continues to increase, it is difficult using and manage
        this multitude of databases. Red Hat commissioned Forrester Consulting found that simplified managed database access solutions help decision-makers maximize the impact, mitigate the complexity
        of their database-as-a-service (DBaaS) solutions, and manage their diverse landscape of databases.</p>
      <rh-cta>
        <a href="https://www.redhat.com/rhdc/managed-files/cl-maximize-with-a-managed-database-access-solution-analyst-material-f32006-202208-en.pdf">Download</a>
      </rh-cta>
    </rh-tab-panel>
    <rh-tab slot="tab">Servers</rh-tab>
    <rh-tab-panel>Servers</rh-tab-panel>
    <rh-tab slot="tab">Cloud</rh-tab>
    <rh-tab-panel>Cloud</rh-tab-panel>
  </rh-tabs>
</rh-context-demo>
{{< /raw >}}


### deprecation

{{< raw >}}
<rh-tabs theme="base">
  <rh-tab id="users"
          slot="tab"
          active>Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Containers</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
<rh-tabs id="new">
  <rh-tab id="users"
          slot="tab">Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Containers</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
{{< /raw >}}


### icons and text

{{< raw >}}
<rh-tabs>
  <rh-tab slot="tab"
          active> Users <rh-icon slot="icon"
             icon="users"
             set="ui"></rh-icon></rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab"> Containers <rh-icon slot="icon"
             icon="container"></rh-icon></rh-tab>
  <rh-tab-panel>Containers</rh-tab-panel>
  <rh-tab slot="tab"> Database: Long SQL Statement <rh-icon slot="icon"
             icon="datacenter"></rh-icon></rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab"> Server <rh-icon slot="icon"
             icon="server"></rh-icon></rh-tab>
  <rh-tab-panel>Server</rh-tab-panel>
  <rh-tab slot="tab"> No Icon </rh-tab>
  <rh-tab-panel>System</rh-tab-panel>
</rh-tabs>
{{< /raw >}}


### index

{{< raw >}}
<rh-tabs>
  <rh-tab id="users"
          slot="tab"
          active>Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Containers</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
{{< /raw >}}


### long tab content

{{< raw >}}
<rh-tabs>
  <rh-tab id="users"
          slot="tab"
          active>Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Red Hat OpenShift Service on AWS</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
<rh-tabs vertical>
  <rh-tab id="users"
          slot="tab">Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Red Hat OpenShift Service on AWS</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
{{< /raw >}}


### manual activation

{{< raw >}}
<rh-tabs manual>
  <rh-tab id="users" slot="tab">Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Containers</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
{{< /raw >}}


### nested

{{< raw >}}
<rh-tabs>
  <rh-tab id="users"
          slot="tab"
          active>Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Containers</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Nested</rh-tab>
  <rh-tab-panel>
    <rh-tabs>
      <rh-tab slot="tab">Nested Users</rh-tab>
      <rh-tab-panel>Users</rh-tab-panel>
      <rh-tab slot="tab">Nested Containers</rh-tab>
      <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
      <rh-tab slot="tab">Nested Database</rh-tab>
      <rh-tab-panel>Database</rh-tab-panel>
      <rh-tab slot="tab">Nested Servers</rh-tab>
      <rh-tab-panel>Servers</rh-tab-panel>
      <rh-tab slot="tab">Nested Cloud</rh-tab>
      <rh-tab-panel>Cloud</rh-tab-panel>
    </rh-tabs>
  </rh-tab-panel>
</rh-tabs>
{{< /raw >}}


### right to left

{{< raw >}}
<section dir="rtl"
         lang="he"
         id="rtl">
  <p>מאת <a href="https://www.gnu.org/philosophy/shouldbefree.he.html">"מדוע על תוכנה להיות חופשית" </a>- המוסד לתוכנה חופשית</p>
  <rh-tabs>
    <rh-tab slot="tab"
            active>הקדמה</rh-tab>
    <rh-tab-panel>
      <p>קיום מוצר ה”תוכנה“ מעלה את השאלה על דרך קבלת ההחלטות בקשר אליה. לדוגמא, נניח כי אדם מסוים, בעל עותק של תוכנה, פוגש אדם אחר המעונין בעותק שלה. האם הם יכולים להעתיק את התוכנה? מי צריך להחליט על
        כך? האם האנשים המעורבים, או שמא ישות נוספת, ”הבעלים“ של המוצר?</p>
      <p>מפתחי תוכנה בדרך כלל מתייחסים לשאלות אלו בהנחה שהקריטריון לתשובה הוא מיקסום רווחים למפתח. הכוח הפוליטי של עולם העסקים הוביל את השלטון לאמץ הן את הקריטריון והן את התשובה שהוצעה : למוצר תוכנה
        יש בעלים, שהם בד"כ חברה המזוהה עם פיתוח המוצר. </p>
      <p>במאמר זה ננסה לבחון את אותה שאלה, אבל עם קריטריון שונה: השגשוג והחופש של הציבור. </p>
      <p>התשובה אינה יכולה להתבסס על החקיקה הנוכחית – החוק צריך לעקוב אחרי האתיקה ולא להפך, ובכלל, החקיקה הנוכחית אינה פותרת את השאלה, למרות שהיא מציעה מספר פתרונות אפשריים. הדרך היחידה לשפוט בנושא
        היא לגלות מי נפגע ומי מרוויח בהכרה בבעלות על תוכנה, למה, ובאיזו מידה. במילים אחרות, עלינו לבצע אנליזת מחיר/תועלת בשם החברה כישות אחת, תוך שאנו לוקחים בחשבון את החופש האישי ובנוסף את היצור של
        מוצרים ברי קיימא. </p>
      <p>במאמר זה יתוארו האפקטים של ”בעלות“, ונראה כי התוצאה היא מזיקה בעליל. המסקנה המתבקשת היא כי למתכנתים יש את החובה לעודד אחרים לחלוק, להפיץ מחדש, ללמוד ולשפר תוכנה שהם כותבים, או במילים אחרות,
        לכתוב תוכנה חופשית. </p>
    </rh-tab-panel>
    <rh-tab slot="tab">כיצד בעלי קניין מצדיקים את כוחם</rh-tab>
    <rh-tab-panel>
      <p>אלו הנהנים מהמערכת הנוכחית, בה תוכנות הנן רכוש, מציעים שני נימוקים לתמיכה
        בדרישתם לבעלות על תוכנה: הנימוק הרגשי והנימוק הכלכלי.</p>
      <p>הנימוק הרגשי פשוט: ”השקעתי מאמצים, זמן ונשמה במוצר, המוצר בא <em>ממני</em>,
        הוא <em>שלי</em>!“</p>
      <p>נימוק זה אינה דורש מאמץ רב להפרכתו. תחושת ה-”קשר הרגשי“ אינה בלתי-נמנעת אלא
        כזו שמתכנתים יכולים לטפח כאשר זה מתאים להם. ניקח לדוגמא, כיצד תחושת הקשר
        נעלמת במסתוריות במקרה הנפוץ בו מתכנתים מוכנים, בד"כ, להעביר את כל זכויותיהם
        לתאגיד גדול תמורת שכר. בניגוד לכך, שיקלו את גישתם של האמנים והאומנים הגדולים
        של ימי הביניים, שאפילו לא חתמו על עבודותיהם. להם, שם האמן לא היה חשוב. מה
        שכן היה חשוב היה הצורך שהעבודה תשרת והעובדה שהיא נעשתה. גישה זה שלטה מאות
        שנים.</p>
      <p>הנימוק הכלכלי מתואר כך : ”אני רוצה להתעשר (או, בתיאור הפופולרי אך לא מדויק,
        להתקיים), אם לא אוכל על ידי פיתוח תוכנה, אז לא אפתח תוכנה. כל האחרים הם
        כמוני, ולכן אף אחד לא יתכנת. ואז מה תעשו?“ – איום זה בד"כ מוסווה כעצה
        ידידותית.</p>
      <p>בהמשך נראה מדוע איום זה הנו איום סרק. אבל תחילה נפנה לטפל בהנחה סמויה שניתן
        לגלות ביתר קלות בניסוח אחר של הטיעון הנ"ל.</p>
      <p>ניסוח זה מתחיל בהשוואת התועלת הציבורית של תוכנה קניינית לעומת אי קיום תוכנה,
        ומסיים במסקנה כי פיתוח תוכנה קניינית הנו, למרות הכל, משתלם, ויש
        לעודדו. הטעות פה הנה ההשוואה של שתי תוצאות בלבד, תוכנה קניינית
        ואי-תוכנה. וההנחה שאין אפשרויות נוספות.</p>
      <p>בהינתן מערכת של זכויות יוצרים, פיתוח תוכנה בד"כ מקושר עם קיום בעלים השולט על
        השימוש בתוכנה. כל עוד קישור זה קיים, אנו לרוב עומדים בפני הבחירה בין תוכנה
        קניינית וכלום. בכל אופן, הקישור הזה אינו טבעי או בלתי-נמנע, הוא תוצאה של
        מדיניות משפטית/חברתית שעליה אנו מערערים: ההחלטה על קנייניות. הגדרת הברירה
        כבחירה בין תוכנה קניינית לאי קיום תוכנה ממש דורש מאיתנו לשאול</p>
    </rh-tab-panel>
    <rh-tab slot="tab">הטיעון נגד קנייניות</rh-tab>
    <rh-tab-panel>
      <p>השאלה שעל הפרק – ”האם פיתוח תוכנה צריך להיות קשור לקיום בעלות והגבלת שימוש?“</p>
      <p>על מנת לפסוק בסוגיה זו, אנו צריכים למדוד את ההשפעה של כל אחת מהפעילויות על
        החברה באופן בלתי תלוי. ההשפעה של פיתוח תוכנה (ללא קשר לשיטות ההפצה שלה)
        וההשפעה של הגבלת השימוש (בהנחה שהתוכנה פותחה). אם אחת מהפעילויות הללו הנה
        חיובית והשניה הרסנית, מוטב לנו לבטל את הקישור בינן ולבצע רק את החיובית.</p>
      <p>או במילים אחרות – אם הגבלת הפצה של תוכנה פוגע בחברה, אזי מתכנת מוסרי ימנע
        מאפשרות זו.</p>
      <p>למדידת התוצאה של הגבלת ההפצה, אנו צריכים להשוות את הערך לחברה של תוכנה קניינית (מוגבלת הפצה) והערך של אותה תוכנה, הזמינה לכולם – כלומר השוואה של
        שני עולמות שונים.</p>
      <p>השוואה זו מטפלת גם בנימוק הנגדי הפשוט שעולה מפעם לפעם - ”התועלת לשכן שלו נתת
        עותק תוכנה מבטלת את התועלת לבעלי התוכנה“ – נימוק זה מניח שהנזק והתועלת שווים
        בערכם. בהשוואה שנבצע נשווה גם את הערכים, ונראה כי התועלת גדולה לעין שיעור
        מהנזק.</p>
      <p>להבהרת הטיעון, בואו ננסה להפעיל אותו על נושא אחר – בניית כבישים.</p>
      <p>ניתן לממן את בניית כל הכבישים על ידי אגרת-מעבר. מימון כזה ידרוש הקמת נקודות
        גבייה בכל צומת רחובות. מערכת כזו תהווה תמריץ גדול לשיפור הדרכים ותגרום לכל
        משתמש בדרך לשלם על השימוש בה, אולם, נקודת גבייה היא מכשול מלאכותי המפריע
        לזרימת התנועה. מלאכותי מכיוון שאינו תוצאה של דרך פעולתן של מכוניות או
        כבישים.</p>
      <p>בהשוואת דרכים חופשיות לכבישי אגרה אנו מוצאים (בהנחה כי כל שאר הפרמטרים זהים)
        כי כבישים רגילים זולים יותר לבנייה ולתחזוקה, בטוחים יותר ויותר יעילים
        בשימוש. <a href="#f2">(2)</a> במדינה ענייה, אגרות מעבר מונעות את השימוש
        בכבישים מאזרחים רבים. יוצא מזה כי כבישים נטולי אגרה (חופשיים) מציעים יותר
        תועלת לחברה בעלות קטנה יותר והם מועדפים על ידי החברה. לכן, על החברה לבחור
        לממן סלילת כבישים בדרכים אחרות, ולא על ידי נקודות גביית אגרה. השימוש בכביש,
        אחרי בנייתו, צריך להיות חופשי.</p>
      <p>כאשר הסניגורים של כבישי האגרה טוענים כי הם רק דרך לגיוס כספים, הם מסתירים את
        העובדה שיש חלופות. נקודות גביית אגרה אכן מגייסות כספים, אבל הן גם מורידות את
        רמת הכביש. כבישי האגרה אינם טובים כמו כבישים חופשיים.</p>
      <p>כמובן, בניית כביש חופשי עולה כסף, שעל הציבור לשלם. אולם, אין משמעות הדבר
        מעבר לשיטת כבישי האגרה. אנו, הנדרשים לשלם בכל מקרה, נקבל תמורה טובה יותר
        לכספנו בבניית כבישים חופשיים.</p>
      <p>איננו טוענים כי דרך אגרה גרועה מאי קיום דרך, למרות שטענה זו תתקיים אם האגרה
        תהיה גבוהה במידה כזאת שלא נוכל להשתמש בדרך (זו אינה מדיניות הגיונית לבעלי
        הכביש). אולם, כל עוד כבישי אגרה גורמים לבזבוז ואי נוחות, עדיף לגייס את הכסף
        בדרכים פחות פוגעניות.</p>
      <p>ניתן להפעיל אותם טיעונים לגבי פיתוח תוכנה, נראה כי קיום ”נקודות גבייה“
        לתוכנה שימושית עולה רבות לחברה – גורם לתוכנה להיות יקרה יותר לפיתוח, יקרה
        יותר להפצה ופחות מספקת ויעילה לשימוש. מזה יוצא כי יש לעודד פיתוח תוכנה
        בדרכים אחרות, ויתוארו גישות אחדות לעידוד ומימון (בהיקף הנדרש באמת) של פיתוח
        תוכנה.</p>
    </rh-tab-panel>
  </rh-tabs>
</section>
{{< /raw >}}


### vertical

{{< raw >}}
<rh-tabs vertical>
  <rh-tab id="users"
          slot="tab"
          active>Users</rh-tab>
  <rh-tab-panel>Users</rh-tab-panel>
  <rh-tab slot="tab">Containers</rh-tab>
  <rh-tab-panel>Containers <a href="#">Focusable element</a></rh-tab-panel>
  <rh-tab slot="tab">Database</rh-tab>
  <rh-tab-panel>Database</rh-tab-panel>
  <rh-tab slot="tab">Servers</rh-tab>
  <rh-tab-panel>Servers</rh-tab-panel>
  <rh-tab slot="tab">Cloud</rh-tab>
  <rh-tab-panel>Cloud</rh-tab-panel>
</rh-tabs>
{{< /raw >}}

