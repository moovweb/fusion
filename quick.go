package fusion

import (
	"errors"
	"golog"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/* Quick Bundler Instance */
/*
 * TODO(SJ): 
 *     - everytime I can throw a panic, wrap in a function so I can handle it w a defer nicely
 *     - move base functions into a base class and promote the the base interface
 */

type QuickBundlerInstance struct {
	bundlerInstance
	Log *golog.Logger
}

func NewQuickBundler(bundlesPath string, logger *golog.Logger) (*QuickBundlerInstance, error) {
	bundles, projectPath, err := getBundles(bundlesPath)

	if err != nil {
		return nil, err
	}

	b := &QuickBundlerInstance{}
	b.Bundles = bundles
	b.ProjectPath = *projectPath
	b.Log = logger

	return b, nil
}

func (qb *QuickBundlerInstance) GetConfig(conf string) []interface{} {
	results := make([]interface{}, 0)
	for _, rawConfig := range qb.Bundles {
		config := rawConfig.(map[interface{}]interface{})
		if config[conf] != nil {
			results = append(results, config[conf])
		}
	}
	return results
}

func (qb *QuickBundlerInstance) Run() []error {

	var data string

	errors := make([]error, 0)

	for _, config := range qb.Bundles {

		inputFiles, err := qb.gatherFiles(config)

		if err != nil {
			qb.Log.Info("Couldn't create bundle: %v", err)
			errors = append(errors, err)
		}

		data = ""

		for _, inputFile := range inputFiles {

			rawJS, err := ioutil.ReadFile(inputFile)

			if err != nil {
				errors = append(errors, err)
				continue
			}

			relativeFilename := strings.Replace(inputFile, qb.ProjectPath, "", -1) // remove everything including $project/assets/javascript/
			if strings.HasPrefix(relativeFilename, ".remote/") {
				relativeFilename = strings.Replace(relativeFilename, ".remote/", "", 1) // remove the .remote/ prefix from remote files, since they already show the URL
				relativeFilename = strings.Replace(relativeFilename, "__colon__", ":", -1)
				relativeFilename = strings.Replace(relativeFilename, "__slash__", "/", -1)
			}
			data += "\n\n/*\n * File: " + relativeFilename + "\n */\n" + string(rawJS) + "\n"

		}

		outputFile := qb.getOutputFile(config)

		err = ioutil.WriteFile(outputFile, []uint8(data), os.FileMode(0644))

		if err != nil {
			qb.Log.Info("Couldn't write output js (%v): %v", outputFile, err)
			errors = append(errors, err)
		}

	}

	return errors
}

/****** TODO(SJ): Put these in a base struct and inherit these methods ******/

func (qb *QuickBundlerInstance) gatherFiles(rawConfig interface{}) (filenames []string, error error) {
	config := rawConfig.(map[interface{}]interface{})

	var files []interface{}

	if config[":input_files"] != nil {
		files = config[":input_files"].([]interface{})
	}

	for _, rawInputFile := range files {
		inputFile := rawInputFile.(string)

		if isURL(inputFile) {
			var path string
			path, error = qb.getRemoteFile(inputFile)
			if error != nil {
				return
			}
			filenames = append(filenames, path)
		} else {
			absolutePath := qb.Absolutize(inputFile)

			if strings.Contains(absolutePath, qb.ProjectPath) {
				filenames = append(filenames, absolutePath)
			} else {
				qb.Log.Warning("'%s' lies outside of the '/assets/javascript' directory; skipped.", inputFile)
			}
		}
	}

	var inputDirectory string

	if config[":input_directory"] != nil {
		inputDirectory = config[":input_directory"].(string)
	}

	var entries []os.FileInfo
	var absoluteDirectoryPath string

	if len(inputDirectory) != 0 {
		absoluteDirectoryPath = qb.Absolutize(inputDirectory)

		if !strings.Contains(absoluteDirectoryPath, qb.ProjectPath) {
			qb.Log.Warning("'%s' lies outside of the '/assets/javascript' directory; skipped.", inputDirectory)
		} else {
			newEntries, err := ioutil.ReadDir(absoluteDirectoryPath)
			if err != nil {
				qb.Log.Warning("Cannot read input directory: " + absoluteDirectoryPath)
			} else {
				entries = newEntries
			}
		}
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			qb.Log.Info("Skipped file " + entry.Name())
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".js") {
			qb.Log.Info("Skipped file " + entry.Name())
			continue
		}

		filenames = append(filenames, filepath.Join(absoluteDirectoryPath, entry.Name()))
	}

	return filenames, nil
}

func (qb *QuickBundlerInstance) getRemoteFile(url string) (string, error) {

	filename := strings.Replace(url, "/", "__slash__", -1)
	filename = strings.Replace(filename, ":", "__colon__", -1)

	remoteDirectory := filepath.Join(qb.ProjectPath, ".remote")
	path := filepath.Join(remoteDirectory, filename)

	_, err := os.Stat(path)

	if err == nil {
		return path, nil
	}

	response, err := http.Get(url)

	if err != nil {
		return path, errors.New("Error fetching file: " + url + ":" + err.Error())
	}

	if response.Status != "200 OK" {
		return path, errors.New("Error fetching file:" + url + ":: Status Code : " + response.Status)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return path, errors.New("Error reading response: " + err.Error())
	}

	_, err = os.Stat(remoteDirectory)

	if err != nil {
		err = os.Mkdir(remoteDirectory, os.FileMode(0755))

		if err != nil {
			return path, errors.New("Couldn't create directory: " + remoteDirectory + "(" + err.Error() + ")")
		}
	}

	ioutil.WriteFile(path, data, os.FileMode(0644))

	return path, nil
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

func isURL(path string) bool {
	thisURL, err := url.Parse(path)

	if err != nil {
		return false
	}

	realURL := false

	for _, prefix := range [4]string{"//", "http://", "https://", "ftp://"} {
		valid := strings.HasPrefix(thisURL.String(), prefix)
		if valid {
			realURL = true
		}
	}

	return realURL
}

func (qb *QuickBundlerInstance) Absolutize(path string) string {
	absolutePath := filepath.Join(qb.ProjectPath, path)

	return absolutePath
}
