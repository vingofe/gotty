package gotty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


type Object struct {
	Value interface{}
}

func TestGetNonStruct(t *testing.T) {
	_, err := Get(nil, "Value")
	assert.NotNil(t, err)
	_, err = Get("foo", "Value")
	assert.NotNil(t, err)
	_, err = Get(65535, "Value")
	assert.NotNil(t, err)
}


func TestGetString(t *testing.T) {
	expected := "foo"
	actual, err := Get(Object{
		Value: expected,
	}, "Value")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}


func TestGetBool(t *testing.T) {
	expected := true
	actual, err := Get(Object{
		Value: expected,
	}, "Value")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}


func TestGetInt(t *testing.T) {
	expected := 1
	actual, err := Get(Object{
		Value: expected,
	}, "Value")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}


func TestGetCustomType(t *testing.T) {
	type CustomType string
	var expected CustomType = "foo"
	actual, err := Get(Object{
		Value: expected,
	}, "Value")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}


func TestGetNil(t *testing.T) {
	actual, err := Get(Object{
		Value: nil,
	}, "Value")
	assert.Nil(t, err)
	assert.Nil(t, actual)
}


func TestGetNotExist(t *testing.T) {
	expected := "foo"
	object := Object{
		Value: expected,
	}
	_, err := Get(object, "Value1")
	assert.NotNil(t, err)


	_, err = Get(object, "Value.AnotherValue")
	assert.NotNil(t, err)
}


func TestGetNested(t *testing.T) {
	expected := "foo"
	type A struct {
		AValue string
	}
	type B struct {
		BValue    A
		BValuePtr *A
		BValueNil *A
	}
	type C struct {
		CValue B
	}
	object1 := A{
		AValue: expected,
	}
	object2 := B{
		BValue:    object1,
		BValuePtr: &object1,
	}
	object3 := C{
		CValue: object2,
	}


	actual, err := Get(object3, "CValue.BValue.AValue")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)


	actual, err = Get(object3, "CValue.BValueNil")
	assert.Nil(t, err)
	assert.True(t, actual == nil)


	actual, err = Get(object3, "CValue.BValuePtr.AValue")
	assert.Nil(t, err)
	assert.True(t, actual == expected)


	actual, err = Get(object3, "CValue.BValue")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual.(A).AValue)
}


func TestGetNestedPtr(t *testing.T) {
	expected := "foo"
	type A struct {
		AValue  *string
		AValue2 *int
	}
	type B struct {
		BValue *A
	}
	type C struct {
		CValue *B
	}
	object1 := A{
		AValue: &expected,
	}
	object2 := B{
		BValue: &object1,
	}
	object3 := C{
		CValue: &object2,
	}


	actual, err := Get(object3, "CValue.BValue.AValue")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)


	object1.AValue = nil
	actual, err = Get(object3, "CValue.BValue.AValue")
	assert.Nil(t, err)
	assert.Nil(t, actual)


	actual, err = Get(object3, "CValue.BValue.AValue2")
	assert.Nil(t, err)
	assert.Nil(t, actual)


	actual, err = Get(object3, "CValue..BValue.AValue2")
	assert.NotNil(t, err)
}


func TestGetArraySlice(t *testing.T) {
	expected := "foo"
	var arr [1]string
	arr[0] = expected


	actual, err := Get(arr, "0")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)


	actual, err = Get(arr, "1")
	assert.NotNil(t, err)
	assert.Nil(t, actual)


	arr1 := []*string{&expected, nil}
	actual, err = Get(arr1, "0")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)


	actual, err = Get(arr1, "1")
	assert.Nil(t, err)
	assert.Nil(t, actual)


	slice1 := arr1[0:1]
	actual, err = Get(slice1, "0")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)


	inner := Object{
		Value: expected,
	}
	arr2 := []Object{inner}
	wrapper := Object{
		Value: arr2,
	}


	actual, err = Get(wrapper, "Value.0.Value")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}


func TestGetValueByKeyMap(t *testing.T) {
	expected := "foo"
	o := Object{
		Value: expected,
	}
	m := map[string]Object {
		"key1": o,
	}


	actual, err := Get(m, "key1.Value")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)


	m2 := map[int]string {
		1: expected,
	}


	actual, err = Get(m2, "1")
	assert.NotNil(t, err)
	assert.Nil(t, actual)
}


func TestGetValueByKeyInterface(t *testing.T) {
	var expected interface{}
	o := &Object{
		Value: expected,
	}
	o1 := Object{
		Value: o.Value,
	}


	actual, err := Get(o1, "Value.Value")
	assert.NotNil(t, err)
	assert.Nil(t, actual)
}
