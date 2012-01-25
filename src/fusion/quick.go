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
}

func NewQuickBundler(bundlesPath string) (*QuickBundlerInstance) {
	bundles, projectPath := getBundles(bundlesPath)
	
	b := &QuickBundlerInstance{}
	b.Bundles = bundles
	b.ProjectPath = projectPath

	return b
}

func (qb *QuickBundlerInstance) Run() {

	var inputFiles []string
	var data string
	
	for _, config := range(qb.Bundles) {
		
		inputFiles = qb.gatherFiles(config)
		data = ""
		
		for _, inputFile := range(inputFiles) {
			
			rawJS, err := ioutil.ReadFile(inputFile)
			
			if err != nil {
				panic("Couldn't open file:" + inputFile)
			}
			
			data += "\n" + string(rawJS) + "\n"
			
		}
		
		outputFile := qb.getOutputFile(config)		
		
		err := ioutil.WriteFile(outputFile, []uint8(data), uint32(0777) )
		
		if err != nil {
			panic("Couldn't write file:" + outputFile)
		}
		
		
	}
	
}


/****** TODO(SJ): Put these in a base struct and inherit these methods ******/

func (qb *QuickBundlerInstance) gatherFiles(rawConfig interface{}) (filenames []string) {
//	config := rawConfig.(map[string]interface{})	
	config := rawConfig.(map[interface{}]interface{})	
	
	var files []interface{}
	
	if config[":input_files"] != nil {
		files = config[":input_files"].([]interface{})
	}
	
	for _, rawInputFile := range(files) {
		inputFile := rawInputFile.(string)


		if isURL(inputFile) {
			filenames = append(filenames, qb.getRemoteFile(inputFile) )
						
		} else {
			absolutePath := qb.Absolutize(inputFile)

			filenames = append(filenames, absolutePath)
		}			
	}
	
	var inputDirectory string
	
	if config[":input_directory"] != nil {
		inputDirectory = config[":input_directory"].(string)
	}
	
	if len(inputDirectory) != 0 {

		entries, err := ioutil.ReadDir(inputDirectory)

		if err != nil {
			panic("Cannot read input directory:" + inputDirectory)
		}
		
		for _, entry := range(entries) {
			filenames = append(filenames, qb.Absolutize(entry.Name) )
		}
		
	}	
	
	return filenames
}

func (qb *QuickBundlerInstance) getRemoteFile(url string) (path string) {

	filename := strings.Replace(url, "/", "_", -1)
	filename = strings.Replace(filename, ":", "", -1)

	remoteDirectory := filepath.Join(qb.ProjectPath, ".remote")
	path = filepath.Join(remoteDirectory, filename)
		
	_, err := os.Stat(path)
	
	if err == nil {
		return path
	}
	
	response, err := http.Get(url)	
	
	if err != nil {
		panic("Error fetching file:" + url + ":" + err.String() )
	}
	
	if response.Status != "200 OK" {
		panic("Error fetching file:" + url + ":: Status Code : " + response.Status )
	}
	
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body);

	if err != nil {
		println("Error reading response")
		panic(err)
	}
		
	_, err = os.Stat(remoteDirectory)
	
	if err != nil {
		err = os.Mkdir(remoteDirectory, uint32(0777) )

		if err != nil {
			println("Couldn't create directory:" + remoteDirectory)
			panic(err)
		}
	}

	ioutil.WriteFile(path, data, uint32(0777) )

	return path
}

func (qb *QuickBundlerInstance) getOutputFile(rawConfig interface{}) (path string) {
	config := rawConfig.(map[interface{}]interface{})

	outputFile := config[":output_file"].(string)

	if len(outputFile) == 0 {
		panic("Bundle missing output file.")
	}
	
	return filepath.Join(qb.ProjectPath, outputFile)
}

/* Helper Functions */

func isURL(path string) (bool) {
	thisURL, err := url.Parse(path)
	
	if err != nil {		
		return false
	}
	
	realURL := false

	for _, prefix := range( [4]string{"//","http://","https://","ftp://"} ) {
			valid := strings.HasPrefix(thisURL.String(), prefix)
			if valid {
				realURL = true
			}
	}
	
	return realURL
}

func (qb *QuickBundlerInstance) Absolutize(path string) (string) {
	absolutePath := filepath.Join(qb.ProjectPath, path)
		
	return absolutePath
}