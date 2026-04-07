---
title: "Menu Dropdown"
imports:
  - rh-icon
  - rh-menu-dropdown
---

<p>12 demos for <code>&lt;rh-menu-dropdown&gt;</code></p>


### basic toggle with anchor links


<div id="menu-dropdown-container">
  <rh-menu-dropdown>
    <span slot="toggle-label">Basic toggle with links</span>
    <rh-menu-item href="#">Link one</rh-menu-item>
    <rh-menu-item href="#">Link two</rh-menu-item>
    <rh-menu-item href="#">Link three</rh-menu-item>
    <rh-menu-item disabled href="#">Disabled link</rh-menu-item>
    <hr/>
    <rh-menu-item href="#">Separated link</rh-menu-item>
    <rh-menu-item href="#" external>Separated, external link</rh-menu-item>
  </rh-menu-dropdown>
</div>



### basic toggle with fit text


<div id="menu-dropdown-container">
  <rh-menu-dropdown>
    <span slot="toggle-label">Basic toggle</span>
    <rh-menu-item>Action one</rh-menu-item>
    <rh-menu-item>Action two</rh-menu-item>
    <rh-menu-item>Action three</rh-menu-item>
    <rh-menu-item disabled>Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action</rh-menu-item>
  </rh-menu-dropdown>
</div>



### basic toggle with icon on left


<div id="menu-dropdown-container">
  <rh-menu-dropdown>
    <rh-icon slot="toggle-label" set="ui" icon="auto-light-dark-mode"></rh-icon>
    <span slot="toggle-label">Basic toggle</span>
    <rh-menu-item>Action one</rh-menu-item>
    <rh-menu-item>Action two</rh-menu-item>
    <rh-menu-item>Action three</rh-menu-item>
    <rh-menu-item disabled>Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action</rh-menu-item>
  </rh-menu-dropdown>
</div>



### color context


<rh-context-demo>
  <rh-menu-dropdown>
    <span slot="toggle-label">Basic toggle</span>
    <rh-menu-item>Action one</rh-menu-item>
    <rh-menu-item>Action two</rh-menu-item>
    <rh-menu-item>Action three</rh-menu-item>
    <rh-menu-item disabled>Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action</rh-menu-item>
  </rh-menu-dropdown>
</rh-context-demo>



### compact borderless variant


<div id="menu-dropdown-container">
  <rh-menu-dropdown variant="borderless" layout="compact" accessible-label="Toggle menu">
    <rh-menu-item>Action one</rh-menu-item>
    <rh-menu-item>Action two</rh-menu-item>
    <rh-menu-item>Action three</rh-menu-item>
    <rh-menu-item disabled>Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action</rh-menu-item>
  </rh-menu-dropdown>
</div>



### compact


<div id="menu-dropdown-container">
  <rh-menu-dropdown layout="compact" accessible-label="Toggle menu">
    <rh-menu-item>Action one</rh-menu-item>
    <rh-menu-item>Action two</rh-menu-item>
    <rh-menu-item>Action three</rh-menu-item>
    <rh-menu-item disabled>Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action</rh-menu-item>
  </rh-menu-dropdown>
</div>



### disabled


<div id="menu-dropdown-container">
  <rh-menu-dropdown disabled>
    <span slot="toggle-label">Basic toggle</span>
    <rh-menu-item>Action one</rh-menu-item>
    <rh-menu-item>Action two</rh-menu-item>
    <rh-menu-item>Action three</rh-menu-item>
    <rh-menu-item disabled>Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action</rh-menu-item>
  </rh-menu-dropdown>
</div>



### drodown items with icons


<div id="menu-dropdown-container">
  <rh-menu-dropdown>
    <span slot="toggle-label">Basic toggle</span>
    <rh-menu-item><rh-icon slot="icon" set="ui" icon="profile"></rh-icon> Action one</rh-menu-item>
    <rh-menu-item><rh-icon slot="icon" set="ui" icon="profile"></rh-icon> Action two</rh-menu-item>
    <rh-menu-item><rh-icon slot="icon" set="ui" icon="profile"></rh-icon> Action three</rh-menu-item>
    <rh-menu-item disabled><rh-icon slot="icon" set="ui" icon="profile"></rh-icon> Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item><rh-icon slot="icon" set="ui" icon="document"></rh-icon> Separated action</rh-menu-item>
  </rh-menu-dropdown>
</div>



### dropdown items with item descriptions


<div id="menu-dropdown-container">
  <rh-menu-dropdown>
    <span slot="toggle-label">Basic toggle</span>
    <rh-menu-item>Action one <p slot="description">Description for the 1st action</p></rh-menu-item>
    <rh-menu-item>Action two <p slot="description">Description for the 2nd action</p></rh-menu-item>
    <rh-menu-item>Action three <p slot="description">Description for the 3rd action</p></rh-menu-item>
    <rh-menu-item disabled>Disabled Action <p slot="description">Description for a disabled action</p></rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action <p slot="description">Description for a separate action</p></rh-menu-item>
  </rh-menu-dropdown>
</div>



### dropdown with group headings


<div id="menu-dropdown-container">
  <rh-menu-dropdown>
    <span slot="toggle-label">Basic toggle</span>
    <rh-menu-item-group group-heading="Group heading">
      <rh-menu-item>Action one</rh-menu-item>
      <rh-menu-item>Action two</rh-menu-item>
      <rh-menu-item>Action three</rh-menu-item>
      <rh-menu-item disabled>Disabled Action</rh-menu-item>
    </rh-menu-item-group>
    <hr/>
    <rh-menu-item-group group-heading="Group heading">
      <rh-menu-item>Separated action</rh-menu-item>
      <rh-menu-item>Separated action two</rh-menu-item>
    </rh-menu-item-group>
  </rh-menu-dropdown>
</div>



### index


<div id="menu-dropdown-container">
  <rh-menu-dropdown>
    <span slot="toggle-label">Basic toggle</span>
    <rh-menu-item>Action one</rh-menu-item>
    <rh-menu-item>Action two</rh-menu-item>
    <rh-menu-item>Action three</rh-menu-item>
    <rh-menu-item disabled>Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action</rh-menu-item>
  </rh-menu-dropdown>
</div>



### info action icon only


<div id="menu-dropdown-container">
  <rh-menu-dropdown>
    <rh-icon accessible-label="Toggle menu" slot="toggle-label" set="ui" icon="auto-light-dark-mode"></rh-icon>
    <rh-menu-item>Action one</rh-menu-item>
    <rh-menu-item>Action two</rh-menu-item>
    <rh-menu-item>Action three</rh-menu-item>
    <rh-menu-item disabled>Disabled Action</rh-menu-item>
    <hr/>
    <rh-menu-item>Separated action</rh-menu-item>
  </rh-menu-dropdown>
</div>


