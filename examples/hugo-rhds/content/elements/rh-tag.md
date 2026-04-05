---
title: "Tag"
imports:
  - rh-tag
---

<p>12 demos for <code>&lt;rh-tag&gt;</code></p>


### color context

{{< raw >}}
<rh-context-demo>
  <dl>
    <dt>Filled</dt>
    <dd>
      <rh-tag color="red">Red</rh-tag>
      <rh-tag color="red-orange">Red Orange</rh-tag>
      <rh-tag color="orange">Orange</rh-tag>
      <rh-tag color="yellow">Yellow</rh-tag>
      <rh-tag color="green">Green</rh-tag>
      <rh-tag color="teal">Teal</rh-tag>
      <rh-tag color="blue">Blue</rh-tag>
      <rh-tag color="purple">Purple</rh-tag>
      <rh-tag color="gray">Gray</rh-tag>
    </dd>
    <dt>Filled with icon</dt>
    <dd>
      <rh-tag color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Filled with slotted icon</dt>
    <dd>
      <rh-tag color="red"> Red                        <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag color="red-orange"> Red Orange          <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag icon="information">Gray                 <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Filled with link</dt>
    <dd>
      <rh-tag href="#" color="red">Red</rh-tag>
      <rh-tag href="#" color="red-orange">Red Orange</rh-tag>
      <rh-tag href="#" color="orange">Orange</rh-tag>
      <rh-tag href="#" color="yellow">Yellow</rh-tag>
      <rh-tag href="#" color="green">Green</rh-tag>
      <rh-tag href="#" color="teal">Teal</rh-tag>
      <rh-tag href="#" color="blue">Blue</rh-tag>
      <rh-tag href="#" color="purple">Purple</rh-tag>
      <rh-tag href="#" color="gray">Gray</rh-tag>
    </dd>
    <dt>Filled with icon and link</dt>
    <dd>
      <rh-tag href="#" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag href="#" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag href="#" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag href="#" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag href="#" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag href="#" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag href="#" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag href="#" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag href="#" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Filled with slotted icon and link</dt>
    <dd>
      <rh-tag href="#" color="red"> Red                        <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" color="red-orange"> Red Orange          <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" icon="information">Gray                 <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Outline</dt>
    <dd>
      <rh-tag variant="outline" color="red">Red</rh-tag>
      <rh-tag variant="outline" color="red-orange">Red Orange</rh-tag>
      <rh-tag variant="outline" color="orange">Orange</rh-tag>
      <rh-tag variant="outline" color="yellow">Yellow</rh-tag>
      <rh-tag variant="outline" color="green">Green</rh-tag>
      <rh-tag variant="outline" color="teal">Teal</rh-tag>
      <rh-tag variant="outline" color="blue">Blue</rh-tag>
      <rh-tag variant="outline" color="purple">Purple</rh-tag>
      <rh-tag variant="outline" color="gray">Gray</rh-tag>
    </dd>
    <dt>Outlined with icon</dt>
    <dd>
      <rh-tag variant="outline" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag variant="outline" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag variant="outline" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag variant="outline" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag variant="outline" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag variant="outline" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag variant="outline" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag variant="outline" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag variant="outline" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Outlined with slotted icon</dt>
    <dd>
      <rh-tag variant="outline" color="red">Red                         <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="outline" color="red-orange">Red Orange           <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="outline" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="outline" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="outline" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="outline" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="outline" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="outline" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="outline" color="gray icon="information">Gray     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Outline with link</dt>
    <dd>
      <rh-tag href="#" variant="outline" color="red">Red</rh-tag>
      <rh-tag href="#" variant="outline" color="red-orange">Red Orange</rh-tag>
      <rh-tag href="#" variant="outline" color="orange">Orange</rh-tag>
      <rh-tag href="#" variant="outline" color="yellow">Yellow</rh-tag>
      <rh-tag href="#" variant="outline" color="green">Green</rh-tag>
      <rh-tag href="#" variant="outline" color="teal">Teal</rh-tag>
      <rh-tag href="#" variant="outline" color="blue">Blue</rh-tag>
      <rh-tag href="#" variant="outline" color="purple">Purple</rh-tag>
      <rh-tag href="#" variant="outline" color="gray">Gray</rh-tag>
    </dd>
    <dt>Outlined with icon and link</dt>
    <dd>
      <rh-tag href="#" variant="outline" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag href="#" variant="outline" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag href="#" variant="outline" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag href="#" variant="outline" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag href="#" variant="outline" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag href="#" variant="outline" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag href="#" variant="outline" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag href="#" variant="outline" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag href="#" variant="outline" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Outlined with slotted icon and link</dt>
    <dd>
      <rh-tag href="#" variant="outline" color="red">Red                         <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="outline" color="red-orange">Red Orange           <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="outline" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="outline" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="outline" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="outline" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="outline" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="outline" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="outline" color="gray icon="information">Gray     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Desaturated</dt>
    <dd>
      <rh-tag variant="desaturated" color="red">Red</rh-tag>
      <rh-tag variant="desaturated" color="red-orange">Red Orange</rh-tag>
      <rh-tag variant="desaturated" color="orange">Orange</rh-tag>
      <rh-tag variant="desaturated" color="yellow">Yellow</rh-tag>
      <rh-tag variant="desaturated" color="green">Green</rh-tag>
      <rh-tag variant="desaturated" color="teal">Teal</rh-tag>
      <rh-tag variant="desaturated" color="blue">Blue</rh-tag>
      <rh-tag variant="desaturated" color="purple">Purple</rh-tag>
      <rh-tag variant="desaturated" color="gray">Gray</rh-tag>
    </dd>
    <dt>Desaturated with icon</dt>
    <dd>
      <rh-tag variant="desaturated" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag variant="desaturated" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag variant="desaturated" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag variant="desaturated" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag variant="desaturated" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag variant="desaturated" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag variant="desaturated" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag variant="desaturated" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag variant="desaturated" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Desaturated with slotted icon</dt>
    <dd>
      <rh-tag variant="desaturated" color="red">Red                         <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="desaturated" color="red-orange">Red Orange           <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="desaturated" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="desaturated" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="desaturated" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="desaturated" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="desaturated" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="desaturated" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag variant="desaturated" color="gray icon="information">Gray     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Desaturated with link</dt>
    <dd>
      <rh-tag href="#" variant="desaturated" color="red">Red</rh-tag>
      <rh-tag href="#" variant="desaturated" color="red-orange">Red Orange</rh-tag>
      <rh-tag href="#" variant="desaturated" color="orange">Orange</rh-tag>
      <rh-tag href="#" variant="desaturated" color="yellow">Yellow</rh-tag>
      <rh-tag href="#" variant="desaturated" color="green">Green</rh-tag>
      <rh-tag href="#" variant="desaturated" color="teal">Teal</rh-tag>
      <rh-tag href="#" variant="desaturated" color="blue">Blue</rh-tag>
      <rh-tag href="#" variant="desaturated" color="purple">Purple</rh-tag>
      <rh-tag href="#" variant="desaturated" color="gray">Gray</rh-tag>
    </dd>
    <dt>Desaturated with icon and link</dt>
    <dd>
      <rh-tag href="#" variant="desaturated" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag href="#" variant="desaturated" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag href="#" variant="desaturated" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag href="#" variant="desaturated" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag href="#" variant="desaturated" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag href="#" variant="desaturated" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag href="#" variant="desaturated" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag href="#" variant="desaturated" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag href="#" variant="desaturated" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Desaturated with slotted icon and link</dt>
    <dd>
      <rh-tag href="#" variant="desaturated" color="red">Red                         <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="desaturated" color="red-orange">Red Orange           <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="desaturated" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="desaturated" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="desaturated" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="desaturated" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="desaturated" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="desaturated" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag href="#" variant="desaturated" color="gray icon="information">Gray     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Compact filled</dt>
    <dd>
      <rh-tag size="compact" color="red">Red</rh-tag>
      <rh-tag size="compact" color="red-orange">Red Orange</rh-tag>
      <rh-tag size="compact" color="orange">Orange</rh-tag>
      <rh-tag size="compact" color="yellow">Yellow</rh-tag>
      <rh-tag size="compact" color="green">Green</rh-tag>
      <rh-tag size="compact" color="teal">Teal</rh-tag>
      <rh-tag size="compact" color="blue">Blue</rh-tag>
      <rh-tag size="compact" color="purple">Purple</rh-tag>
      <rh-tag size="compact" color="gray">Gray</rh-tag>
    </dd>
    <dt>Compact filled with icon</dt>
    <dd>
      <rh-tag size="compact" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag size="compact" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag size="compact" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag size="compact" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag size="compact" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag size="compact" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag size="compact" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag size="compact" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag size="compact" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Compact filled with slotted icon</dt>
    <dd>
      <rh-tag size="compact" color="red"> Red                        <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" color="red-orange"> Red Orange          <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" icon="information">Gray                 <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Compact filled with link</dt>
    <dd>
      <rh-tag size="compact" href="#" color="red">Red</rh-tag>
      <rh-tag size="compact" href="#" color="red-orange">Red Orange</rh-tag>
      <rh-tag size="compact" href="#" color="orange">Orange</rh-tag>
      <rh-tag size="compact" href="#" color="yellow">Yellow</rh-tag>
      <rh-tag size="compact" href="#" color="green">Green</rh-tag>
      <rh-tag size="compact" href="#" color="teal">Teal</rh-tag>
      <rh-tag size="compact" href="#" color="blue">Blue</rh-tag>
      <rh-tag size="compact" href="#" color="purple">Purple</rh-tag>
      <rh-tag size="compact" href="#" color="gray">Gray</rh-tag>
    </dd>
    <dt>Compact filled with icon and link</dt>
    <dd>
      <rh-tag size="compact" href="#" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag size="compact" href="#" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag size="compact" href="#" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag size="compact" href="#" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag size="compact" href="#" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag size="compact" href="#" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag size="compact" href="#" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag size="compact" href="#" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag size="compact" href="#" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Compact filled with slotted icon and link</dt>
    <dd>
      <rh-tag size="compact" href="#" color="red"> Red                        <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" color="red-orange"> Red Orange          <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" icon="information">Gray                 <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Compact outline</dt>
    <dd>
      <rh-tag size="compact" variant="outline" color="red">Red</rh-tag>
      <rh-tag size="compact" variant="outline" color="red-orange">Red Orange</rh-tag>
      <rh-tag size="compact" variant="outline" color="orange">Orange</rh-tag>
      <rh-tag size="compact" variant="outline" color="yellow">Yellow</rh-tag>
      <rh-tag size="compact" variant="outline" color="green">Green</rh-tag>
      <rh-tag size="compact" variant="outline" color="teal">Teal</rh-tag>
      <rh-tag size="compact" variant="outline" color="blue">Blue</rh-tag>
      <rh-tag size="compact" variant="outline" color="purple">Purple</rh-tag>
      <rh-tag size="compact" variant="outline" color="gray">Gray</rh-tag>
    </dd>
    <dt>Compact outlined with icon</dt>
    <dd>
      <rh-tag size="compact" variant="outline" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag size="compact" variant="outline" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag size="compact" variant="outline" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag size="compact" variant="outline" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag size="compact" variant="outline" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag size="compact" variant="outline" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag size="compact" variant="outline" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag size="compact" variant="outline" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag size="compact" variant="outline" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Compact outlined with slotted icon</dt>
    <dd>
      <rh-tag size="compact" variant="outline" color="red">Red                         <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="outline" color="red-orange">Red Orange           <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="outline" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="outline" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="outline" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="outline" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="outline" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="outline" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="outline" color="gray icon="information">Gray     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Compact outline with link</dt>
    <dd>
      <rh-tag size="compact" href="#" variant="outline" color="red">Red</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="red-orange">Red Orange</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="orange">Orange</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="yellow">Yellow</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="green">Green</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="teal">Teal</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="blue">Blue</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="purple">Purple</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="gray">Gray</rh-tag>
    </dd>
    <dt>Compact outlined with icon and link</dt>
    <dd>
      <rh-tag size="compact" href="#" variant="outline" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Compact outlined with slotted icon and link</dt>
    <dd>
      <rh-tag size="compact" href="#" variant="outline" color="red">Red                         <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="red-orange">Red Orange           <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="outline" color="gray icon="information">Gray     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Compact desaturated</dt>
    <dd>
      <rh-tag size="compact" variant="desaturated" color="red">Red</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="red-orange">Red Orange</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="orange">Orange</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="yellow">Yellow</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="green">Green</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="teal">Teal</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="blue">Blue</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="purple">Purple</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="gray">Gray</rh-tag>
    </dd>
    <dt>Compact desaturated with icon</dt>
    <dd>
      <rh-tag size="compact" variant="desaturated" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag size="compact" variant="desaturated" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Compact desaturated with slotted icon</dt>
    <dd>
      <rh-tag size="compact" variant="desaturated" color="red">Red                         <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="desaturated" color="red-orange">Red Orange           <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="desaturated" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="desaturated" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="desaturated" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="desaturated" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="desaturated" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="desaturated" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" variant="desaturated" color="gray icon="information">Gray     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
    <dt>Compact desaturated with link</dt>
    <dd>
      <rh-tag size="compact" href="#" variant="desaturated" color="red">Red</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="red-orange">Red Orange</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="orange">Orange</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="yellow">Yellow</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="green">Green</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="teal">Teal</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="blue">Blue</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="purple">Purple</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="gray">Gray</rh-tag>
    </dd>
    <dt>Compact desaturated with icon and link</dt>
    <dd>
      <rh-tag size="compact" href="#" variant="desaturated" color="red"        icon="information-fill">Red</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="red-orange" icon="information-fill">Red Orange</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="orange"     icon="information-fill">Orange</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="yellow"     icon="information-fill">Yellow</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="green"      icon="information-fill">Green</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="teal"       icon="information-fill">Teal</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="blue"       icon="information-fill">Blue</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="purple"     icon="information-fill">Purple</rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="gray"       icon="information-fill">Gray</rh-tag>
    </dd>
    <dt>Compact desaturated with slotted icon and link</dt>
    <dd>
      <rh-tag size="compact" href="#" variant="desaturated" color="red">Red                         <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="red-orange">Red Orange           <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="orange">Orange                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="yellow">Yellow                   <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="green">Green                     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="teal" icon="information">Teal    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="blue" icon="information">Blue    <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="purple" icon="information">Purple<svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
      <rh-tag size="compact" href="#" variant="desaturated" color="gray icon="information">Gray     <svg slot="icon"><use href="#svg-icon"></use></svg></rh-tag>
    </dd>
  </dl>
