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

    if @options[:bundle_file_path].nil? && @options[:bundle_configs].nil?
      raise Exception("Configuration error -- must specify #{:bundle_file_path} when configuring Fusion javascript bundler")
    end

    @options[:project_path] = File.join(@options[:bundle_file_path].split("/")[0..-2]) if @options[:bundle_file_path]
    
  end

  # So that the bundler can be configured and run in two different places ... like Sass
  module_function :configure

  class Basic

    def initialize
      @bundle_options = Fusion.instance_variable_get('@options') 
      @log = @bundle_options[:logger] || Logger.new(STDOUT)

      if @bundle_options[:bundle_configs]
        @bundle_configs = @bundle_options[:bundle_configs]
      else
        @bundle_configs = YAML::load(File.open(@bundle_options[:bundle_file_path]))
      end

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

        unless File.exists?(directory)
          @log.debug "Path #{directory} does not exist"
          FileUtils.mkpath(directory)
          @log.debug "Created path: #{directory}"
        end

        file_names = Dir.open(directory).entries.sort.find_all {|filename| filename.end_with?(".js") }

        input_files += file_names.collect do |file_name|
          File.join(directory, file_name)
        end
      end

      input_files.uniq
    end

    def get_output_file(config)
      raise Exception.new("Undefined js bundler output file") if config[:output_file].nil?
      output_file = File.join(@bundle_options[:project_path], config[:output_file])
      path_directories = output_file.split("/")

      if path_directories.size > 1
        path = File.join(File.expand_path("."), path_directories[0..-2].join("/"))
        FileUtils::mkpath(File.join(path,"/"))
      end

      output_file
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
      options << ["language_in","ECMASCRIPT5"] # This will be compatible w all newer browsers, and helps us avoid old IE quirks

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

  class DebugOptimized < Optimized
    def gather_files(config)
      @log.debug "Warning ... using Debug compiler."
      input_files = []

      if(config[:input_files])
        config[:input_files].each do |input_file|
          @log.debug "Remote file? #{!(input_file =~ URI::regexp).nil?}"

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

      input_files.uniq!

      # Now wrap each file in a try/catch block and update the input_files list

      FileUtils::mkpath(File.join(@bundle_options[:project_path],".debug"))

      input_files.collect do |input_file|
        contents = File.open(input_file).read
        new_input_file =""
        file_name = input_file.split("/").last

        if input_file.include?(".remote")
          new_input_file = input_file.gsub(".remote",".debug")
        else
          new_input_file = File.join(@bundle_options[:project_path], ".debug", file_name)
        end

        new_contents = "///////////////////\n//mw_bundle: #{input_file}\n///////////////////\n\n try{\n#{contents}\n}catch(e){\nconsole.log('Error (' + e + 'generated in : #{input_file}');\n}"
        File.open(new_input_file,"w") {|f| f << new_contents}

        new_input_file
      end

    end
  end

end
