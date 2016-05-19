package fusion

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type bundlerInstance struct {
	ProjectPath string
	OutputPath  string
	Version     string
	Bundles     []interface{}
}

func getBundles(bundlesPath string) (bundles []interface{}, projectPath *string, err error) {
	path, _ := filepath.Split(bundlesPath)
	data, err := ioutil.ReadFile(bundlesPath)
	if err != nil {
		return nil, nil, err
	}

	someBundles := make([]interface{}, 0)
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
	Run() []error
	GetConfig(config string) []interface{}
}
