package auth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/banyejiu/ruoyi-go/internal/modules/user"
	"github.com/banyejiu/ruoyi-go/pkg/database"
	"github.com/banyejiu/ruoyi-go/pkg/jwtutil"
	"github.com/banyejiu/ruoyi-go/pkg/logger"
)

var (
	integrationConfigOnce sync.Once
	integrationConfigErr  error
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

func TestLoginRequiresValidCaptcha(t *testing.T) {
	ctx := context.Background()
	db, redisClient := integrationDependencies(t)
	service := NewService(redisClient, user.NewRepository(db))

	_, err := service.Login(ctx, LoginRequest{
		Username: "admin",
		Password: "secret123",
		Code:     "8",
		UUID:     "captcha-missing",
	})
	if !errors.Is(err, ErrCaptchaInvalid) {
		t.Fatalf("expected captcha invalid, got %v", err)
	}
}

func TestLoginSuccessAfterCaptchaVerification(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	ctx := context.Background()
	db, redisClient := integrationDependencies(t)
	username := "auth_it_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	testUser := user.SysUser{
		UserName: username,
		NickName: "Integration User",
		Password: string(passwordHash),
		Status:   "0",
		DelFlag:  "0",
	}
	if err := db.WithContext(ctx).Create(&testUser).Error; err != nil {
		t.Fatalf("create test user: %v", err)
	}
	t.Cleanup(func() {
		if err := db.WithContext(context.Background()).Delete(&user.SysUser{}, testUser.UserID).Error; err != nil {
			t.Fatalf("cleanup test user: %v", err)
		}
	})

	uuid := "captcha-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := redisClient.Set(ctx, captchaCacheKey(uuid), "8", time.Minute).Err(); err != nil {
		t.Fatalf("seed captcha: %v", err)
	}
	t.Cleanup(func() {
		if err := redisClient.Del(context.Background(), captchaCacheKey(uuid)).Err(); err != nil {
			t.Fatalf("cleanup captcha: %v", err)
		}
	})

	now := time.Unix(1893456000, 0)
	service := &Service{
		redis:          redisClient,
		userRepository: user.NewRepository(db),
		now:            func() time.Time { return now },
	}

	resp, err := service.Login(context.Background(), LoginRequest{
		Username: username,
		Password: "secret123",
		Code:     "8",
		UUID:     uuid,
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if resp == nil || resp.Token == "" {
		t.Fatalf("expected login token, got %+v", resp)
	}
	if resp.UserName != username {
		t.Fatalf("expected user name %q, got %q", username, resp.UserName)
	}
	if _, err := redisClient.Get(ctx, captchaCacheKey(uuid)).Result(); !errors.Is(err, redis.Nil) {
		t.Fatalf("expected captcha to be consumed, got err=%v", err)
	}
}

func integrationDependencies(t *testing.T) (*gorm.DB, *redis.Client) {
	t.Helper()

	loadIntegrationConfig(t)
	logger.Init("dev")

	db, err := database.InitDB()
	if err != nil {
		t.Skipf("skip integration test, mysql unavailable: %v", err)
	}
	if !db.Migrator().HasTable(&user.SysUser{}) {
		t.Skip("skip integration test, sys_user table not found")
	}

	redisClient, err := database.InitRedis()
	if err != nil {
		t.Skipf("skip integration test, redis unavailable: %v", err)
	}

	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		_ = redisClient.Close()
	})

	return db, redisClient
}

func loadIntegrationConfig(t *testing.T) {
	t.Helper()

	integrationConfigOnce.Do(func() {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			integrationConfigErr = errors.New("resolve test file path failed")
			return
		}

		configPath := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../config.yml"))
		if _, err := os.Stat(configPath); err != nil {
			integrationConfigErr = err
			return
		}

		viper.Reset()
		viper.SetConfigFile(configPath)
		integrationConfigErr = viper.ReadInConfig()
	})

	if integrationConfigErr != nil {
		t.Skipf("skip integration test, config unavailable: %v", integrationConfigErr)
	}
}
