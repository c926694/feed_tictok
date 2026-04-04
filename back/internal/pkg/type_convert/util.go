package type_convert

import (
	"errors"
)

const (
	err = "类型转换异常"
)

func AnyToUint64(number any) (uint64, error) {
	if num, ok := number.(uint64); ok {
		return num, nil
	}
	return 0, errors.New(err)
}

func AnyToFloat64(number any) (float64, error) {
	if num, ok := number.(float64); ok {
		return num, nil
	}
	return 0, errors.New(err)
}
