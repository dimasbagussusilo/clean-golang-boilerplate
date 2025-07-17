package helper

import (
	"fmt"
	"reflect"

	"appsku-golang/app/global-utils/model"
)

func GenerateMustActive(table string, reqField string, id int, status string) *model.MustActiveRequest {
	return &model.MustActiveRequest{
		Table:    table,
		ReqField: reqField,
		Clause:   fmt.Sprintf("id = %d AND status = '%s'", id, status),
		Id:       id,
	}
}

func IntContainArray(array []int64, target int64) bool {
	for _, value := range array {
		if value == target {
			return true
		}
	}
	return false
}

func StringContainArray(array []string, target string) bool {
	for _, value := range array {
		if value == target {
			return true
		}
	}
	return false
}

func IsArrayNilValues(arr ...interface{}) bool {
	arrLength := len(arr)
	nilLength := 0

	for _, ar := range arr {
		v := reflect.ValueOf(ar)

		switch v.Kind() {
		case reflect.Map, reflect.Array, reflect.String:
			if v.Len() < 1 {
				nilLength += 1
			}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Float32, reflect.Float64:
			if v.IsZero() {
				nilLength += 1
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v.IsZero() {
				nilLength += 1
			}
		}
	}

	return nilLength >= arrLength
}
