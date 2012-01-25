package fusion

import(
	"io/ioutil"
)

/* Quick Bundler Instance */

type QuickBundlerInstance struct{
	Bundles []BundleConfig
}

func (qb *QuickBundlerInstance) Run() {

	var inputFiles []string
	var data string
	
	for _, config := range(qb.Bundles) {
		
		inputFiles = qb.gatherFiles(&config)
		data = ""
		
		for _, inputFile := range(inputFiles) {
			rawJS, err := ioutil.ReadFile(inputFile)
			
			if err != nil {
				panic("Couldn't open file:" + inputFile)
			}
			
			data += "\n" + string(rawJS) + "\n"
			
		}
		
		outputFile := qb.getOutputFile(&config)		
		err := ioutil.WriteFile(outputFile, []uint8(data), uint32(0777) )
		
		if err != nil {
			panic("Couldn't write file:" + outputFile)
		}
		
		
	}
	
}


/****** TODO(SJ): Put these in a base struct and inherit these methods ******/

func (qb *QuickBundlerInstance) gatherFiles(config *BundleConfig) (filenames []string) {
	
	return filenames
}

func (qb *QuickBundlerInstance) getRemoteFile(url string) (filepath string) {
	return "heyo"
}

func (qb *QuickBundlerInstance) getOutputFile(config *BundleConfig) (filepath string) {
	if len(config.OutputFile) == 0 {
		panic("Bundle missing output file.")
	}
	
	
	return "hey"
}