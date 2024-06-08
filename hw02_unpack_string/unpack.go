package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func isDigit(ch byte) bool {
	_, err := strconv.Atoi(string(ch))
	return err == nil
}

func Unpack(s string) (string, error) {
	if s == "" {
		return s, nil
	}
	if isDigit(s[0]) {
		return "", ErrInvalidString
	}

	var sb strings.Builder
	escape := false // Флаг для обработки символа '\'

	for i := range s {
		if escape {
			sb.WriteByte(s[i])
			escape = false
			continue
		}
		if s[i] == '\\' {
			escape = true
			continue
		}

		if isDigit(s[i]) && i+1 < len(s) && isDigit(s[i+1]) {
			return "", ErrInvalidString
		}

		if isDigit(s[i]) {
			n, err := strconv.Atoi(string(s[i]))
			if err != nil {
				return "", err
			}
			if n == 0 {
				tmp := sb.String()
				sb.Reset()
				sb.WriteString(tmp[:len(tmp)-1])
			} else {
				sb.WriteString(strings.Repeat(string(s[i-1]), n-1))
			}
		} else {
			sb.WriteByte(s[i])
		}
	}
	return sb.String(), nil
}
