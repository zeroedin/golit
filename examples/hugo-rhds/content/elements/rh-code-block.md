---
title: "Code block"
imports:
  - rh-badge
  - rh-code-block
  - rh-cta
lightdom:
  - rh-code-block-lightdom.css
  - rh-cta-lightdom-shim.css
---

17 demos for `<rh-code-block>`

### Default (index)

{{< raw >}}
```html
<!DOCTYPE html>
<title>Title</title>
<style>body {width: 500px;}</style>
<body>
  <p checked class="title" id="title">Title</p>
  <!-- here goes the rest of the page -->
</body>
```
{{< /raw >}}

### Actions

{{< raw >}}
```css
#content {
  display: block;
  background-color: var(--rh-color-surface-lighter, #f2f2f2);
  border: var(--rh-border-width-sm, 1px) solid var(--rh-color-border-subtle-on-light, #c7c7c7);
  border-block-start-width: var(--rh-code-block-border-block-start-width, var(--rh-border-width-sm, 1px));
  font-family: var(--rh-font-family-code, RedHatMono, "Red Hat Mono", "Courier New", Courier, monospace);
  color: var(--rh-color-text-primary-on-light, #151515);
  padding: var(--rh-space-xl, 24px);
  height: calc(100% - 2 * var(--rh-space-xl, 24px));
  border-radius: var(--rh-border-radius-default, 3px);
  max-width: 1000px;
  max-height: 640px;
  overflow-y: auto;
}
```
{{< /raw >}}

### Actions i18n

{{< raw >}}
{{< raw >}}
<rh-code-block actions="wrap copy">
  <span slot="action-label-copy">העתק</span>
  <span slot="action-label-copy" hidden data-code-block-state="active">הועתק!</span>
  <span slot="action-label-wrap">לעבור לגלישת שורות</span>
  <span slot="action-label-wrap" hidden data-code-block-state="active">לעבור להצפת שורות</span>
  <script type="text/css">#content {
  display: block;
  background-color: var(--rh-color-surface-lighter, #f2f2f2);
  border: var(--rh-border-width-sm, 1px) solid var(--rh-color-border-subtle-on-light, #c7c7c7);
  font-family: var(--rh-font-family-code, RedHatMono, "Red Hat Mono", "Courier New", Courier, monospace);
  padding: var(--rh-space-xl, 24px);
  max-width: 1000px;
  max-height: 640px;
  overflow-y: auto;
}</script>
</rh-code-block>
{{< /raw >}}
{{< /raw >}}

### Callout badges

{{< raw >}}
{{< raw >}}
<rh-code-block>
  <script type="text/html"><p>Script tags in HTML must be escaped</p></script>
  <rh-badge state="info">1</rh-badge>
</rh-code-block>
{{< /raw >}}
{{< /raw >}}

### Hide line numbers

{{< raw >}}
{{< raw >}}
<rh-code-block line-numbers="hidden">
  <script type="text/html"><!DOCTYPE html>
<title>Title</title>
<style>body {width: 500px;}</style>
<body>
  <p checked class="title" id="title">Title</p>
</body></script>
</rh-code-block>
{{< /raw >}}
{{< /raw >}}

### Client-side highlighting (HTML)

{{< raw >}}
```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <title>Cards Galore!</title>
  </head>
  <body>
    <main>
      <rh-card>
        <h2 slot="header">Card</h2>
        <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>
        <rh-cta slot="footer" priority="primary">
          <a href="#">Call to action</a>
        </rh-cta>
      </rh-card>
    </main>
  </body>
</html>
```
{{< /raw >}}

### Client-side highlighting (CSS)

{{< raw >}}
```css
rh-card.avatar-card {
  width: 360px;
  &::part(body) {
    margin-block-start: var(--rh-space-lg, 16px);
  }
  & p {
    margin-block-start: 0;
  }
  & h4 {
    font-weight: var(--rh-font-weight-heading-regular, 300);
    font-size: var(--rh-font-size-body-text-md, 1rem);
    font-family: var(--rh-font-family-body-text);
    line-height: var(--rh-line-height-body-text, 1.5);
  }
}
```
{{< /raw >}}

### Client-side highlighting (YAML)

{{< raw >}}
```yaml
extends:
  - stylelint-config-standard
  - '@stylistic/stylelint-config'
plugins:
  - ./node_modules/@rhds/tokens/plugins/stylelint.js
  - '@stylistic/stylelint-plugin'
rules:
  rhds/token-values: true
  rhds/no-unknown-token-name:
    - true
    - allowed:
      - --rh-icon-size
```
{{< /raw >}}

### Client-side highlighting (JavaScript)

{{< raw >}}
```javascript
import '@rhds/elements/rh-card/rh-card.js';

const card = document.querySelector('rh-card');
card.addEventListener('click', () => {
  console.log('Card clicked!');
});
```
{{< /raw >}}

### Resizable

{{< raw >}}
{{< raw >}}
<rh-code-block resizable>
  <script type="text/text">Error: Error creating network Load Balancer: AccessDenied: User: arn:aws:sts::970xxxxxxxxx:assumed-role/ManagedOpenShift-Installer-Role/163xxxxxxxxxxxxxxxx is not authorized to perform: iam:CreateServiceLinkedRole on resource: arn:aws:iam::970xxxxxxxxx:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing</script>
</rh-code-block>
{{< /raw >}}
{{< /raw >}}

### Sizes — Standard

{{< raw >}}
```text
oc apply -f ostoy-microservice-deployment.yaml
```
{{< /raw >}}

### Sizes — Multi-line

{{< raw >}}
```bash
$ podman login -u flozanorht quay.io
Password:
Login Succeeded!
$ skopeo copy docker://registry.access.redhat.com/ubi8/ubi:8.0-122 \
docker://quay.io/flozanorht/ubi:8
...
Writing manifest to image destination
Storing signatures
```
{{< /raw >}}

### Below the fold

{{< raw >}}
{{< raw >}}
<section style="height: 100vh; display: flex; align-items: center; justify-content: center;">
  <p><strong>Scroll down to view code block</strong></p>
</section>
{{< /raw >}}

```html
<!DOCTYPE html>
<title>Title</title>
<style>body {width: 500px;}</style>
<body>
  <p checked class="title" id="title">Title</p>
  <!-- here goes the rest of the page -->
</body>
```
{{< /raw >}}

### Bash example

{{< raw >}}
```bash
for node in $(oc get nodes -o jsonpath='{.items[*].metadata.name}'); do
  echo ${node}
  oc adm cordon ${node}
done
```
{{< /raw >}}

### JSON example

{{< raw >}}
```json
{
  "apiVersion": "v1",
  "kind": "ConfigMap",
  "metadata": {
    "name": "cluster-monitoring-config",
    "namespace": "openshift-monitoring"
  },
  "data": {
    "config.yaml": "enableUserWorkload: true"
  }
}
```
{{< /raw >}}

### TypeScript example

{{< raw >}}
```typescript
import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('my-element')
export class MyElement extends LitElement {
  @property() name = 'World';

  static styles = css`:host { display: block; }`;

  render() {
    return html`<p>Hello, ${this.name}!</p>`;
  }
}
```
{{< /raw >}}

### Ruby example

{{< raw >}}
```ruby
require 'sinatra'

get '/' do
  'Hello from Red Hat!'
end

post '/deploy' do
  content_type :json
  { status: 'deployed', timestamp: Time.now }.to_json
end
```
{{< /raw >}}
