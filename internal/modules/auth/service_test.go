package auth

import (
	"context"
	"errors"
	"net/http/httptest"
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
	"github.com/banyejiu/ruoyi-go/internal/security"
	"github.com/banyejiu/ruoyi-go/pkg/constants"
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

	redisClient := redisClientForUnitTest(t)
	tokenStore := security.NewTokenStore(redisClient)
	now := time.Unix(1893456000, 0)
	tokenStore.SetNowFunc(func() time.Time { return now })
	service := &Service{
		redis:          redisClient,
		tokenStore:     tokenStore,
		sessionService: security.NewSessionService(tokenStore),
		now:            func() time.Time { return now },
	}

	resp, err := service.loginUser(context.Background(), LoginRequest{
		Password: "secret123",
	}, &user.SysUser{
		UserID:   1,
		UserName: "admin",
		NickName: "Admin",
		Password: string(passwordHash),
		Status:   "0",
	}, LoginMeta{
		IPAddr:      "127.0.0.1",
		Browser:     "Chrome",
		OS:          "macOS",
		Roles:       []string{"admin"},
		Permissions: []string{"system:user:list"},
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
	if claims.Token == "" {
		t.Fatal("expected login token claim")
	}

	loginUser, err := service.GetLoginUser(context.Background(), claims.Token)
	if err != nil {
		t.Fatalf("get login user: %v", err)
	}
	if loginUser.UserID != 1 || loginUser.User == nil || loginUser.User.UserName != "admin" {
		t.Fatalf("unexpected login user: %+v", loginUser)
	}
	if len(loginUser.Permissions) != 1 || loginUser.Permissions[0] != "system:user:list" {
		t.Fatalf("unexpected permissions: %+v", loginUser.Permissions)
	}
	if len(loginUser.Roles) != 1 || loginUser.Roles[0] != "admin" {
		t.Fatalf("unexpected roles: %+v", loginUser.Roles)
	}
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	redisClient := redisClientForUnitTest(t)
	tokenStore := security.NewTokenStore(redisClient)
	service := &Service{
		redis:          redisClient,
		tokenStore:     tokenStore,
		sessionService: security.NewSessionService(tokenStore),
		now:            time.Now,
	}

	_, err = service.loginUser(context.Background(), LoginRequest{
		Password: "wrong-password",
	}, &user.SysUser{
		UserID:   1,
		UserName: "admin",
		Password: string(passwordHash),
		Status:   "0",
	}, LoginMeta{})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestLoginUserDisabledAccount(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	redisClient := redisClientForUnitTest(t)
	tokenStore := security.NewTokenStore(redisClient)
	service := &Service{
		redis:          redisClient,
		tokenStore:     tokenStore,
		sessionService: security.NewSessionService(tokenStore),
		now:            time.Now,
	}

	_, err = service.loginUser(context.Background(), LoginRequest{
		Password: "secret123",
	}, &user.SysUser{
		UserID:   1,
		UserName: "admin",
		Password: string(passwordHash),
		Status:   "1",
	}, LoginMeta{})
	if !errors.Is(err, ErrAccountDisabled) {
		t.Fatalf("expected account disabled, got %v", err)
	}
}

func TestLogoutDeletesLoginSession(t *testing.T) {
	ctx := context.Background()
	db, redisClient := integrationDependencies(t)

	now := time.Unix(1893456000, 0)
	tokenStore := security.NewTokenStore(redisClient)
	tokenStore.SetNowFunc(func() time.Time { return now })
	service := &Service{
		redis:          redisClient,
		userAuthReader: user.NewService(user.NewRepository(db)),
		tokenStore:     tokenStore,
		sessionService: security.NewSessionService(tokenStore),
		now:            func() time.Time { return now },
	}

	loginUser := &security.LoginUser{
		UserID:     1,
		Token:      "session-token",
		LoginTime:  now.UnixMilli(),
		ExpireTime: jwtutil.ExpiresAt(now).UnixMilli(),
		User: &security.UserInfo{
			UserID:   1,
			UserName: "admin",
		},
	}
	if err := service.tokenStore.Save(ctx, loginUser); err != nil {
		t.Fatalf("save login session: %v", err)
	}

	if err := service.Logout(ctx, loginUser.Token); err != nil {
		t.Fatalf("logout: %v", err)
	}

	_, err := service.GetLoginUser(ctx, loginUser.Token)
	if !errors.Is(err, security.ErrTokenNotFound) {
		t.Fatalf("expected token not found, got %v", err)
	}
}

func TestLogoutRejectsMissingToken(t *testing.T) {
	redisClient := redisClientForUnitTest(t)
	tokenStore := security.NewTokenStore(redisClient)
	service := &Service{
		redis:          redisClient,
		tokenStore:     tokenStore,
		sessionService: security.NewSessionService(tokenStore),
		now:            time.Now,
	}

	err := service.Logout(context.Background(), "")
	if !errors.Is(err, security.ErrTokenInvalid) {
		t.Fatalf("expected token invalid, got %v", err)
	}
}

func TestGetInfoReturnsRolesAndPermissions(t *testing.T) {
	redisClient := redisClientForUnitTest(t)
	tokenStore := security.NewTokenStore(redisClient)
	service := &Service{
		redis:          redisClient,
		tokenStore:     tokenStore,
		sessionService: security.NewSessionService(tokenStore),
		now:            time.Now,
	}

	now := time.Unix(1893456000, 0)
	loginUser := &security.LoginUser{
		UserID:      1,
		Token:       "get-info-token",
		LoginTime:   now.UnixMilli(),
		ExpireTime:  jwtutil.ExpiresAt(now).UnixMilli(),
		Roles:       []string{"admin"},
		Permissions: []string{"system:user:list"},
		User: &security.UserInfo{
			UserID:   1,
			UserName: "admin",
			NickName: "Admin",
		},
	}
	if err := tokenStore.Save(context.Background(), loginUser); err != nil {
		t.Fatalf("save login session: %v", err)
	}

	data, err := service.GetInfo(context.Background(), loginUser.Token)
	if err != nil {
		t.Fatalf("get info: %v", err)
	}
	if data.User == nil || data.User.UserName != "admin" {
		t.Fatalf("unexpected user info: %+v", data.User)
	}
	if len(data.Roles) != 1 || data.Roles[0] != "admin" {
		t.Fatalf("unexpected roles: %+v", data.Roles)
	}
	if len(data.Permissions) != 1 || data.Permissions[0] != "system:user:list" {
		t.Fatalf("unexpected permissions: %+v", data.Permissions)
	}
}

func TestListOnlineUsersAndForceLogout(t *testing.T) {
	redisClient := redisClientForUnitTest(t)
	tokenStore := security.NewTokenStore(redisClient)
	service := &Service{
		redis:          redisClient,
		tokenStore:     tokenStore,
		sessionService: security.NewSessionService(tokenStore),
		now:            time.Now,
	}

	now := time.Unix(1893456000, 0)
	tokenStore.SetNowFunc(func() time.Time { return now })
	loginUser := &security.LoginUser{
		UserID:      1,
		Token:       "online-token",
		LoginTime:   now.UnixMilli(),
		ExpireTime:  jwtutil.ExpiresAt(now).UnixMilli(),
		IPAddr:      "127.0.0.1",
		Browser:     "Chrome",
		OS:          "macOS",
		Permissions: []string{"monitor:online:query"},
		User: &security.UserInfo{
			UserID:   1,
			UserName: "admin",
			NickName: "Admin",
		},
	}
	if err := tokenStore.Save(context.Background(), loginUser); err != nil {
		t.Fatalf("save login session: %v", err)
	}

	items, err := service.ListOnlineUsers(context.Background())
	if err != nil {
		t.Fatalf("list online users: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected online users")
	}

	found := false
	for _, item := range items {
		if item.TokenID == loginUser.Token {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected token %q in online users: %+v", loginUser.Token, items)
	}

	if err := service.ForceLogout(context.Background(), loginUser.Token); err != nil {
		t.Fatalf("force logout: %v", err)
	}
	if _, err := service.GetLoginUser(context.Background(), loginUser.Token); !errors.Is(err, security.ErrTokenNotFound) {
		t.Fatalf("expected login user removed, got %v", err)
	}
}

func TestLoginRequiresValidCaptcha(t *testing.T) {
	ctx := context.Background()
	db, redisClient := integrationDependencies(t)
	service := NewService(redisClient, user.NewService(user.NewRepository(db)), nil)

	_, err := service.Login(ctx, LoginRequest{
		Username: "admin",
		Password: "secret123",
		Code:     "8",
		UUID:     "captcha-missing",
	}, LoginMeta{})
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
	if err := redisClient.Set(ctx, constants.CaptchaCacheKey(uuid), "8", time.Minute).Err(); err != nil {
		t.Fatalf("seed captcha: %v", err)
	}
	t.Cleanup(func() {
		if err := redisClient.Del(context.Background(), constants.CaptchaCacheKey(uuid)).Err(); err != nil {
			t.Fatalf("cleanup captcha: %v", err)
		}
	})

	now := time.Unix(1893456000, 0)
	tokenStore := security.NewTokenStore(redisClient)
	tokenStore.SetNowFunc(func() time.Time { return now })
	service := &Service{
		redis:          redisClient,
		userAuthReader: user.NewService(user.NewRepository(db)),
		tokenStore:     tokenStore,
		sessionService: security.NewSessionService(tokenStore),
		now:            func() time.Time { return now },
	}

	req := httptest.NewRequest("POST", "/api/public/auth/login", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X) AppleWebKit/537.36 Chrome/136.0.0.0 Safari/537.36")

	resp, err := service.Login(context.Background(), LoginRequest{
		Username: username,
		Password: "secret123",
		Code:     "8",
		UUID:     uuid,
	}, BuildLoginMeta(req, nil))
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if resp == nil || resp.Token == "" {
		t.Fatalf("expected login token, got %+v", resp)
	}
	if resp.UserName != username {
		t.Fatalf("expected user name %q, got %q", username, resp.UserName)
	}
	if resp.ExpiresAt <= now.UnixMilli() {
		t.Fatalf("expected expiresAt after login time, got %d", resp.ExpiresAt)
	}
	if _, err := redisClient.Get(ctx, constants.CaptchaCacheKey(uuid)).Result(); !errors.Is(err, redis.Nil) {
		t.Fatalf("expected captcha to be consumed, got err=%v", err)
	}

	claims, err := jwtutil.Parse(resp.Token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	loginUser, err := service.GetLoginUser(ctx, claims.Token)
	if err != nil {
		t.Fatalf("get login user: %v", err)
	}
	if loginUser.IPAddr != "127.0.0.1" || loginUser.Browser != "Chrome" || loginUser.OS != "macOS" {
		t.Fatalf("unexpected login meta: %+v", loginUser)
	}
	if loginUser.User == nil || loginUser.User.UserName != username {
		t.Fatalf("unexpected login user payload: %+v", loginUser)
	}
}

func TestBuildLoginMeta(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/public/auth/login", nil)
	req.RemoteAddr = "10.0.0.1:8080"
	req.Header.Set("X-Forwarded-For", "203.0.113.5, 10.0.0.1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/136.0.0.0 Safari/537.36")

	meta := BuildLoginMeta(req, []string{"system:user:list"})
	if meta.IPAddr != "203.0.113.5" {
		t.Fatalf("expected forwarded ip, got %q", meta.IPAddr)
	}
	if meta.Browser != "Chrome" {
		t.Fatalf("expected Chrome browser, got %q", meta.Browser)
	}
	if meta.OS != "Windows" {
		t.Fatalf("expected Windows os, got %q", meta.OS)
	}
	if len(meta.Permissions) != 1 || meta.Permissions[0] != "system:user:list" {
		t.Fatalf("unexpected permissions: %+v", meta.Permissions)
	}
}

func TestTokenStoreRoundTrip(t *testing.T) {
	ctx := context.Background()
	store := security.NewTokenStore(redisClientForUnitTest(t))
	now := time.Unix(1893456000, 0)
	store.SetNowFunc(func() time.Time { return now })

	loginUser := &security.LoginUser{
		UserID:     1,
		Token:      "round-trip-token",
		LoginTime:  now.UnixMilli(),
		ExpireTime: jwtutil.ExpiresAt(now).UnixMilli(),
		IPAddr:     "127.0.0.1",
		Permissions: []string{
			"system:user:list",
		},
		User: &security.UserInfo{
			UserID:   1,
			UserName: "admin",
		},
	}

	if err := store.Save(ctx, loginUser); err != nil {
		t.Fatalf("save token store: %v", err)
	}

	saved, err := store.Get(ctx, loginUser.Token)
	if err != nil {
		t.Fatalf("get token store: %v", err)
	}
	if saved.UserID != loginUser.UserID || saved.User == nil || saved.User.UserName != "admin" {
		t.Fatalf("unexpected saved login user: %+v", saved)
	}
	if err := store.Delete(ctx, loginUser.Token); err != nil {
		t.Fatalf("delete token store: %v", err)
	}
}

func redisClientForUnitTest(t *testing.T) *redis.Client {
	t.Helper()

	loadIntegrationConfig(t)
	logger.Init("dev")

	redisClient, err := database.InitRedis()
	if err != nil {
		t.Skipf("skip test, redis unavailable: %v", err)
	}

	t.Cleanup(func() {
		_ = redisClient.Close()
	})

	return redisClient
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
