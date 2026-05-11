package security

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/banyejiu/ruoyi-go/pkg/jwtutil"
	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenInvalid  = errors.New("token invalid")
	ErrTokenNotFound = errors.New("token not found")
)

const sessionRefreshWindow = 120 * time.Minute

type SessionService struct {
	tokenStore *TokenStore
	now        func() time.Time
}

func NewSessionService(tokenStore *TokenStore) *SessionService {
	return &SessionService{
		tokenStore: tokenStore,
		now:        time.Now,
	}
}

func (s *SessionService) GetLoginUser(ctx context.Context, token string) (*LoginUser, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, ErrTokenInvalid
	}

	loginUser, err := s.tokenStore.Get(ctx, token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrTokenNotFound
		}

		return nil, err
	}

	return loginUser, nil
}

func (s *SessionService) DeleteLoginUser(ctx context.Context, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return ErrTokenInvalid
	}

	if _, err := s.GetLoginUser(ctx, token); err != nil {
		return err
	}

	return s.tokenStore.Delete(ctx, token)
}

func (s *SessionService) VerifyAndRefresh(ctx context.Context, token string) (*LoginUser, bool, error) {
	loginUser, err := s.GetLoginUser(ctx, token)
	if err != nil {
		return nil, false, err
	}

	remaining := time.Until(loginUser.ExpireAt())
	if s.now != nil {
		remaining = loginUser.ExpireAt().Sub(s.now())
	}

	if remaining > sessionRefreshWindow {
		return loginUser, false, nil
	}

	now := time.Now()
	if s.now != nil {
		now = s.now()
	}
	loginUser.ExpireTime = jwtutil.ExpiresAt(now).UnixMilli()

	if err := s.tokenStore.Save(ctx, loginUser); err != nil {
		return nil, false, err
	}

	return loginUser, true, nil
}

func (s *SessionService) ListLoginUsers(ctx context.Context) ([]*LoginUser, error) {
	keys, err := s.tokenStore.Keys(ctx)
	if err != nil {
		return nil, err
	}

	loginUsers := make([]*LoginUser, 0, len(keys))
	for _, key := range keys {
		loginUser, err := s.tokenStore.Get(ctx, tokenFromCacheKey(key))
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}

			return nil, err
		}

		loginUsers = append(loginUsers, loginUser)
	}

	return loginUsers, nil
}
