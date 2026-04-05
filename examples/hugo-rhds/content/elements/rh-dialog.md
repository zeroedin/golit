---
title: "Dialog"
imports:
  - rh-button
  - rh-cta
  - rh-dialog
lightdom:
  - rh-cta-lightdom-shim.css
---

<p>8 demos for <code>&lt;rh-dialog&gt;</code></p>


### color context

{{< raw >}}
<rh-context-demo>
  <rh-dialog trigger="context-trigger">
    <h2 slot="header">Color Context</h2>
    <p>Lorem ipsum dolor sit amet, <a href="#foo">consectetur adipisicing</a> elit, sed do eiusmod tempor
      incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
      laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
      esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
      qui officia deserunt mollit anim id est laborum.</p>
    <p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
      consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
    <p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
      parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
      consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
      erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
    <p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
      odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
      condimentum nibh, ut fermentum massa justo sit amet risus.</p>
    <rh-cta slot="footer" href="#bar">Call to Action</rh-cta>
  </rh-dialog>
  <rh-button id="context-trigger">Open Dialog</rh-button>
</rh-context-demo>
{{< /raw >}}


### events

{{< raw >}}
<form id="dialog-events">
  <rh-dialog id="dialog" trigger="trigger">
    <h2 slot="header">Modal dialog with a header</h2>
    <p>Lorem ipsum dolor sit amet, <a href="#foo">consectetur adipisicing</a> elit, sed do eiusmod tempor incididunt
      ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
      aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
      fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit
      anim id est laborum.</p>
    <rh-cta slot="footer" href="#bar">Call to Action</rh-cta>
  </rh-dialog>
  <rh-button id="trigger">Open</rh-button>
  <fieldset>
    <legend>Events Fired</legend>
    <output name="events">No events yet</output>
  </fieldset>
</form>
{{< /raw >}}


### form

{{< raw >}}
<rh-dialog id="dialog-element"  trigger="trigger">
    <form id="form-element"> 
        <p>
          <label>
            Favorite RHDS Token Value:
            <select id="selectElement">
              <option value="default">Choose…</option>
              <option>--rh-color-brand-red</option>
              <option>--rh-color-red-50</option>
              <option>--rh-color-status-note</option>
            </select>
          </label>
        </p>
        <div>
            <button id="cancellationButton" value="cancel" formmethod="dialog">Cancel</button>
            <button type="submit" id="confirmation-button" value="default">Submit</button>
          </div>
        </form>
        </rh-dialog>
<rh-button id="trigger">Open Dialog</rh-button>
<output></output>
{{< /raw >}}


### index

{{< raw >}}
<rh-button id="first-modal-trigger">Open</rh-button>
<rh-dialog trigger="first-modal-trigger">
  <h2 slot="header">Modal dialog with a header</h2>
  <p>Lorem ipsum dolor sit amet, <a href="#foo">consectetur adipisicing</a> elit, sed do eiusmod tempor incididunt
    ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
    aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
    fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit
    anim id est laborum.</p>
  <rh-cta slot="footer" href="#bar">Call to Action</rh-cta>
</rh-dialog>
{{< /raw >}}


### lots of content

{{< raw >}}
<p>Scroll down to find the dialog's trigger.</p>
<p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
  consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
<p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
  parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
  consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
  erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
<p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
  odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
  condimentum nibh, ut fermentum massa justo sit amet risus.</p>
<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
  laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
  qui officia deserunt mollit anim id est laborum.</p>
<p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
  consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
<p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
  parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
  consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
  erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
<p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
  odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
  condimentum nibh, ut fermentum massa justo sit amet risus.</p>
<p>Lorem ipsum dolor sit amet, <a href="#foo">consectetur adipisicing</a> elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
  laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
  qui officia deserunt mollit anim id est laborum.</p>
<p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
  consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
<p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
  parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
  consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
  erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
<p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
  odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
  condimentum nibh, ut fermentum massa justo sit amet risus.</p>
<p>Lorem ipsum dolor sit amet, <a href="#foo">enim ad minim veniam</a> elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
  laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
  qui officia deserunt mollit anim id est laborum.</p>
<p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
  consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
<p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
  parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
  consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
  erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
<p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
  odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
  condimentum nibh, ut fermentum massa justo sit amet risus.</p>
