package security

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const tokenCachePrefix = "login_tokens:"

type TokenStore struct {
	redis *redis.Client
	now   func() time.Time
}

func NewTokenStore(redisClient *redis.Client) *TokenStore {
	return &TokenStore{
		redis: redisClient,
		now:   time.Now,
	}
}

func (s *TokenStore) SetNowFunc(now func() time.Time) {
	if s == nil || now == nil {
		return
	}

	s.now = now
}

func (s *TokenStore) NewToken() string {
	return uuid.NewString()
}

func (s *TokenStore) Save(ctx context.Context, loginUser *LoginUser) error {
	if s == nil || s.redis == nil || loginUser == nil || loginUser.Token == "" {
		return errors.New("invalid token store state")
	}

	payload, err := json.Marshal(loginUser)
	if err != nil {
		return err
	}

	now := time.Now()
	if s.now != nil {
		now = s.now()
	}

	ttl := loginUser.ExpireAt().Sub(now)
	if ttl <= 0 {
		return errors.New("login session expired")
	}

	return s.redis.Set(ctx, tokenCacheKey(loginUser.Token), payload, ttl).Err()
}

func (s *TokenStore) Get(ctx context.Context, token string) (*LoginUser, error) {
	if s == nil || s.redis == nil {
		return nil, errors.New("token store unavailable")
	}

	payload, err := s.redis.Get(ctx, tokenCacheKey(token)).Bytes()
	if err != nil {
		return nil, err
	}

	var loginUser LoginUser
	if err := json.Unmarshal(payload, &loginUser); err != nil {
		return nil, err
	}

	return &loginUser, nil
}

func (s *TokenStore) Delete(ctx context.Context, token string) error {
	if s == nil || s.redis == nil {
		return errors.New("token store unavailable")
	}

	return s.redis.Del(ctx, tokenCacheKey(token)).Err()
}

func (s *TokenStore) Keys(ctx context.Context) ([]string, error) {
	if s == nil || s.redis == nil {
		return nil, errors.New("token store unavailable")
	}

	keys, err := s.redis.Keys(ctx, tokenCachePrefix+"*").Result()
	if err != nil {
		return nil, err
	}

	sort.Strings(keys)
	return keys, nil
}

func tokenCacheKey(token string) string {
	return tokenCachePrefix + token
}

func tokenFromCacheKey(key string) string {
	return strings.TrimPrefix(key, tokenCachePrefix)
}
