---
title: "Navigation (vertical)"
imports:
  - rh-icon
  - rh-navigation-link
  - rh-navigation-vertical
lightdom:
  - rh-navigation-vertical-lightdom.css
---

<p>7 demos for <code>&lt;rh-navigation-vertical&gt;</code></p>


### color context

{{< raw >}}
<rh-context-demo>
  <rh-navigation-vertical bordered="inline-start">
    <rh-navigation-link href="#manage-applications" current-page>Manage applications</rh-navigation-link>
    <rh-navigation-vertical-list summary="1. Managing applications">
      <rh-navigation-link>
        <a class="index-link" href="#managing-applications">Managing applications</a>
      </rh-navigation-link>
      <rh-navigation-vertical-list summary="1.1 Application management lifecycle">
        <rh-navigation-link>
          <a class="index-link" href="#application-management-lifecycle">Application Management Lifecycle</a>
        </rh-navigation-link>
        <rh-navigation-link href="#application-model-definitions">1.1.1 Application Model Definitions</rh-navigation-link>
        <rh-navigation-link href="#managing-applications-with-the-console">1.1.2 Managing applications with the console</rh-navigation-link>
        <rh-navigation-link href="#channels">1.1.3 Channels</rh-navigation-link>
        <rh-navigation-link href="#subscriptions">1.1.4 Subscriptions</rh-navigation-link>
        <rh-navigation-link href="#placement-rules">1.1.5 Placement rules</rh-navigation-link>
        <rh-navigation-vertical-list summary="1.1.6 Application">
          <rh-navigation-link>
            <a class="index-link" href="#application">Application</a>
          </rh-navigation-link>
          <rh-navigation-vertical-list summary="1.1.6.1 Creating and managing channels">
            <rh-navigation-link>
              <a class="index-link" href="#creating-and-managing-channels">Creating and managing channels</a>
            </rh-navigation-link>
            <rh-navigation-link href="#creating-channels">1.1.6.1.1 Creating and managing channels</rh-navigation-link>
            <rh-navigation-vertical-list summary="1.1.6.1.2 Updating channel">
              <rh-navigation-link>
                <a class="index-link" href="#updating-channel">Updating channel</a>
              </rh-navigation-link>
              <rh-navigation-link href="#deleting-channel">1.1.6.1.2.1 Deleting channel</rh-navigation-link>
            </rh-navigation-vertical-list>
            <rh-navigation-link href="#managing-deployments-with-channels">1.1.6.1.3 Managing deployments with channels</rh-navigation-link>
          </rh-navigation-vertical-list>
        </rh-navigation-vertical-list>
      </rh-navigation-vertical-list>
    </rh-navigation-vertical-list>
    <rh-navigation-link href="#legal-notice">Legal Notice</rh-navigation-link>
  </rh-navigation-vertical>
</rh-context-demo>
{{< /raw >}}


### contained

