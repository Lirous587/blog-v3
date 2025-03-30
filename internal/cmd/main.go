package cmd

import (
	"blog/pkg/httpserver"
	"blog/pkg/validator"
	"github.com/gin-gonic/gin"
)

type TestJSON struct {
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"required"`
	Phone string `json:"phone" binding:"mobile_cn"`
}

func Main() {
	// 创建服务器
	s := httpserver.New(8080)
	r := s.Router

	s2 := httpserver.New(8081)

	r2 := s2.Router

	r.POST("/test", func(c *gin.Context) {
		var test TestJSON
		if err := c.ShouldBindJSON(&test); err != nil {
			// 翻译错误
			errMsg := validator.TranslateError(err, "zh")
			c.AbortWithStatusJSON(400, gin.H{
				"msg":   "参数错误",
				"error": errMsg,
			})
			return
		}
		c.JSON(200, gin.H{
			"msg": "hello",
		})
	})
	r2.POST("/test", func(c *gin.Context) {
		var test TestJSON
		if err := c.ShouldBindJSON(&test); err != nil {
			// 翻译错误
			errMsg := validator.TranslateError(err, "zh")
			c.AbortWithStatusJSON(400, gin.H{
				"msg":   "参数错误",
				"error": errMsg,
			})
			return
		}
		c.JSON(200, gin.H{
			"msg": "hello",
		})
	})

	go s.Run()

	s2.Run()
}
