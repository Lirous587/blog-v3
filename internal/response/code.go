package response

type code int

type Code = code

// 成功码
const codeSuccess code = 2000

// 参数错误 4000-4099
const (
	CodeParamInvalid code = 4000 + iota
	CodeParamFormat
	CodeParamMissing
)

// 认证授权错误 4100-4199
const (
	CodeAuthFailed code = 4100 + iota
	CodeTokenInvalid
	CodeTokenExpired
	CodeRefreshInvalid
	CodeRefreshExpired
)

// 资源错误 4200-4299
const (
	CodeResourceNotFound code = 4200 + iota
	CodeResourceExists
)

// 服务器错误 5000-5999
const (
	CodeServerError code = 5000 + iota
	CodeDatabaseError
	CodeInternalError
	CodeAdminExist
)
