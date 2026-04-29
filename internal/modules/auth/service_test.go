package auth

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"

	"github.com/banyejiu/ruoyi-go/internal/modules/user"
	"github.com/banyejiu/ruoyi-go/pkg/jwtutil"
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

func TestLoginUserSuccess(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	now := time.Unix(1893456000, 0)
	service := &Service{
		now: func() time.Time { return now },
	}

	resp, err := service.loginUser(LoginRequest{
		Password: "secret123",
	}, &user.SysUser{
		UserID:   1,
		UserName: "admin",
		NickName: "Admin",
		Password: string(passwordHash),
		Status:   "0",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token")
	}
	if resp.TokenType != "Bearer" {
		t.Fatalf("expected Bearer token type, got %q", resp.TokenType)
	}
	if resp.UserID != 1 || resp.UserName != "admin" || resp.NickName != "Admin" {
		t.Fatalf("unexpected login response: %+v", resp)
	}

	claims, err := jwtutil.Parse(resp.Token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.UserID != 1 {
		t.Fatalf("expected user_id claim 1, got %d", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Fatalf("expected username claim admin, got %q", claims.Username)
	}
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	service := &Service{
		now: time.Now,
	}

	_, err = service.loginUser(LoginRequest{
		Password: "wrong-password",
	}, &user.SysUser{
		UserID:   1,
		UserName: "admin",
		Password: string(passwordHash),
		Status:   "0",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestLoginUserDisabledAccount(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	service := &Service{
		now: time.Now,
	}

	_, err = service.loginUser(LoginRequest{
		Password: "secret123",
	}, &user.SysUser{
		UserID:   1,
		UserName: "admin",
		Password: string(passwordHash),
		Status:   "1",
	})
	if !errors.Is(err, ErrAccountDisabled) {
		t.Fatalf("expected account disabled, got %v", err)
	}
}
