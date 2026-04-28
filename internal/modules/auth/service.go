package auth

import (
	"context"
	"image/color"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

const captchaTTL = 5 * time.Minute

type Service struct {
	redis *redis.Client
}

type CaptchaResponse struct {
	CaptchaEnabled bool   `json:"captchaEnabled"`
	Img            string `json:"img"`
	UUID           string `json:"uuid"`
}

func NewService(redisClient *redis.Client) *Service {
	return &Service{redis: redisClient}
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

func captchaCacheKey(uuid string) string {
	return "captcha_codes:" + uuid
}
