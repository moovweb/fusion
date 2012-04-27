package fusion


import(
  "io/ioutil"
  "testing"
)

func TestYAMLLoading(t *testing.T) {
  rawBundles, _, _ := getBundles("test/example-bundle.yml")
	bundles := rawBundles.(map[interface{}]interface{})

	var files []interface{}

	if bundles[":input_files"] != nil {
		files = bundles[":input_files"].([]interface{})
	}
	
	if len(files) != 1 {
	  t.FailNow("Input files array is empty!")
	}

  inputFile := files[0].(string)
	expected := "http://d1topzp4nao5hp.cloudfront.net/uranium-upload/0.1.23/uranium.js"
	
  if inputFile != expected {
    t.FailNow("Invalid input file. Got (" + inputFile + "), expected (" + expected + ")")
  }
	
}

