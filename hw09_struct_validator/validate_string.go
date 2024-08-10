package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func validateString(f reflect.StructField,
	val reflect.Value,
	vErr ValidationErrors,
	fieldName string) (ValidationErrors, error) {
	validators := strings.Split(f.Tag.Get("validate"), "|")
	for _, item := range validators {
		validator := strings.Split(item, ":")

		switch validator[0] {
		case LEN:
			valLen, pErr := strconv.Atoi(validator[1])
			if pErr != nil {
				return vErr, pErr // Возвращаем системную ошибку
			}

			if len(val.String()) != valLen {
				vErr = append(vErr, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("expected string len %d, got %d", valLen, len(val.String())),
				})
			}
		case IN:
			if !strings.Contains(validator[1], val.String()) {
				vErr = append(vErr, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value %s does not contains into %v", val.String(), validator[1]),
				})
			}
		case REGEXP:
			re, pErr := regexp.Compile(validator[1])
			if pErr != nil {
				return vErr, pErr
			}
			if !re.MatchString(val.String()) {
				vErr = append(vErr, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("value %s does not match regexp %s", val.String(), validator[1]),
				})
			}
		}
	}
	return vErr, nil // Возвращаем ошибки валидации и nil для системных ошибок
}
