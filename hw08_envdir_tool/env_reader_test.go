package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestFile(dir, name, content string) error {
	filePath := filepath.Join(dir, name)
	return os.WriteFile(filePath, []byte(content), 0o644)
}

func TestReadDir_EmptyDirectory(t *testing.T) {
	dir, err := os.MkdirTemp("", "envdir_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	env, err := ReadDir(dir)
	require.NoError(t, err)

	assert.Empty(t, env, "expected empty environment map")
}

func TestReadDir_FilesWithEnvVariables(t *testing.T) {
	dir, err := os.MkdirTemp("", "envdir_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = createTestFile(dir, "VAR1", "value1\n")
	require.NoError(t, err)
	err = createTestFile(dir, "VAR2", "value2\n")
	require.NoError(t, err)

	env, err := ReadDir(dir)
	require.NoError(t, err)

	expected := Environment{
		"VAR1": {Value: "value1"},
		"VAR2": {Value: "value2"},
	}

	for key, val := range expected {
		assert.Equal(t, val.Value, env[key].Value, "mismatch for key "+key)
	}
}

func TestReadDir_FilesWithEqualsInName(t *testing.T) {
	dir, err := os.MkdirTemp("", "envdir_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = createTestFile(dir, "VAR=BAD", "value\n")
	require.NoError(t, err)

	env, err := ReadDir(dir)
	require.NoError(t, err)

	_, exists := env["VAR=BAD"]
	assert.False(t, exists, "expected VAR=BAD to be ignored")
}

func TestReadDir_ZeroLengthFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "envdir_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = createTestFile(dir, "EMPTY", "")
	require.NoError(t, err)

	env, err := ReadDir(dir)
	require.NoError(t, err)

	assert.True(t, env["EMPTY"].NeedRemove, "expected EMPTY to be marked for removal")
}

func TestReadDir_FilesWithTrailingSpaces(t *testing.T) {
	dir, err := os.MkdirTemp("", "envdir_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = createTestFile(dir, "VAR1", "value1   \n")
	require.NoError(t, err)
	err = createTestFile(dir, "VAR2", "value2\t\t\n")
	require.NoError(t, err)

	env, err := ReadDir(dir)
	require.NoError(t, err)

	assert.Equal(t, "value1", env["VAR1"].Value, "mismatch for VAR1")
	assert.Equal(t, "value2", env["VAR2"].Value, "mismatch for VAR2")
}

func TestReadDir_FilesWithNullTerminator(t *testing.T) {
	dir, err := os.MkdirTemp("", "envdir_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = createTestFile(dir, "VAR1", "value1\x00value2\n")
	require.NoError(t, err)

	env, err := ReadDir(dir)
	require.NoError(t, err)

	assert.Equal(t, "value1\nvalue2", env["VAR1"].Value, "mismatch for VAR1")
}
