package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Environment struct {
	env map[string]string
}

func ReadEnvironment() *Environment {
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		keyValue := strings.Split(entry, "=")
		env[keyValue[0]] = keyValue[1]
	}
	return &Environment{env: env}
}

func (env *Environment) String(key string) string {
	return env.env[key]
}

func (env *Environment) Int(key string) (int, error) {
	value := env.env[key]
	if value == "" {
		return -1, fmt.Errorf("No configuration found for key %s", key)
	}
	return strconv.Atoi(value)
}

func (env *Environment) Base64(key string) ([]byte, error) {
	value := env.env[key]
	if value == "" {
		return nil, fmt.Errorf("No configuration found for key %s", key)
	}
	decoded, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err == nil {
		return decoded, err
	}
	return base64.StdEncoding.DecodeString(value)
}

func FindFile(filename, defaultFileName string) string {
	if filename != "" {
		if fileExists(filename) {
			return filename
		}
	}
	etcPath := fmt.Sprintf("/etc/%s", defaultFileName)
	if fileExists(etcPath) {
		return filename
	}
	home, _ := homedir.Dir()
	homepath := fmt.Sprintf("%s/.%s", home, defaultFileName)
	if fileExists(homepath) {
		return homepath
	}
	if fileExists(defaultFileName) {
		return defaultFileName
	}
	fmt.Printf("Config file %s not found", defaultFileName)
	os.Exit(1)
	return ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
