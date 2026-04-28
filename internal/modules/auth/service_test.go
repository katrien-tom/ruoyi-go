package auth

import (
	"strconv"
	"strings"
	"testing"

	"github.com/mojocn/base64Captcha"
)

func TestCaptchaDriverGenerate(t *testing.T) {
	driver := base64Captcha.NewDriverMath(
		48,
		130,
		2,
		base64Captcha.OptionShowSlimeLine|base64Captcha.OptionShowSineLine,
		nil,
		nil,
		nil,
	)
	id, question, answer := driver.GenerateIdQuestionAnswer()
	if id == "" {
		t.Fatal("expected non-empty captcha id")
	}
	if !strings.HasSuffix(question, "=?") {
		t.Fatalf("expected math question to end with =?, got %q", question)
	}
	if _, err := strconv.Atoi(answer); err != nil {
		t.Fatalf("expected numeric answer, got %q", answer)
	}
	if !strings.ContainsAny(question, "+-x") {
		t.Fatalf("expected math operator in question, got %q", question)
	}
}

func TestCaptchaDriverDraw(t *testing.T) {
	driver := base64Captcha.NewDriverMath(
		48,
		130,
		2,
		base64Captcha.OptionShowSlimeLine|base64Captcha.OptionShowSineLine,
		nil,
		nil,
		nil,
	)
	item, err := driver.DrawCaptcha("3+5=?")
	if err != nil {
		t.Fatalf("draw captcha: %v", err)
	}
	img := item.EncodeB64string()
	const prefix = "data:image/png;base64,"
	if !strings.HasPrefix(img, prefix) {
		t.Fatalf("expected prefix %q, got %q", prefix, img)
	}
	if len(img) <= len(prefix) {
		t.Fatalf("expected image payload after prefix, got %q", img)
	}
}
