package response

// 业务模块错误

// admin相关错误 10000-10099
const (
	CodeAdminExist code = 10000 + iota
	CodeAdminXX
)

// label模块错误 100100-100199
const (
	CodeLabelNotFound      code = 100100 + iota
	CodeLabelNameDuplicate      // 标签name重复
)

// essay模块错误 100200-100299
const (
	CodeEssayNotFound code = 10200 + iota
	CodeEssayXX
)

// maxim模块错误 100300-100399
const (
	CodeMaximNotFound code = 10300 + iota
	CodeMaximXX
)

// 友链模块错误 100400-100499
const (
	CodeFriendLinkNotFound     code = 10400 + iota
	CodeFriendLinkUrlDuplicate      // 友链URL重复
)
