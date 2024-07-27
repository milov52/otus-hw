package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		offset      int64
		limit       int64
		expectError error
		expected    string
	}{
		{
			name:        "Copy entire file",
			content:     "Hello, World!",
			offset:      0,
			limit:       0,
			expectError: nil,
			expected:    "Hello, World!",
		},
		{
			name:        "Copy with offset",
			content:     "Hello, World!",
			offset:      7,
			limit:       0,
			expectError: nil,
			expected:    "World!",
		},
		{
			name:        "Copy with limit",
			content:     "Hello, World!",
			offset:      0,
			limit:       5,
			expectError: nil,
			expected:    "Hello",
		},
		{
			name:        "Copy with offset and limit",
			content:     "Hello, World!",
			offset:      7,
			limit:       5,
			expectError: nil,
			expected:    "World",
		},
		{
			name:        "Offset exceeds file size",
			content:     "Hello, World!",
			offset:      20,
			limit:       0,
			expectError: ErrOffsetExceedsFileSize,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcFile, err := os.CreateTemp("", "src")
			require.NoError(t, err, "Failed to create temp source file")
			defer os.Remove(srcFile.Name())

			dstFile, err := os.CreateTemp("", "dst")
			require.NoError(t, err, "Failed to create temp destination file")
			defer os.Remove(dstFile.Name())

			_, err = srcFile.Write([]byte(tt.content))
			require.NoError(t, err, "Failed to write to source file")

			_, err = srcFile.Seek(0, io.SeekStart)
			require.NoError(t, err, "Failed to seek in source file")

			err = Copy(srcFile.Name(), dstFile.Name(), tt.offset, tt.limit)
			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err, "Unexpected error during copy")

				result, err := os.ReadFile(dstFile.Name())
				require.NoError(t, err, "Failed to read destination file")

				assert.Equal(t, tt.expected, string(result), "Unexpected content in destination file")
			}
		})
	}
}

func TestCopySameFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		offset      int64
		limit       int64
		expectError error
		expected    string
	}{
		{
			name:        "Copy entire file",
			content:     "Hello, World!",
			offset:      0,
			limit:       0,
			expectError: nil,
			expected:    "Hello, World!",
		},
		{
			name:        "Copy with offset",
			content:     "Hello, World!",
			offset:      7,
			limit:       0,
			expectError: nil,
			expected:    "World!",
		},
		{
			name:        "Copy with limit",
			content:     "Hello, World!",
			offset:      0,
			limit:       5,
			expectError: nil,
			expected:    "Hello",
		},
		{
			name:        "Copy with offset and limit",
			content:     "Hello, World!",
			offset:      7,
			limit:       5,
			expectError: nil,
			expected:    "World",
		},
		{
			name:        "Offset exceeds file size",
			content:     "Hello, World!",
			offset:      20,
			limit:       0,
			expectError: ErrOffsetExceedsFileSize,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcFile, err := os.CreateTemp("", "src")
			require.NoError(t, err, "Failed to create temp source file")
			defer os.Remove(srcFile.Name())

			_, err = srcFile.Write([]byte(tt.content))
			require.NoError(t, err, "Failed to write to source file")

			_, err = srcFile.Seek(0, io.SeekStart)
			require.NoError(t, err, "Failed to seek in source file")

			err = Copy(srcFile.Name(), srcFile.Name(), tt.offset, tt.limit)
			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err, "Unexpected error during copy")

				result, err := os.ReadFile(srcFile.Name())
				require.NoError(t, err, "Failed to read destination file")

				assert.Equal(t, tt.expected, string(result), "Unexpected content in destination file")
			}
		})
	}
}
