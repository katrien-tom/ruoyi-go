package validation

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type testRequest struct {
	Mobile   string `json:"mobile" binding:"required,mobile"`
	Username string `json:"username" binding:"required,username"`
}

func TestCustomValidators(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("init validation: %v", err)
	}

	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("gin validator engine is unavailable")
	}

	tests := []struct {
		name      string
		req       testRequest
		wantError bool
	}{
		{
			name:      "valid mobile and username",
			req:       testRequest{Mobile: "13800138000", Username: "test_user"},
			wantError: false,
		},
		{
			name:      "invalid mobile",
			req:       testRequest{Mobile: "12345", Username: "test_user"},
			wantError: true,
		},
		{
			name:      "invalid username - too short",
			req:       testRequest{Mobile: "13800138000", Username: "abc"},
			wantError: true,
		},
		{
			name:      "invalid username - special chars",
			req:       testRequest{Mobile: "13800138000", Username: "test@user"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if (err != nil) != tt.wantError {
				t.Errorf("validation error = %v, wantError = %v", err, tt.wantError)
			}
		})
	}
}
