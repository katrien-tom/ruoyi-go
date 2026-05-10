package auth

import (
	"context"
	"errors"
	"image/color"
	"sort"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/banyejiu/ruoyi-go/internal/modules/menu"
	"github.com/banyejiu/ruoyi-go/internal/modules/user"
	"github.com/banyejiu/ruoyi-go/internal/security"
	"github.com/banyejiu/ruoyi-go/pkg/jwtutil"
)

const captchaTTL = 5 * time.Minute

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountDisabled    = errors.New("account disabled")
	ErrCaptchaInvalid     = errors.New("captcha invalid")
)

type Service struct {
	redis          *redis.Client
	userAuthReader UserAuthReader
	menuReader     MenuReader
	tokenStore     *security.TokenStore
	sessionService *security.SessionService
	now            func() time.Time
}

type UserAuthReader interface {
	FindByUserName(ctx context.Context, userName string) (*user.SysUser, error)
	FindPermissionsByUserID(ctx context.Context, userID int64) ([]string, error)
	FindRoleKeysByUserID(ctx context.Context, userID int64) ([]string, error)
}

type MenuReader interface {
	FindByUserID(ctx context.Context, userID int64) ([]menu.SysMenu, error)
}

func NewService(redisClient *redis.Client, userAuthReader UserAuthReader, menuReader MenuReader) *Service {
	return &Service{
		redis:          redisClient,
		userAuthReader: userAuthReader,
		menuReader:     menuReader,
		tokenStore:     security.NewTokenStore(redisClient),
		sessionService: security.NewSessionService(security.NewTokenStore(redisClient)),
		now:            time.Now,
	}
}

func (s *Service) GetCaptcha(ctx context.Context) (*CaptchaResponse, error) {
	driver := base64Captcha.NewDriverMath(
		48,
		130,
		1,
		base64Captcha.OptionShowSlimeLine|base64Captcha.OptionShowSineLine,
		&color.RGBA{R: 244, G: 246, B: 251, A: 255},
		nil,
		nil,
	)
	id, question, answer := driver.GenerateIdQuestionAnswer()
	item, err := driver.DrawCaptcha(question)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Set(ctx, captchaCacheKey(id), answer, captchaTTL).Err(); err != nil {
		return nil, err
	}

	return &CaptchaResponse{
		CaptchaEnabled: true,
		Img:            item.EncodeB64string(),
		UUID:           id,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest, meta LoginMeta) (*LoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	code := strings.TrimSpace(req.Code)
	uuid := strings.TrimSpace(req.UUID)
	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}
	if err := s.verifyCaptcha(ctx, uuid, code); err != nil {
		return nil, err
	}

	authUser, err := s.userAuthReader.FindByUserName(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	permissions, err := s.userAuthReader.FindPermissionsByUserID(ctx, authUser.UserID)
	if err != nil {
		return nil, err
	}
	meta.Permissions = permissions

	roles, err := s.userAuthReader.FindRoleKeysByUserID(ctx, authUser.UserID)
	if err != nil {
		return nil, err
	}
	meta.Roles = roles

	return s.loginUser(ctx, req, authUser, meta)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.sessionService.DeleteLoginUser(ctx, token)
}

func (s *Service) GetLoginUser(ctx context.Context, token string) (*security.LoginUser, error) {
	return s.sessionService.GetLoginUser(ctx, token)
}

func (s *Service) GetInfo(ctx context.Context, token string) (*GetInfoResponse, error) {
	loginUser, err := s.GetLoginUser(ctx, token)
	if err != nil {
		return nil, err
	}

	return &GetInfoResponse{
		User:        loginUser.User,
		Roles:       defaultRoles(loginUser.Roles),
		Permissions: defaultPermissions(loginUser.Permissions),
	}, nil
}

func (s *Service) GetRouters(ctx context.Context, token string) ([]RouterResponse, error) {
	loginUser, err := s.GetLoginUser(ctx, token)
	if err != nil {
		return nil, err
	}

	menus, err := s.menuReader.FindByUserID(ctx, loginUser.UserID)
	if err != nil {
		return nil, err
	}

	return buildRouterTree(menus, 0), nil
}

func buildRouterTree(menus []menu.SysMenu, parentID int64) []RouterResponse {
	var routers []RouterResponse
	for _, m := range menus {
		if m.ParentID != parentID {
			continue
		}

		router := RouterResponse{
			Name:      m.RouteName,
			Path:      m.Path,
			Hidden:    m.Visible == "1",
			Redirect:  "noRedirect",
			Component: deptr(m.Component),
			Query:     deptr(m.Query),
			Meta: RouterMetaResponse{
				Title:   m.MenuName,
				Icon:    m.Icon,
				NoCache: m.IsCache == 1,
			},
			Children: buildRouterTree(menus, m.MenuID),
		}

		if m.IsFrame == 0 {
			router.Meta.Link = m.Path
		}

		if router.Component == "" {
			router.Component = "Layout"
		}

		routers = append(routers, router)
	}

	sort.Slice(routers, func(i, j int) bool {
		return routers[i].Name < routers[j].Name
	})

	return routers
}

func deptr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *Service) ListOnlineUsers(ctx context.Context) ([]OnlineUserResponse, error) {
	loginUsers, err := s.sessionService.ListLoginUsers(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]OnlineUserResponse, 0, len(loginUsers))
	for _, loginUser := range loginUsers {
		if loginUser == nil || loginUser.User == nil {
			continue
		}

		items = append(items, OnlineUserResponse{
			TokenID:    loginUser.Token,
			UserID:     loginUser.UserID,
			UserName:   loginUser.User.UserName,
			NickName:   loginUser.User.NickName,
			DeptID:     loginUser.DeptID,
			IPAddr:     loginUser.IPAddr,
			LoginIP:    loginUser.LoginIP,
			Browser:    loginUser.Browser,
			OS:         loginUser.OS,
			LoginTime:  loginUser.LoginTime,
			ExpireTime: loginUser.ExpireTime,
		})
	}

	return items, nil
}

