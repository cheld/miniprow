package config

import (
	"testing"
)

func TestDestinationCtx(t *testing.T) {

	event := EventCtx{}
	event.Objectiv = ""
	event.Input = map[string]interface{}{
		"inputkey1": "inputvalue1",
		"inputkey2": "inputvalue2",
	}
	t.Run("mytest", func(t *testing.T) {

		rule := Rule{
			Destination: "some-destination",
			Values: map[string]interface{}{
				"target": "test",
				"nested": map[string]string{
					"nestedkey": "nestedvalue",
				},
				"template": "String with {{ .Input.inputkey1 }}",
			},
		}
		trigger := rule.DestinationCtx(event)
		if trigger.Name != "some-destination" {
			t.Errorf("got %s, want %s", trigger.Name, "some-destination")
		}
		if trigger.Values["target"] != "test" {
			t.Errorf("got %s, want %s", trigger.Values["target"], "test")
		}
		if (trigger.Values["nested"]).(map[string]string)["nestedkey"] != "nestedvalue" {
			t.Errorf("got %s, want %s", "", "nestedvalue")
		}
		if trigger.Values["template"] != "String with inputvalue1" {
			t.Errorf("got %s, want %s", trigger.Values["template"], "String with inputvalue1")
		}
	})
}
