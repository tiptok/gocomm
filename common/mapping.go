package common

import (
	"reflect"
	"strings"
)

func CamelCase(name string, firstUpper bool) string {
	array := []byte(name)
	if len(array) == 0 {
		return ""
	}
	rspArray := make([]byte, len(array))
	if firstUpper {
		copy(rspArray[:1], strings.ToUpper(string(array[:1])))
	} else {
		copy(rspArray[:1], strings.ToLower(string(array[:1])))
	}
	copy(rspArray[1:], array[1:])
	return string(rspArray)
}

func ObjectToMap(o interface{}) map[string]interface{} {
	if o == nil {
		return nil
	}
	value := reflect.ValueOf(o)
	if value.Kind() != reflect.Ptr {
		return nil
	}
	elem := value.Elem()
	relType := elem.Type()
	m := make(map[string]interface{})
	for i := 0; i < relType.NumField(); i++ {
		field := relType.Field(i)
		if elem.Field(i).IsZero() {
			continue
		}
		m[CamelCase(field.Name, false)] = elem.Field(i).Interface()
	}
	return m
}
