package helper

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// Must returns an error if string variable is empty
func Must(key string) (string, error) {
	result := strings.Trim(viper.GetString(key), "\t ")
	if "" == result {
		return "", fmt.Errorf("error: required key %s not present", key)
	}
	return result, nil
}
