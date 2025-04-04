package utils

func ComputeOffset(page int, size int) int {
	if page < 1 {
		page = 1
	}
	if size < 0 {
		size = 0
	}
	return (page - 1) * size
}

func ComputePages(count int64, pageSize int) int {
	if count <= 0 {
		return 0 // 或返回1，取决于业务需求
	}
	if pageSize <= 0 {
		return 0 // 避免除零错误
	}
	sizeInt64 := int64(pageSize)
	return int((count + sizeInt64 - 1) / sizeInt64)
}

// BuildLikeQuery 构建用于SQL LIKE查询的模式字符串
// 参数:
//   - keyword: 要查询的关键词，如为空则返回"%"（匹配所有）
//   - matchType: 可选的匹配类型，支持以下值:
//   - "start": 前缀匹配，返回"keyword%"
//   - "end": 后缀匹配，返回"%keyword"
//   - "exact": 精确匹配，返回"keyword"
//   - 其他或不提供: 包含匹配(默认)，返回"%keyword%"
//
// 返回值:
//
//	格式化后的LIKE查询模式字符串
//
// 示例:
//
//	BuildLikeQuery("test")         // 返回 "%test%"
//	BuildLikeQuery("test", "start") // 返回 "test%"
//	BuildLikeQuery("", "exact")    // 返回 "%"
func BuildLikeQuery(keyword string, matchType ...string) string {
	if keyword == "" {
		return "%"
	}

	var match string
	if len(matchType) > 0 {
		match = matchType[0]
	}

	switch match {
	case "start":
		return keyword + "%"
	case "end":
		return "%" + keyword
	case "exact":
		return keyword
	default: // "contains"
		return "%" + keyword + "%"
	}
}
