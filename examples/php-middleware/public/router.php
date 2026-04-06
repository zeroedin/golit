<?php
/**
 * Front controller with golit SSR middleware.
 *
 * Usage: php -S 0.0.0.0:8080 -t public public/router.php
 *
 * SSR modes (see also GOLIT_DISABLED):
 * - GOLIT_SERVE_URL set (e.g. http://127.0.0.1:9777): POST HTML to the warm
 *   `golit serve` process (fast path, no per-request process spawn).
 * - Otherwise: run `golit transform` on a temp dir (cold path; needs GOLIT_BIN).
 */

$uri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);

// Let PHP's built-in server handle static assets directly.
if (preg_match('/\.(js|css|png|jpg|gif|svg|ico|woff2?)$/', $uri)) {
    return false;
}

$page = match ($uri) {
    '/', '/index.php' => 'home',
    '/about'          => 'about',
    default           => null,
};

if ($page === null) {
    http_response_code(404);
    echo '<!DOCTYPE html><html><body><h1>404 Not Found</h1></body></html>';
    exit;
}

// --- Render the page template into the layout ---
ob_start();
$pageTitle = ucfirst($page);
require __DIR__ . "/pages/{$page}.php";
$content = ob_get_clean();

ob_start();
require __DIR__ . '/pages/_layout.php';
$html = ob_get_clean();

// --- golit SSR middleware ---
if (getenv('GOLIT_DISABLED')) {
    echo $html;
    exit;
}

$serveUrl = getenv('GOLIT_SERVE_URL');
if ($serveUrl !== false && $serveUrl !== '') {
    $endpoint = rtrim($serveUrl, '/') . '/render';
    $ctx = stream_context_create([
        'http' => [
            'method'  => 'POST',
            'header'  => "Content-Type: text/html; charset=utf-8\r\n",
            'content' => $html,
            'timeout' => 60,
            'ignore_errors' => true,
        ],
    ]);
    $response = @file_get_contents($endpoint, false, $ctx);
    $status   = 0;
    if (isset($http_response_header)) {
        foreach ($http_response_header as $h) {
            if (preg_match('#^HTTP/\S+\s+(\d+)#', $h, $m)) {
                $status = (int) $m[1];
                break;
            }
        }
    }
    if ($response !== false && $status >= 200 && $status < 300) {
        echo $response;
    } else {
        echo $html;
        error_log("golit serve POST failed (HTTP {$status})");
    }
    exit;
}

$golitBin = getenv('GOLIT_BIN') ?: dirname(__DIR__, 3) . '/dist/golit';
$defsDir  = getenv('GOLIT_DEFS') ?: dirname(__DIR__) . '/bundles';

if (!is_file($golitBin)) {
    // No golit binary available — serve untransformed HTML.
    echo $html;
    exit;
}

$tmpDir = sys_get_temp_dir() . '/golit-php-' . uniqid();
mkdir($tmpDir, 0755, true);
file_put_contents("{$tmpDir}/page.html", $html);

$cmd = implode(' ', array_map('escapeshellarg', [
    $golitBin, 'transform', $tmpDir, '--defs', $defsDir,
])) . ' 2>&1';

exec($cmd, $output, $exitCode);

if ($exitCode === 0 && is_file("{$tmpDir}/page.html")) {
    echo file_get_contents("{$tmpDir}/page.html");
} else {
    echo $html;
    error_log("golit transform failed (exit {$exitCode}): " . implode("\n", $output));
}

foreach (glob("{$tmpDir}/*") ?: [] as $path) {
    if (is_file($path)) {
        @unlink($path);
    }
}
if (is_dir($tmpDir)) {
    @rmdir($tmpDir);
}
