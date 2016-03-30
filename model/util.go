package model

import "reflect"

func prepareToUpdate(s interface{}) map[string]interface{} {
	v := reflect.ValueOf(s).Elem() // the struct variable
	kvMap := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		name := fieldInfo.Name
		if name == "Id" {
			//skip to update id
			continue
		}
		value := v.FieldByName(name)
		if !isEmptyOrSkipValue(value) {
			kvMap[fieldInfo.Tag.Get("redis")] = value
		}
	}
	return kvMap
}

func isEmptyOrSkipValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		//return v.Len() == 0
		return true
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
