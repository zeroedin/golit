---
title: "Health index"
imports:
  - rh-health-index
---

<p>4 demos for <code>&lt;rh-health-index&gt;</code></p>


### color context

<rh-context-demo>
  <ul class="demo">
    <li>
      <h2>SM</h2>
      <div><rh-health-index grade="A" size="sm">A</rh-health-index></div>
      <div><rh-health-index grade="B" size="sm">B</rh-health-index></div>
      <div><rh-health-index grade="C" size="sm">C</rh-health-index></div>
      <div><rh-health-index grade="D" size="sm">D</rh-health-index></div>
      <div><rh-health-index grade="E" size="sm">E</rh-health-index></div>
      <div><rh-health-index grade="F" size="sm">F</rh-health-index></div>
    </li>
    <li>
      <h2>Default</h2>
      <div><rh-health-index grade="A">A</rh-health-index></div>
      <div><rh-health-index grade="B">B</rh-health-index></div>
      <div><rh-health-index grade="C">C</rh-health-index></div>
      <div><rh-health-index grade="D">D</rh-health-index></div>
      <div><rh-health-index grade="E">E</rh-health-index></div>
      <div><rh-health-index grade="F">F</rh-health-index></div>
    </li>
    <li>
      <h2>LG</h2>
      <div><rh-health-index grade="A" size="lg">A</rh-health-index></div>
      <div><rh-health-index grade="B" size="lg">B</rh-health-index></div>
      <div><rh-health-index grade="C" size="lg">C</rh-health-index></div>
      <div><rh-health-index grade="D" size="lg">D</rh-health-index></div>
      <div><rh-health-index grade="E" size="lg">E</rh-health-index></div>
      <div><rh-health-index grade="F" size="lg">F</rh-health-index></div>
    </li>
    <li>
      <h2>XL</h2>
      <div><rh-health-index grade="A" size="xl">A</rh-health-index></div>
      <div><rh-health-index grade="B" size="xl">B</rh-health-index></div>
      <div><rh-health-index grade="C" size="xl">C</rh-health-index></div>
      <div><rh-health-index grade="D" size="xl">D</rh-health-index></div>
      <div><rh-health-index grade="E" size="xl">E</rh-health-index></div>
      <div><rh-health-index grade="F" size="xl">F</rh-health-index></div>
    </li>
  </ul>
</rh-context-demo>


### index

<rh-health-index grade="C"
                 size="lg">C</rh-health-index>


### screen readers

<section>
  <h2>&lt;rh-health-index> element</h2>
  <rh-health-index grade="A" size="sm">A</rh-health-index>
  <rh-health-index grade="B" size="sm">B</rh-health-index>
  <rh-health-index grade="C" size="sm">C</rh-health-index>
  <rh-health-index grade="D" size="md">D</rh-health-index>
  <rh-health-index grade="E" size="lg">E</rh-health-index>
  <rh-health-index grade="F" size="xl">F</rh-health-index>
</section>
<section>
  <h2>ARIA attrs for <abbr title="localization">l10n</abbr></h2>
  <rh-health-index grade="A"
      aria-label="בריאות מדורג מדרגות א עד ו"
      aria-valuetext="דרגה א"
      size="sm">A</rh-health-index>
</section>
<section>
  <h2>Native &gt;meter> element</h2>
  <meter min="1" max="6" value="1" aria-valuetext="Grade A" aria-label="Health graded A to F"></meter>
  <meter min="1" max="6" value="2" aria-valuetext="Grade B" aria-label="Health graded A to F"></meter>
  <meter min="1" max="6" value="3" aria-valuetext="Grade C" aria-label="Health graded A to F"></meter>
  <meter min="1" max="6" value="4" aria-valuetext="Grade D" aria-label="Health graded A to F"></meter>
  <meter min="1" max="6" value="5" aria-valuetext="Grade E" aria-label="Health graded A to F"></meter>
  <meter min="1" max="6" value="6" aria-valuetext="Grade F" aria-label="Health graded A to F"></meter>
</section>
<section>
  <h2>ARIA-attributes &gt;meter> element</h2>
  <div role="meter" aria-valuemin="1" aria-valuemax="6" aria-valuenow="1" aria-valuetext="Grade A" aria-label="Health graded A to F">A</div>
  <div role="meter" aria-valuemin="1" aria-valuemax="6" aria-valuenow="1" aria-valuetext="Grade B" aria-label="Health graded A to F">B</div>
  <div role="meter" aria-valuemin="1" aria-valuemax="6" aria-valuenow="1" aria-valuetext="Grade C" aria-label="Health graded A to F">C</div>
  <div role="meter" aria-valuemin="1" aria-valuemax="6" aria-valuenow="1" aria-valuetext="Grade D" aria-label="Health graded A to F">D</div>
  <div role="meter" aria-valuemin="1" aria-valuemax="6" aria-valuenow="1" aria-valuetext="Grade E" aria-label="Health graded A to F">E</div>
  <div role="meter" aria-valuemin="1" aria-valuemax="6" aria-valuenow="1" aria-valuetext="Grade F" aria-label="Health graded A to F">F</div>
</section>


### variants

<ul class="demo">
  <li>
    <h2>SM</h2>
    <rh-health-index grade="A" size="sm">A</rh-health-index>
    <rh-health-index grade="B" size="sm">B</rh-health-index>
    <rh-health-index grade="C" size="sm">C</rh-health-index>
    <rh-health-index grade="D" size="sm">D</rh-health-index>
    <rh-health-index grade="E" size="sm">E</rh-health-index>
    <rh-health-index grade="F" size="sm">F</rh-health-index>
  </li>
  <li>
    <h2>Default</h2>
    <rh-health-index grade="A">A</rh-health-index>
    <rh-health-index grade="B">B</rh-health-index>
    <rh-health-index grade="C">C</rh-health-index>
    <rh-health-index grade="D">D</rh-health-index>
    <rh-health-index grade="E">E</rh-health-index>
    <rh-health-index grade="F">F</rh-health-index>
  </li>
  <li>
    <h2>LG</h2>
    <rh-health-index grade="A" size="lg">A</rh-health-index>
    <rh-health-index grade="B" size="lg">B</rh-health-index>
    <rh-health-index grade="C" size="lg">C</rh-health-index>
    <rh-health-index grade="D" size="lg">D</rh-health-index>
    <rh-health-index grade="E" size="lg">E</rh-health-index>
    <rh-health-index grade="F" size="lg">F</rh-health-index>
  </li>
  <li>
    <h2>XL</h2>
    <rh-health-index grade="A" size="xl">A</rh-health-index>
    <rh-health-index grade="B" size="xl">B</rh-health-index>
    <rh-health-index grade="C" size="xl">C</rh-health-index>
    <rh-health-index grade="D" size="xl">D</rh-health-index>
    <rh-health-index grade="E" size="xl">E</rh-health-index>
    <rh-health-index grade="F" size="xl">F</rh-health-index>
  </li>
</ul>

