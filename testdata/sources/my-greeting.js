import { LitElement, html, css } from 'lit';

class MyGreeting extends LitElement {
  static properties = {
    name: { type: String, reflect: true },
    active: { type: Boolean },
  };

  static styles = css`
    :host { display: block; }
    p { color: blue; }
  `;

  constructor() {
    super();
    this.name = 'World';
    this.active = false;
  }

  render() {
    return html`<p class="${this.active ? 'highlight' : ''}">Hello, ${this.name}!</p>`;
  }
}

customElements.define('my-greeting', MyGreeting);
