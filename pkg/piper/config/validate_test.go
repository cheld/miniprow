package config

import "testing"

func TestValidate(t *testing.T) {
	type testcase struct {
		name     string
		events   []Event
		triggers []Trigger
		err      bool
	}
	testcases := []testcase{
		{
			name:     "No events",
			events:   []Event{},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "Event: no source",
			events:   []Event{{}},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "Event: no type",
			events:   []Event{{Source: "mysource"}},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "Event: no rules",
			events:   []Event{{Source: "mysource", Type: "mytype"}},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "Event: no trigger",
			events:   []Event{{Source: "mysource", Type: "mytype", If_contains: "test"}},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "Event: reference to trigger does not exist",
			events:   []Event{{Source: "mysource", Type: "mytype", If_contains: "test", Trigger: "mytrigger"}},
			triggers: []Trigger{},
			err:      true,
		},
		{
			name:     "Trigger: no type",
			events:   []Event{{Source: "mysource", Type: "mytype", If_contains: "test", Trigger: "mytrigger"}},
			triggers: []Trigger{{Name: "mytrigger"}},
			err:      true,
		},
		{
			name:     "Trigger: wrong type",
			events:   []Event{{Source: "mysource", Type: "mytype", If_contains: "test", Trigger: "mytrigger"}},
			triggers: []Trigger{{Name: "mytrigger", Type: "mytype"}},
			err:      true,
		},
		{
			name:     "Trigger: no name",
			events:   []Event{{Source: "mysource", Type: "mytype", If_contains: "test", Trigger: "mytrigger"}},
			triggers: []Trigger{{Name: "mytrigger", Type: "http"}, {}},
			err:      true,
		},
		{
			name:     "Ingnore case",
			events:   []Event{{Source: "mysource", Type: "mytype", If_contains: "test", Trigger: "myTrigger"}},
			triggers: []Trigger{{Name: "mytrigger", Type: "hTTp"}},
			err:      false,
		},
	}
	cfg := Configuration{}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cfg.Triggers = tc.triggers
			cfg.Events = tc.events
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
