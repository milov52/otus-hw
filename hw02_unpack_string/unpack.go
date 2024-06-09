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

func isSymbol(ch rune) bool {
	return !isDigit(ch) && ch != '\\'
}

func Unpack(s string) (string, error) {
	runeArr := []rune(s)
	if s == "" {
		return s, nil
	}
	// Ошибка если 1 число или последнее \
	if isDigit(runeArr[0]) {
		return "", ErrInvalidString
	}

	var sb strings.Builder
	escape := false // Флаг для обработки символа '\'

	for i, ch := range runeArr {
		// Ошибка - экранирование символа
		if ch == '\\' && i+1 < len(runeArr) && isSymbol(runeArr[i+1]) {
			return "", ErrInvalidString
		}

		// Ошибка если 2 подряд числа
		if isDigit(ch) && i+1 < len(runeArr) && isDigit(runeArr[i+1]) && !escape {
			return "", ErrInvalidString
		}

		if escape {
			sb.WriteRune(ch)
			escape = false
			continue
		}
		if ch == '\\' {
			escape = true
			continue
		}

		if isDigit(ch) {
			n, err := strconv.Atoi(string(ch))
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
			sb.WriteRune(ch)
		}
	}
	if escape {
		return "", ErrInvalidString
	}
	return sb.String(), nil
}
