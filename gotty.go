package gotty

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)


func IsNilValue(value reflect.Value) bool {
	kind := value.Kind()
	isNilable := kind == reflect.Chan || kind == reflect.Func ||
		kind == reflect.Interface || kind == reflect.Map ||
		kind == reflect.Ptr || kind == reflect.Slice


	if isNilable && value.IsNil() {
		return true
	}
	return false
}


func FieldByNameForMapValue(value reflect.Value, key string) reflect.Value {
	if value.Kind() == reflect.Map {
		switch value.Type().Key().Name() {
		case "string":
			return value.MapIndex(reflect.ValueOf(key))
		}
	}
	return reflect.Value{}
}


func FieldByIndexForArrayValue(value reflect.Value, key string) reflect.Value {
	if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
		i, err := strconv.ParseInt(key, 10, 32)
		if err == nil && i >= 0 && int(i) < value.Len() {
			return value.Index(int(i))
		}
	}
	return reflect.Value{}
}


func GetValueByKey(value reflect.Value, key string) reflect.Value {
	//fmt.Println(value.Kind())
	switch value.Kind() {
	case reflect.Ptr:
		v := value.Elem()
		return GetValueByKey(v, key)
	case reflect.Interface:
		v := reflect.ValueOf(value.Interface())
		return GetValueByKey(v, key)
	case reflect.Map:
		return FieldByNameForMapValue(value, key)
	case reflect.Array:
		return FieldByIndexForArrayValue(value, key)
	case reflect.Slice:
		return FieldByIndexForArrayValue(value, key)
	case reflect.Struct:
		return value.FieldByName(key)
	}
	return reflect.Value{}
}


// Get gets the nested value from object by give path str, eg. 'Get(object, "a.b.c")'
func Get(object interface{}, path string) (interface{}, error) {
	keySlice := strings.Split(path, ".")
	value := reflect.ValueOf(object)


	for _, key := range keySlice {
		value = GetValueByKey(value, key)


		if !value.IsValid() {
			return nil, errors.New("invalid path")
		}
	}


	if IsNilValue(value) {
		return nil, nil
	}
	return reflect.Indirect(value).Interface(), nil
}