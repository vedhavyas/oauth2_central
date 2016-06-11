package config

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestLoadConfigFile(t *testing.T) {

	cases := []struct {
		fileName       string
		expectedResult bool
	}{
		{fileName: "no_file", expectedResult: false},
		{fileName: "./config.json", expectedResult: false},
	}

	for _, test := range cases {
		err := LoadConfigFile(test.fileName)
		result := true
		if err != nil {
			result = false
		}
		assert.Equal(t, result, test.expectedResult)
	}
}
