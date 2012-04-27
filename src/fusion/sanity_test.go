package fusion


import(
  "testing"
)

func TestYAMLLoading(t *testing.T) {
	bundles, _, _ := getBundles("test/example-bundle.yml")
	if len(bundles) != 1 {
		t.Fatalf("there should be one bundle entry\n")
	}
	b := bundles[0].(map[interface{}]interface{})
	input_files := b[":input_files"].([]interface{})
	if len(input_files) != 1 {
		t.Fatalf("there should be 1 input file: %d\n", len(input_files))
	}
	name := input_files[0].(string)
	if name != "http://d1topzp4nao5hp.cloudfront.net/uranium-upload/0.1.23/uranium.js" {
		t.Fatalf("the input file name does not match\n")
	}
}

