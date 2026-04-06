import { LitElement, html, css } from 'lit';

class MyCounter extends LitElement {
  static properties = {
    count: { type: Number },
  };

  static styles = css`
    :host {
      display: block;
      font-family: system-ui, -apple-system, sans-serif;
      padding: 1.5rem;
      border: 2px solid #e2e8f0;
      border-radius: 12px;
      text-align: center;
      max-width: 280px;
      background: #f8fafc;
    }
    h3 {
      margin: 0 0 0.75rem;
      font-size: 0.875rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: #64748b;
    }
    .count {
      font-size: 3rem;
      font-weight: 700;
      color: #1e40af;
      margin: 0.5rem 0;
      line-height: 1;
    }
    .controls {
      display: flex;
      gap: 0.5rem;
      justify-content: center;
      margin-top: 1rem;
    }
    button {
      font-size: 1.25rem;
      width: 2.5rem;
      height: 2.5rem;
      border: 1px solid #cbd5e1;
      border-radius: 8px;
      cursor: pointer;
      background: white;
      color: #334155;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: all 0.15s ease;
    }
    button:hover {
      background: #e2e8f0;
      border-color: #94a3b8;
    }
    button:active {
      transform: scale(0.95);
    }
  `;

  constructor() {
    super();
    this.count = 0;
  }

  render() {
    return html`
      <h3>Counter</h3>
      <div class="count">${this.count}</div>
      <div class="controls">
        <button @click=${() => this.count--} aria-label="Decrement">&minus;</button>
        <button @click=${() => this.count++} aria-label="Increment">&plus;</button>
      </div>
    `;
  }
}

customElements.define('my-counter', MyCounter);
