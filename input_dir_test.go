package fusion

import (
	"reflect"
	"testing"

	"github.com/moovweb/golog"
)

var (
	CONFIG_SINGULAR_SUCCESS = map[interface{}]interface{}{
		"input_directory": "main",
	}
	CONFIG_SINGULAR_FAIL = map[interface{}]interface{}{
		"input_directory": []string{"main"},
	}
	CONFIG_PLURAL_SUCCESS = map[interface{}]interface{}{
		"input_directories": []string{"main", "vendor", "main/login"},
	}
	CONFIG_PLURAL_FAIL = map[interface{}]interface{}{
		"input_directories": "main",
	}
	CONFIG_BOTH_SUCCESS = map[interface{}]interface{}{
		"input_directory":   "main",
		"input_directories": []string{"main", "vendor", "main/login"},
	}
	CONFIG_EMPTY_SUCCESS = map[interface{}]interface{}{}

	SINGULAR_INPUT_DIRS = []string{"main"}
	PLURAL_INPUT_DIRS   = []string{"main", "vendor", "main/login"}
	BOTH_INPUT_DIRS     = []string{"main", "vendor", "main/login", "main"}
	EMPTY_INPUT_DIRS    = []string{}
)

func helper(t *testing.T, config map[interface{}]interface{}, errFlag bool, errMsg string, expectedInputDirs []string) {
	qb := &QuickBundlerInstance{}
	qb.Log = golog.NewLogger("")
	inputDirs, err := qb.getInputDirs(config)

	if errFlag {
		if err == nil {
			t.Fatalf("Expected an error: %s. But received no error.", errMsg)
		} else if err.Error() != errMsg {
			t.Fatalf("Expected an error: %s. But received an error: %s", errMsg, err.Error())
		}
	} else {
		if err != nil {
			t.Fatalf("Expected no error. But received an error: %s", err.Error())
		} else {
			resultInputDirs := make([]string, len(inputDirs))
			for i := 0; i < len(inputDirs); i++ {
				if str, ok := inputDirs[i].(string); ok {
					resultInputDirs[i] = str
				} else {
					t.Fatal("Failed to convert at least one element in InputDirs to string!")
					return
				}
			}
			if !reflect.DeepEqual(resultInputDirs, expectedInputDirs) {
				t.Fatalf("Expected an array of InputDirs as: %v. But received: %v", expectedInputDirs, resultInputDirs)
			}
		}
	}
}

func TestGetInputDirsSingularSuccess(t *testing.T) {
	helper(t, CONFIG_SINGULAR_SUCCESS, false, "", SINGULAR_INPUT_DIRS)
}

func TestGetInputDirsSingularFail(t *testing.T) {
	helper(t, CONFIG_SINGULAR_FAIL, true, SINGULAR_CONV_ERR, []string{})
}

func TestGetInputDirsPluralSuccess(t *testing.T) {
	helper(t, CONFIG_PLURAL_SUCCESS, false, "", PLURAL_INPUT_DIRS)
}

func TestGetInputDirsPluralFail(t *testing.T) {
	helper(t, CONFIG_PLURAL_FAIL, true, PLURAL_ARRAY_ERR, []string{})
}

func TestGetInputDirsBothSuccess(t *testing.T) {
	helper(t, CONFIG_BOTH_SUCCESS, false, "", BOTH_INPUT_DIRS)
}

func TestGetInputDirsEmptySuccess(t *testing.T) {
	helper(t, CONFIG_EMPTY_SUCCESS, false, "", EMPTY_INPUT_DIRS)
}
