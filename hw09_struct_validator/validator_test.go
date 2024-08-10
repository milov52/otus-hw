package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte `validate:"len:16"`
		Payload   []byte `validate:"len:32"`
		Signature []byte `validate:"len:64"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			name: "invalid ID and Role",
			in: User{
				ID:     "invalid",
				Name:   "Jon Snow",
				Age:    25,
				Email:  "jonsnow@got.com",
				Role:   "admin1",
				Phones: []string{"12345678901", "09876543210"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   fmt.Errorf("expected string len 36, got 7"),
				},
				{
					Field: "Role",
					Err:   fmt.Errorf("value admin1 does not contains into [admin stuff]"),
				},
			},
		},
		{
			name: "valid user",
			in: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Jon Snow",
				Age:    25,
				Email:  "jonsnow@got.com",
				Role:   "admin",
				Phones: []string{"12345678901", "09876543210"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			name: "invalid email and age",
			in: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Arya Stark",
				Age:    16,
				Email:  "invalid-email",
				Role:   "stuff",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
					Err:   fmt.Errorf("value 16 is less than min 18"),
				},
				{
					Field: "Email",
					Err:   fmt.Errorf("value invalid-email does not match regexp ^\\w+@\\w+\\.\\w+$"),
				},
			},
		},
		{
			name: "invalid phone length",
			in: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Tyrion Lannister",
				Age:    30,
				Email:  "tyrion@got.com",
				Role:   "admin",
				Phones: []string{"1234567"}, // Invalid phone length
			},
			expectedErr: ValidationErrors{
				{
					Field: "Phones[0]",
					Err:   fmt.Errorf("expected string len 11, got 7"),
				},
			},
		},

		{
			name: "non-struct type",
			in:   54,
			expectedErr: ValidationErrors{
				{
					Field: "int",
					Err:   fmt.Errorf("expected struct, got int"),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			if err != nil {
				var vErr ValidationErrors

				if !errors.As(err, &vErr) {
					t.Fatalf("expected error %q, got %q", tt.expectedErr, err)
				}

				for i, expectedErr := range tt.expectedErr {
					if vErr[i].Field != expectedErr.Field {
						t.Errorf("expected field %q, got %q", expectedErr.Field, vErr[i].Field)
					}
					if vErr[i].Err.Error() != expectedErr.Err.Error() {
						t.Errorf("expected error at index %d: %v, got: %v", i, expectedErr.Err, vErr[i].Err)
					}
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			name: "valid token",
			in: Token{
				Header:    make([]byte, 16), // Valid length
				Payload:   make([]byte, 32), // Valid length
				Signature: make([]byte, 64), // Valid length
			},
			expectedErr: ValidationErrors{},
		},
		{
			name: "invalid token header length",
			in: Token{
				Header:    make([]byte, 15), // Invalid length
				Payload:   make([]byte, 32),
				Signature: make([]byte, 64),
			},
			expectedErr: ValidationErrors{
				{
					Field: "Header",
					Err:   fmt.Errorf("expected length 16, got 15"),
				},
			},
		},
		{
			name: "invalid token payload length",
			in: Token{
				Header:    make([]byte, 16),
				Payload:   make([]byte, 31), // Invalid length
				Signature: make([]byte, 64),
			},
			expectedErr: ValidationErrors{
				{
					Field: "Payload",
					Err:   fmt.Errorf("expected length 32, got 31"),
				},
			},
		},
		{
			name: "invalid token signature length",
			in: Token{
				Header:    make([]byte, 16),
				Payload:   make([]byte, 32),
				Signature: make([]byte, 63), // Invalid length
			},
			expectedErr: ValidationErrors{
				{
					Field: "Signature",
					Err:   fmt.Errorf("expected length 64, got 63"),
				},
			},
		},
		{
			name: "non-struct type",
			in:   54,
			expectedErr: ValidationErrors{
				{
					Field: "int",
					Err:   fmt.Errorf("expected struct, got int"),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			if err != nil {
				var vErr ValidationErrors

				if !errors.As(err, &vErr) {
					t.Fatalf("expected error %q, got %q", tt.expectedErr, err)
				}

				for i, expectedErr := range tt.expectedErr {
					if vErr[i].Field != expectedErr.Field {
						t.Errorf("expected field %q, got %q", expectedErr.Field, vErr[i].Field)
					}
					if vErr[i].Err.Error() != expectedErr.Err.Error() {
						t.Errorf("expected error at index %d: %v, got: %v", i, expectedErr.Err, vErr[i].Err)
					}
				}
			}
		})
	}
}

func TestValidateApp(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			name: "valid version",
			in: App{
				Version: "1.0.0", // Valid length (5 characters)
			},
			expectedErr: ValidationErrors{},
		},
		{
			name: "too short version",
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   fmt.Errorf("expected string len 5, got 3"),
				},
			},
		},
		{
			name: "too long version",
			in: App{
				Version: "1.0.00",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   fmt.Errorf("expected string len 5, got 6"),
				},
			},
		},
		{
			name: "non-struct type",
			in:   54,
			expectedErr: ValidationErrors{
				{
					Field: "int",
					Err:   fmt.Errorf("expected struct, got int"),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			if err != nil {
				var vErr ValidationErrors
				if !errors.As(err, &vErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}

				if len(vErr) != len(tt.expectedErr) {
					t.Fatalf("expected %d errors, got %d", len(tt.expectedErr), len(vErr))
				}

				for i, expectedErr := range tt.expectedErr {
					if vErr[i].Field != expectedErr.Field {
						t.Errorf("expected field %q, got %q", expectedErr.Field, vErr[i].Field)
					}
					if vErr[i].Err.Error() != expectedErr.Err.Error() {
						t.Errorf("expected error at index %d: %v, got: %v", i, expectedErr.Err, vErr[i].Err)
					}
				}
			}
		})
	}
}

func TestValidateResponse(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			name: "valid response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: ValidationErrors{},
		},
		{
			name: "invalid response code",
			in: Response{
				Code: 403, // Invalid code
				Body: "Forbidden",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   fmt.Errorf("value 403 does not contains into [200 404 500]"),
				},
			},
		},
		{
			name: "valid response with empty body",
			in: Response{
				Code: 404, // Valid code
				Body: "",  // Body can be empty
			},
			expectedErr: ValidationErrors{},
		},
		{
			name: "valid response with empty body but invalid code",
			in: Response{
				Code: 500, // Valid code
				Body: "",  // Body can be empty
			},
			expectedErr: ValidationErrors{},
		},
		{
			name: "missing body field",
			in: Response{
				Code: 200, // Valid code
			},
			expectedErr: ValidationErrors{}, // Body is omitted, no error
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			if err != nil {
				var vErr ValidationErrors
				if !errors.As(err, &vErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}

				if len(vErr) != len(tt.expectedErr) {
					t.Fatalf("expected %d errors, got %d", len(tt.expectedErr), len(vErr))
				}

				for i, expectedErr := range tt.expectedErr {
					if vErr[i].Field != expectedErr.Field {
						t.Errorf("expected field %q, got %q", expectedErr.Field, vErr[i].Field)
					}

					if vErr[i].Err.Error() != expectedErr.Err.Error() {
						t.Errorf("expected error at index %d: %v, got: %v", i, expectedErr.Err, vErr[i].Err)
					}
				}
			}
		})
	}
}
