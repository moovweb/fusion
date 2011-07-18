Gem::Specification.new do |s|
  s.name = "fusion"
  s.version = "0.0.5"
  s.platform = Gem::Platform::RUBY

  s.authors = ["Sean Jezewski"]
  s.email = ["sean@moovweb.com"]

  s.homepage = "http://github.com/moovweb/fusion"
  s.summary = "Simple javascript bundler plugin."
  s.description = "Fusion bundles and re-bundles your javascript in two modes - quick (dumb concatenation) and optimized (google closure compiler's SIMPLE_OPTIMIZATIONS level"

  s.files = Dir['README.md','lib/*','compiler/*']
  s.require_path = "lib"
  
  s.add_dependency "open4"
  s.add_dependency "mechanize"

end
