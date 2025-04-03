package response

var (
	errCodeMsgMap = map[code]string{
		// 参数错误
		CodeParamInvalid: "参数无效",
		//CodeParamFormat:  "参数格式错误",
		//CodeParamMissing: "缺少必要参数",

		// 认证错误
		CodeAuthFailed:     "认证失败",
		CodeTokenInvalid:   "无效的令牌",
		CodeTokenExpired:   "令牌已过期",
		CodeRefreshInvalid: "无效的refreshToken",
		// 资源错误
		CodeResourceNotFound: "资源未找到",
		CodeResourceExists:   "资源已存在",

		// 服务端错误消息映射
		CodeAdminExist:     "管理员已初始化",
		CodeServerError:    "服务器错误",
		CodeDatabaseError:  "数据库错误",
		CodeInternalError:  "内部服务器错误",
		CodeRecordNotFound: "该记录不存在",
	}
)
