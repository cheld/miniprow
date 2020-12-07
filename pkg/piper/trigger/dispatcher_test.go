package trigger

import (
	"testing"

	"github.com/cheld/miniprow/pkg/piper/config"
)

func TestExecute_errorhandling(t *testing.T) {
	cfg := config.Configuration{}
	cfg.Triggers = []config.Trigger{
		{
			Name: "some-trigger-definition",
			Type: "not-implemented",
		},
		{
			Name: "debug-trigger",
			Type: "debug",
		},
	}
	dispatcher := NewDispatcher(cfg)
	type testcase struct {
		name    string
		trigger string
		err     bool
	}
	testcases := []testcase{
		{
			name:    "Task correlates to trigger definition, but no implentation exists",
			trigger: "some-trigger-definition",
			err:     true,
		},
		{
			name:    "No trigger matches this task",
			trigger: "no-trigger",
			err:     true,
		},
		{
			name:    "Trigger found",
			trigger: "debug-trigger",
			err:     false,
		},
		{
			name:    "Trigger with different case writings",
			trigger: "DeBug-Trigger",
			err:     false,
		},
		{
			name: "Missing trigger name",
			err:  true,
		},
	}

	// Test error
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tasks := []config.Task{
				config.Task{
					Trigger: tc.trigger,
				},
			}
			dispatcher.Execute(tasks)

		})
	}
}
