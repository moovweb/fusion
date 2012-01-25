package fusion

import(
	"url"
	"io/ioutil"
	"path/filepath"
)

/* Quick Bundler Instance */
/*
 * TODO(SJ): 
 *     - everytime I can throw a panic, wrap in a function so I can handle it w a defer nicely
 *     - move base functions into a base class and promote the the base interface
 */


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
	
	for _, inputFile := range(config.InputFiles) {
		if isURL(inputFile) {
			filenames = append(filenames, qb.getRemoteFile(inputFile) )
						
		} else {
			absolutePath := absolutize(inputFile)
			
			filenames = append(filenames, absolutePath)
		}			
	}
	
	if len(config.InputDirectory) != 0 {

		entries, err := ioutil.ReadDir(config.InputDirectory)

		if err != nil {
			panic("Cannot read input directory:" + config.InputDirectory)
		}
		
		for _, entry := range(entries) {
			filenames = append(filenames, absolutize(entry.Name) )
		}
		
	}	
	
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

/* Helper Functions */

func isURL(path string) (bool) {
	_, err := url.Parse(path)
	
	if err != nil {
		return false
	}
	
	return true
}

func absolutize(path string) (string) {
	absolutePath, err := filepath.Abs(path)
	
	if err != nil {
		panic("Cannot get absolute filepath for file:" + path)
	}
	
	return absolutePath
}