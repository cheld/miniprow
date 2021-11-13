package config

import "testing"

func TestValidate(t *testing.T) {
	type testcase struct {
		name     string
		rules    []Rule
		triggers []Trigger
		err      bool
	}
	testcases := []testcase{
		{
			name:     "No rules",
			rules:    []Rule{},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "rule: no event",
			rules:    []Rule{{Event: "myevent"}},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "rule: no rules",
			rules:    []Rule{{Event: "mysource"}},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "rule: no trigger",
			rules:    []Rule{{Event: "mysource", If_contains: "test"}},
			triggers: []Trigger{},
			err:      true,
		},
		/*{
			name:     "rule: reference to trigger does not exist",
			rules:    []Rule{{Event: "mysource", If_contains: "test", Trigger: "mytrigger"}},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "Trigger: no type",
			rules:    []Rule{{Event: "mysource", If_contains: "test", Trigger: "mytrigger"}},
			triggers: []Trigger{{Name: "mytrigger"}},
			err:      true,
		},
		{
			name:     "Trigger: wrong type",
			rules:    []Rule{{Event: "mysource", If_contains: "test", Trigger: "mytrigger"}},
			triggers: []Trigger{{Name: "mytrigger", Type: "mytype"}},
			err:      true,
		},
		{
			name:     "Trigger: no name",
			rules:    []Rule{{Event: "mysource", If_contains: "test", Trigger: "mytrigger"}},
			triggers: []Trigger{{Name: "mytrigger", Type: "http"}, {}},
			err:      true,
		},
		{
			name:     "Ingnore case",
			rules:    []Rule{{Event: "mysource", If_contains: "test", Trigger: "myTrigger"}},
			triggers: []Trigger{{Name: "mytrigger", Type: "hTTp"}},
			err:      false,
		},
		*/}
	cfg := Configuration{}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cfg.Triggers = tc.triggers
			cfg.Rules = tc.rules
			err := Validate(cfg)
			if err == nil && tc.err {
				t.Errorf("Error expected but not received")
			}
			if err != nil && !tc.err {
				t.Errorf("Error not exptected: %v", err)
			}
		})
	}
}
