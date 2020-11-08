package config

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"time"
)

func ProcessTemplate(tpl string, data interface{}) (string, error) {
	if !strings.Contains(tpl, "{{") {
		return tpl, nil
	}
	if strings.HasPrefix(tpl, "'") && strings.HasSuffix(tpl, "'") {
		tpl = tpl[1 : len(tpl)-1]
	}
	tpl = strings.ReplaceAll(tpl, "${{", "{{")
	var result bytes.Buffer
	t, err := template.New("tmp").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("Cannot parse template: %v. Error: %v", tpl, err)
	}
	err = t.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("Cannot process template: %v. Error: %v", tpl, err)
	}
	return result.String(), nil
}

func ProcessAllTemplates(templates, input interface{}) (interface{}, error) {
	ctx := context{
		input: input,
	}
	result := ctx.copy(templates)
	if ctx.err != nil {
		return nil, ctx.err
	}
	return result, nil
}

type context struct {
	input interface{}
	err   error
}

// Copy creates a deep copy of whatever is passed to it and returns the copy
// in an interface{}.  The returned value will need to be asserted to the
// correct type.
func (ctx *context) copy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	// Make the interface a reflect.Value
	original := reflect.ValueOf(src)

	// Make a copy of the same type as the original.
	cpy := reflect.New(original.Type()).Elem()

	// Recursively copy the original.
	ctx.copyRecursive(original, cpy)

	// Return the copy as an interface.
	return cpy.Interface()
}

// copyRecursive does the actual copying of the interface. It currently has
// limited support for what it can handle. Add as needed.
func (ctx *context) copyRecursive(original, cpy reflect.Value) {

	// handle according to original's Kind
	switch original.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		originalValue := original.Elem()

		// if  it isn't valid, return.
		if !originalValue.IsValid() {
			return
		}
		cpy.Set(reflect.New(originalValue.Type()))
		ctx.copyRecursive(originalValue, cpy.Elem())

	case reflect.Interface:
		// If this is a nil, don't do anything
		if original.IsNil() {
			return
		}
		// Get the value for the interface, not the pointer.
		originalValue := original.Elem()

		// Get the value by calling Elem().
		copyValue := reflect.New(originalValue.Type()).Elem()
		ctx.copyRecursive(originalValue, copyValue)
		cpy.Set(copyValue)

	case reflect.Struct:
		t, ok := original.Interface().(time.Time)
		if ok {
			cpy.Set(reflect.ValueOf(t))
			return
		}
		// Go through each field of the struct and copy it.
		for i := 0; i < original.NumField(); i++ {
			// The Type's StructField for a given field is checked to see if StructField.PkgPath
			// is set to determine if the field is exported or not because CanSet() returns false
			// for settable fields.  I'm not sure why.  -mohae
			if original.Type().Field(i).PkgPath != "" {
				continue
			}
			ctx.copyRecursive(original.Field(i), cpy.Field(i))
		}

	case reflect.Slice:
		if original.IsNil() {
			return
		}
		// Make a new slice and copy each element.
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			ctx.copyRecursive(original.Index(i), cpy.Index(i))
		}

	case reflect.Map:
		if original.IsNil() {
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			ctx.copyRecursive(originalValue, copyValue)
			copyKey := ctx.copy(key.Interface())
			cpy.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}

	default:
		if original.Kind() == reflect.String {
			n, err := ProcessTemplate(original.String(), ctx.input)
			if err != nil {
				ctx.err = err
				n = original.String()
			}
			newValue := reflect.ValueOf(n)
			cpy.Set(newValue)
		} else {
			cpy.Set(original)
		}
	}
}
