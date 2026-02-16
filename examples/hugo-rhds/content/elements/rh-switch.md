---
title: "Switch"
imports:
  - rh-button
  - rh-switch
---

<p>10 demos for <code>&lt;rh-switch&gt;</code></p>


### check icon

<section>
  <form>
    <fieldset>
      <legend>Checked with label</legend>
      <rh-switch id="switch-a"
                 accessible-label="Switch A"
                 message-on="Message when on"
                 message-off="Message when off"
                 checked
                 show-check-icon></rh-switch>
    </fieldset>
  </form>
</section>


### color context

<rh-context-demo>
  <rh-switch accessible-label="Checked"
             message-on="Checked on"
             message-off="Checked off"
             checked></rh-switch>
  <rh-switch accessible-label="Unchecked"
             message-on="Unchecked on"
             message-off="Unchecked off"></rh-switch>
  <rh-switch accessible-label="Checked with icon"
             message-on="Checked with icon on"
             message-off="Checked with icon off"
             checked
             show-check-icon></rh-switch>
  <rh-switch accessible-label="Unchecked with icon"
             message-on="Unchecked with icon on"
             message-off="Unchecked with icon off"
             show-check-icon></rh-switch>
  <rh-switch accessible-label="Checked disabled"
             message-on="Disabled checked on"
             message-off="Disabled checked off"
             checked
             disabled></rh-switch>
  <rh-switch accessible-label="Unchecked disabled"
             message-on="Disabled unchecked on"
             message-off="Disabled unchecked off"
             disabled></rh-switch>
  <rh-switch accessible-label="Checked disabled with icon"
             message-on="Disabled checked with icon on"
             message-off="Disabled checked with icon off"
             show-check-icon
             checked
             disabled></rh-switch>
  <rh-switch accessible-label="Unchecked disabled with icon"
             message-on="Disabled unchecked with icon on"
             message-off="Disabled unchecked with icon off"
             show-check-icon
             disabled></rh-switch>
</rh-context-demo>


### disabled

<section>
  <form>
    <fieldset>
      <legend>Checked and Disabled</legend>
      <rh-switch id="switch-a"
                 disabled
                 message-on="Message when on"
                 message-off="Message when off"
                 accessible-label="Switch A"
                 checked></rh-switch>
    </fieldset>
    <fieldset>
      <label for="switch-b">Switch B</label>
      <rh-switch disabled
                 id="switch-b"
                 message-on="Message when on"
                 message-off="Message when off"
                 checked></rh-switch>
    </fieldset>
  </form>
</section>


### fieldset

<section>
  <p>A switch toggles the state of a setting (between on and off). Switches provide a more explicit, visible representation on a setting.</p>
  <form>
    <fieldset id="fieldset-a">
      <legend>Option A</legend>
      <rh-switch id="switch-a"
                 accessible-label="Switch A"
                 message-on="Message when on"
                 message-off="Message when off"
                 checked></rh-switch>
    </fieldset>
    <fieldset id="fieldset-b">
      <legend>Option B (Explicit label)</legend>
      <label for="switch-b">Switch B</label>
      <rh-switch id="switch-b"
                 message-on="Message when on"
                 message-off="Message when off"
                 checked></rh-switch>
    </fieldset>
    <fieldset id="form-disabled">
      <legend>Form Disabled State</legend>
      <label for="disable-fieldset-a">Disable Fieldset A</label>
      <rh-switch id="disable-fieldset-a"
                 message-on="Fieldset A is disabled"
                 message-off="ieldset A is enabled"
                 aria-controls="fieldset-a"></rh-switch>
      <label for="disable-switch-a">Disable Switch A</label>
      <rh-switch id="disable-switch-a"
                 message-on="Switch A is disabled"
                 message-off="Switch A is enabled"
                 aria-controls="switch-a"></rh-switch>
      <label for="disable-fieldset-b">Disable Fieldset B</label>
      <rh-switch id="disable-fieldset-b"
                 message-on="Fieldset A is disabled"
                 message-off="Fieldset A is enabled"
                 aria-controls="fieldset-b"></rh-switch>
      <label for="disable-switch-b">Disable Switch B</label>
      <rh-switch id="disable-switch-b"
                 message-on="Switch A is disabled"
                 message-off="Switch A is enabled"
                 aria-controls="switch-b"></rh-switch>
    </fieldset>
  </form>
</section>


### index

<rh-switch id="switch-a"
           accessible-label="Switch A"
           message-on="Message when on"
           message-off="Message when off"
           checked></rh-switch>


### nested in label

<section>
  <form id="nested-label">
    <fieldset>
      <legend>Nested in a label</legend>
      <label> Dark Mode
        <rh-switch name="status"
                   message-on="On"
                   message-off="Off"
                   checked></rh-switch>
      </label>
      <rh-button type="submit">Submit</rh-button>
      <output>Submit to read status</output>
    </fieldset>
  </form>
</section>


### reversed

