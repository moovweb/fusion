require 'yaml'
require 'open4'
require 'uri'
require 'cgi'
require 'rubygems' #I'm getting an error loading mechanize ... do I need to load this?
require 'mechanize'
require 'logger'

module Fusion

  def configure(options)
    @options ||= {}
    @options.update(options)

    raise Exception("Configuration error -- must specify #{:bundle_file_path} when configuring Fusion javascript bundler") if @options[:bundle_file_path].nil?

    @options[:project_path] = File.join(@options[:bundle_file_path].split("/")[0..-2])
  end

  # So that the bundler can be configured and run in two different places ... like Sass
  module_function :configure

  class Basic

    def initialize
      @bundle_options = Fusion.instance_variable_get('@options') 
      @log = @bundle_options[:logger] || Logger.new(STDOUT)

      @bundle_configs = YAML::load(File.open(@bundle_options[:bundle_file_path]))
    end

    def run
      start = Time.now

      @bundle_configs.each do |config|
        bundle(config)
      end

      @log.debug "Javascript Reloaded #{@bundle_configs.size} bundle(s) (#{Time.now - start}s)"
    end

    def gather_files(config)
      input_files = []

      if(config[:input_files])
        config[:input_files].each do |input_file|
          if (input_file =~ URI::regexp).nil?
            # Not a URL
            input_files << File.join(@bundle_options[:project_path], input_file)
          else
            # This is a remote file, if we don't have it, get it
            input_files << get_remote_file(input_file)
          end          
        end
      end

      if (config[:input_directory])
        directory = File.join(@bundle_options[:project_path],config[:input_directory])

        file_names = Dir.open(directory).entries.sort.find_all {|filename| filename.end_with?(".js") }

        input_files += file_names.collect do |file_name|
          File.join(directory, file_name)
        end
      end

      input_files.uniq
    end

    def get_output_file(config)
      raise Exception.new("Undefined js bundler output file") if config[:output_file].nil?
      File.join(@bundle_options[:project_path], config[:output_file])
    end

    def get_remote_file(url)
      filename = CGI.escape(url)
      local_directory = File.join(@bundle_options[:project_path], ".remote")
      local_file_path = File.join(local_directory, filename)
      
      return local_file_path if File.exists?(local_file_path)
      
      @log.debug "Fetching remote file (#{url})"

      m = Mechanize.new
      response = m.get(url)

      raise Exception.new("Error downloading file (#{url}) -- returned code #{repsonse.code}") unless response.code == "200"
      
      @log.debug "Got file (#{url})"

      unless Dir.exists?(local_directory)
        Dir.mkdir(local_directory)
      end
                  
      File.open(local_file_path,"w") {|f| f << response.body}
      
      local_file_path
    end
    

  end

  class Quick < Basic

    def bundle(config)
      js = []

      input_files = gather_files(config)      

      input_files.each do |input_file|
        js << File.open(input_file, "r").read
      end
      
      js = js.join("\n")
      
      File.open(get_output_file(config), "w") { |f| f << js }
    end

  end

  class Optimized < Basic
    
    def bundle(config)
      options = []

      options << ["js_output_file", get_output_file(config)]
      options << ["compilation_level", "SIMPLE_OPTIMIZATIONS"]

      gather_files(config).each do |input_file|
        options << ["js", input_file]
      end

      options.collect! do |option|
        "--#{option[0]} #{option[1]}"
      end

      jar_file = File.join(__FILE__.split("/")[0..-3].join("/"), "/compiler/compiler.jar")
      cmd = "java -jar #{jar_file} #{options.join(" ")}"
      io = IO.popen(cmd, "w")
      io.close

      raise Exception.new("Error creating bundle: #{get_output_file(config)}") unless $?.exitstatus == 0
    end
    
  end

end
