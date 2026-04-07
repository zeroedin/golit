import { spawn } from 'node:child_process';

async function waitForHealth(url, { timeout = 30000, interval = 200 } = {}) {
  const start = Date.now();
  while (Date.now() - start < timeout) {
    try {
      const res = await fetch(url);
      if (res.ok) return;
    } catch {
      // server not ready yet
    }
    await new Promise(r => setTimeout(r, interval));
  }
  throw new Error(`golit serve did not become healthy at ${url} within ${timeout}ms`);
}

export default function golitSSRPlugin(eleventyConfig, opts = {}) {
  const {
    binary = '../../dist/golit',
    defs = 'bundles/',
    port = 9777,
    ignore = [],
    concurrency,
    preload = [],
  } = opts;

  let proc;
  const baseURL = `http://127.0.0.1:${port}`;

  eleventyConfig.on('eleventy.before', async () => {
    const args = ['serve', '--defs', defs, '--listen', `127.0.0.1:${port}`];
    for (const tag of ignore) {
      args.push('--ignore', tag);
    }
    for (const mod of preload) {
      args.push('--preload', mod);
    }
    if (concurrency) {
      args.push('-j', String(concurrency));
    }

    proc = spawn(binary, args, {
      stdio: ['ignore', 'inherit', 'inherit'],
    });

    proc.on('error', (err) => {
      console.error(`golit serve failed to start: ${err.message}`);
    });

    await waitForHealth(`${baseURL}/health`);
    console.log(`[golit-ssr] golit serve ready on ${baseURL}`);
  });

  eleventyConfig.on('eleventy.after', async () => {
    if (proc) {
      proc.kill('SIGTERM');
      await new Promise(resolve => proc.on('close', resolve));
      console.log('[golit-ssr] golit serve stopped');
    }
  });

  eleventyConfig.addTransform('render-golit', async function(content) {
    const { outputPath } = this.page;
    if (!outputPath?.endsWith('.html')) return content;

    try {
      const res = await fetch(`${baseURL}/render`, {
        method: 'POST',
        headers: { 'Content-Type': 'text/html; charset=utf-8' },
        body: content,
      });
      if (!res.ok) {
        const err = await res.text();
        console.error(`[golit-ssr] render failed for ${outputPath}: ${err}`);
        return content;
      }
      return await res.text();
    } catch (err) {
      console.error(`[golit-ssr] render error for ${outputPath}: ${err.message}`);
      return content;
    }
  });
}
