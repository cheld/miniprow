package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
)

func ReadConfigFile(fileName, environmentKey string) *[]byte {
	cfg, err := Environment.Value(environmentKey).Base64()
	if err == nil {
		return &cfg
	}
	filePath := FindExistingFile(DefaultConfigLocations(fileName))
	if filePath == "" {
		fmt.Printf("Config file %s not found\n", fileName)
		os.Exit(1)
	}
	cfg, err = ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Config file %s could not be loaded\n", fileName)
		os.Exit(1)
	}
	return &cfg
}

func DefaultConfigLocations(fileName string) []string {
	etcDir := fmt.Sprintf("/etc/%s", fileName)
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Cannot access home dir, %s", err)
	}
	homeDir := fmt.Sprintf("%s/.%s", home, fileName)
	currentDir := fileName
	locations := []string{currentDir, homeDir, etcDir}
	return locations
}

func FindExistingFile(filenames []string) string {
	for _, filename := range filenames {
		if fileExists(filename) {
			return filename
		}
	}
	return ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
