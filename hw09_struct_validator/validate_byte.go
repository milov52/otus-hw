package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ByteValidators struct {
	Len *int
}

func parseByteValidators(tag string) (*ByteValidators, error) {
	bVal := ByteValidators{}

	for _, item := range strings.Split(tag, "|") {
		validator := strings.Split(item, ":")
		if validator[0] == LEN {
			lenVal, pErr := strconv.Atoi(validator[1])
			if pErr != nil {
				return nil, pErr
			}
			bVal.Len = &lenVal
		}
	}
	return &bVal, nil
}

func validateByte(bVal ByteValidators, val reflect.Value, fieldName string, vErr ValidationErrors) ValidationErrors {
	if bVal.Len != nil && val.Len() != *bVal.Len {
		vErr = append(vErr, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("expected length %d, got %d", *bVal.Len, val.Len()),
		})
	}
	return vErr
}
