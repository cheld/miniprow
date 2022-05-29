package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

func ReadConfiguration(fileName, environmentKey string) (*[]byte, error) {
	cfg, err := Environment.Value(environmentKey).Base64()
	if err == nil {
		logrus.Infof("Configuration loaded from environemnt with key %s\n", environmentKey)
		return &cfg, nil
	}
	if fileName == "" {
		return nil, fmt.Errorf("No configuration found")
	}
	logrus.Infof("Loading configuration from %s\n", fileName)
	cfg, err = ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Config file %s could not be loaded, %s", fileName, err)
	}
	return &cfg, nil
}

func DefaultConfigLocations(fileName string) []string {
	etcDir := fmt.Sprintf("/etc/%s", fileName)
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Cannot access home dir, %s", err)
	}
	currentDir := fileName
	homeDir := fmt.Sprintf("%s/.%s", home, fileName)
	locations := []string{currentDir, homeDir, etcDir}
	return locations
}

func FindExistingFile(filenames []string) string {
	for _, filename := range filenames {
		if FileExists(filename) {
			return filename
		}
	}
	return ""
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
