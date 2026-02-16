---
title: "Back to top"
imports:
  - rh-back-to-top
---

<p>5 demos for <code>&lt;rh-back-to-top&gt;</code></p>


### always visible

<rh-back-to-top visible="always" href="#top">Back to top</rh-back-to-top>


### color context

<rh-context-demo>
  <div id="overflow">
    <a id="top" href="#bottom">Go to bottom</a>
    <p>Scroll down or press tab to see the back to top button</p>
    <a id="bottom" href="#top">Go to top</a>
  <div>
  <rh-back-to-top href="#top">Back to top</rh-back-to-top>
</rh-context-demo>


### index

<div id="overflow">
  <p>Scroll down to reveal the back to top element</p>
  <rh-back-to-top href="#">Back to top</rh-back-to-top>
</div>


### no slotted text

<div id="overflow">
  <a id="top" href="#bottom">Go to bottom</a>
  <p>Scroll down or press tab to see the back to top button</p>
  <a id="bottom" href="#top">Go to top</a>
  <rh-back-to-top href="#top" label="Return to top"></rh-back-to-top>
</div>


### scroll distance

<div id="overflow">
  <a id="top" href="#bottom">Go to bottom</a>
  <p>Scroll down (50px) or press tab to see the back to top button</p>
  <a id="bottom" href="#top">Go to top</a>
  <rh-back-to-top scroll-distance="50" href="#top">Back to top</rh-back-to-top>
</div>

