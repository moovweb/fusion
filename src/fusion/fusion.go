package fusion

import(
	"path/filepath"
	"io/ioutil"
	yaml "launchpad.net/goyaml"
)

/* -- LAME : go-yaml doesn't change the names to camel case
type BundleConfig struct {
	OutputFile string
	InputDirectory string
	InputFiles []string
}
*/

type BundleConfig struct {
	Output_file string
	Input_directory string
	Input_files []string
}


type bundlerInstance struct{
	ProjectPath string
	Bundles []interface{}
}

func GetBundles(bundlesPath string) (bundles []interface{}, projectPath string) {	
	projectPath, _ = filepath.Split(bundlesPath)

	data, err := ioutil.ReadFile(bundlesPath)
	
	if err != nil {
		panic("Couldn't read file:" + bundlesPath)
	}
	
	someBundles := make([]interface{},0)
	err = yaml.Unmarshal(data, &someBundles)
		
	if err != nil {
		panic("Bad bundle format. Couldn't unmarshal: " + bundlesPath)
	}
	
	return someBundles, projectPath
}


func getBundles(bundlesPath string) (bundles []BundleConfig, projectPath string) {	
	projectPath, _ = filepath.Split(bundlesPath)

	data, err := ioutil.ReadFile(bundlesPath)
	
	if err != nil {
		panic("Couldn't read file:" + bundlesPath)
	}
	
	someBundles := make([]BundleConfig,0)
	err = yaml.Unmarshal(data, &someBundles)
		
	if err != nil {
		panic("Bad bundle format. Couldn't unmarshal: " + bundlesPath)
	}
	
	return someBundles, projectPath
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