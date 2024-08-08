package main

import (
	"errors"
	"os"
	"os/exec"
)

func updateEnv(env Environment) {
	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		os.Setenv(k, v.Value)
	}
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	args := cmd[1:]
	command := cmd[0]
	c := exec.Command(command, args...)

	updateEnv(env)
	// Получение текущих переменных окружения
	currentEnv := os.Environ()
	// Копирование текущих переменных окружения в новый срез
	c.Env = append([]string{}, currentEnv...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err := c.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return 1
	}
	return 0
}
