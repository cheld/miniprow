package config

import (
	"os"
	"strings"
)

func Environ(overrideVariables map[string]string) map[string]string {
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		keyValue := strings.Split(entry, "=")
		env[keyValue[0]] = keyValue[1]
	}
	if overrideVariables != nil {
		for k, v := range overrideVariables {
			env[k] = v
		}
	}
	return env
}
