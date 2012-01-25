package fusion

import(
	"url"
	"io/ioutil"
	"path/filepath"
	"http"
	"strings"
	"os"
)

/* Quick Bundler Instance */
/*
 * TODO(SJ): 
 *     - everytime I can throw a panic, wrap in a function so I can handle it w a defer nicely
 *     - move base functions into a base class and promote the the base interface
 */

type QuickBundlerInstance struct{
	bundlerInstance
	client *http.Client
}

func NewQuickBundler(bundlesPath string) (*QuickBundlerInstance) {
//	bundles, projectPath := getBundles(bundlesPath)
	getBundles(bundlesPath)	
	
	b := &QuickBundlerInstance{}
//	b.Bundles = bundles
//	b.ProjectPath = projectPath

	return b
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
	
	for _, inputFile := range(config.Input_files) {
		if isURL(inputFile) {
			filenames = append(filenames, qb.getRemoteFile(inputFile) )
						
		} else {
			absolutePath := absolutize(inputFile)
			
			filenames = append(filenames, absolutePath)
		}			
	}
	
	if len(config.Input_directory) != 0 {

		entries, err := ioutil.ReadDir(config.Input_directory)

		if err != nil {
			panic("Cannot read input directory:" + config.Input_directory)
		}
		
		for _, entry := range(entries) {
			filenames = append(filenames, absolutize(entry.Name) )
		}
		
	}	
	
	return filenames
}

func (qb *QuickBundlerInstance) getRemoteFile(url string) (path string) {
	if qb.client == nil {
		// Initialize http client
		qb.client = &http.Client{}
	}
	
	response, err := qb.client.Get(url)
	
	if err != nil || response.Status != "200" {
		panic("Error fetching file:" + url + ":" + err.String() )
	}
	
	var data []byte

	_, err = response.Body.Read(data)
	
	if err != nil {
		println("Error reading response")
		panic(err)
	}
	
	remoteDirectory := filepath.Join(qb.ProjectPath, ".remote")
	
	_, err = os.Stat(remoteDirectory)
	
	if err != nil {
		err = os.Mkdir(remoteDirectory, uint32(0777) )

		if err != nil {
			println("Couldn't create directory:" + remoteDirectory)
			panic(err)
		}
	}
	
	filename := strings.Replace(url, "/", "_", -1)
	filename = strings.Replace(filename, ":", "", -1)
	path = filepath.Join(remoteDirectory, filename)

	ioutil.WriteFile(path, data, uint32(0777) )

	return path
}

func (qb *QuickBundlerInstance) getOutputFile(config *BundleConfig) (path string) {
	if len(config.Output_file) == 0 {
		panic("Bundle missing output file.")
	}
	
	return filepath.Join(qb.ProjectPath, config.Output_file)
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