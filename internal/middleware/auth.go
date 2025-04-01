package middleware

// import (
// 	"errors"
// 	"github.com/gin-gonic/gin"
// 	"net/http"
// 	"strings"
// )

// const (
// 	needLogin        = "需要登录"
// 	needRefreshToken = "需要refreshToken"
// )

// func validateTokenFormat(c *gin.Context) (token string, err error) {
// 	authHeader := c.Request.Header.Get("Authorization")
// 	if authHeader == "" {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"msg": "token为空",
// 		})
// 		err = errors.New("token为空")
// 		return
// 	}
// 	parts := strings.SplitN(authHeader, " ", 2)
// 	if !(len(parts) == 2 && parts[0] == "Bearer") {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"msg": "token格式错误",
// 		})
// 		err = errors.New("token格式错误")
// 		return
// 	}
// 	token = parts[1]
// 	return
// }

// func validateRefreshTokenFormat(c *gin.Context) (refreshToken string, err error) {
// 	refreshToken = c.Request.Header.Get("Refresh-Token")
// 	if refreshToken == "" {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"msg": needRefreshToken,
// 		})
// 		err = errors.New(needRefreshToken)
// 		return
// 	}
// 	return
// }
