package response

var (
	errCodeMsgMap = map[code]string{

		// 服务端错误消息映射
		CodeServerError:      "服务器错误",
		CodeIllegalOperation: "非法操作",

		// 认证错误
		CodeAuthFailed:     "认证失败",
		CodeTokenInvalid:   "无效的令牌",
		CodeTokenExpired:   "令牌已过期",
		CodeRefreshInvalid: "无效的refreshToken",

		// admin模块
		CodeAdminExist: "管理员已初始化",
		// label模块
		CodeLabelNotFound:      "该标签不存在",
		CodeLabelNameDuplicate: "标签名重复",

		CodeFriendLinkNotFound:     "该友链不存在",
		CodeFriendLinkUrlDuplicate: "友链URL重复",
	}
)
