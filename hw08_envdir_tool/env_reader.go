package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		if fileInfo, _ := os.Stat(filePath); fileInfo.Size() == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		envItem := EnvValue{}
		scanner := bufio.NewScanner(strings.NewReader(string(fileContent)))
		if scanner.Scan() {
			firstLine := scanner.Text()
			// Удаление пробелов и табуляций в конце
			firstLine = strings.TrimRight(firstLine, " \t")
			// Замена терминальных нулей на новую строку
			firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")
			envItem.Value = firstLine
		}

		env[file.Name()] = envItem
	}
	return env, nil
}
