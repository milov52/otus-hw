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
	// Ошибка если 1 число или последнее \
	if isDigit(runeArr[0]) || (runeArr[len(runeArr)-1] == '\\' && runeArr[len(runeArr)-2] != '\\') {
		return "", ErrInvalidString
	}

	var sb strings.Builder
	escape := false // Флаг для обработки символа '\'

	for i := range runeArr {
		if escape {
			if isDigit(runeArr[i]) || runeArr[i] == '\\' {
				sb.WriteRune(runeArr[i])
			} else {
				return "", ErrInvalidString
			}
			escape = false
			continue
		}
		if runeArr[i] == '\\' {
			escape = true
			continue
		}

		// Ошибка если 2 подряд числа
		if isDigit(runeArr[i]) && i+1 < len(runeArr) && isDigit(runeArr[i+1]) {
			return "", ErrInvalidString
		}

		if isDigit(runeArr[i]) {
			n, err := strconv.Atoi(string(runeArr[i]))
			if err != nil {
				return "", err
			}
			if n == 0 {
				tmpArr := []rune(sb.String())
				sb.Reset()
				sb.WriteString(string(tmpArr[:len(tmpArr)-1]))
			} else {
				sb.WriteString(strings.Repeat(string(runeArr[i-1]), n-1))
			}
		} else {
			sb.WriteRune(runeArr[i])
		}
	}

	return sb.String(), nil
}
