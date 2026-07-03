package models

import (
	"strings"
	"testing"
)

func TestCreateUserInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateUserInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid input",
			input: CreateUserInput{
				Username:        "admin",
				Password:        "secret",
				ConfirmPassword: "secret",
			},
			wantErr: false,
		},
		{
			name: "Passwords do not match",
			input: CreateUserInput{
				Username:        "admin",
				Password:        "secret",
				ConfirmPassword: "wrong",
			},
			wantErr: true,
			errMsg:  "confirm password", // Should mention confirm password field
		},
		{
			name: "Empty username",
			input: CreateUserInput{
				Username:        "",
				Password:        "secret",
				ConfirmPassword: "secret",
			},
			wantErr: true,
			errMsg:  "login", // Should mention login field
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" {
				formatted := FormatValidationError(err)
				if !strings.Contains(formatted, tt.errMsg) {
					t.Errorf("FormatValidationError() = %v, want it to contain %v", formatted, tt.errMsg)
				}
			}
		})
	}
}

func TestUpdateUserInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   UpdateUserInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid username only",
			input: UpdateUserInput{
				Username: "admin",
			},
			wantErr: false,
		},
		{
			name: "Valid with password change",
			input: UpdateUserInput{
				Username:        "admin",
				Password:        "newsecret",
				ConfirmPassword: "newsecret",
			},
			wantErr: false,
		},
		{
			name: "Empty username",
			input: UpdateUserInput{
				Username: "",
			},
			wantErr: true,
			errMsg:  "login",
		},
		{
			name: "Passwords do not match",
			input: UpdateUserInput{
				Username:        "admin",
				Password:        "secret",
				ConfirmPassword: "wrong",
			},
			wantErr: true,
			errMsg:  "confirm password",
		},
		{
			name: "Password set without confirmation",
			input: UpdateUserInput{
				Username: "admin",
				Password: "secret",
			},
			wantErr: true,
			errMsg:  "confirm password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" {
				formatted := FormatValidationError(err)
				if !strings.Contains(formatted, tt.errMsg) {
					t.Errorf("FormatValidationError() = %v, want it to contain %v", formatted, tt.errMsg)
				}
			}
		})
	}
}
