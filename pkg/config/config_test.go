package config

import (
	"testing"
)

func TestDestinationCtx(t *testing.T) {

	inputData := InputData{}
	inputData.Objectiv = ""
	inputData.Input = map[string]interface{}{
		"inputkey1": "inputvalue1",
		"inputkey2": "inputvalue2",
	}
	t.Run("mytest", func(t *testing.T) {

		eventSpec := EventSpec{
			Trigger: "some-destination",
			Values: map[string]interface{}{
				"target": "test",
				"nested": map[string]string{
					"nestedkey": "nestedvalue",
				},
				"template": "String with {{ .Input.inputkey1 }}",
			},
		}
		eventData := eventSpec.Process(inputData)
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
