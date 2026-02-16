import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('my-card')
export class MyCard extends LitElement {
  @property({ type: String })
  title: string = 'Untitled';

  @property({ type: String })
  subtitle?: string;

  @property({ type: Number })
  count: number = 0;

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

  render() {
    return html`
      <h2>${this.title}</h2>
      <p class="subtitle">${this.subtitle ?? ''}</p>
      <p class="count">Count: ${this.count}</p>
      <slot></slot>
    `;
  }
}
