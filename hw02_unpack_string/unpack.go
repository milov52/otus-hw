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

	if isDigit(runeArr[0]) || (runeArr[len(runeArr)-1] == '\\' && runeArr[len(runeArr)-2] != '\\') {
		return "", ErrInvalidString
	}

	var sb strings.Builder
	escape := false // Флаг для обработки символа '\'
	for i, ch := range runeArr {
		if escape {
			if isDigit(ch) || ch == '\\' {
				sb.WriteRune(ch)
			} else {
				return "", ErrInvalidString
			}
			escape = false
			continue
		}

		if ch == '\\' {
			escape = true
			continue
		}

		if isDigit(ch) {
			if i == 0 || (i < len(runeArr)-1 && isDigit(runeArr[i+1])) {
				return "", ErrInvalidString
			}

			count, err := strconv.Atoi(string(ch))
			if err != nil {
				return "", err
			}

			if count == 0 {
				if sb.Len() > 0 {
					tmpArr := []rune(sb.String())
					sb.Reset()
					sb.WriteString(string(tmpArr[:len(tmpArr)-1]))
				}
			} else {
				sb.WriteString(strings.Repeat(string(runeArr[i-1]), count-1))
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
