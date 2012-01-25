package fusion

import(
	"path/filepath"
	"io/ioutil"
	yaml "launchpad.net/goyaml"
)

type BundleConfig struct {
	OutputFile string
	InputDirectory string
	InputFiles []string
}

type bundlerInstance struct{
	projectPath string
	Bundles []BundleConfig
}

func NewBundlerFromFile(bundlesPath string) (*bundlerInstance) {
	projectPath, _ := filepath.Split(bundlesPath)

	data, err := ioutil.ReadFile(bundlesPath)
	
	if err != nil {
		panic("Couldn't read file:" + bundlesPath)
	}
	
	bundles := make([]BundleConfig,0)
	
	err = yaml.Unmarshal(data, &bundles)
	
	if err != nil {
		panic("Bad bundle format. Couldn't unmarshal: " + bundlesPath)
	}
	
	return NewBundlerFromBundleConfig(bundles, projectPath)
}

func NewBundlerFromBundleConfig(bundles []BundleConfig, projectPath string) (*bundlerInstance) {
	return &bundlerInstance{
		Bundles: bundles,
	}
}

type Bundler interface {
	Run() bool
}