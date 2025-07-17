package helper

import (
	"errors"
	"fmt"
	"reflect"
)

func ValidateRequiredFields(data interface{}) error {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	if v.Kind() != reflect.Struct {
		return errors.New("expected a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check if the field has the `required:"true"` tag
		if required, ok := field.Tag.Lookup("required"); ok && required == "true" {
			// Check if the field is zero (empty/unset)
			if isZeroValue(value) {
				return fmt.Errorf("field '%s' is required but empty", field.Name)
			}
		}
	}

	return nil
}

func isZeroValue(v reflect.Value) bool {
	// Handle pointer values by checking what they point to
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	// Use the zero value of the type to compare
	zero := reflect.Zero(v.Type()).Interface()
	return reflect.DeepEqual(v.Interface(), zero)
}
