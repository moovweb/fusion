module Fusion
  class Quick < Fusion::Basic

    def bundle(config)
      js = []
      input_files = gather_files(config)      

      input_files.each do |input_file|
        js << File.open(input_file, "r").read
      end
      
      js = js.join("\n")
      
      File.open(get_output_file(config), "w") { |f| f << js }

      js
    end

  end
end
