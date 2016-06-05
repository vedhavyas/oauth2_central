package config

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestLoadConfigFile(t *testing.T) {

	cases := []struct {
		fileName string
		result   bool
	}{
		{fileName: "no_file", result: false},
		{fileName: "../config_file.json", result: true},
	}

	for _, test := range cases {
		err := LoadConfigFile(test.fileName)
		result := true
		if err != nil {
			result = false
		}
		assert.Equal(t, result, test.result)
	}
}
