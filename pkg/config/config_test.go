package config

import (
	"fmt"
	"testing"
)

func TestIsMatching(t *testing.T) {

	type testCase struct {
		name string

		ifContains string
		ifEquals   string
		ifTrue     string
		input      string

		expected bool
	}
	testcases := []testCase{
		{
			name:       "Contains: Simple test",
			ifContains: "/test",
			input:      "Please run more /test",
			expected:   true,
		},
		{
			name:       "Contains: No match",
			ifContains: "/test",
			input:      "Text without match",
			expected:   false,
		},
		{
			name:       "Contains: Twice",
			ifContains: "/test",
			input:      "Please run more /test and /test again",
			expected:   true,
		},
		{
			name:       "Contains: Empty input",
			ifContains: "/test",
			input:      "",
			expected:   false,
		},
		{
			name:       "Contains: Ignore case",
			ifContains: "/TEST",
			input:      "Please run more /test",
			expected:   true,
		},
		{
			name:       "Contains: space",
			ifContains: "/TEST all",
			input:      "Run /test all",
			expected:   true,
		},
		{
			name:     "Equals: Simple test",
			ifEquals: "/test",
			input:    "/test",
			expected: true,
		},
		{
			name:     "Equals: No match",
			ifEquals: "/test",
			input:    "Text /test",
			expected: false,
		},
		{
			name:     "Equals: Empty input",
			ifEquals: "/test",
			input:    "",
			expected: false,
		},
		{
			name:     "Equals: Ignore case",
			ifEquals: "/TEST",
			input:    "/test",
			expected: true,
		},
		{
			name:     "Equals: Space",
			ifEquals: "/TEST all",
			input:    "/test all",
			expected: true,
		},
		{
			name:     "IfTrue: Simple test",
			ifTrue:   "true",
			input:    "Hello world",
			expected: true,
		},
		{
			name:     "IfTrue: No match",
			ifTrue:   "false",
			input:    "Hello world",
			expected: false,
		},
		{
			name:     "IfTrue: Empty input",
			ifTrue:   "true",
			input:    "",
			expected: true,
		},
		{
			name:     "IfTrue: Ignore case",
			ifTrue:   "TRUE",
			input:    "Hello world",
			expected: true,
		},
		{
			name:     "IfTrue: Condition",
			ifTrue:   "${{ eq .Value \"test\" }}",
			input:    "test",
			expected: true,
		},
		{
			name:     "IfTrue: Error",
			ifTrue:   "${{ eq .",
			input:    "test",
			expected: false,
		},
		{
			name:     "Multi rule with one false",
			ifEquals: "test",
			ifTrue:   "false",
			input:    "test",
			expected: false,
		},
		{
			name:     "Multi rule with all true",
			ifEquals: "test",
			ifTrue:   "true",
			input:    "test",
			expected: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			event := Event{
				If_contains: tc.ifContains,
				If_equals:   tc.ifEquals,
				If_true:     tc.ifTrue,
			}
			source := Source{
				Value: tc.input,
			}
			result := event.IsMatching(source)
			if result != tc.expected {
				t.Fatalf("Error expected %v, got %v.", tc.expected, result)
			}
		})
	}
}

func TestTrigger(t *testing.T) {
	cfg := Configuration{}
	cfg.Triggers = []Trigger{
		{
			Name: "some-trigger-definition",
		},
	}
	type testcase struct {
		Name    string
		Trigger string
		Found   bool
	}
	testcases := []testcase{
		{
			Name:    "Simple match",
			Trigger: "some-trigger-definition",
			Found:   true,
		},
		{
			Name:    "Ignore case",
			Trigger: "sOme-Trigger-definition",
			Found:   true,
		},
		{
			Name:    "Not found",
			Trigger: "not available",
			Found:   false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			trigger := cfg.Trigger(tc.Trigger)
			if trigger == nil && tc.Found {
				t.Fatalf("Trigger not found")
			}
			if trigger != nil && !tc.Found {
				t.Fatalf("Trigger found, but not expected")
			}
		})
	}
}

func TestEvent(t *testing.T) {
	cfg := Configuration{}
	cfg.Events = []Event{
		{
			Source:    "github",
			Type:      "comment",
			If_equals: "/test",
		},
	}
	type testcase struct {
		Name string

		Source string
		Type   string
		Value  string

		Found bool
	}
	testcases := []testcase{
		{
			Name:   "Simple match",
			Source: "github",
			Type:   "comment",
			Value:  "/test",
			Found:  true,
		},
		{
			Name:   "Ignore case",
			Source: "GithuB",
			Type:   "COmmEnt",
			Value:  "/test",
			Found:  true,
		},
		{
			Name:   "No match",
			Source: "github",
			Type:   "comment",
			Value:  "/not-equals",
			Found:  false,
		},
		{
			Name:  "Empty source",
			Type:  "comment",
			Value: "/test",
			Found: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			source := Source{Value: tc.Value}
			event := cfg.Event(tc.Source, tc.Type, source)
			if event == nil && tc.Found {
				t.Fatalf("Event not found")
			}
			if event != nil && !tc.Found {
				t.Fatalf("Trigger found, but not expected")
			}
		})
	}
}

func TestDestinationCtx(t *testing.T) {

	source := Source{}
	source.Value = ""
	source.Payload = map[string]interface{}{
		"inputkey1": "inputvalue1",
		"inputkey2": "inputvalue2",
	}
	t.Run("mytest", func(t *testing.T) {

		event := Event{
			Trigger: "some-destination",
			Values: map[string]interface{}{
				"target": "test",
				"nested": map[string]string{
					"nestedkey": "nestedvalue",
				},
				"template": "String with {{ .Payload.inputkey1 }}",
			},
		}
		task, err := event.NewTask(source)
		if err != nil {
			//fmt.Println(err)
			t.Errorf("Error not exptected: %v", err)
		}
		if task.Trigger != "some-destination" {
			t.Errorf("got %s, want %s", task.Trigger, "some-destination")
		}
		if task.Values["target"] != "test" {
			t.Errorf("got %s, want %s", task.Values["target"], "test")
		}
		fmt.Println(task.Values["nested"])
		//if (task.Values["nested"]).(map[string]string)["nestedkey"] != "nestedvalue" {
		//	t.Errorf("got %s, want %s", "", "nestedvalue")
		//}
		if task.Values["template"] != "String with inputvalue1" {
			t.Errorf("got %s, want %s", task.Values["template"], "String with inputvalue1")
		}
	})
}
