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
			rule := Rule{
				If_contains: tc.ifContains,
				If_equals:   tc.ifEquals,
				If_true:     tc.ifTrue,
			}
			request := Request{
				Value: tc.input,
			}
			ctx := Ctx{
				Request: request,
			}
			result := rule.IsMatching(ctx)
			if result != tc.expected {
				t.Fatalf("Error expected %v, got %v.", tc.expected, result)
			}
		})
	}
}

func TestFindTrigger(t *testing.T) {
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
			trigger := cfg.FindTrigger(tc.Trigger)
			if trigger == nil && tc.Found {
				t.Fatalf("Trigger not found")
			}
			if trigger != nil && !tc.Found {
				t.Fatalf("Trigger found, but not expected")
			}
		})
	}
}

func TestFindEvent(t *testing.T) {
	cfg := Configuration{}
	cfg.Rules = []Rule{
		{
			Event:     "github",
			If_equals: "/test",
		},
	}
	type testcase struct {
		Name string

		Event string
		Value string

		Found bool
	}
	testcases := []testcase{
		{
			Name:  "Simple match",
			Event: "github",
			Value: "/test",
			Found: true,
		},
		{
			Name:  "Ignore case",
			Event: "GithuB",
			Value: "/test",
			Found: true,
		},
		{
			Name:  "No match",
			Event: "github",
			Value: "/not-equals",
			Found: false,
		},
		{
			Name:  "Empty",
			Value: "/test",
			Found: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			request := Request{Value: tc.Value}
			source := Ctx{Request: request}
			event := cfg.GetMatchingRule(tc.Event, source)
			if event == nil && tc.Found {
				t.Fatalf("Event not found")
			}
			if event != nil && !tc.Found {
				t.Fatalf("Trigger found, but not expected")
			}
		})
	}
}

func TestBuildTask(t *testing.T) {
	event := Rule{
		Trigger: map[string]interface{}{
			"action": "mytrigger",
			"simple": "test",
			"nested": map[string]string{
				"nestedkey": "nestedvalue",
			},
			"template": "Hello {{ .Payload.sourcekey }}",
		},
	}
	type testcase struct {
		Name string

		Payload map[string]string

		Expected Task
	}
	testcases := []testcase{
		{
			Name: "Simple",

			Payload: map[string]string{
				"sourcekey": "World",
			},

			Expected: Task{
				Trigger: "mytrigger",
				Values: map[string]interface{}{
					"simple": "test",
					"nested": map[string]string{
						"nestedkey": "nestedvalue",
					},
					"template": "Hello World",
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := Ctx{}
			ctx.Request.Payload = tc.Payload

			task, err := event.BuildTask(ctx)
			if err != nil {
				t.Errorf("Error not exptected: %v", err)
			}
			if task.Trigger != tc.Expected.Trigger {
				t.Errorf("got %s, want %s", task.Trigger, tc.Expected.Trigger)
			}
			if task.Values["simple"] != tc.Expected.Values["simple"] {
				t.Errorf("got %s, want %s", task.Values["simple"], tc.Expected.Values["simple"])
			}
			if task.Values["template"] != tc.Expected.Values["template"] {
				t.Errorf("got %s, want %s", task.Values["template"], tc.Expected.Values["template"])
			}
			if task.Values["template"] != tc.Expected.Values["template"] {
				t.Errorf("got %s, want %s", task.Values["template"], tc.Expected.Values["template"])
			}
			nestedValue := task.Values["nested"].(map[string]string)["nestedkey"]
			expectedNestedValue := tc.Expected.Values["nested"].(map[string]string)["nestedkey"]
			if nestedValue != expectedNestedValue {
				t.Errorf("got %s, want %s", nestedValue, expectedNestedValue)
			}
		})
	}
}
