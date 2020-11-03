package config

import (
	"fmt"
	"testing"
)

func TestProcessTemplate(t *testing.T) {

	type testCase struct {
		Name string

		Template string
		Input    string

		Expected string
	}

	testcases := []testCase{
		{
			Name:     "Simple Template",
			Template: "Hello {{ . }}",
			Input:    "World",
			Expected: "Hello World",
		},
		{
			Name:     "No Template",
			Template: "Hello World",
			Input:    "Test",
			Expected: "Hello World",
		},
		{
			Name:     "Only Template",
			Template: "{{.}}",
			Input:    "Hello World",
			Expected: "Hello World",
		},
		{
			Name:     "Start with dollar",
			Template: "${{.}}",
			Input:    "Hello World",
			Expected: "Hello World",
		},
		{
			Name:     "Dollar in middle",
			Template: "Hello ${{.}}",
			Input:    "World",
			Expected: "Hello World",
		},
		{
			Name:     "Dollar in middle",
			Template: "Hello ${{.",
			Input:    "World",
			Expected: "Hello World",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := ProcessTemplate(tc.Template, tc.Input)
			if result != tc.Expected && err == nil {
				t.Fatalf("Error expected %s, got %s.", tc.Expected, result)
			}
		})
	}

	t.Run("Syntax Error", func(t *testing.T) {
		_, err := ProcessTemplate("Hello ${{.", "World")
		if err != nil {
			t.Fatalf("Error expected error.")
		}
	})

}

func TestProcessAllTemplates(t *testing.T) {

	type testCase struct {
		Name string

		Templates map[string]string
		Inputs    map[string]string

		Expected map[string]string
	}

	testcases := []testCase{
		{
			Name: "Simple",
			Templates: map[string]string{
				"tpl1": "Hello {{ .Inputs.value1 }}",
			},
			Inputs: map[string]string{
				"value1": "World",
			},
			Expected: map[string]string{
				"tpl1": "Hello Wold",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			result := ProcessAllTemplates(tc.Templates, tc)
			fmt.Println(result)
			for k, v := range result.(map[string]string) {
				if tc.Expected[k] != v {
					t.Fatalf("Error expected %s, got %s.", tc.Expected[k], v)
				}
			}
		})
	}
}
