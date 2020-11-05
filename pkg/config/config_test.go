package config

import (
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
			ifTrue:   "${{ eq .Objectiv \"test\" }}",
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
			eventInput := EventInput{
				Objectiv: tc.input,
			}
			result := event.IsMatching(eventInput)
			if result != tc.expected {
				t.Fatalf("Error expected %v, got %v.", tc.expected, result)
			}
		})
	}
}

func TestDestinationCtx(t *testing.T) {

	eventInpt := EventInput{}
	eventInpt.Objectiv = ""
	eventInpt.Input = map[string]interface{}{
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
				"template": "String with {{ .Input.inputkey1 }}",
			},
		}
		eventData := event.Handle(eventInpt)
		if eventData.Name != "some-destination" {
			t.Errorf("got %s, want %s", eventData.Name, "some-destination")
		}
		if eventData.Values["target"] != "test" {
			t.Errorf("got %s, want %s", eventData.Values["target"], "test")
		}
		if (eventData.Values["nested"]).(map[string]string)["nestedkey"] != "nestedvalue" {
			t.Errorf("got %s, want %s", "", "nestedvalue")
		}
		if eventData.Values["template"] != "String with inputvalue1" {
			t.Errorf("got %s, want %s", eventData.Values["template"], "String with inputvalue1")
		}
	})
}
