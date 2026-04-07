---
title: "Timestamp"
imports:
  - rh-timestamp
  - rh-tooltip
---

<p>6 demos for <code>&lt;rh-timestamp&gt;</code></p>


### custom format


<rh-timestamp id="timestamp-custom-format"
              date="Sat Jan 01 2022 00:00:00 GMT-0500"></rh-timestamp>



### fallback


<rh-timestamp date="Tue Aug 09 2006 14:57:00 GMT-0400">Tue Aug 09 2006 14:57:00 GMT-0400</rh-timestamp>



### formats


<rh-timestamp date-format="full"
              time-format="full"></rh-timestamp>
<rh-timestamp date-format="full"></rh-timestamp>
<rh-timestamp time-format="full"></rh-timestamp>
<rh-timestamp date-format="medium"
              time-format="short"
              display-suffix="US Eastern"></rh-timestamp>
<rh-timestamp date-format="full"
              locale="es"></rh-timestamp>



### index


<rh-timestamp date="Tue Aug 09 2006 14:57:00 GMT-0400"></rh-timestamp>



### relative


<rh-timestamp date="Tue Aug 09 2006 14:57:00 GMT-0400 (Eastern Daylight Time)" relative></rh-timestamp>
<rh-timestamp date="Tue Aug 09 2006 14:57:00 GMT-0400 (Eastern Daylight Time)" locale="es" relative></rh-timestamp>
<rh-timestamp date="Tue Aug 09 2022 14:57:00 GMT-0400 (Eastern Daylight Time)" relative></rh-timestamp>
<rh-timestamp date="Tue Aug 09 2022 14:57:00 GMT-0400 (Eastern Daylight Time)" locale="es" relative></rh-timestamp>
<rh-timestamp date="Tue Aug 09 2099 14:57:00 GMT-0400 (Eastern Daylight Time)" relative></rh-timestamp>
<rh-timestamp date="Tue Aug 09 2099 14:57:00 GMT-0400 (Eastern Daylight Time)" locale="es" relative></rh-timestamp>



### tooltip


<section>
  <h2>Default tooltip</h2>
  <p>
    <rh-tooltip>
      <rh-timestamp></rh-timestamp>
      <rh-timestamp slot="content" utc></rh-timestamp>
    </rh-tooltip>
  </p>
  <p>
    <rh-tooltip>
      <rh-timestamp></rh-timestamp>
      <rh-timestamp slot="content" utc display-suffix="Coordinated Universal Time"></rh-timestamp>
    </rh-tooltip>
  </p>
</section>
<section>
  <h2>Custom tooltip</h2>
  <p>
    <rh-tooltip>
      <rh-timestamp date="Tue Aug 09 2022 14:57:00 GMT-0400 (Eastern Daylight Time)"></rh-timestamp>
      <span slot="content">Last updated on <rh-timestamp date="Tue Aug 09 2022 14:57:00 GMT-0400 (Eastern Daylight Time)" date-format="long" time-format="short" utc></rh-timestamp></span>
    </rh-tooltip>
  </p>
  <p>
    <rh-tooltip>
      Halloween
      <rh-timestamp slot="content" date="Mon Oct 31 2022 00:00:00 GMT-0400 (Eastern Daylight Time)" date-format="medium"></rh-timestamp>
    </rh-tooltip>
  </p>
</section>
<section>
  <h2>Relative with tooltip</h2>
  <p>
    <rh-tooltip>
      <rh-timestamp date="Tue Aug 09 2022 14:57:00 GMT-0400 (Eastern Daylight Time)" relative></rh-timestamp>
      <rh-timestamp slot="content" date="Tue Aug 09 2022 14:57:00 GMT-0400 (Eastern Daylight Time)"></rh-timestamp>
    </rh-tooltip>
  </p>
  <p>
    <rh-tooltip>
      <rh-timestamp date="Aug 09 2024 14:57:00 GMT-0400 (Eastern Daylight Time)" relative></rh-timestamp>
      <rh-timestamp slot="content" date="Aug 09 2024 14:57:00 GMT-0400 (Eastern Daylight Time)"></rh-timestamp>
    </rh-tooltip>
  </p>
</section>


