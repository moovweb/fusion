package fusion

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/moovweb/golog"
)

/* Quick Bundler Instance */
/*
 * TODO(SJ):
 *     - everytime I can throw a panic, wrap in a function so I can handle it w a defer nicely
 *     - move base functions into a base class and promote the the base interface
 */
const (
	SINGULAR_CONV_ERR = "Couldn't convert input directory into string. Make sure the field in bundles.yml is a line rather than a list."
	PLURAL_ARRAY_ERR  = "Input directories should be an array. Make sure the field in bundles.yml is a list rather than a line."
	NO_SINGULAR_WARN  = "No input directory specified, this might be an error. Please check to make sure it's not."
)

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
	b.Version = "0.1"

	return b, nil
}

func (qb *QuickBundlerInstance) GetConfig(conf string) []interface{} {
	results := make([]interface{}, 0)
	for _, rawConfig := range qb.Bundles {
		config := rawConfig.(map[interface{}]interface{})
		if config[conf] != nil {
			results = append(results, config[conf])
		} else if conf[0] == ':' {
			conf = conf[1:]
			if config[conf] != nil {
				results = append(results, config[conf])
			}
		} else {
			conf = ":" + conf
			if config[conf] != nil {
				results = append(results, config[conf])
			}
		}
	}
	return results
}

func (qb *QuickBundlerInstance) Run() []error {

	var data string

	errArr := make([]error, 0)

	for _, rawConfig := range qb.Bundles {
		config := rawConfig.(map[interface{}]interface{})

		inputFiles, err := qb.gatherFiles(config)

		if err != nil {
			qb.Log.Infof("Couldn't create bundle: %v", err)
			errArr = append(errArr, err)
		}

		data = ""

		includeVersion := true
		if config[":include_version"] != nil {
			if config[":include_version"].(bool) == false {
				includeVersion = false
			}
		} else if config["include_version"] != nil {
			if config["include_version"].(bool) == false {
				includeVersion = false
			}
		}
		if includeVersion == true {
			data += "// Bundled with Fusion v" + qb.Version + "\n\n"
		}

		protected := false
		if config[":protected"] != nil {
			if config[":protected"].(bool) == true {
				protected = true
				data += "(function() {"
			}
		} else if config["protected"] != nil {
			if config["protected"].(bool) == true {
				protected = true
				data += "(function() {"
			}
		}

		for _, inputFile := range inputFiles {

			rawJS, err := ioutil.ReadFile(inputFile)

			if err != nil {
				errArr = append(errArr, err)
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

		if protected == true {
			data += "\n})();"
		}

		outputFile, err := qb.getOutputFile(config)
		if err != nil {
			errArr = append(errArr, err)
			return errArr
		}

		err = ioutil.WriteFile(outputFile, []uint8(data), os.FileMode(0644))

		if err != nil {
			qb.Log.Infof("Couldn't write output js (%v): %v", outputFile, err)
			errArr = append(errArr, err)
		}

	}

	return errArr
}

/****** TODO(SJ): Put these in a base struct and inherit these methods ******/

func (qb *QuickBundlerInstance) gatherFiles(config map[interface{}]interface{}) (filenames []string, err error) {

	var files []interface{}

	if config[":input_files"] != nil {
		files = config[":input_files"].([]interface{})
	} else if config["input_files"] != nil {
		files = config["input_files"].([]interface{})
	} else {
		qb.Log.Warningf("No input files specified, this might be an error. Please check to make sure it's not.")
	}

	for _, rawInputFile := range files {
		inputFile := rawInputFile.(string)

		if isURL(inputFile) {
			var path string
			path, err = qb.getRemoteFile(inputFile)
			if err != nil {
				return
			}
			filenames = append(filenames, path)
		} else {
			absolutePath := qb.Absolutize(inputFile)

			if strings.Contains(absolutePath, qb.ProjectPath) {
				filenames = append(filenames, absolutePath)
			} else {
				qb.Log.Warningf("'%s' lies outside of the '/assets/javascript' directory; skipped.", inputFile)
			}
		}
	}

	inputDirs, err := qb.getInputDirs(config)
	if err != nil {
		return
	}

	var entries []string

	for _, rawDir := range inputDirs {
		dir := rawDir.(string)
		if len(dir) != 0 {
			absDir := qb.Absolutize(dir)

			if !strings.Contains(absDir, qb.ProjectPath) {
				qb.Log.Warningf("'%s' lies outside of the '/assets/javascript' directory; skipped.", absDir)
			} else {
				newEntries, err := ioutil.ReadDir(absDir)
				if err != nil {
					qb.Log.Warningf("Cannot read input directory: " + absDir)
				} else {
					for _, entry := range newEntries {
						entries = append(entries, filepath.Join(absDir, entry.Name()))
					}
				}
			}
		}
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry, ".") {
			qb.Log.Infof("Skipped file " + entry)
			continue
		}
		if !strings.HasSuffix(entry, ".js") {
			qb.Log.Infof("Skipped file " + entry)
			continue
		}

		filenames = append(filenames, entry)
	}

	return filenames, nil
}

func (qb *QuickBundlerInstance) getInputDirs(config map[interface{}]interface{}) (inputDirs []interface{}, err error) {
	// for backwards compatibility
	var inputDirectory string

	inputDirOld := config[":input_directory"]
	if inputDirOld == nil {
		inputDirOld = config["input_directory"]
	}
	if inputDirOld != nil {
		if str, ok := inputDirOld.(string); ok {
			inputDirectory = str
		} else {
			err = fmt.Errorf(SINGULAR_CONV_ERR)
			return inputDirs, err
		}
	} else {
		qb.Log.Warningf(NO_SINGULAR_WARN)
	}

	inputDirNew := config[":input_directories"]
	if inputDirNew == nil {
		inputDirNew = config["input_directories"]
	}
	if inputDirNew != nil {
		v := reflect.ValueOf(inputDirNew)
		if v.Kind() != reflect.Slice {
			err = fmt.Errorf(PLURAL_ARRAY_ERR)
			return inputDirs, err
		}
		for i := 0; i < v.Len(); i++ {
			inputDirs = append(inputDirs, v.Index(i).Interface())
		}
	}

	if inputDirectory != "" {
		inputDirs = append(inputDirs, interface{}(inputDirectory))
	}

	return inputDirs, err
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

func (qb *QuickBundlerInstance) getOutputFile(config map[interface{}]interface{}) (string, error) {
	var outputFile string
	if config[":output_file"] != nil {
		outputFile = config[":output_file"].(string)
	} else if config["output_file"] != nil {
		outputFile = config["output_file"].(string)
	} else {
		return "", errors.New("No output file specified, please specify an :output_file in bundles.yml.")
	}

	if len(outputFile) == 0 {
		return "", errors.New("Bundle missing output file.")
	}

	return filepath.Join(qb.ProjectPath, outputFile), nil
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
