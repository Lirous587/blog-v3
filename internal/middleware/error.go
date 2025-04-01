package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		// 处理错误
		if len(ctx.Errors) > 0 {
			for _, e := range ctx.Errors {
				// 记录详细错误日志
				log.Printf("Error: %+v\n", e.Err)
			}
		}
	}
}