{{< raw >}}
<div id="container">
  <rh-navigation-vertical>
    <rh-navigation-link href="#" current-page>Home</rh-navigation-link>
    <rh-navigation-vertical-list summary="About">
      <rh-navigation-link href="#about">About the Design System</rh-navigation-link>
      <rh-navigation-link href="#about/roadmap">Roadmap</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-vertical-list summary="Get Started">
      <rh-navigation-link href="#get-started">About the Design System</rh-navigation-link>
      <rh-navigation-link href="#get-started/designers">Designers</rh-navigation-link>
      <rh-navigation-link href="#get-started/developers">Developers</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-vertical-list summary="Foundations">
      <rh-navigation-link href="#foundations">Overview</rh-navigation-link>
      <rh-navigation-link href="#foundations/color">Color</rh-navigation-link>
      <rh-navigation-link href="#foundations/grid">Grid</rh-navigation-link>
      <rh-navigation-link href="#foundations/iconography">Iconography</rh-navigation-link>
      <rh-navigation-link href="#foundations/interactions">Interactions</rh-navigation-link>
      <rh-navigation-link href="#foundations/spacing">Spacing</rh-navigation-link>
      <rh-navigation-link href="#foundations/typography">Typography</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-vertical-list summary="Tokens">
      <rh-navigation-link href="#tokens">Overview</rh-navigation-link>
      <rh-navigation-link href="#tokens/color">Global Color</rh-navigation-link>
      <rh-navigation-link href="#tokens/box-shadow">Box Shadow</rh-navigation-link>
      <rh-navigation-link href="#tokens/typography">Typography</rh-navigation-link>
      <rh-navigation-link href="#tokens/border">Border</rh-navigation-link>
      <rh-navigation-link href="#tokens/opacity">Opacity</rh-navigation-link>
      <rh-navigation-link href="#tokens/space">Space</rh-navigation-link>
      <rh-navigation-link href="#tokens/length">Length</rh-navigation-link>
      <rh-navigation-link href="#tokens/icon">Icon</rh-navigation-link>
      <rh-navigation-link href="#tokens/breakpoints">Breakpoints</rh-navigation-link>
      <rh-navigation-link href="#tokens/media-queries">Media Queries</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-vertical-list summary="Elements">
      <rh-navigation-link href="#all-elements">All Elements</rh-navigation-link>
      <rh-navigation-link href="#elements/accordion">Accordion</rh-navigation-link>
      <rh-navigation-link href="#elements/alert">Alert</rh-navigation-link>
      <rh-navigation-link href="#elements/announcement">Announcement</rh-navigation-link>
      <rh-navigation-link href="#elements/audio-player">Audio Player</rh-navigation-link>
      <rh-navigation-link href="#elements/avatar">Avatar</rh-navigation-link>
      <rh-navigation-link href="#elements/back-to-top">Back to Top</rh-navigation-link>
      <rh-navigation-link href="#elements/badge">Badge</rh-navigation-link>
      <rh-navigation-link href="#elements/breadcrumb">Breadcrumb</rh-navigation-link>
      <rh-navigation-link href="#elements/button">Button</rh-navigation-link>
      <rh-navigation-link href="#elements/card">Card</rh-navigation-link>
      <rh-navigation-link href="#elements/chip">Chip</rh-navigation-link>
      <rh-navigation-link href="#elements/code-block">Code Block</rh-navigation-link>
      <rh-navigation-link href="#elements/call-to-action">Call to Action</rh-navigation-link>
      <rh-navigation-link href="#elements/dialog">Dialog</rh-navigation-link>
      <rh-navigation-link href="#elements/disclosure">Disclosure</rh-navigation-link>
      <rh-navigation-link href="#elements/footer">Footer</rh-navigation-link>
      <rh-navigation-link href="#elements/health-index">Health Index</rh-navigation-link>
      <rh-navigation-link href="#elements/icon">Icon</rh-navigation-link>
      <rh-navigation-link href="#elements/jump-links">Jump Links</rh-navigation-link>
      <rh-navigation-link href="#elements/navigation-primary">Navigation Primary</rh-navigation-link>
      <rh-navigation-link href="#elements/navigation-secondary">Navigation Secondary</rh-navigation-link>
      <rh-navigation-link href="#elements/navigation-vertical">Navigation Vertical</rh-navigation-link>
      <rh-navigation-link href="#elements/pagination">Pagination</rh-navigation-link>
      <rh-navigation-link href="#elements/progress-indicator">Progress Indicator</rh-navigation-link>
      <rh-navigation-link href="#elements/popover">Popover</rh-navigation-link>
      <rh-navigation-link href="#elements/progress-steps">Progress Steps</rh-navigation-link>
      <rh-navigation-link href="#elements/site-status">Site Status</rh-navigation-link>
      <rh-navigation-link href="#elements/skip-link">Skip Link</rh-navigation-link>
      <rh-navigation-link href="#elements/spinner">Spinner</rh-navigation-link>
      <rh-navigation-link href="#elements/statistic">Statistic</rh-navigation-link>
      <rh-navigation-link href="#elements/subnavigation">Subnavigation</rh-navigation-link>
      <rh-navigation-link href="#elements/surface">Surface</rh-navigation-link>
      <rh-navigation-link href="#elements/switch">Switch</rh-navigation-link>
      <rh-navigation-link href="#elements/table">Table</rh-navigation-link>
      <rh-navigation-link href="#elements/tabs">Tabs</rh-navigation-link>
      <rh-navigation-link href="#elements/tag">Tag</rh-navigation-link>
      <rh-navigation-link href="#elements/tile">Tile</rh-navigation-link>
      <rh-navigation-link href="#elements/timestamp">Timestamp</rh-navigation-link>
      <rh-navigation-link href="#elements/tooltip">Tooltip</rh-navigation-link>
      <rh-navigation-link href="#elements/video-embed">Video Embed</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-vertical-list summary="Theming">
      <rh-navigation-link href="#theming/overview">Overview</rh-navigation-link>
      <rh-navigation-link href="#theming/color-palettes">Color Palettes</rh-navigation-link>
      <rh-navigation-link href="#theming/customizing">Customizing</rh-navigation-link>
      <rh-navigation-link href="#theming/developers">Developers</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-vertical-list summary="Patterns">
      <rh-navigation-link href="#patterns/all-patterns">All Patterns</rh-navigation-link>
      <rh-navigation-link href="#patterns/card">Card</rh-navigation-link>
      <rh-navigation-link href="#patterns/tabs">Tabs</rh-navigation-link>
      <rh-navigation-link href="#patterns/filter">Filter</rh-navigation-link>
      <rh-navigation-link href="#patterns/form">Form</rh-navigation-link>
      <rh-navigation-link href="#patterns/link-with-icon">Link with Icon</rh-navigation-link>
      <rh-navigation-link href="#patterns/search-bar">Search Bar</rh-navigation-link>
      <rh-navigation-link href="#patterns/sticky-banner">Sticky Banner</rh-navigation-link>
      <rh-navigation-link href="#patterns/sticky-card">Sticky Card</rh-navigation-link>
      <rh-navigation-link href="#patterns/tile">Tile</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-vertical-list summary="Personalization">
      <rh-navigation-link href="#personalization/all-personalization-patterns">All Personalization Patterns</rh-navigation-link>
      <rh-navigation-link href="#personalization/announcement">Announcement</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-vertical-list summary="Accessibility">
      <rh-navigation-link href="#accessibility/fundamentals">Fundamentals</rh-navigation-link>
      <rh-navigation-link href="#accessibility/accessibility-tools">Accessibility Tools</rh-navigation-link>
      <rh-navigation-link href="#accessibility/assistive-technologies">Assistive Technologies</rh-navigation-link>
      <rh-navigation-link href="#accessibility/ci-cd">CI/CD</rh-navigation-link>
      <rh-navigation-link href="#accessibility/content">Content</rh-navigation-link>
      <rh-navigation-link href="#accessibility/contributors">Contributors</rh-navigation-link>
      <rh-navigation-link href="#accessibility/design">Design</rh-navigation-link>
      <rh-navigation-link href="#accessibility/development">Development</rh-navigation-link>
      <rh-navigation-link href="#accessibility/manual-testing">Manual Testing</rh-navigation-link>
      <rh-navigation-link href="#accessibility/resources">Resources</rh-navigation-link>
      <rh-navigation-link href="#accessibility/screen-readers">Screen Readers</rh-navigation-link>
    </rh-navigation-vertical-list>
    <rh-navigation-link href="#design-code-status">Design & Code Status</rh-navigation-link>
    <rh-navigation-link href="#release-notes">Release Notes</rh-navigation-link>
    <rh-navigation-link href="#get-support">Get Support</rh-navigation-link>
  </rh-navigation-vertical>
