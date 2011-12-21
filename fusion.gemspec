version = File.read("VERSION").strip
if File.exists?"JENKINS"
  version += "."
  version += File.read("JENKINS").strip
end

buildf = File.open("BUILD_VERSION", 'w')
buildf.puts version
buildf.close

Gem::Specification.new do |s|
  s.name = "fusion"
  s.version = version
  s.platform = Gem::Platform::RUBY

  s.authors = ["Sean Jezewski"]
  s.email = ["sean@moovweb.com"]

  s.homepage = "http://github.com/moovweb/fusion"
  s.summary = "Simple javascript bundler plugin."
  s.description = "Fusion bundles and re-bundles your javascript in two modes - quick (dumb concatenation) and optimized (google closure compiler's SIMPLE_OPTIMIZATIONS level"

  s.files = Dir['README.md', 'BUILD_VERSION', 'Gemfile', 'Gemfile.lock', 'Rakefile', 'lib/*', 'compiler/*']
  s.require_path = "lib"

  s.add_development_dependency('moov_build_tasks', ['~> 0.2.0'])
  s.add_dependency "open4"
  s.add_dependency "mechanize"
end

