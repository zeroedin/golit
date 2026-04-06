require 'tmpdir'
require 'fileutils'
require 'net/http'
require 'uri'

# Rack middleware that pipes HTML responses through golit SSR.
#
# Fast path: set GOLIT_SERVE_URL (e.g. http://127.0.0.1:9777) to POST HTML to
# a long-running `golit serve` process.
# Cold path: uses GOLIT_BIN + temp dir + `golit transform` (no extra process).
#
# This middleware works with any Rack-compatible framework:
# Rails, Sinatra, Hanami, Roda, or plain Rack apps.
#
#   # Rails example (config/application.rb):
#   config.middleware.use GolitMiddleware,
#     golit_bin: '/usr/local/bin/golit',
#     defs_dir:  Rails.root.join('bundles').to_s
#
class GolitMiddleware
  def initialize(app, golit_bin:, defs_dir:)
    @app      = app
    @golit    = golit_bin
    @defs_dir = defs_dir
  end

  def call(env)
    status, headers, body = @app.call(env)
    return [status, headers, body] if ENV['GOLIT_DISABLED']

    content_type = rack_header(headers, 'content-type').to_s
    return [status, headers, body] unless content_type.include?('text/html')

    html = +""
    body.each { |chunk| html << chunk }
    body.close if body.respond_to?(:close)

    transformed = transform(html)

    rack_delete_header!(headers, 'content-length')
    headers['content-length'] = transformed.bytesize.to_s
    [status, headers, [transformed]]
  end

  private

  # Rack allows mixed-case header keys depending on server and framework.
  def rack_header(headers, name)
    needle = name.downcase
    headers.each { |k, v| return v if k.to_s.downcase == needle }
    nil
  end

  def rack_delete_header!(headers, name)
    needle = name.downcase
    headers.reject! { |k, _| k.to_s.downcase == needle }
  end

  def transform(html)
    base = ENV['GOLIT_SERVE_URL'].to_s.strip
    unless base.empty?
      uri = URI("#{base.chomp('/')}/render")
      res = Net::HTTP.start(
        uri.host, uri.port,
        use_ssl: uri.scheme == 'https',
        open_timeout: 5,
        read_timeout: 60
      ) do |http|
        http.post(uri.path, html, { 'Content-Type' => 'text/html; charset=utf-8' })
      end
      return res.body if res.is_a?(Net::HTTPSuccess)

      $stderr.puts "golit serve failed (#{res.code}) — serving untransformed HTML"
      return html
    end

    Dir.mktmpdir('golit-rack-') do |dir|
      path = File.join(dir, 'page.html')
      File.write(path, html)

      success = system(@golit, 'transform', dir, '--defs', @defs_dir,
                        out: File::NULL, err: :out)

      if success && File.exist?(path)
        File.read(path)
      else
        $stderr.puts "golit transform failed — serving untransformed HTML"
        html
      end
    end
  end
end
