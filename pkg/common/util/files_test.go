package util

import (
	"testing"
)

func TestDefaultConfigLocations(t *testing.T) {

	// test if default config locations are added to file
	locations := DefaultConfigLocations("test.yaml")
	if locations[0] != "test.yaml" {
		t.Errorf("Expected current dir, but is %s", locations[0])
	}
	if locations[2] != "/etc/test.yaml" {
		t.Errorf("Expected current dir, but is %s", locations[0])
	}

	// assert that no error thrown
	locations = DefaultConfigLocations("")
	if len(locations[0]) != 0 {
		t.Error("Should handle empty string gracefully")
	}
}

func TestFindExistingFile(t *testing.T) {
	givenFiles := []string{"non-existing.yaml", "/etc/environment"}
	existingFile := FindExistingFile(givenFiles)
	if existingFile != "/etc/environment" {
		t.Error("File not found")
	}
}

func TestReadConfiguration(t *testing.T) {

	// test existing file
	content, _ := ReadConfiguration("/etc/environment", "")
	if len(*content) == 0 {
		t.Error("File could not be loaded")
	}

}
