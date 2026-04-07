---
title: "Tooltip"
imports:
  - rh-button
  - rh-icon
  - rh-tooltip
---

<p>9 demos for <code>&lt;rh-tooltip&gt;</code></p>


### bottom


<div class="tooltip-container">
  <rh-tooltip position="bottom">
    <rh-button>Bottom Tooltip</rh-button>
    <span slot="content">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
      labore et dolore magna aliqua. Mi eget mauris pharetra et ultrices.</span>
  </rh-tooltip>
</div>



### color context


<rh-context-demo>
  <rh-tooltip>
    <rh-button>Tooltip</rh-button>
    <span slot="content">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
      labore et dolore magna aliqua. Mi eget mauris pharetra et ultrices.</span>
  </rh-tooltip>
</rh-context-demo>



### content attributes


<rh-tooltip content="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
    labore et dolore magna aliqua. Mi eget mauris pharetra et ultrices.">
  <rh-button>Tooltip</rh-button>
</rh-tooltip>



### index


<rh-tooltip>
  <rh-button>Tooltip</rh-button>
  <span slot="content">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
    labore et dolore magna aliqua. Mi eget mauris pharetra et ultrices.</span>
</rh-tooltip>



### left


<div class="tooltip-container">
  <rh-tooltip position="left">
    <rh-button>Left Tooltip</rh-button>
    <span slot="content">Some tooltip content</span>
  </rh-tooltip>
</div>



### right


<div class="tooltip-container">
<rh-tooltip position="right">
  <rh-button>Right Tooltip</rh-button>
    <span slot="content">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
      labore et dolore magna aliqua. Mi eget mauris pharetra et ultrices.</span>
  </rh-tooltip>
</div>



### rtl


<div class="tooltip-container" dir="rtl" lang="he">
  <rh-tooltip position="right">
    <rh-button>עם ישראל חי</rh-button>
    <span slot="content">
  בְּאֶרֶץ-יִשְׂרָאֵל קָם הָעָם הַיְּהוּדִי, בָּהּ עֻצְּבָה דְּמוּתוֹ הָרוּחָנִית, הַדָּתִית וְהַמְּדִינִית, בָּהּ חַי חַיֵּי קוֹמְמִיּוּת מַמְלַכְתִּית, בָּהּ יָצַר נִכְסֵי תַּרְבּוּת לְאֻמִּיִּים וּכְלַל-אֱנוֹשִׁיִּים וְהוֹרִישׁ לָעוֹלָם כֻּלּוֹ אֶת סֵפֶר הַסְּפָרִים הַנִּצְחִי.
  לְאַחַר שֶׁהֻגְלָה הָעָם מֵאַרְצוֹ בְּכֹחַ הַזְּרוֹעַ שָׁמַר לָהּ אֱמוּנִים בְּכָל אַרְצוֹת פְּזוּרָיו, וְלֹא חָדַל מִתְּפִלָּה וּמִתִּקְוָה לָשׁוּב לְאַרְצוֹ וּלְחַדֵּשׁ בְּתוֹכָהּ אֶת חֵרוּתוֹ הַמְּדִינִית. </span>
  </rh-tooltip>
  <rh-tooltip position="left">
    <rh-button>עם ישראל חי</rh-button>
    <span slot="content">
  בְּאֶרֶץ-יִשְׂרָאֵל קָם הָעָם הַיְּהוּדִי, בָּהּ עֻצְּבָה דְּמוּתוֹ הָרוּחָנִית, הַדָּתִית וְהַמְּדִינִית, בָּהּ חַי חַיֵּי קוֹמְמִיּוּת מַמְלַכְתִּית, בָּהּ יָצַר נִכְסֵי תַּרְבּוּת לְאֻמִּיִּים וּכְלַל-אֱנוֹשִׁיִּים וְהוֹרִישׁ לָעוֹלָם כֻּלּוֹ אֶת סֵפֶר הַסְּפָרִים הַנִּצְחִי.
  לְאַחַר שֶׁהֻגְלָה הָעָם מֵאַרְצוֹ בְּכֹחַ הַזְּרוֹעַ שָׁמַר לָהּ אֱמוּנִים בְּכָל אַרְצוֹת פְּזוּרָיו, וְלֹא חָדַל מִתְּפִלָּה וּמִתִּקְוָה לָשׁוּב לְאַרְצוֹ וּלְחַדֵּשׁ בְּתוֹכָהּ אֶת חֵרוּתוֹ הַמְּדִינִית. </span>
  </rh-tooltip>
</div>



### silent


<p>Adding the <code>silent</code> attribute makes <code>&lt;rh-tooltip&gt;</code>&rsquo;s tooltip content inaccessible. Only use when providing another means of accessibility for content in the <code>content</code> slot.</p>
<rh-tooltip silent>
  <rh-button variant="secondary">
    <rh-icon set="ui" icon="copy" accessible-label="Copy to Clipboard"></rh-icon>
  </rh-button>
  <span slot="content">Copy to Clipboard</span>
</rh-tooltip>



### top


<div class="tooltip-container">
  <rh-tooltip position="top">
    <rh-button>Top Tooltip</rh-button>
    <span slot="content">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
      labore et dolore magna aliqua. Mi eget mauris pharetra et ultrices.</span>
  </rh-tooltip>
</div>


