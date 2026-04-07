---
title: "Disclosure"
imports:
  - rh-disclosure
  - rh-icon
  - rh-jump-links
lightdom:
  - rh-disclosure-lightdom-shim.css
---

<p>7 demos for <code>&lt;rh-disclosure&gt;</code></p>


### color context


<rh-context-demo>
  <rh-disclosure summary="Collapsed panel title">
    <p>Lorem ipsum dolor sit amet consectetur adipisicing, elit. Velit distinctio, nesciunt nobis sit.</p>
  </rh-disclosure>
</rh-context-demo>



### events


<form id="disclosure-events">
  <rh-disclosure summary="Collapsed panel title">
    <p>Lorem ipsum dolor <a href="#">sit amet consectetur</a> adipisicing, elit. Velit distinctio, nesciunt nobis sit, a dolor, non numquam rerum recusandae, deserunt enim assumenda quidem. Id impedit necessitatibus obcaecati ratione reprehenderit laborum?</p>
  </rh-disclosure>
  <fieldset>
    <legend>Events Fired</legend>
    <output name="events">No events yet</output>
  </fieldset>
</form>



### index


<rh-disclosure summary="Collapsed panel title">
  <p>Lorem ipsum dolor sit amet consectetur adipisicing, elit. Velit distinctio, nesciunt nobis sit, a dolor, non numquam rerum recusandae, deserunt enim assumenda quidem. Id impedit necessitatibus obcaecati ratione reprehenderit laborum?</p>
</rh-disclosure>



### nested disclosures


<rh-disclosure summary="This is the top level disclosure">
  <p>Try not to nest disclosures. If nesting, consider using an accordion instead. Lorem ipsum <a href="#">dolor</a>.</p>
  <rh-disclosure summary="This is the second level disclosure">
    <p>Test using the ESC key + focus on an element when nesting disclosures. <a href="#">Fake link 2</a> and more text.</p>
    <form action="#" method="get" class="form-example" style="margin-block: var(--rh-space-md, 8px);">
      <div class="form-example">
        <label for="name">Enter your name: </label>
        <input type="text" name="name" id="name" required />
      </div>
      <div class="form-example">
        <label for="favcity">Which is your favorite city?</label>
        <select id="favcity" name="select">
        <option value="1">Amsterdam</option>
        <option value="2">Buenos Aires</option>
        </select>
      </div>
      <div class="form-example">
        <fieldset>
        <legend>Choose a shipping method:</legend>
        <input id="overnight" type="radio" name="shipping" value="overnight">
        <label for="overnight">Overnight</label><br>
        <input id="twoday" type="radio" name="shipping" value="twoday">
        <label for="twoday">Two day</label><br>
        </fieldset>
      </div>
      <div class="form-example">
        <fieldset>
        <legend>Select your pizza toppings:</legend>
        <input id="ham" type="checkbox" name="toppings" value="ham">
        <label for="ham">Ham</label><br>
        <input id="pepperoni" type="checkbox" name="toppings" value="pepperoni">
        <label for="pepperoni">Pepperoni</label><br>
        </fieldset>
      </div>
      <div class="form-example">
        <input type="submit" value="Subscribe!" />
      </div>
    </form>
    <p>This is a sentence with <a href="#">a link</a>.</p>
  </rh-disclosure>
</rh-disclosure>



### nested jump links


<rh-disclosure>
  <h2 slot="summary" id="sections">Sections</h2>
  <rh-jump-links aria-labelledby="sections">
    <rh-jump-link href="#section-1">Section 1</rh-jump-link>
    <rh-jump-link href="#section-2">Section 2</rh-jump-link>
    <rh-jump-link href="#section-3">Section 3</rh-jump-link>
    <rh-jump-link href="#section-4">Section 4</rh-jump-link>
    <rh-jump-link href="#section-5">Section 5</rh-jump-link>
  </rh-jump-links>
</rh-disclosure>



### slotted summary


<rh-disclosure>
  <span slot="summary" class="icon">
    This is a slotted summary with extra markup <rh-icon set="ui" icon="like"></rh-icon>
  </span>
  <p>Instead of using <code>&lt;rh-disclosure summary="Hello world"&gt;</code>, users can slot content into a <code>summary</code> slot and include additional HTML if needed.</p>
  <p>Also note that slotted <code>summary</code> content will render on the page if/when JavaScript fails to load.</p>
</rh-disclosure>



### variants


<div class="container">
  <rh-disclosure summary="Compact variant disclosure" variant="compact">
    <p>Lorem ipsum dolor sit amet consectetur adipisicing, elit. Velit distinctio, nesciunt nobis sit, a dolor, non numquam rerum recusandae, deserunt enim assumenda quidem. Id impedit necessitatibus obcaecati ratione reprehenderit laborum?</p>
  </rh-disclosure>
  <rh-disclosure summary="Borderless variant disclosure" variant="borderless">
    <p>Lorem ipsum dolor sit amet consectetur adipisicing, elit. Velit distinctio, nesciunt nobis sit, a dolor, non numquam rerum recusandae, deserunt enim assumenda quidem. Id impedit necessitatibus obcaecati ratione reprehenderit laborum?</p>
  </rh-disclosure>
  <rh-disclosure summary="Borderless compact disclosure" variant="borderless compact">
    <p>Lorem ipsum dolor sit amet consectetur adipisicing, elit. Velit distinctio, nesciunt nobis sit, a dolor, non numquam rerum recusandae, deserunt enim assumenda quidem. Id impedit necessitatibus obcaecati ratione reprehenderit laborum?</p>
  </rh-disclosure>
</div>


