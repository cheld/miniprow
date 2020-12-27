package util

import (
	"testing"
)

func TestGetValue(t *testing.T) {
	testData := map[string]string{
		"a_string": "myvalue",
		"a_int":    "5",
	}
	Environment = &Env{
		env: testData,
	}

	// string
	str := Environment.Value("a_string").String()
	if str != "myvalue" {
		t.Errorf("Expected myvalue, but received %s", str)
	}
	str = Environment.Value("non_existing").String()
	if str != "" {
		t.Errorf("Expected myvalue, but received %s", str)
	}

	// int
	i, err := Environment.Value("a_int").Int()
	if i != 5 && err == nil {
		t.Errorf("Expected myvalue, but received %d", i)
	}
	i, err = Environment.Value("non_existing").Int()
	if i != -1 && err == nil {
		t.Errorf("Expected myvalue, but received %d", i)
	}

}

func TestUpdate(t *testing.T) {
	testData := map[string]string{
		"a_string": "myvalue",
		"a_int":    "5",
	}
	Environment = &Env{
		env: testData,
	}

	//int
	i := 1
	Environment.Value("a_int").Update(&i)
	if i != 5 {
		t.Error("int was not updated")
	}
	i = 1
	Environment.Value("non_existing").Update(&i)
	if i != 1 {
		t.Errorf("int was updated to %d, but should not", i)
	}
	i = 1
	Environment.Value("a_string").Update(&i)
	if i != 1 {
		t.Errorf("int was updated to %d, but should not", i)
	}

	//string
	s := "init"
	Environment.Value("a_string").Update(&s)
	if s != "myvalue" {
		t.Error("int was not updated")
	}
	s = "init"
	Environment.Value("non_existing").Update(&s)
	if s != "init" {
		t.Errorf("string was updated to %s, but should not", s)
	}
}
