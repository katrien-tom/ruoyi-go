package auth

import (
	"context"
	"errors"
	"image/color"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/banyejiu/ruoyi-go/internal/modules/user"
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
	userRepository *user.Repository
	now            func() time.Time
}

func NewService(redisClient *redis.Client, userRepository *user.Repository) *Service {
	return &Service{
		redis:          redisClient,
		userRepository: userRepository,
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

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
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

	authUser, err := s.userRepository.FindByUserName(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	return s.loginUser(req, authUser)
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

func (s *Service) loginUser(req LoginRequest, authUser *user.SysUser) (*LoginResponse, error) {
	if authUser.Status != "0" {
		return nil, ErrAccountDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authUser.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := jwtutil.Sign(authUser.UserID, authUser.UserName, s.now())
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		UserID:    authUser.UserID,
		UserName:  authUser.UserName,
		NickName:  authUser.NickName,
	}, nil
}

func captchaCacheKey(uuid string) string {
	return "captcha_codes:" + uuid
}