<div class="wrap">
  <rh-button id="lots-of-content-trigger">Open</rh-button>
  <rh-dialog trigger="lots-of-content-trigger">
    <h2 slot="header">Modal with a header with a truly excessive super duper long title containing entirely too many words and in addition to the remarkably verbose title it also contains and a lot of content</h2>
    <p>Lorem ipsum dolor sit amet, <a href="#foo">nulla pariatur</a> elit, sed do eiusmod tempor
      incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
      laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
      esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
      qui officia deserunt mollit anim id est laborum.</p>
    <p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
      consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
    <p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
      parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
      consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
      erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
    <p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
      odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
      condimentum nibh, ut fermentum massa justo sit amet risus.</p>
    <p>Lorem ipsum dolor sit amet, <a href="#foo">justo odio</a> elit, sed do eiusmod tempor
      incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
      laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
      esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
      qui officia deserunt mollit anim id est laborum.</p>
    <p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
      consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
    <p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
      parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
      consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
      erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
    <p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
      odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
      condimentum nibh, ut fermentum massa justo sit amet risus.</p>
    <p>Lorem ipsum dolor sit amet, <a href="#foo">nisl consectetur</a> elit, sed do eiusmod tempor
      incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
      laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
      esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
      qui officia deserunt mollit anim id est laborum.</p>
    <p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
      consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
    <p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
      parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
      consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
      erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
    <p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
      odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
      condimentum nibh, ut fermentum massa justo sit amet risus.</p>
    <p>Lorem ipsum dolor sit amet, <a href="#foo">ridiculus mus</a> elit, sed do eiusmod tempor
      incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
      laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
      esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
      qui officia deserunt mollit anim id est laborum.</p>
    <p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
      consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
    <p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
      parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
      consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
      erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
    <p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
      odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
      condimentum nibh, ut fermentum massa justo sit amet risus.</p>
    <rh-cta slot="footer" href="#bar">Call to Action</rh-cta>
  </rh-dialog>
</div>
<p>Lorem ipsum dolor sit amet, <a href="#foo">Duis aute</a> elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
  laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
  qui officia deserunt mollit anim id est laborum.</p>
<p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
  consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
<p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
  parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
  consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
  erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
<p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
  odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
  condimentum nibh, ut fermentum massa justo sit amet risus.</p>
<p>Lorem ipsum dolor sit amet, <a href="#foo">Curabitur blandit</a> elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
  laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
  qui officia deserunt mollit anim id est laborum.</p>
<p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
  consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
<p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
  parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
  consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
  erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
<p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
  odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
  condimentum nibh, ut fermentum massa justo sit amet risus.</p>
<p>Lorem ipsum dolor sit amet, <a href="#foo">natoque penatibus</a> elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
  laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
  qui officia deserunt mollit anim id est laborum.</p>
<p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
  consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
<p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
  parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
  consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
  erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
<p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
  odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
  condimentum nibh, ut fermentum massa justo sit amet risus.</p>
<p>Lorem ipsum dolor sit amet, <a href="#foo">consectetur adipisicing</a> elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
  laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
  qui officia deserunt mollit anim id est laborum.</p>
<p>Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Donec sed odio dui. Sed posuere
  consectetur est at lobortis. Donec sed odio dui. Donec ullamcorper nulla non metus auctor fringilla.</p>
<p>Cras justo odio, dapibus ac facilisis in, egestas eget quam. Cum sociis natoque penatibus et magnis dis
  parturient montes, nascetur ridiculus mus. Curabitur blandit tempus porttitor. Lorem ipsum dolor sit amet,
  consectetur adipiscing elit. Curabitur blandit tempus porttitor. Duis mollis, est non commodo luctus, nisi
  erat porttitor ligula, eget lacinia odio sem nec elit. Curabitur blandit tempus porttitor.</p>
<p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula, eget lacinia odio sem nec elit. Donec sed
  odio dui. Maecenas faucibus mollis interdum. Fusce dapibus, tellus ac cursus commodo, tortor mauris
  condimentum nibh, ut fermentum massa justo sit amet risus.</p>
{{< /raw >}}


### no header content

{{< raw >}}
<rh-button id="no-header-content-trigger">Open</rh-button>
<rh-dialog trigger="no-header-content-trigger">
  <h2>This dialog has no slotted header content</h2>
  <p>All this content exists in the default slot.</p>
  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore
    magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
    consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
    pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est
    laborum.</p>
  <rh-cta slot="footer" href="#bar">Call to Action</rh-cta>
</rh-dialog>
{{< /raw >}}


### no headings

{{< raw >}}
<rh-button id="no-headings-trigger">Open</rh-button>
<rh-dialog accessible-label="This dialog's accessible name" trigger="no-headings-trigger">
  <p>This dialog uses the <code>accessible-label</code> attribute to define its accessible name for assistive technology.
     If the dialog includes a heading, the heading automatically becomes its accessible name—as long as the
     <code>accessible-label</code> attribute does not exist. If neither a heading nor the <code>accessible-label</code>
     attribute are provided, the dialog's accessible name defaults to the text content of its trigger element.
     Learn more about these attributes in the <a href="https://ux.redhat.com/elements/dialog/code/">rh-dialog docs</a>.
  </p>
  <rh-cta slot="footer" href="#">Call to Action</rh-cta>
</rh-dialog>
{{< /raw >}}


### video modal

{{< raw >}}
<rh-button id="video-modal-trigger">Open Video Dialog</rh-button>
<rh-dialog id="video-modal"
           type="video"
           trigger="video-modal-trigger">
  <video controls src="https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4"></video>
</rh-dialog>
<rh-button id="youtube-modal-trigger">Open YouTube Dialog</rh-button>
<rh-dialog id="youtube-modal"
           type="video"
           trigger="youtube-modal-trigger">
  <iframe src="https://www.youtube.com/embed/aqz-KE-bpKQ?enablejsapi=1" title="YouTube video player" frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowfullscreen></iframe>
</rh-dialog>
{{< /raw >}}