</div>
{{< /raw >}}


### index

{{< raw >}}
<rh-navigation-vertical>
  <rh-navigation-link href="#" current-page>Home</rh-navigation-link>
  <rh-navigation-vertical-list summary="About">
    <rh-navigation-link href="#about">About the Design System</rh-navigation-link>
    <rh-navigation-link href="#about/roadmap">Roadmap</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Get Started">
    <rh-navigation-link href="#get-started">About the Design System</rh-navigation-link>
    <rh-navigation-link href="#get-started/designers">Designers</rh-navigation-link>
    <rh-navigation-link href="#get-started/developers">Developers</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Foundations">
    <rh-navigation-link href="#foundations">Overview</rh-navigation-link>
    <rh-navigation-link href="#foundations/color">Color</rh-navigation-link>
    <rh-navigation-link href="#foundations/grid">Grid</rh-navigation-link>
    <rh-navigation-link href="#foundations/iconography">Iconography</rh-navigation-link>
    <rh-navigation-link href="#foundations/interactions">Interactions</rh-navigation-link>
    <rh-navigation-link href="#foundations/spacing">Spacing</rh-navigation-link>
    <rh-navigation-link href="#foundations/typography">Typography</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Tokens">
    <rh-navigation-link href="#tokens">Overview</rh-navigation-link>
    <rh-navigation-link href="#tokens/color">Global Color</rh-navigation-link>
    <rh-navigation-link href="#tokens/box-shadow">Box Shadow</rh-navigation-link>
    <rh-navigation-link href="#tokens/typography">Typography</rh-navigation-link>
    <rh-navigation-link href="#tokens/border">Border</rh-navigation-link>
    <rh-navigation-link href="#tokens/opacity">Opacity</rh-navigation-link>
    <rh-navigation-link href="#tokens/space">Space</rh-navigation-link>
    <rh-navigation-link href="#tokens/length">Length</rh-navigation-link>
    <rh-navigation-link href="#tokens/icon">Icon</rh-navigation-link>
    <rh-navigation-link href="#tokens/breakpoints">Breakpoints</rh-navigation-link>
    <rh-navigation-link href="#tokens/media-queries">Media Queries</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Elements">
    <rh-navigation-link href="#all-elements">All Elements</rh-navigation-link>
    <rh-navigation-link href="#elements/accordion">Accordion</rh-navigation-link>
    <rh-navigation-link href="#elements/alert">Alert</rh-navigation-link>
    <rh-navigation-link href="#elements/announcement">Announcement</rh-navigation-link>
    <rh-navigation-link href="#elements/audio-player">Audio Player</rh-navigation-link>
    <rh-navigation-link href="#elements/avatar">Avatar</rh-navigation-link>
    <rh-navigation-link href="#elements/back-to-top">Back to Top</rh-navigation-link>
    <rh-navigation-link href="#elements/badge">Badge</rh-navigation-link>
    <rh-navigation-link href="#elements/breadcrumb">Breadcrumb</rh-navigation-link>
    <rh-navigation-link href="#elements/button">Button</rh-navigation-link>
    <rh-navigation-link href="#elements/card">Card</rh-navigation-link>
    <rh-navigation-link href="#elements/chip">Chip</rh-navigation-link>
    <rh-navigation-link href="#elements/code-block">Code Block</rh-navigation-link>
    <rh-navigation-link href="#elements/call-to-action">Call to Action</rh-navigation-link>
    <rh-navigation-link href="#elements/dialog">Dialog</rh-navigation-link>
    <rh-navigation-link href="#elements/disclosure">Disclosure</rh-navigation-link>
    <rh-navigation-link href="#elements/footer">Footer</rh-navigation-link>
    <rh-navigation-link href="#elements/health-index">Health Index</rh-navigation-link>
    <rh-navigation-link href="#elements/icon">Icon</rh-navigation-link>
    <rh-navigation-link href="#elements/jump-links">Jump Links</rh-navigation-link>
    <rh-navigation-link href="#elements/navigation-primary">Navigation Primary</rh-navigation-link>
    <rh-navigation-link href="#elements/navigation-secondary">Navigation Secondary</rh-navigation-link>
    <rh-navigation-link href="#elements/navigation-vertical">Navigation Vertical</rh-navigation-link>
    <rh-navigation-link href="#elements/pagination">Pagination</rh-navigation-link>
    <rh-navigation-link href="#elements/progress-indicator">Progress Indicator</rh-navigation-link>
    <rh-navigation-link href="#elements/popover">Popover</rh-navigation-link>
    <rh-navigation-link href="#elements/progress-steps">Progress Steps</rh-navigation-link>
    <rh-navigation-link href="#elements/site-status">Site Status</rh-navigation-link>
    <rh-navigation-link href="#elements/skip-link">Skip Link</rh-navigation-link>
    <rh-navigation-link href="#elements/spinner">Spinner</rh-navigation-link>
    <rh-navigation-link href="#elements/statistic">Statistic</rh-navigation-link>
    <rh-navigation-link href="#elements/subnavigation">Subnavigation</rh-navigation-link>
    <rh-navigation-link href="#elements/surface">Surface</rh-navigation-link>
    <rh-navigation-link href="#elements/switch">Switch</rh-navigation-link>
    <rh-navigation-link href="#elements/table">Table</rh-navigation-link>
    <rh-navigation-link href="#elements/tabs">Tabs</rh-navigation-link>
    <rh-navigation-link href="#elements/tag">Tag</rh-navigation-link>
    <rh-navigation-link href="#elements/tile">Tile</rh-navigation-link>
    <rh-navigation-link href="#elements/timestamp">Timestamp</rh-navigation-link>
    <rh-navigation-link href="#elements/tooltip">Tooltip</rh-navigation-link>
    <rh-navigation-link href="#elements/video-embed">Video Embed</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Theming">
    <rh-navigation-link href="#theming/overview">Overview</rh-navigation-link>
    <rh-navigation-link href="#theming/color-palettes">Color Palettes</rh-navigation-link>
    <rh-navigation-link href="#theming/customizing">Customizing</rh-navigation-link>
    <rh-navigation-link href="#theming/developers">Developers</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Patterns">
    <rh-navigation-link href="#patterns/all-patterns">All Patterns</rh-navigation-link>
    <rh-navigation-link href="#patterns/card">Card</rh-navigation-link>
    <rh-navigation-link href="#patterns/tabs">Tabs</rh-navigation-link>
    <rh-navigation-link href="#patterns/filter">Filter</rh-navigation-link>
    <rh-navigation-link href="#patterns/form">Form</rh-navigation-link>
    <rh-navigation-link href="#patterns/link-with-icon">Link with Icon</rh-navigation-link>
    <rh-navigation-link href="#patterns/search-bar">Search Bar</rh-navigation-link>
    <rh-navigation-link href="#patterns/sticky-banner">Sticky Banner</rh-navigation-link>
    <rh-navigation-link href="#patterns/sticky-card">Sticky Card</rh-navigation-link>
    <rh-navigation-link href="#patterns/tile">Tile</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Personalization">
    <rh-navigation-link href="#personalization/all-personalization-patterns">All Personalization Patterns</rh-navigation-link>
    <rh-navigation-link href="#personalization/announcement">Announcement</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Accessibility">
    <rh-navigation-link href="#accessibility/fundamentals">Fundamentals</rh-navigation-link>
    <rh-navigation-link href="#accessibility/accessibility-tools">Accessibility Tools</rh-navigation-link>
    <rh-navigation-link href="#accessibility/assistive-technologies">Assistive Technologies</rh-navigation-link>
    <rh-navigation-link href="#accessibility/ci-cd">CI/CD</rh-navigation-link>
    <rh-navigation-link href="#accessibility/content">Content</rh-navigation-link>
    <rh-navigation-link href="#accessibility/contributors">Contributors</rh-navigation-link>
    <rh-navigation-link href="#accessibility/design">Design</rh-navigation-link>
    <rh-navigation-link href="#accessibility/development">Development</rh-navigation-link>
    <rh-navigation-link href="#accessibility/manual-testing">Manual Testing</rh-navigation-link>
    <rh-navigation-link href="#accessibility/resources">Resources</rh-navigation-link>
    <rh-navigation-link href="#accessibility/screen-readers">Screen Readers</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-link href="#design-code-status">Design & Code Status</rh-navigation-link>
  <rh-navigation-link href="#release-notes">Release Notes</rh-navigation-link>
  <rh-navigation-link href="#get-support">Get Support</rh-navigation-link>
