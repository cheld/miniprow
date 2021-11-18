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
			ifTrue:   "${{ eq .Request.Value \"test\" }}",
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
			request := RequestCtx{
				Value: tc.input,
			}
			ctx := &Ctx{
				Request: request,
			}
			result := rule.IsMatching(ctx)
			if result != tc.expected {
				t.Fatalf("Error expected %v, got %v.", tc.expected, result)
			}
		})
	}
}

func TestGetTrigger(t *testing.T) {
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
			trigger := cfg.GetTrigger(tc.Trigger)
			if trigger == nil && tc.Found {
				t.Fatalf("Trigger not found")
			}
			if trigger != nil && !tc.Found {
				t.Fatalf("Trigger found, but not expected")
			}
		})
	}
}

func TestGetRule(t *testing.T) {
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
			request := RequestCtx{
				Event: tc.Event,
				Value: tc.Value,
			}
			ctx := &Ctx{Request: request}
			rule := cfg.GetFirstMatchingRule(ctx)
			if rule == nil && tc.Found {
				t.Fatalf("Event not found")
			}
			if rule != nil && !tc.Found {
				t.Fatalf("Trigger found, but not expected")
			}
		})
	}
}

func TestApply(t *testing.T) {
	rule := Rule{
		Then: []Then{{
			Apply: "mytrigger",
			With: map[string]string{
				"param":   "value",
				"template": "Hello {{ .Request.Payload.sourcekey }}",
			},
		},
		},
	}
	type testcase struct {
		Name     string
		Payload  map[string]string
		Expected TriggerCtx
	}
	testcases := []testcase{
		{
			Name: "Simple",
			Payload: map[string]string{
				"sourcekey": "World",
			},
			Expected: TriggerCtx{
				Input: map[string]string{
					"param":   "value",
					"template": "Hello World",
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := &Ctx{}
			ctx.Request.Payload = tc.Payload
			err := rule.Apply(ctx)
			if err != nil {
				t.Errorf("Error not exptected: %v", err)
			}
			if ctx.Trigger["mytrigger"].Input["simple"] != tc.Expected.Input["simple"] {
				t.Errorf("got %s, want %s", ctx.Trigger["mytrigger"].Input["simple"], tc.Expected.Input["simple"])
			}
			if ctx.Trigger["mytrigger"].Input["template"] != tc.Expected.Input["template"] {
				t.Errorf("got %s, want %s", ctx.Trigger["mytrigger"].Input["template"], tc.Expected.Input["template"])
			}

		})
	}
}
