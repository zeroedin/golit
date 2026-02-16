---
title: "Alert"
imports:
  - rh-alert
  - rh-button
---

<p>8 demos for <code>&lt;rh-alert&gt;</code></p>


### alternate

<section id="alert-variant-alternate">
  <rh-alert variant="alternate">
    <h3 slot="header">Neutral</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="info" variant="alternate">
    <h3 slot="header">Info</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="success" variant="alternate">
    <h3 slot="header">Success</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="warning" variant="alternate">
    <h3 slot="header">Warning</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="danger" variant="alternate">
    <h3 slot="header">Danger</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="caution" variant="alternate">
    <h3 slot="header">Caution</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
</section>


### color context

<rh-context-demo>
  <section id="alert-context">
    <rh-alert>
      <h3 slot="header">Default</h3>
      <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
        egestas, a sollicitudin mauris tincidunt.</p>
      <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
      <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
    </rh-alert>
    <rh-alert variant="alternate">
      <h3 slot="header">Default</h3>
      <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
        egestas, a sollicitudin mauris tincidunt.</p>
      <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
      <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
    </rh-alert>
  </section>
</rh-context-demo>


### deprecated states

<section id="alert-states">
  <rh-alert state="error">
    <h3 slot="header">Error - alias of Danger</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="link" data-action="dismiss">Dismiss</rh-button>
    <rh-button slot="actions" variant="link" data-action="confirm">Confirm</rh-button>
  </rh-alert>
  <rh-alert state="default">
    <h3 slot="header">Default - alias of Neutral</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="link" data-action="dismiss">Dismiss</rh-button>
    <rh-button slot="actions" variant="link" data-action="confirm">Confirm</rh-button>
  </rh-alert>
  <rh-alert state="note">
    <h3 slot="header">Note - alias of Info</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="link" data-action="dismiss">Dismiss</rh-button>
    <rh-button slot="actions" variant="link" data-action="confirm">Confirm</rh-button>
  </rh-alert>
</section>


### dismissable

<section id="alert-dismissable">
  <rh-alert dismissable>
    <h3 slot="header">Default dismissable</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert variant="alternate" dismissable>
    <h3 slot="header">Inline dismissable</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert toast dismissable>
    <h3 slot="header">Toast dismissable</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert dismissable data-on-close="prevent-default">
    <h3 slot="header">Dismissable With Prevent Default</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
</section>


### index

<rh-alert>
  <h3 slot="header">Default</h3>
  <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
    egestas, a sollicitudin mauris tincidunt.</p>
  <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
  <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
</rh-alert>


### inline

<section id="alert-variant-inline">
  <rh-alert variant="inline">
    <h3 slot="header">Default</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="info" variant="inline">
    <h3 slot="header">Info</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="success" variant="inline">
    <h3 slot="header">Success</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="warning" variant="inline">
    <h3 slot="header">Warning</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="danger" variant="inline">
    <h3 slot="header">Danger</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
</section>


### states

<section id="alert-states">
  <rh-alert state="danger">
    <h3 slot="header">Danger</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="warning">
    <h3 slot="header">Warning</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="link" data-action="dismiss">Dismiss</rh-button>
    <rh-button slot="actions" variant="link" data-action="confirm">Confirm</rh-button>
  </rh-alert>
  <rh-alert state="caution">
    <h3 slot="header">Caution</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="link" data-action="dismiss">Dismiss</rh-button>
    <rh-button slot="actions" variant="link" data-action="confirm">Confirm</rh-button>
  </rh-alert>
  <rh-alert state="neutral">
    <h3 slot="header">Neutral</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="info">
    <h3 slot="header">Info</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
  <rh-alert state="success">
    <h3 slot="header">Success</h3>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eleifend elit sed est
      egestas, a sollicitudin mauris tincidunt.</p>
    <rh-button slot="actions" variant="secondary" data-action="confirm">Confirm</rh-button>
    <rh-button slot="actions" variant="link" data-action="dismiss">Cancel</rh-button>
  </rh-alert>
</section>


### toast

<form id="alert-variant-toast">
  <fieldset>
    <legend>Alert State</legend>
    <label><input type="radio" name="state" value="neutral">Neutral</label>
    <label><input type="radio" name="state" value="info">Info</label>
    <label><input type="radio" name="state" value="success">Success</label>
    <label><input type="radio" name="state" value="caution">Caution</label>
    <label><input type="radio" name="state" value="warning">Warning</label>
    <label><input type="radio" name="state" value="danger">Danger</label>
  </fieldset>
  <fieldset>
    <legend>Persistence</legend>
    <label><input type="checkbox" name="persistent">Persistent</label>
  </fieldset>
  <fieldset>
    <legend>Actions</legend>
    <label><input type="radio" name="actions" value="none">No actions</label>
    <label><input type="radio" name="actions" value="primary">Single action</label>
    <label><input type="radio" name="actions" value="secondary">Secondary action</label>
  </fieldset>
  <rh-button>Toast alert</rh-button>
</form>

