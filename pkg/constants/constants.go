package constants

// Database field values
const (
	StatusNormal   = "0"
	StatusDisabled = "1"

	DelFlagNormal  = "0"
	DelFlagDeleted = "2"

	VisibleShow   = "0"
	VisibleHidden = "1"

	IsCacheNo  = 0
	IsCacheYes = 1

	IsFrameNo  = 0
	IsFrameYes = 1
)

// Redis key prefixes
const (
	RedisPrefixCaptcha    = "captcha_codes:"
	RedisPrefixLoginToken = "login_tokens:"
)

// CaptchaCacheKey returns the Redis key for captcha storage
func CaptchaCacheKey(uuid string) string {
	return RedisPrefixCaptcha + uuid
}

// LoginTokenCacheKey returns the Redis key for login token storage
func LoginTokenCacheKey(token string) string {
	return RedisPrefixLoginToken + token
}
