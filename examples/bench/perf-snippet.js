// Performance API self-reporting snippet for golit benchmarks.
// Only activates when the URL contains ?bench — zero impact on normal usage.
// After window.onload + settle time, writes metrics into <pre id="golit-perf">
// so Chrome headless --dump-dom can extract them.
(function () {
  if (!location.search.includes('bench')) return;

  var metrics = {};

  try {
    new PerformanceObserver(function (list) {
      list.getEntries().forEach(function (entry) {
        if (entry.name === 'first-contentful-paint') {
          metrics.fcp = Math.round(entry.startTime * 100) / 100;
        }
      });
    }).observe({ type: 'paint', buffered: true });
  } catch (_) {}

  try {
    new PerformanceObserver(function (list) {
      var entries = list.getEntries();
      if (entries.length) {
        metrics.lcp = Math.round(entries[entries.length - 1].startTime * 100) / 100;
      }
    }).observe({ type: 'largest-contentful-paint', buffered: true });
  } catch (_) {}

  window.addEventListener('load', function () {
    setTimeout(function () {
      var nav = performance.getEntriesByType('navigation')[0];
      if (nav) {
        metrics.ttfb = Math.round(nav.responseStart * 100) / 100;
        metrics.domContentLoaded = Math.round(nav.domContentLoadedEventEnd * 100) / 100;
        metrics.loadEvent = Math.round(nav.loadEventEnd * 100) / 100;
        metrics.domInteractive = Math.round(nav.domInteractive * 100) / 100;
        metrics.responseEnd = Math.round(nav.responseEnd * 100) / 100;
      }

      var el = document.createElement('pre');
      el.id = 'golit-perf';
      el.style.display = 'none';
      el.textContent = JSON.stringify(metrics);
      document.body.appendChild(el);
    }, 500);
  });
})();
