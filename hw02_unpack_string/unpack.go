package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func isDigit(ch rune) bool {
	_, err := strconv.Atoi(string(ch))
	return err == nil
}

func Unpack(s string) (string, error) {
	runeArr := []rune(s)
	if s == "" {
		return s, nil
	}
	if isDigit(runeArr[0]) {
		return "", ErrInvalidString
	}

	var sb strings.Builder
	escape := false // Флаг для обработки символа '\'

	for i := range runeArr {
		if escape {
			sb.WriteRune(runeArr[i])
			escape = false
			continue
		}
		if s[i] == '\\' {
			escape = true
			continue
		}

		if isDigit(runeArr[i]) && i+1 < len(runeArr) && isDigit(runeArr[i+1]) {
			return "", ErrInvalidString
		}

		if isDigit(runeArr[i]) {
			n, err := strconv.Atoi(string(runeArr[i]))
			if err != nil {
				return "", err
			}
			if n == 0 {
				tmp := sb.String()
				sb.Reset()
				sb.WriteString(tmp[:len(tmp)-1])
			} else {
				sb.WriteString(strings.Repeat(string(runeArr[i-1]), n-1))
			}
		} else {
			sb.WriteRune(runeArr[i])
		}
	}
	return sb.String(), nil
}