</rh-context-demo>
<svg inert aria-hidden="true">
  <defs>
    <svg id="svg-icon" fill="currentColor" viewBox="0 0 512 512" aria-hidden="true" role="img">
      <path d="M256 8C119.043 8 8 119.083 8 256c0 136.997 111.043 248 248 248s248-111.003 248-248C504 119.083 392.957 8 256 8zm0 110c23.196 0 42 18.804 42 42s-18.804 42-42 42-42-18.804-42-42 18.804-42 42-42zm56 254c0 6.627-5.373 12-12 12h-88c-6.627 0-12-5.373-12-12v-24c0-6.627 5.373-12 12-12h12v-64h-12c-6.627 0-12-5.373-12-12v-24c0-6.627 5.373-12 12-12h64c6.627 0 12 5.373 12 12v100h12c6.627 0 12 5.373 12 12v24z"></path>
    </svg>
  </defs>
</svg>
{{< /raw >}}


### colors

{{< raw >}}
<rh-tag color="red">Red</rh-tag>
<rh-tag color="orange">Orange</rh-tag>
<rh-tag color="green">Green</rh-tag>
<rh-tag color="cyan">Cyan</rh-tag>
<rh-tag color="blue">Blue</rh-tag>
<rh-tag color="purple">Purple</rh-tag>
<rh-tag>Gray</rh-tag>
{{< /raw >}}


