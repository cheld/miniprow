package config

import (
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
			Name:     "Single quotes",
			Template: "'{{.}}'",
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
		{
			Name:     "Empty",
			Template: " ",
			Input:    "World",
			Expected: " ",
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

	// Test error
	t.Run("Syntax Error", func(t *testing.T) {
		_, err := ProcessTemplate("Hello ${{.", "World")
		if err == nil {
			t.Fatalf("Error expected error.")
		}
	})

}

func TestProcessAllTemplates_mapInput(t *testing.T) {

	type mapTestCase struct {
		Name string

		Templates map[string]string
		Inputs    map[string]string

		Expected map[string]string
	}
	testcases := []mapTestCase{
		{
			Name: "Simple",
			Templates: map[string]string{
				"tpl1": "Hello {{ .Input }}",
			},
			Inputs: map[string]string{
				"Input": "World",
			},
			Expected: map[string]string{
				"tpl1": "Hello World",
			},
		},
		{
			Name: "Simple with Dollar",
			Templates: map[string]string{
				"tpl1": "${{ .Input }}",
			},
			Inputs: map[string]string{
				"Input": "Hello World",
			},
			Expected: map[string]string{
				"tpl1": "Hello World",
			},
		},
		{
			Name: "No Template",
			Templates: map[string]string{
				"tpl1": "Hello World",
			},
			Inputs: map[string]string{
				"Input": "Hello World",
			},
			Expected: map[string]string{
				"tpl1": "Hello World",
			},
		},
		{
			Name: "Two Templates",
			Templates: map[string]string{
				"tpl1": "Hello {{ .Input }}",
				"tpl2": "Hello again {{ .Input }}",
			},
			Inputs: map[string]string{
				"Input": "World",
			},
			Expected: map[string]string{
				"tpl1": "Hello World",
				"tpl2": "Hello again World",
			},
		},
		{
			Name:      "Empty",
			Templates: map[string]string{},
			Inputs: map[string]string{
				"Input": "World",
			},
			Expected: map[string]string{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := ProcessAllTemplates(tc.Templates, tc.Inputs)
			if err != nil {
				t.Fatalf("Error occured %s.", err)
			}
			for k, v := range result.(map[string]string) {
				if tc.Expected[k] != v {
					t.Fatalf("Error expected %s, got %s.", tc.Expected[k], v)
				}
			}
		})
	}
}

func TestProcessAllTemplates_listInput(t *testing.T) {

	type mapTestCase struct {
		Name string

		Templates []string
		Inputs    map[string]string

		Expected []string
	}
	testcases := []mapTestCase{
		{
			Name: "Simple",
			Templates: []string{
				"Hello {{ .Input }}",
			},
			Inputs: map[string]string{
				"Input": "World",
			},
			Expected: []string{
				"Hello World",
			},
		},
		{
			Name: "Two",
			Templates: []string{
				"Hello {{ .Input }}",
				"Hello again {{ .Input }}",
			},
			Inputs: map[string]string{
				"Input": "World",
			},
			Expected: []string{
				"Hello World",
				"Hello again World",
			},
		},
		{
			Name:      "Empty",
			Templates: []string{},
			Inputs: map[string]string{
				"Input": "World",
			},
			Expected: []string{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := ProcessAllTemplates(tc.Templates, tc.Inputs)
			if err != nil {
				t.Fatalf("Error occured %s.", err)
			}
			for i, v := range result.([]string) {
				if tc.Expected[i] != v {
					t.Fatalf("Error expected %s, got %s.", tc.Expected[i], v)
				}
			}
		})
	}
}

func TestProcessAllTemplates_error(t *testing.T) {

	type mapTestCase struct {
		Name string

		Templates []string
		Inputs    map[string]string

		Expected []string
	}
	testcases := []mapTestCase{
		{
			Name: "error",
			Templates: []string{
				"Hello {{ .Input }}",
				"Error {{.",
			},
			Inputs: map[string]string{
				"Input": "World",
			},
			Expected: []string{},
		},
	}
	for _, tc := range testcases {
		t.Run("Process all with error", func(t *testing.T) {
			_, err := ProcessAllTemplates(tc.Templates, tc.Inputs)
			if err == nil {
				t.Fatalf("Error expected error.")
			}
		})
	}
}