</rh-navigation-vertical>
{{< /raw >}}


### nested

{{< raw >}}
<rh-navigation-vertical bordered="inline-start">
  <rh-navigation-link href="#manage-applications">Manage applications</rh-navigation-link>
  <rh-navigation-vertical-list highlight summary="1. Managing applications">
    <rh-navigation-link current-page>
      <a class="index-link" aria-current="page" href="#managing-applications">Managing applications</a>
    </rh-navigation-link>
    <rh-navigation-vertical-list highlight summary="1.1 Application management lifecycle">
      <rh-navigation-link>
        <a class="index-link" href="#application-management-lifecycle">Application Management Lifecycle</a>
      </rh-navigation-link>
      <rh-navigation-link href="#application-model-definitions">1.1.1 Application Model Definitions</rh-navigation-link>
      <rh-navigation-link href="#managing-applications-with-the-console">1.1.2 Managing applications with the console</rh-navigation-link>
      <rh-navigation-link href="#channels">1.1.3 Channels</rh-navigation-link>
      <rh-navigation-link href="#subscriptions">1.1.4 Subscriptions</rh-navigation-link>
      <rh-navigation-link href="#placement-rules">1.1.5 Placement rules</rh-navigation-link>
      <rh-navigation-vertical-list highlight summary="1.1.6 Application">
        <rh-navigation-link>
          <a class="index-link" href="#application">Application</a>
        </rh-navigation-link>
        <rh-navigation-vertical-list highlight summary="1.1.6.1 Creating and managing channels">
          <rh-navigation-link>
            <a class="index-link" href="#creating-and-managing-channels">Creating and managing channels</a>
          </rh-navigation-link>
          <rh-navigation-link href="#creating-channels">1.1.6.1.1 Creating and managing channels</rh-navigation-link>
          <rh-navigation-vertical-list highlight summary="1.1.6.1.2 Updating channel">
            <rh-navigation-link>
              <a class="index-link" href="#updating-channel">Updating channel</a>
            </rh-navigation-link>
            <rh-navigation-link href="#deleting-channel">1.1.6.1.2.1 Deleting channel</rh-navigation-link>
          </rh-navigation-vertical-list>
          <rh-navigation-link href="#managing-deployments-with-channels">1.1.6.1.3 Managing deployments with channels</rh-navigation-link>
        </rh-navigation-vertical-list>
      </rh-navigation-vertical-list>
    </rh-navigation-vertical-list>
  </rh-navigation-vertical-list>
  <rh-navigation-link href="#legal-notice">Legal Notice</rh-navigation-link>
