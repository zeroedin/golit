---
title: "Button"
imports:
  - rh-button
---

<p>5 demos for <code>&lt;rh-button&gt;</code></p>


### color context


<rh-context-demo>
  <rh-button danger>Danger</rh-button>
  <rh-button>Primary</rh-button>
  <rh-button variant="link">Link</rh-button>
  <rh-button variant="secondary">Secondary</rh-button>
  <rh-button variant="secondary" danger>Secondary Danger</rh-button>
  <rh-button variant="tertiary">Tertiary</rh-button>
  <rh-button variant="close">Close</rh-button>
  <rh-button variant="play">Play</rh-button>
  <rh-button disabled>Disabled</rh-button>
  <rh-button danger icon="information-fill">Danger</rh-button>
  <rh-button icon="information-fill">Primary</rh-button>
  <rh-button variant="link" icon="information-fill">Link</rh-button>
  <rh-button variant="secondary" icon="information-fill">Secondary</rh-button>
  <rh-button variant="secondary" danger icon="information-fill">Secondary Danger</rh-button>
  <rh-button variant="tertiary" icon="information-fill">Tertiary</rh-button>
  <rh-button variant="close" icon="information-fill">Close</rh-button>
  <rh-button variant="play" icon="information-fill">Play</rh-button>
  <rh-button disabled icon="information-fill">Disabled</rh-button>
</rh-context-demo>



### form control


<form id="form">
  <fieldset id="fieldset">
    <legend>
      rh-button in a <code>&lt;fieldset&gt;</code> element;
      clicking this button must submit the form
    </legend>
    <rh-button id="button" type="submit">Submit</rh-button>
  </fieldset>
  <fieldset id="checkboxes">
    <legend>Use these checkboxes to toggle disabled state</legend>
    <input id="fst" type="checkbox" data-controls="fieldset" type="checkbox">
    <label for="fst">Disable fieldset</label>
    <input id="btn" type="checkbox" data-controls="button">
    <label for="btn">Disable rh-button</label>
  </fieldset>
  <fieldset id="outputs">
    <legend>Observe and reset form state</legend>
    <rh-button type="reset">Reset</rh-button>
    <label for="output">Form status:</label>
    <output id="output" name="output">Pending</output>
  </fieldset>
</form>



### icon


<section id="button-with-icon">
  <rh-button icon="error-fill" danger>Danger</rh-button>
  <rh-button icon="information-fill">Primary</rh-button>
  <rh-button icon="external-link" icon-set="microns" variant="link">Link</rh-button>
  <rh-button icon="bug-fill" variant="secondary">Secondary</rh-button>
  <rh-button icon="bug-fill" variant="secondary" danger>Secondary Danger</rh-button>
  <rh-button icon="close" icon-set="microns" variant="tertiary">Tertiary</rh-button>
  <rh-button icon="close" icon-set="microns" disabled>Disabled</rh-button>
</section>



### index


<rh-button>Primary</rh-button>



### variants


<section id="button-variants">
  <rh-button danger>Danger</rh-button>
  <rh-button>Primary</rh-button>
  <rh-button variant="link">Link</rh-button>
  <rh-button variant="secondary">Secondary</rh-button>
  <rh-button variant="secondary" danger>Secondary Danger</rh-button>
  <rh-button variant="tertiary">Tertiary</rh-button>
  <rh-button variant="close">Close</rh-button>
  <rh-button variant="play">Play</rh-button>
  <rh-button disabled>Disabled</rh-button>
</section>


