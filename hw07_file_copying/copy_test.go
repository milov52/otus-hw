package main

import (
	"errors"
	"io"
	"os"
	"testing"
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
			if err != nil {
				t.Fatalf("Failed to create temp source file: %v", err)
			}
			defer os.Remove(srcFile.Name())

			dstFile, err := os.CreateTemp("", "dst")
			if err != nil {
				t.Fatalf("Failed to create temp destination file: %v", err)
			}
			defer os.Remove(dstFile.Name())

			if _, err := srcFile.Write([]byte(tt.content)); err != nil {
				t.Fatalf("Failed to write to source file: %v", err)
			}

			// Reset file pointer to start
			if _, err := srcFile.Seek(0, io.SeekStart); err != nil {
				t.Fatalf("Failed to seek in source file: %v", err)
			}

			err = Copy(srcFile.Name(), dstFile.Name(), tt.offset, tt.limit)
			if !errors.Is(err, tt.expectError) {
				t.Fatalf("Expected error %v, got %v", tt.expectError, err)
			}

			if tt.expectError == nil {
				result, err := io.ReadAll(dstFile)
				if err != nil {
					t.Fatalf("Failed to read destination file: %v", err)
				}

				if string(result) != tt.expected {
					t.Errorf("Expected content %q, got %q", tt.expected, string(result))
				}
			}
		})
	}
}