func (s *Service) ForceLogout(ctx context.Context, token string) error {
	return s.sessionService.DeleteLoginUser(ctx, token)
}

func (s *Service) verifyCaptcha(ctx context.Context, uuid, code string) error {
	if uuid == "" || code == "" {
		return ErrCaptchaInvalid
	}

	expected, err := s.redis.GetDel(ctx, captchaCacheKey(uuid)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCaptchaInvalid
		}

		return err
	}

	if strings.TrimSpace(expected) != strings.TrimSpace(code) {
		return ErrCaptchaInvalid
	}

	return nil
}

func (s *Service) loginUser(ctx context.Context, req LoginRequest, authUser *user.SysUser, meta LoginMeta) (*LoginResponse, error) {
	if authUser.Status != "0" {
		return nil, ErrAccountDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authUser.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := s.now()
	loginToken := s.tokenStore.NewToken()
	loginUser := &security.LoginUser{
		UserID:      authUser.UserID,
		DeptID:      authUser.DeptID,
		Token:       loginToken,
		LoginTime:   now.UnixMilli(),
		ExpireTime:  jwtutil.ExpiresAt(now).UnixMilli(),
		IPAddr:      meta.IPAddr,
		LoginIP:     meta.IPAddr,
		Browser:     meta.Browser,
		OS:          meta.OS,
		Roles:       append([]string(nil), meta.Roles...),
		Permissions: append([]string(nil), meta.Permissions...),
		User: &security.UserInfo{
			UserID:      authUser.UserID,
			DeptID:      authUser.DeptID,
			UserName:    authUser.UserName,
			NickName:    authUser.NickName,
			UserType:    authUser.UserType,
			Email:       authUser.Email,
			Phonenumber: authUser.Phonenumber,
			Sex:         authUser.Sex,
			Avatar:      authUser.Avatar,
			Status:      authUser.Status,
			DelFlag:     authUser.DelFlag,
			Remark:      authUser.Remark,
		},
	}
	if err := s.tokenStore.Save(ctx, loginUser); err != nil {
		return nil, err
	}

	token, err := jwtutil.Sign(authUser.UserID, authUser.UserName, loginToken, now)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		UserID:    authUser.UserID,
		UserName:  authUser.UserName,
		NickName:  authUser.NickName,
		DeptID:    authUser.DeptID,
		ExpiresAt: loginUser.ExpireTime,
	}, nil
}

type LoginMeta struct {
	IPAddr      string
	Browser     string
	OS          string
	Roles       []string
	Permissions []string
}

func BuildLoginMeta(req *http.Request, permissions []string) LoginMeta {
	userAgent := req.UserAgent()

	return LoginMeta{
		IPAddr:      clientIP(req),
		Browser:     detectBrowser(userAgent),
		OS:          detectOS(userAgent),
		Permissions: permissions,
	}
}

func defaultRoles(roles []string) []string {
	if len(roles) == 0 {
		return []string{"default"}
	}

	return append([]string(nil), roles...)
}

func defaultPermissions(permissions []string) []string {
	if len(permissions) == 0 {
		return []string{"*:*:*"}
	}

	return append([]string(nil), permissions...)
}

func clientIP(req *http.Request) string {
	if req == nil {
		return ""
	}

	for _, header := range []string{"X-Forwarded-For", "X-Real-IP"} {
		value := strings.TrimSpace(req.Header.Get(header))
		if value == "" {
			continue
		}

		if header == "X-Forwarded-For" {
			parts := strings.Split(value, ",")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0])
			}
		}

		return value
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr))
	if err == nil {
		return host
	}

	return strings.TrimSpace(req.RemoteAddr)
}

func detectBrowser(userAgent string) string {
	switch {
	case strings.Contains(userAgent, "Edg/"):
		return "Edge"
	case strings.Contains(userAgent, "Chrome/"):
		return "Chrome"
	case strings.Contains(userAgent, "Firefox/"):
		return "Firefox"
	case strings.Contains(userAgent, "Safari/") && strings.Contains(userAgent, "Version/"):
		return "Safari"
	case strings.Contains(userAgent, "MSIE") || strings.Contains(userAgent, "Trident/"):
		return "IE"
	default:
		return "Unknown"
	}
}

func detectOS(userAgent string) string {
	switch {
	case strings.Contains(userAgent, "Windows"):
		return "Windows"
	case strings.Contains(userAgent, "Mac OS X"):
		return "macOS"
	case strings.Contains(userAgent, "Android"):
		return "Android"
	case strings.Contains(userAgent, "iPhone"), strings.Contains(userAgent, "iPad"), strings.Contains(userAgent, "iOS"):
		return "iOS"
	case strings.Contains(userAgent, "Linux"):
		return "Linux"
	default:
		return "Unknown"
	}
}

func captchaCacheKey(uuid string) string {
	return "captcha_codes:" + uuid
}
