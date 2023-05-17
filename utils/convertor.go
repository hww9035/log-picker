package utils

import (
	"fmt"
	"reflect"
	"regexp"
)

// StructToMap 结构体转map
func StructToMap(value any) (map[string]any, error) {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data type %T not support, shuld be struct or pointer to struct", value)
	}
	result := make(map[string]any)
	fieldNum := t.NumField()
	pattern := `^[A-Z]`
	regex := regexp.MustCompile(pattern)
	for i := 0; i < fieldNum; i++ {
		name := t.Field(i).Name
		tag := t.Field(i).Tag.Get("json")
		if regex.MatchString(name) && tag != "" {
			if v.Kind() == reflect.Ptr { // 指针类型
				result[tag] = v.Elem().Field(i).Interface()
			} else {
				result[tag] = v.Field(i).Interface()
			}
		}
	}

	return result, nil
}
