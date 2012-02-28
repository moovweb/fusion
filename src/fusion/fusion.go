package fusion

import(
	"path/filepath"
	"io/ioutil"
	yaml "goyaml"
	"os"
)

type bundlerInstance struct{
	ProjectPath string
	Bundles []interface{}
}

func getBundles(bundlesPath string) (bundles []interface{}, projectPath *string, thisError *os.Error) {	

	path, _ := filepath.Split(bundlesPath)
	data, err := ioutil.ReadFile(bundlesPath)

	if err != nil {
		err = os.NewError("Couldn't read file:" + bundlesPath)
		return nil, nil, &err
	}
	
	someBundles := make([]interface{},0)
	err = yaml.Unmarshal(data, &someBundles)
		
	if err != nil {
		panic("Bad bundle format. Couldn't unmarshal: " + bundlesPath)
	}
	
	return someBundles, &path, nil
}

/*
This would be nice w reflection:

func NewBundlerFromBundleConfig(bundles []BundleConfig, projectPath string) (*bundlerInstance) {
	return &bundlerInstance{
		Bundles: bundles,
		ProjectPath: projectPath,
	}
}
*/

type Bundler interface {
	Run() bool
}
