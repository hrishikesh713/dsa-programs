package utils

import (
	"fmt"
	"reflect"
)

func IsValueNil(v any) (bool, error) {
	if v == nil {
		return false, fmt.Errorf("value passed is a nil interface value")
	}
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Pointer {
		return false, fmt.Errorf("this function only checks the validity of a pointer value stored in an interface")
	}
	return vv.IsNil(), nil
}