### desaturated with icon

{{< raw >}}
<rh-tag variant="desaturated" icon="information-fill" color="red"    >Red</rh-tag>
<rh-tag variant="desaturated" icon="information-fill" color="orange" >Orange</rh-tag>
<rh-tag variant="desaturated" icon="information-fill" color="green"  >Green</rh-tag>
<rh-tag variant="desaturated" icon="information-fill" color="cyan"   >Cyan</rh-tag>
<rh-tag variant="desaturated" icon="information-fill" color="blue"   >Blue</rh-tag>
<rh-tag variant="desaturated" icon="information-fill" color="purple" >Purple</rh-tag>
<rh-tag variant="desaturated" icon="information-fill" color="gray"   >Gray</rh-tag>
{{< /raw >}}


### desaturated

{{< raw >}}
<rh-tag variant="desaturated" color="red">Red</rh-tag>
<rh-tag variant="desaturated" color="orange">Orange</rh-tag>
<rh-tag variant="desaturated" color="green">Green</rh-tag>
<rh-tag variant="desaturated" color="cyan">Cyan</rh-tag>
<rh-tag variant="desaturated" color="blue">Blue</rh-tag>
<rh-tag variant="desaturated" color="purple">Purple</rh-tag>
<rh-tag variant="desaturated" color="gray">Gray</rh-tag>
{{< /raw >}}


