package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func validateInt(f reflect.StructField, val reflect.Value, vErr ValidationErrors, fieldName string) (ValidationErrors, error) {
	validators := strings.Split(f.Tag.Get("validate"), "|")
	for _, item := range validators {
		validator := strings.Split(item, ":")
		switch validator[0] {
		case
			MIN:
			valMin, pErr := strconv.Atoi(validator[1])
			if pErr != nil {
				return vErr, pErr
			}
			if val.Int() < int64(valMin) {
				vErr = append(vErr, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value %d is less than min %d", val.Int(), valMin),
				})
			}
		case MAX:
			valMax, pErr := strconv.Atoi(validator[1])
			if pErr != nil {
				return vErr, pErr
			}
			if val.Int() > int64(valMax) {
				vErr = append(vErr, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value %d is more than min %d", val.Int(), valMax),
				})
			}
		case IN:
			if !strings.Contains(validator[1], strconv.FormatInt(val.Int(), 10)) {
				vErr = append(vErr, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value %d does not contains into %v", val.Int(), validator[1]),
				})
			}
		default:
			return vErr, fmt.Errorf("unknown validator %s", validator[0])
		}

	}
	return vErr, nil
}
