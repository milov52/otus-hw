package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

const (
	LEN    = "len"
	REGEXP = "regexp"
	IN     = "in"
	MIN    = "min"
	MAX    = "max"
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(err.Field)
		sb.WriteString(": ")
		sb.WriteString(err.Err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func validateField(f reflect.StructField, val reflect.Value, vErr ValidationErrors) (ValidationErrors, error) {
	fieldName := f.Name
	switch f.Type.Kind() {
	case reflect.String:
		return validateString(f, val, vErr, fieldName)
	case reflect.Int:
		return validateInt(f, val, vErr, fieldName)
	case reflect.Slice:
		elemKind := f.Type.Elem().Kind()
		for i := 0; i < val.Len(); i++ {
			elemVal := val.Index(i)
			switch elemKind {
			case reflect.String:
				fieldName = fmt.Sprintf("%s[%d]", fieldName, i)
				vErr, _ = validateString(f, elemVal, vErr, fieldName)
			case reflect.Int:
				fieldName = fmt.Sprintf("%s[%d]", fieldName, i)
				vErr, _ = validateInt(f, elemVal, vErr, fieldName)
			case reflect.Uint8:
				return validateByte(f, val, vErr, fieldName)
			default:
				vErr = append(vErr, ValidationError{
					Field: "Slice",
					Err:   fmt.Errorf("not supported slice type"),
				})
			}
		}
	}
	return vErr, nil
}

func Validate(v any) error {
	ve := make(ValidationErrors, 0)

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		ve = append(ve, ValidationError{
			Field: rv.Type().Name(),
			Err:   fmt.Errorf("expected struct, got %s", rv.Type().Name()),
		})
		return ve
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		var sysErr error
		ve, sysErr = validateField(field, value, ve)
		if sysErr != nil {
			return sysErr
		}
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}