### filled color with icon

{{< raw >}}
<rh-tag icon="information-fill" color="red">Red</rh-tag>
<rh-tag icon="information-fill" color="red-orange">Red Orange</rh-tag>
<rh-tag icon="information-fill" color="orange">Orange</rh-tag>
<rh-tag icon="information-fill" color="yellow">Yellow</rh-tag>
<rh-tag icon="information-fill" color="green">Green</rh-tag>
<rh-tag icon="information-fill" color="teal">Teal</rh-tag>
<rh-tag icon="information-fill" color="blue">Blue</rh-tag>
<rh-tag icon="information-fill" color="purple">Purple</rh-tag>
<rh-tag icon="information-fill" color="gray">Gray</rh-tag>
{{< /raw >}}


### filled color with slotted icon

{{< raw >}}
<rh-tag color="red">Red<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag color="red-orange">Red Orange<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag color="orange">Orange<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag color="yellow">Yellow<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag color="green">Green<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag color="teal">Teal<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag color="blue">Blue<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag color="purple">Purple<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag color="gray">Gray<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<svg aria-hidden="true" hidden inert>
  <defs>
    <svg id="slotted-icon"
         fill="currentColor"
         xmlns="http://www.w3.org/2000/svg"
         viewBox="0 0 24 24">
      <path d="M17,8a9,9,0,1,0,9,9A9,9,0,0,0,17,8Zm0,3.992a1.524,1.524,0,1,1-1.524,1.524A1.524,1.524,0,0,1,17,11.992Zm2.032,9.218a.436.436,0,0,1-.435.435H15.4a.436.436,0,0,1-.435-.435v-.871A.436.436,0,0,1,15.4,19.9h.435V17.581H15.4a.436.436,0,0,1-.435-.435v-.871a.436.436,0,0,1,.435-.435h2.323a.436.436,0,0,1,.435.435V19.9H18.6a.436.436,0,0,1,.435.435Z" transform="translate(-5 -5)"/>
    </svg>
  </defs>
