package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type StringValidators struct {
	Len    *int
	In     []string
	RegExp *regexp.Regexp
}

func NewStringValidators() *StringValidators {
	return &StringValidators{
		In: make([]string, 0),
	}
}

func parseStringValidators(tag string) (*StringValidators, error) {
	strValidators := NewStringValidators()
	for _, item := range strings.Split(tag, "|") {
		validator := strings.Split(item, ":")
		switch validator[0] {
		case LEN:
			valLen, pErr := strconv.Atoi(validator[1])
			if pErr != nil {
				return nil, pErr
			}
			strValidators.Len = &valLen
		case REGEXP:
			re, pErr := regexp.Compile(validator[1])
			if pErr != nil {
				return nil, pErr
			}
			strValidators.RegExp = re
		case IN:
			strValidators.In = strings.Split(validator[1], ",")
		}
	}
	return strValidators, nil
}

func validateString(sVal StringValidators, val string, fieldName string, vErr ValidationErrors) ValidationErrors {
	if sVal.Len != nil && len(val) != *sVal.Len {
		vErr = append(vErr, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("expected string len %d, got %d", *sVal.Len, len(val)),
		})
	}

	if sVal.RegExp != nil && !sVal.RegExp.MatchString(val) {
		vErr = append(vErr, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("value %s does not match regexp %s", val, sVal.RegExp.String()),
		})
	}

	if len(sVal.In) > 0 && !is_contains(sVal.In, val) {
		vErr = append(vErr, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("value %s does not contains into %v", val, sVal.In),
		})
	}

	return vErr
}