<form>
  <fieldset>
    <legend>Reversed message</legend>
    <rh-switch id="switch-a"
               accessible-label="Switch A"
               message-on="Message when on"
               message-off="Message when off"
               checked
               reversed></rh-switch>
  </fieldset>
  <fieldset>
    <legend>Reversed with label</legend>
    <label for="switch-b">Switch B</label>
    <rh-switch id="switch-b"
               message-on="Message when on"
               message-off="Message when off"
               checked
               reversed></rh-switch>
  </fieldset>
  <fieldset>
    <legend>Reversed with slotted message</legend>
    <rh-switch id="switch-a"
               accessible-label="Switch A"
               checked
               reversed>
      <span slot="message-on">Message when on</span>
      <span slot="message-off">Message when off</span>
    </rh-switch>
  </fieldset>
  <fieldset>
    <legend>Reversed with label and slotted message</legend>
    <label for="switch-b">Switch B</label>
    <rh-switch id="switch-b"
               checked
               reversed>
      <span slot="message-on">Message when on</span>
      <span slot="message-off">Message when off</span>
    </rh-switch>
  </fieldset>
</form>


### rich messages

<rh-switch id="switch-a"
           accessible-label="Switch A"
           checked>
  <span slot="message-on">Message when <strong>on</strong></span>
  <span slot="message-off">Message when <strong>off</strong></span>
</rh-switch>


### right to left

<form dir="rtl">
  <fieldset>
    <legend>Right To Left - No Label</legend>
    <div>
      <rh-switch accessible-label="RTL Switch No Messages"></rh-switch>
    </div>
    <div>
      <rh-switch accessible-label="RTL Switch No Messages, Checked" checked></rh-switch>
    </div>
    <div>
      <rh-switch id="switch-a"
                 accessible-label="RTL Switch with Messages, Checked"
                 message-on="Message when on"
                 message-off="Message when off"
                 checked></rh-switch>
    </div>
    <div>
      <rh-switch id="switch-b"
                 message-on="Message when on"
                 message-off="Message when off"
                 accessible-label="RTL Switch with Messages, Show Checked Icon"
                 show-check-icon></rh-switch>
    </div>
    <div>
      <rh-switch id="switch-c"
                 accessible-label="RTL with Reversed Messages, Show Checked Icon"
                 message-on="Reversed message when on"
                 message-off="Reversed message when off"
                 show-check-icon
                 reversed></rh-switch>
    </div>
  </fieldset>
  <fieldset disabled>
    <legend>Right to Left - Disabled</legend>
    <div>
      <rh-switch id="switch-w" accessible-label="Disabled RTL Switch No Label"></rh-switch>
    </div>
    <div>
      <rh-switch id="switch-x" accessible-label="Disabled RTL Switch No Label, Checked" checked></rh-switch>
    </div>
    <div>
      <rh-switch id="switch-y" accessible-label="Disabled RTL Switch No Label, Checked, Show Checked Icon" checked show-check-icon></rh-switch>
    </div>
    <div>
      <label for="switch-z">Disabled RTL Switch with Label and Messages</label>
      <rh-switch id="switch-z"
                 message-on="Disabled Message when on"
                 message-off="Disabled Message when off"></rh-switch>
    </div>
    <div>
      <label for="switch-1">Disabled RTL Switch with Label and with Messages, Show Checked Icon</label>
      <rh-switch id="switch-1"
                 message-on="Disabled Message when on"
                 message-off="Disabled Message when off"
                 checked
                 show-check-icon></rh-switch>
    </div>
  </fieldset>
  <fieldset>
    <legend>Right to Left - With Sibling Label</legend>
    <div>
      <label for="switch-e">RTL Switch</label>
      <rh-switch id="switch-e"
                 message-on="Message when on"
                 message-off="Message when off"
                 checked></rh-switch>
    </div>
    <div>
      <label for="switch-f">RTL Switch</label>
      <rh-switch id="switch-f"
                 message-on="Message when on"
                 message-off="Message when off"></rh-switch>
    </div>
    <div>
      <label for="switch-g">RTL Switch, Show Checked Icon</label>
      <rh-switch id="switch-g"
                 checked
                 show-check-icon
                 message-on="Message when on"
                 message-off="Message when off"></rh-switch>
    </div>
    <div>
      <rh-switch id="switch-h"
                 message-on="Reversed message when on"
                 message-off="Reversed message when off"
                 checked
                 show-check-icon
                 reversed></rh-switch>
      <label for="switch-h">Switch with Reversed Messages, Show Checked Icon</label>
    </div>
  </fieldset>
  <fieldset>
    <legend>Right to Left - With Rich Messages</legend>
    <div>
      <rh-switch accessible-label="RTL Switch" checked>
        <div slot="messages">Message when
          <span data-state="on">On</span>
          <span data-state="of" hidden>Off</span>
        </div>
      </rh-switch>
    </div>
  </fieldset>
</form>


### without messages

<section>
  <form>
    <fieldset>
      <legend>Without messages</legend>
      <rh-switch id="switch-a" accessible-label="Switch A" checked></rh-switch>
      <rh-switch id="switch-b" accessible-label="Switch B" checked show-check-icon></rh-switch>
    </fieldset>
  </form>
</section>