</rh-navigation-vertical>
{{< /raw >}}


### open

{{< raw >}}
<rh-navigation-vertical>
  <rh-navigation-link href="#">Home</rh-navigation-link>
  <rh-navigation-vertical-list summary="About">
    <rh-navigation-link href="#about">About the Design System</rh-navigation-link>
    <rh-navigation-link href="#about/roadmap">Roadmap</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Get Started">
    <rh-navigation-link href="#get-started">About the Design System</rh-navigation-link>
    <rh-navigation-link href="#get-started/designers">Designers</rh-navigation-link>
    <rh-navigation-link href="#get-started/developers">Developers</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Foundations">
    <rh-navigation-link href="#foundations">Overview</rh-navigation-link>
    <rh-navigation-link href="#foundations/color">Color</rh-navigation-link>
    <rh-navigation-link href="#foundations/grid">Grid</rh-navigation-link>
    <rh-navigation-link href="#foundations/iconography">Iconography</rh-navigation-link>
    <rh-navigation-link href="#foundations/interactions">Interactions</rh-navigation-link>
    <rh-navigation-link href="#foundations/spacing">Spacing</rh-navigation-link>
    <rh-navigation-link href="#foundations/typography">Typography</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Tokens">
    <rh-navigation-link href="#tokens">Overview</rh-navigation-link>
    <rh-navigation-link href="#tokens/color">Global Color</rh-navigation-link>
    <rh-navigation-link href="#tokens/box-shadow">Box Shadow</rh-navigation-link>
    <rh-navigation-link href="#tokens/typography">Typography</rh-navigation-link>
    <rh-navigation-link href="#tokens/border">Border</rh-navigation-link>
    <rh-navigation-link href="#tokens/opacity">Opacity</rh-navigation-link>
    <rh-navigation-link href="#tokens/space">Space</rh-navigation-link>
    <rh-navigation-link href="#tokens/length">Length</rh-navigation-link>
    <rh-navigation-link href="#tokens/icon">Icon</rh-navigation-link>
    <rh-navigation-link href="#tokens/breakpoints">Breakpoints</rh-navigation-link>
    <rh-navigation-link href="#tokens/media-queries">Media Queries</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Elements" open>
    <rh-navigation-link href="#all-elements" current-page>All Elements</rh-navigation-link>
    <rh-navigation-link href="#elements/accordion">Accordion</rh-navigation-link>
    <rh-navigation-link href="#elements/alert">Alert</rh-navigation-link>
    <rh-navigation-link href="#elements/announcement">Announcement</rh-navigation-link>
    <rh-navigation-link href="#elements/audio-player">Audio Player</rh-navigation-link>
    <rh-navigation-link href="#elements/avatar">Avatar</rh-navigation-link>
    <rh-navigation-link href="#elements/back-to-top">Back to Top</rh-navigation-link>
    <rh-navigation-link href="#elements/badge">Badge</rh-navigation-link>
    <rh-navigation-link href="#elements/breadcrumb">Breadcrumb</rh-navigation-link>
    <rh-navigation-link href="#elements/button">Button</rh-navigation-link>
    <rh-navigation-link href="#elements/card">Card</rh-navigation-link>
    <rh-navigation-link href="#elements/chip">Chip</rh-navigation-link>
    <rh-navigation-link href="#elements/code-block">Code Block</rh-navigation-link>
    <rh-navigation-link href="#elements/call-to-action">Call to Action</rh-navigation-link>
    <rh-navigation-link href="#elements/dialog">Dialog</rh-navigation-link>
    <rh-navigation-link href="#elements/disclosure">Disclosure</rh-navigation-link>
    <rh-navigation-link href="#elements/footer">Footer</rh-navigation-link>
    <rh-navigation-link href="#elements/health-index">Health Index</rh-navigation-link>
    <rh-navigation-link href="#elements/icon">Icon</rh-navigation-link>
    <rh-navigation-link href="#elements/jump-links">Jump Links</rh-navigation-link>
    <rh-navigation-link href="#elements/navigation-primary">Navigation Primary</rh-navigation-link>
    <rh-navigation-link href="#elements/navigation-secondary">Navigation Secondary</rh-navigation-link>
    <rh-navigation-link href="#elements/navigation-vertical">Navigation Vertical</rh-navigation-link>
    <rh-navigation-link href="#elements/pagination">Pagination</rh-navigation-link>
    <rh-navigation-link href="#elements/progress-indicator">Progress Indicator</rh-navigation-link>
    <rh-navigation-link href="#elements/popover">Popover</rh-navigation-link>
    <rh-navigation-link href="#elements/progress-steps">Progress Steps</rh-navigation-link>
    <rh-navigation-link href="#elements/site-status">Site Status</rh-navigation-link>
    <rh-navigation-link href="#elements/skip-link">Skip Link</rh-navigation-link>
    <rh-navigation-link href="#elements/spinner">Spinner</rh-navigation-link>
    <rh-navigation-link href="#elements/statistic">Statistic</rh-navigation-link>
    <rh-navigation-link href="#elements/subnavigation">Subnavigation</rh-navigation-link>
    <rh-navigation-link href="#elements/surface">Surface</rh-navigation-link>
    <rh-navigation-link href="#elements/switch">Switch</rh-navigation-link>
    <rh-navigation-link href="#elements/table">Table</rh-navigation-link>
    <rh-navigation-link href="#elements/tabs">Tabs</rh-navigation-link>
    <rh-navigation-link href="#elements/tag">Tag</rh-navigation-link>
    <rh-navigation-link href="#elements/tile">Tile</rh-navigation-link>
    <rh-navigation-link href="#elements/timestamp">Timestamp</rh-navigation-link>
    <rh-navigation-link href="#elements/tooltip">Tooltip</rh-navigation-link>
    <rh-navigation-link href="#elements/video-embed">Video Embed</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Theming">
    <rh-navigation-link href="#theming/overview">Overview</rh-navigation-link>
    <rh-navigation-link href="#theming/color-palettes">Color Palettes</rh-navigation-link>
    <rh-navigation-link href="#theming/customizing">Customizing</rh-navigation-link>
    <rh-navigation-link href="#theming/developers">Developers</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Patterns">
    <rh-navigation-link href="#patterns/all-patterns">All Patterns</rh-navigation-link>
    <rh-navigation-link href="#patterns/card">Card</rh-navigation-link>
    <rh-navigation-link href="#patterns/tabs">Tabs</rh-navigation-link>
    <rh-navigation-link href="#patterns/filter">Filter</rh-navigation-link>
    <rh-navigation-link href="#patterns/form">Form</rh-navigation-link>
    <rh-navigation-link href="#patterns/link-with-icon">Link with Icon</rh-navigation-link>
    <rh-navigation-link href="#patterns/search-bar">Search Bar</rh-navigation-link>
    <rh-navigation-link href="#patterns/sticky-banner">Sticky Banner</rh-navigation-link>
    <rh-navigation-link href="#patterns/sticky-card">Sticky Card</rh-navigation-link>
    <rh-navigation-link href="#patterns/tile">Tile</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Personalization">
    <rh-navigation-link href="#personalization/all-personalization-patterns">All Personalization Patterns</rh-navigation-link>
    <rh-navigation-link href="#personalization/announcement">Announcement</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list summary="Accessibility">
    <rh-navigation-link href="#accessibility/fundamentals">Fundamentals</rh-navigation-link>
    <rh-navigation-link href="#accessibility/accessibility-tools">Accessibility Tools</rh-navigation-link>
    <rh-navigation-link href="#accessibility/assistive-technologies">Assistive Technologies</rh-navigation-link>
    <rh-navigation-link href="#accessibility/ci-cd">CI/CD</rh-navigation-link>
    <rh-navigation-link href="#accessibility/content">Content</rh-navigation-link>
    <rh-navigation-link href="#accessibility/contributors">Contributors</rh-navigation-link>
    <rh-navigation-link href="#accessibility/design">Design</rh-navigation-link>
    <rh-navigation-link href="#accessibility/development">Development</rh-navigation-link>
    <rh-navigation-link href="#accessibility/manual-testing">Manual Testing</rh-navigation-link>
    <rh-navigation-link href="#accessibility/resources">Resources</rh-navigation-link>
    <rh-navigation-link href="#accessibility/screen-readers">Screen Readers</rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-link href="#design-code-status">Design & Code Status</rh-navigation-link>
  <rh-navigation-link href="#release-notes">Release Notes</rh-navigation-link>
  <rh-navigation-link href="#get-support">Get Support</rh-navigation-link>
