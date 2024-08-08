package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Setup initial environment variables
	os.Setenv("HELLO", "SHOULD_REPLACE")
	os.Setenv("FOO", "SHOULD_REPLACE")
	os.Setenv("UNSET", "SHOULD_REMOVE")
	os.Setenv("ADDED", "from original env")
	os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

	// Define environment to be passed to RunCmd
	env := Environment{
		"HELLO": {Value: "world"},
		"FOO":   {Value: "bar"},
		"UNSET": {NeedRemove: true},
	}

	// Command to print environment variables
	cmd := []string{"env"}

	// Run the command
	returnCode := RunCmd(cmd, env)

	// Check return code
	assert.Equal(t, 0, returnCode, "expected return code 0")

	// Capture the output
	out, err := exec.Command("env").Output()
	require.NoError(t, err, "failed to run command")

	// Convert output to string and split into lines
	output := string(out)
	lines := strings.Split(output, "\n")

	// Check the environment variables
	expectedVars := map[string]string{
		"HELLO": "world",
		"FOO":   "bar",
		"ADDED": "from original env",
		"EMPTY": "SHOULD_BE_EMPTY",
	}
	unexpectedVars := []string{"UNSET"}

	for key, value := range expectedVars {
		found := false
		for _, line := range lines {
			if line == key+"="+value {
				found = true
				break
			}
		}
		assert.True(t, found, "expected to find %s=%s in environment", key, value)
	}

	for _, key := range unexpectedVars {
		for _, line := range lines {
			if strings.HasPrefix(line, key+"=") {
				assert.Fail(t, "expected %s to be unset, but found in environment", key)
			}
		}
	}
}

func TestRunCmd_EmptyFile(t *testing.T) {
	// Setup initial environment variables
	os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

	// Define environment to be passed to RunCmd
	env := Environment{
		"EMPTY": {NeedRemove: true},
	}

	// Command to print environment variables
	cmd := []string{"env"}

	// Run the command
	returnCode := RunCmd(cmd, env)

	// Check return code
	assert.Equal(t, 0, returnCode, "expected return code 0")

	// Capture the output
	out, err := exec.Command("env").Output()
	require.NoError(t, err, "failed to run command")

	// Convert output to string and split into lines
	output := string(out)
	lines := strings.Split(output, "\n")

	// Check that the EMPTY variable is unset
	for _, line := range lines {
		if strings.HasPrefix(line, "EMPTY=") {
			assert.Fail(t, "expected EMPTY to be unset, but found in environment")
		}
	}
}
