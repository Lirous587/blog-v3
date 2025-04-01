package response

var (
	// 客户端错误消息映射
	clientErrCodeMsgMap = map[code]string{
		// 参数错误
		CodeParamInvalid: "参数无效",
		CodeParamFormat:  "参数格式错误",
		CodeParamMissing: "缺少必要参数",

		// 认证错误
		CodeAuthFailed:   "认证失败",
		CodeUnauthorized: "未授权",
		CodeTokenInvalid: "无效的令牌",
		CodeTokenExpired: "令牌已过期",

		// 资源错误
		CodeResourceNotFound: "资源未找到",
		CodeResourceExists:   "资源已存在",
	}

	// 服务端错误消息映射
	serverErrCodeMsgMap = map[code]string{
		CodeServerError:   "服务器错误",
		CodeDatabaseError: "数据库错误",
		CodeInternalError: "内部服务器错误",
	}
)
