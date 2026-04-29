package response

// -----------------------------
// 响应码设计
// -----------------------------
type Code int

const (
	// 成功
	SuccessCode Code = 0

	// 系统通用错误 1000-1999
	ParamError    Code = 1001
	InternalError Code = 1002
	NotFound      Code = 1003

	// 用户/认证 2000-2999
	Unauthorized    Code = 2001
	LoginFailed     Code = 2002
	AccountDisabled Code = 2003
	CaptchaInvalid  Code = 2004

	// 权限 3000-3999
	PermissionDenied Code = 3001

	// 业务逻辑 4000-4999
	DataAlreadyExist Code = 4001
	DataInvalid      Code = 4002

	// 外部依赖 5000-5999
	ThirdPartyFailed Code = 5001

	// 数据库/缓存 6000-6999
	DBWriteFailed Code = 6001
	CacheMiss     Code = 6002

	// 文件/上传下载 7000-7999
	FileTooLarge   Code = 7001
	FileUploadFail Code = 7002
)
