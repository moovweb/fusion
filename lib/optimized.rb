module Fusion
  class Optimized < Fusion::Basic
    
    def bundle(config)
      options = get_options(config)

      jar_file = File.join(__FILE__.split("/")[0..-3].join("/"), "/compiler/compiler.jar")
      cmd = "java -jar #{jar_file} #{options.join(" ")}"
      io = IO.popen(cmd, "w")
      io.close

      raise Exception.new("Error creating bundle: #{get_output_file(config)}") unless $?.exitstatus == 0

      File.read(get_output_file(config))
    end

    private

    def get_options(config)
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

      options
    end

  end

  class AdvancedOptimized < Fusion::Optimized
    def get_options(config)
      options = super(config)

      compilation_level_index = options.find_index {|opt| opt =~ /compilation_level/}
      options[compilation_level_index] = ["--compilation_level  ADVANCED_OPTIMIZATIONS"]

      options      
    end
  end  

end