</svg>
{{< /raw >}}


### filled color

{{< raw >}}
<rh-tag color="red">Red</rh-tag>
<rh-tag color="red-orange">Red Orange</rh-tag>
<rh-tag color="orange">Orange</rh-tag>
<rh-tag color="yellow">Yellow</rh-tag>
<rh-tag color="green">Green</rh-tag>
<rh-tag color="teal">Teal</rh-tag>
<rh-tag color="blue">Blue</rh-tag>
<rh-tag color="purple">Purple</rh-tag>
<rh-tag color="gray">Gray</rh-tag>
{{< /raw >}}


### index

{{< raw >}}
<rh-tag>Tag</rh-tag>
{{< /raw >}}


### links

{{< raw >}}
<dl>
  <dt>Filled with link</dt>
  <dd>
    <rh-tag href="#" color="red">Red</rh-tag>
    <rh-tag href="#" color="red-orange">Red Orange</rh-tag>
    <rh-tag href="#" color="orange">Orange</rh-tag>
    <rh-tag href="#" color="yellow">Yellow</rh-tag>
    <rh-tag href="#" color="green">Green</rh-tag>
    <rh-tag href="#" color="teal">Teal</rh-tag>
    <rh-tag href="#" color="blue">Blue</rh-tag>
    <rh-tag href="#" color="purple">Purple</rh-tag>
    <rh-tag href="#" color="gray">Gray</rh-tag>
  </dd>
  <dt>Filled with icon and link</dt>
  <dd>
    <rh-tag href="#" color="red"        icon="information-fill">Red</rh-tag>
    <rh-tag href="#" color="red-orange" icon="information-fill">Red Orange</rh-tag>
    <rh-tag href="#" color="orange"     icon="information-fill">Orange</rh-tag>
    <rh-tag href="#" color="yellow"     icon="information-fill">Yellow</rh-tag>
    <rh-tag href="#" color="green"      icon="information-fill">Green</rh-tag>
    <rh-tag href="#" color="teal"       icon="information-fill">Teal</rh-tag>
    <rh-tag href="#" color="blue"       icon="information-fill">Blue</rh-tag>
    <rh-tag href="#" color="purple"     icon="information-fill">Purple</rh-tag>
    <rh-tag href="#" color="gray"       icon="information-fill">Gray</rh-tag>
  </dd>
  <dt>Filled with link + disabled</dt>
  <dd>
    <rh-tag href="#" disabled color="red">Red</rh-tag>
    <rh-tag href="#" disabled color="red-orange">Red Orange</rh-tag>
    <rh-tag href="#" disabled color="orange">Orange</rh-tag>
    <rh-tag href="#" disabled color="yellow">Yellow</rh-tag>
    <rh-tag href="#" disabled color="green">Green</rh-tag>
    <rh-tag href="#" disabled color="teal">Teal</rh-tag>
    <rh-tag href="#" disabled color="blue">Blue</rh-tag>
    <rh-tag href="#" disabled color="purple">Purple</rh-tag>
    <rh-tag href="#" disabled color="gray">Gray</rh-tag>
  </dd>
  <dt>Filled with icon and link + disabled</dt>
  <dd>
    <rh-tag href="#" disabled color="red"        icon="information-fill">Red</rh-tag>
    <rh-tag href="#" disabled color="red-orange" icon="information-fill">Red Orange</rh-tag>
    <rh-tag href="#" disabled color="orange"     icon="information-fill">Orange</rh-tag>
    <rh-tag href="#" disabled color="yellow"     icon="information-fill">Yellow</rh-tag>
    <rh-tag href="#" disabled color="green"      icon="information-fill">Green</rh-tag>
    <rh-tag href="#" disabled color="teal"       icon="information-fill">Teal</rh-tag>
    <rh-tag href="#" disabled color="blue"       icon="information-fill">Blue</rh-tag>
    <rh-tag href="#" disabled color="purple"     icon="information-fill">Purple</rh-tag>
    <rh-tag href="#" disabled color="gray"       icon="information-fill">Gray</rh-tag>
  </dd>
