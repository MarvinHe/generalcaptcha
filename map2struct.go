package generalcaptcha

import (
	"reflect"
	"fmt"
	"errors"
)

func SetField(obj interface{}, name string, value interface{}) error {
	// obj must be a pointer
	structValue := reflect.ValueOf(obj).Elem()
	err := setStructField(structValue, name, value)
	return err
}

func SetFields(obj interface{}, data map[string]interface{}) error {
	// obj must be a pointer
	structValue := reflect.ValueOf(obj).Elem()
	for name, value := range data {
		err := setStructField(structValue, name, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func setStructField(structValue reflect.Value, name string, value interface{}) error {
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func EnsureMapString(data map[interface{}]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	var stringKey string
	var ok bool
	for key, value := range data {
		if stringKey, ok = key.(string); !ok {
			fmt.Printf("the key is: %s\n", key)
		}
		result[stringKey] = value
	}
	return result
}
