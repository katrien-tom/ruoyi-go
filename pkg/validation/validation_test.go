package validation

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
}

func TestTranslateValidationError(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("init validation translator: %v", err)
	}

	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("gin validator engine is unavailable")
	}

	err := validate.Struct(loginRequest{})
	if err == nil {
		t.Fatal("expected validation error")
	}

	msg := TranslateError(err)
	if msg != "username为必填字段" {
		t.Fatalf("unexpected translation: %q", msg)
	}
}

func TestTranslateJSONErrors(t *testing.T) {
	msg := TranslateError(&json.SyntaxError{})
	if msg != "请求体不是合法的 JSON 格式" {
		t.Fatalf("unexpected syntax error message: %q", msg)
	}

	msg = TranslateError(&json.UnmarshalTypeError{
		Field: "username",
		Type:  reflect.TypeFor[int](),
	})
	if msg != "username类型错误，应为int" {
		t.Fatalf("unexpected type error message: %q", msg)
	}
}