</dl>
{{< /raw >}}


### outline with icon

{{< raw >}}
<rh-tag variant="outline" icon="information-fill" color="red"    >Red</rh-tag>
<rh-tag variant="outline" icon="information-fill" color="orange" >Orange</rh-tag>
<rh-tag variant="outline" icon="information-fill" color="green"  >Green</rh-tag>
<rh-tag variant="outline" icon="information-fill" color="cyan"   >Cyan</rh-tag>
<rh-tag variant="outline" icon="information-fill" color="blue"   >Blue</rh-tag>
<rh-tag variant="outline" icon="information-fill" color="purple" >Purple</rh-tag>
<rh-tag variant="outline" icon="information-fill" color="gray"   >Gray</rh-tag>
{{< /raw >}}


### outline with slotted icon

{{< raw >}}
<rh-tag variant="outline" color="red">Red<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag variant="outline" color="red-orange">Red Orange<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag variant="outline" color="orange">Orange<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag variant="outline" color="yellow">Yellow<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag variant="outline" color="green">Green<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag variant="outline" color="teal">Teal<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag variant="outline" color="blue">Blue<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag variant="outline" color="purple">Purple<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<rh-tag variant="outline" color="gray">Gray<svg slot="icon"><use href="#slotted-icon"/></svg></rh-tag>
<svg aria-hidden="true" hidden inert>
  <defs>
    <svg id="slotted-icon"
         xmlns="http://www.w3.org/2000/svg"
         viewBox="0 0 24 24">
      <path d="M17,8a9,9,0,1,0,9,9A9,9,0,0,0,17,8Zm0,3.992a1.524,1.524,0,1,1-1.524,1.524A1.524,1.524,0,0,1,17,11.992Zm2.032,9.218a.436.436,0,0,1-.435.435H15.4a.436.436,0,0,1-.435-.435v-.871A.436.436,0,0,1,15.4,19.9h.435V17.581H15.4a.436.436,0,0,1-.435-.435v-.871a.436.436,0,0,1,.435-.435h2.323a.436.436,0,0,1,.435.435V19.9H18.6a.436.436,0,0,1,.435.435Z" transform="translate(-5 -5)"/>
    </svg>
  </defs>
</svg>
{{< /raw >}}


### outline

{{< raw >}}
<rh-tag variant="outline" color="red">Red</rh-tag>
<rh-tag variant="outline" color="orange">Orange</rh-tag>
<rh-tag variant="outline" color="green">Green</rh-tag>
<rh-tag variant="outline" color="cyan">Cyan</rh-tag>
<rh-tag variant="outline" color="blue">Blue</rh-tag>
<rh-tag variant="outline" color="purple">Purple</rh-tag>
<rh-tag variant="outline" color="gray">Gray</rh-tag>
{{< /raw >}}

