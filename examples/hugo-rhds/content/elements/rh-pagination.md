---
title: "Pagination"
imports:
  - rh-button
  - rh-pagination
lightdom:
  - rh-pagination-lightdom.css
---

<p>14 demos for <code>&lt;rh-pagination&gt;</code></p>


### aria current

<rh-pagination>
  <ol>
    <li><a href="#1">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3" aria-current="page">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### color context

<rh-context-demo>
  <rh-pagination>
    <ol>
      <li><a href="#">1</a></li>
      <li><a href="#2">2</a></li>
      <li><a href="#3">3</a></li>
      <li><a href="#4">4</a></li>
      <li><a href="#5">5</a></li>
    </ol>
  </rh-pagination>
  <rh-pagination variant="open">
    <ol>
      <li><a href="#">1</a></li>
      <li><a href="#2">2</a></li>
      <li><a href="#3">3</a></li>
      <li><a href="#4">4</a></li>
      <li><a href="#5">5</a></li>
    </ol>
  </rh-pagination>
  <rh-pagination size="sm">
    <ol>
      <li><a href="#">1</a></li>
      <li><a href="#2">2</a></li>
      <li><a href="#3">3</a></li>
      <li><a href="#4">4</a></li>
      <li><a href="#5">5</a></li>
    </ol>
  </rh-pagination>
</rh-context-demo>


### compact

<rh-pagination id="constrain">
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### index

<rh-pagination>
  <ol>
    <li><a href="?page=1">1</a></li>
    <li><a href="?page=2">2</a></li>
    <li><a href="?page=3">3</a></li>
    <li><a href="?page=4">4</a></li>
    <li><a href="?page=5">5</a></li>
  </ol>
</rh-pagination>


### many pages

<rh-pagination>
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
    <li><a href="#6">6</a></li>
    <li><a href="#7">7</a></li>
    <li><a href="#8">8</a></li>
    <li><a href="#9">9</a></li>
    <li><a href="#10">10</a></li>
    <li><a href="#11">11</a></li>
    <li><a href="#12">12</a></li>
    <li><a href="#13">13</a></li>
    <li><a href="#14">14</a></li>
    <li><a href="#15">15</a></li>
    <li><a href="#16">16</a></li>
    <li><a href="#17">17</a></li>
    <li><a href="#18">18</a></li>
    <li><a href="#19">19</a></li>
    <li><a href="#20">20</a></li>
  </ol>
</rh-pagination>
<p>Paginators with many pages must overflow.</p>


### no numeric control

<rh-pagination>
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### open compact size

<rh-pagination id="constrain" variant="open" size="sm">
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### open compact

<rh-pagination id="constrain" variant="open">
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### open

<rh-pagination variant="open">
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### overflow

<rh-pagination>
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>
<p>Paginators with 5 or fewer pages should not overflow, meaning all links should be visible.
  Once a paginator has more than 5 pages, then it must overflow, meaning some links will be hidden.
  Paginators with 9 or more pages will overflow on boths ends.</p>
<dl>
  <dt>With 5 or fewer pages</dt>
  <dd>No overflow</dd>
  <dt>With more than 5 but fewer than 9 pages</dt>
  <dd>Overflow on one side</dd>
  <dt>With more than 9 pages, active page is less than 6</dt>
  <dd>Overflow end</dd>
  <dt>With more than 9 pages, active page is greater than 6</dt>
  <dd>Overflow both</dd>
  <dt>With more than 9 pages, active page is greater than 5 less than the total (e.g. 16/20)</dt>
  <dd>Overflow start</dd>
</dl>
<fieldset>
  <legend>Adjust pages</legend>
  <rh-button id="add">Add Page</rh-button>
  <rh-button id="remove" danger>Remove Page</rh-button>
</fieldset>


### right to left

<p>צריך להיראות יותר טוב</p>
<rh-pagination id="rtl-pagination" dir="rtl">
  <span slot="go-to-page">עבור לדף</span>
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### size compact

<rh-pagination id="constrain" size="sm">
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### size

<rh-pagination size="sm">
  <ol>
    <li><a href="#">1</a></li>
    <li><a href="#2">2</a></li>
    <li><a href="#3">3</a></li>
    <li><a href="#4">4</a></li>
    <li><a href="#5">5</a></li>
  </ol>
</rh-pagination>


### vue

<div id="app">
  <rh-pagination>
    <ol>
      <li v-for="pageNum in pages" :key="pageNum">
        <a :href="`#${pageNum}`" :aria-current="pageNum == currentPage ? 'page' : null" v-text="pageNum"></a>
      </li>
    </ol>
  </rh-pagination>
  <span v-text="message"></span> <span v-text="currentPage"></span>
</div>

