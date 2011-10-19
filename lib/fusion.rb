require 'basic'
require 'quick'
require 'optimized'

module Fusion

  def configure(options)
    @options ||= {}
    @options.update(options)

    if @options[:bundle_file_path].nil? && @options[:bundle_configs].nil?
      raise Exception("Configuration error -- must specify #{:bundle_file_path} when configuring Fusion javascript bundler")
    end

    @options[:project_path] = File.join(@options[:bundle_file_path].split("/")[0..-2]) if @options[:bundle_file_path]
    
  end

  # So that the bundler can be configured and run in two different places ... like Sass
  module_function :configure

end
