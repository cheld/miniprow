package util

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var Environment = readEnvironment()

type Env struct {
	env map[string]string
}

func (e *Env) Value(k string) *Value {
	return &Value{
		env: e,
		key: k,
	}
}

type Value struct {
	env *Env
	key string
}

func readEnvironment() *Env {
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		keyValue := strings.Split(entry, "=")
		env[keyValue[0]] = keyValue[1]
	}
	return &Env{env: env}
}

func (v *Value) String() string {
	return v.env.env[v.key]
}

func (v *Value) Int() (int, error) {
	value := v.env.env[v.key]
	if value == "" {
		return -1, fmt.Errorf("No configuration found for key %s", v.key)
	}
	return strconv.Atoi(value)
}

func (v *Value) Base64() ([]byte, error) {
	value := v.env.env[v.key]
	if value == "" {
		return nil, fmt.Errorf("No configuration found for key %s", v.key)
	}
	decoded, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(value)
	if err == nil {
		return decoded, err
	}
	return base64.StdEncoding.DecodeString(value)
}

func (v *Value) Update(i interface{}) {
	value := v.env.env[v.key]
	original := reflect.ValueOf(i).Elem()
	switch original.Kind() {
	case reflect.Int:
		valueAsInt, err := strconv.Atoi(value)
		if err == nil {
			original.SetInt(int64(valueAsInt))
		}
	case reflect.String:
		if value != "" {
			original.SetString(value)
		}
	default:
		fmt.Println("Error")
	}
}

func (env *Env) Map() *map[string]string {
	return &env.env
}
