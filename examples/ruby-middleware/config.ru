require 'erb'
require 'rack'
require_relative 'lib/golit_middleware'

# Serve static assets (JS components) from the components/ directory.
map '/components' do
  run Rack::Files.new(File.join(__dir__, 'components'))
end

# Rack app — a minimal router + ERB renderer.
# The same GolitMiddleware class works identically in Rails:
#
#   config.middleware.use GolitMiddleware,
#     golit_bin: '/usr/local/bin/golit',
#     defs_dir:  Rails.root.join('bundles').to_s
#
app = Class.new do
  ROUTES = {
    '/'      => 'home',
    '/about' => 'about',
  }.freeze

  def call(env)
    page = ROUTES[env['PATH_INFO']]
    return not_found unless page

    content    = render(page)
    @title     = page.capitalize
    @content   = content
    @bench     = Rack::Utils.parse_query(env['QUERY_STRING']).key?('bench')
    body       = render('layout')

    [200, { 'content-type' => 'text/html; charset=utf-8' }, [body]]
  end

  private

  def render(template)
    path = File.join(__dir__, 'views', "#{template}.erb")
    ERB.new(File.read(path)).result(binding)
  end

  def not_found
    [404, { 'content-type' => 'text/html' }, ['<h1>404 Not Found</h1>']]
  end
end.new

golit_bin = ENV.fetch('GOLIT_BIN') {
  File.join(__dir__, '..', '..', 'dist', 'golit')
}
defs_dir = ENV.fetch('GOLIT_DEFS') {
  File.join(__dir__, 'bundles')
}

use GolitMiddleware, golit_bin: golit_bin, defs_dir: defs_dir

run app
