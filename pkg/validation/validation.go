package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	zhlocale "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	translator ut.Translator
	once       sync.Once
)

// Init 注册 Gin 默认参数校验器的中文翻译规则。
func Init() error {
	var initErr error

	once.Do(func() {
		engine, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			initErr = errors.New("gin validator engine is unavailable")
			return
		}

		engine.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := field.Tag.Get("json")
			if name == "" {
				return field.Name
			}

			name = strings.Split(name, ",")[0]
			if name == "" || name == "-" {
				return field.Name
			}

			return name
		})

		uni := ut.New(zhlocale.New())
		translator, _ = uni.GetTranslator("zh")
		if translator == nil {
			initErr = errors.New("zh translator is unavailable")
			return
		}

		if err := zhtranslations.RegisterDefaultTranslations(engine, translator); err != nil {
			initErr = err
			return
		}
	})

	return initErr
}

// TranslateError 将绑定或校验错误翻译为前端可直接展示的提示。
func TranslateError(err error) string {
	if err == nil {
		return ""
	}

	if translator == nil {
		_ = Init()
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok && translator != nil {
		messages := make([]string, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			messages = append(messages, fieldErr.Translate(translator))
		}

		return strings.Join(messages, "; ")
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return "请求体不是合法的 JSON 格式"
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := typeErr.Field
		if field == "" {
			field = "参数"
		}

		return fmt.Sprintf("%s类型错误，应为%s", field, typeErr.Type.String())
	}

	if errors.Is(err, io.EOF) {
		return "请求体不能为空"
	}

	return "请求参数不合法"
}
