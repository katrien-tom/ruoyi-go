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

func (s *Service) GetCaptcha(ctx context.Context) (*CaptchaVO, error) {
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

	return &CaptchaVO{
		CaptchaEnabled: true,
		Img:            item.EncodeB64string(),
		UUID:           id,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginDTO) (*LoginVO, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		return nil, ErrInvalidCredentials
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

func (s *Service) loginUser(req LoginDTO, authUser *user.SysUser) (*LoginVO, error) {
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

	return &LoginVO{
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
