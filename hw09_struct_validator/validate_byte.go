package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func validateByte(f reflect.StructField, val reflect.Value, vErr ValidationErrors, fieldName string) (ValidationErrors, error) {
	validators := strings.Split(f.Tag.Get("validate"), "|")
	for _, item := range validators {
		validator := strings.Split(item, ":")

		if validator[0] != LEN {
			return vErr, fmt.Errorf("unknown validator %s", validator[0])
		}

		valLen, pErr := strconv.Atoi(validator[1])
		if pErr != nil {
			return vErr, pErr
		}

		if val.Len() != valLen {
			vErr = append(vErr, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("expected length %d, got %d", valLen, val.Len()),
			})
		}
	}
	return vErr, nil
}
