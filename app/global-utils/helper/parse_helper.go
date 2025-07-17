package helper

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertStringArrayToInt32(strings []string) ([]int32, error) {
	var int32Array []int32

	for _, str := range strings {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			fmt.Printf("Error converting string to int32: %v\n", err)
			continue
		}
		int32Num := int32(num)
		int32Array = append(int32Array, int32Num)
	}

	return int32Array, nil
}

func ConvertStringArrayToInt(strings []string) ([]int, error) {
	var intArray []int

	for _, str := range strings {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			fmt.Printf("Error converting string to int: %v\n", err)
			continue
		}
		intNum := int(num)
		intArray = append(intArray, intNum)
	}

	return intArray, nil
}

func ConvertIntArrayToInt32Array(arr []int) []int32 {
	result := make([]int32, len(arr))
	for i, v := range arr {
		result[i] = int32(v)
	}
	return result
}

func ConvertInt32ArrayToIntArray(arr []int32) []int {
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i] = int(v)
	}
	return result
}

func IntArrayToString(arr []int) string {
	strArr := make([]string, len(arr))
	for i, v := range arr {
		strArr[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strArr, ", ")
}

func Int32ArrayToString(arr []int32) string {
	strArr := make([]string, len(arr))
	for i, v := range arr {
		strArr[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strArr, ", ")
}

func StringToInterfaces(arr ...string) []interface{} {
	arz := make([]interface{}, 0)
	for _, ar := range arr {
		arz = append(arz, ar)
	}

	return arz
}

func ToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
