package response

type code int

const (
	// 成功码
	codeSuccess code = 2000

	// 参数错误 4000-4099
	CodeParamInvalid code = 4000 + iota
	CodeParamFormat
	CodeParamMissing

	// 认证授权错误 4100-4199
	CodeAuthFailed code = 4100 + iota
	CodeUnauthorized
	CodeTokenInvalid
	CodeTokenExpired

	// 资源错误 4200-4299
	CodeResourceNotFound code = 4200 + iota
	CodeResourceExists
)

const (
	// 服务器错误 5000-5999
	CodeServerError code = 5000 + iota
	CodeDatabaseError
	CodeInternalError
)
