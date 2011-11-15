module Fusion
  class Pretty < Fusion::Optimized

    private

    def get_options(config)
      options = super(config)

      compilation_level_index = options.find_index {|opt| opt =~ /compilation_level/}
      options[compilation_level_index] = ["--compilation_level  WHITESPACE_ONLY"]
      options << "--formatting PRETTY_PRINT"

      options      
    end

  end
end
