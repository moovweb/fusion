require 'yaml'
require 'open4'
require 'uri'
require 'cgi'
require 'rubygems' #I'm getting an error loading mechanize ... do I need to load this?
require 'mechanize'
require 'logger'

module Fusion

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

      bundles = @bundle_configs.collect do |config|
        bundle(config)
      end

      @log.debug "Javascript Reloaded #{@bundle_configs.size} bundle(s) (#{Time.now - start}s)"

      bundles
    end

    def gather_files(config)
      input_files = []

      if(config[:input_files])
        config[:input_files].each do |input_file|
          if (input_file =~ URI::regexp).nil?
            # Not a URL
            file_path = input_file

            unless input_file == File.absolute_path(input_file)
              file_path = File.join(@bundle_options[:project_path], input_file)
            end

            input_files << file_path
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

end
