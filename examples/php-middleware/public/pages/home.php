<span class="badge">PHP + golit</span>
<h1>SSR Middleware Demo</h1>
<p>
  This page is rendered by PHP and enhanced with server-side rendered Lit web
  components via <strong>golit</strong>. The counter below arrives pre-rendered
  with Declarative Shadow DOM&mdash;no flash of unstyled content.
</p>

<my-counter count="5"></my-counter>

<p>
  The <code>&lt;my-counter&gt;</code> component above is SSR'd by the PHP front
  controller, which pipes every HTML response through <code>golit transform</code>
  before sending it to the browser.
</p>
