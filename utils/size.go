package utils

import "reflect"

func GetSize[T any](value T) int {
	return int(reflect.TypeOf(value).Size())
}
