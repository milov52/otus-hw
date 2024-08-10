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

type Validator struct {
	IntValidators    *IntValidators
	StringValidators *StringValidators
	ByteValidators   *ByteValidators
}

type ValidateField struct {
	Tag       string
	Name      string
	Validator Validator
}

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

func validateSliceField(f reflect.StructField, val reflect.Value, fieldName string, vErr ValidationErrors) (ValidationErrors, error) {
	elemKind := f.Type.Elem().Kind()

	switch elemKind {
	case reflect.String:
		stringValidator, err := parseStringValidators(f.Tag.Get("validate"))
		if err != nil {
			return nil, err
		}

		for i := 0; i < val.Len(); i++ {
			fieldName := fmt.Sprintf("%s[%d]", fieldName, i)
			vErr = validateString(*stringValidator, val.Index(i).String(), fieldName, vErr)
		}
	case reflect.Int:
		intValidator, err := parseIntValidators(f.Tag.Get("validate"))
		if err != nil {
			return nil, err
		}

		for i := 0; i < val.Len(); i++ {
			fieldName := fmt.Sprintf("%s[%d]", fieldName, i)
			vErr = validateInt(*intValidator, int(val.Index(i).Int()), fieldName, vErr)
		}

	case reflect.Uint8:
		byteValidator, err := parseByteValidators(f.Tag.Get("validate"))
		if err != nil {
			return nil, err
		}
		return validateByte(*byteValidator, val, fieldName, vErr), nil

	default:
		vErr = append(vErr, ValidationError{
			Field: "Slice",
			Err:   fmt.Errorf("not supported slice type"),
		})
	}

	return vErr, nil
}

func validateField(f reflect.StructField, val reflect.Value, vErr ValidationErrors) (ValidationErrors, error) {
	fieldName := f.Name

	switch val.Kind() {
	case reflect.String:
		stringValidator, err := parseStringValidators(f.Tag.Get("validate"))
		if err != nil {
			return nil, err
		}
		return validateString(*stringValidator, val.String(), fieldName, vErr), nil
	case reflect.Int:
		intValidator, err := parseIntValidators(f.Tag.Get("validate"))
		if err != nil {
			return nil, err
		}
		return validateInt(*intValidator, int(val.Int()), fieldName, vErr), nil
	case reflect.Slice:
		return validateSliceField(f, val, fieldName, vErr)
	default:
		vErr = append(vErr, ValidationError{
			Field: fieldName,
			Err:   fmt.Errorf("not supported type"),
		})
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
