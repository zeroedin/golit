<span class="badge">How it works</span>
<h1>About This Example</h1>
<p>
  This example demonstrates <strong>golit</strong> as SSR middleware in a PHP
  application running inside a Podman / Docker container.
</p>

<h2 style="font-size:1.25rem; margin: 1.5rem 0 0.75rem;">The Pipeline</h2>
<ol>
  <li>PHP renders HTML templates containing Lit web components</li>
  <li>The front controller captures the output buffer</li>
  <li>It writes the HTML to a temp file and runs <code>golit transform</code></li>
  <li>golit renders each custom element to Declarative Shadow DOM</li>
  <li>The browser receives pre-rendered HTML that hydrates instantly</li>
</ol>

<my-counter count="42"></my-counter>

<p>
  Try viewing the page source&mdash;you&rsquo;ll see a
  <code>&lt;template shadowrootmode="open"&gt;</code> inside each
  <code>&lt;my-counter&gt;</code>, proving it was server-side rendered.
</p>
