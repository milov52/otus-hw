package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"
)

type IntValidators struct {
	Min *int
	Max *int
	In  []int
}

func NewIntValidators() *IntValidators {
	return &IntValidators{
		In: make([]int, 0),
	}
}

func parseIntValidators(tag string) (*IntValidators, error) {
	intValidator := NewIntValidators()
	for _, item := range strings.Split(tag, "|") {
		validator := strings.Split(item, ":")
		switch validator[0] {
		case MIN:
			valMin, pErr := strconv.Atoi(validator[1])
			if pErr != nil {
				return nil, pErr
			}
			intValidator.Min = &valMin
		case MAX:
			valMax, pErr := strconv.Atoi(validator[1])
			if pErr != nil {
				return nil, pErr
			}
			intValidator.Max = &valMax
		case IN:
			parts := strings.Split(validator[1], ",")
			for _, part := range parts {
				valIn, pErr := strconv.Atoi(part)
				if pErr != nil {
					return nil, pErr
				}
				intValidator.In = append(intValidator.In, valIn)
			}
		}
	}
	return intValidator, nil
}

func validateInt(iVal IntValidators, val int, fieldName string, vErr ValidationErrors) ValidationErrors {
	if iVal.Min != nil && val < *iVal.Min {
		vErr = append(vErr, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("value %d is less than min %d", val, *iVal.Min),
		})
	}

	if iVal.Max != nil && val > *iVal.Max {
		vErr = append(vErr, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("value %d is more than min %d", val, *iVal.Max),
		})
	}

	if len(iVal.In) > 0 && !is_contains(iVal.In, val) {
		vErr = append(vErr, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("value %d does not contains into %v", val, iVal.In),
		})
	}
	return vErr
}