</rh-navigation-vertical>
{{< /raw >}}


### slotted anchor links

{{< raw >}}
<rh-navigation-vertical>
  <rh-navigation-link>
    <a href="#" aria-current="page">Home</a>
  </rh-navigation-link>
  <rh-navigation-vertical-list>
    <span slot="summary">About</span>
    <rh-navigation-link>
      <a href="#about-the-design-system">About the Design System</a>
    </rh-navigation-link>
    <rh-navigation-link>
      <a href="#about/roadmap">Roadmap</a>
    </rh-navigation-link>
  </rh-navigation-vertical-list>
  <rh-navigation-vertical-list>
    <span slot="summary">Get Started</span>
    <rh-navigation-link>
      <a href="#get-started">About the Design System</a>
    </rh-navigation-link>
    <rh-navigation-link>
      <a href="#get-started/designers">Designers</a>
    </rh-navigation-link>
    <rh-navigation-link>
      <a href="#get-started/developers">Developers</a>
    </rh-navigation-link>
  </rh-navigation-vertical-list>
</rh-navigation-vertical>
{{< /raw >}}


### slotted icons

{{< raw >}}
<rh-navigation-vertical>
  <rh-navigation-link href="#" current-page>
    <span><rh-icon set="ui" icon="home-fill"></rh-icon> Icon Example 1</span>
  </rh-navigation-link>
    <rh-navigation-link href="#">
    <span><rh-icon set="ui" icon="language-fill"></rh-icon> Icon Example 2</span>
  </rh-navigation-link>
    <rh-navigation-link href="#">
    <span><rh-icon set="ui" icon="secured-fill"></rh-icon> Icon Example 3</span>
  </rh-navigation-link>
    <rh-navigation-link href="#">
    <span><rh-icon set="ui" icon="bug-fill"></rh-icon> Icon Example 4</span>
  </rh-navigation-link>
  <rh-navigation-vertical-list>
    <span slot="summary"><rh-icon set="ui" icon="check-circle-fill"></rh-icon> Icon Example 5</span>
    <rh-navigation-link href="#">
    <span><rh-icon set="ui" icon="cloud-fill"></rh-icon>Icon Example 6</span>
  </rh-navigation-link>
  <rh-navigation-link href="#">
    <span><rh-icon set="ui" icon="folder-fill"></rh-icon> Icon Example 7</span>
  </rh-navigation-link>
  </rh-navigation-vertical-list>
</rh-navigation-vertical>
{{< /raw >}}

