import { LitElement, html, css } from 'lit';

class MyCard extends LitElement {
  static properties = {
    title: { type: String },
    subtitle: { type: String },
    count: { type: Number },
  };

  static styles = css`
    :host {
      display: block;
      border: 1px solid #ccc;
      padding: 16px;
      border-radius: 8px;
    }
    h2 { margin: 0 0 8px; }
    .subtitle { color: #666; }
    .count { font-weight: bold; }
  `;

  constructor() {
    super();
    this.title = 'Untitled';
    this.subtitle = '';
    this.count = 0;
  }

  render() {
    return html`
      <h2>${this.title}</h2>
      <p class="subtitle">${this.subtitle}</p>
      <p class="count">Count: ${this.count}</p>
      <slot></slot>
    `;
  }
}

customElements.define('my-card', MyCard);
