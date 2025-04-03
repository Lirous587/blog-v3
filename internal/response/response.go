package response

import (
	"blog/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code    code        `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(ctx *gin.Context, data ...any) {
	res := response{
		Code:    codeSuccess,
		Message: "请求成功",
	}
	if len(data) > 0 {
		res.Data = data[0]
	}
	ctx.JSON(http.StatusOK, res)
}

func ClientError(ctx *gin.Context, code code, err error) {
	res := response{
		Code: code,
	}

	msg, ok := clientErrCodeMsgMap[code]
	if ok {
		res.Message = msg
	} else {
		res.Message = "未知客户端错误"
	}

	if code == CodeParamInvalid {
		lang := validator.GetTranslateLang(ctx)
		transErr := validator.TranslateError(err, lang)
		res.Data = transErr.Error()
	}

	if err != nil {
		ctx.Error(err)
	}
	ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
}

func ServerError(ctx *gin.Context, code code, err error) {
	res := response{
		Code: code,
	}
	msg, ok := serverErrCodeMsgMap[code]
	if ok {
		res.Message = msg
	} else {
		res.Message = "未知服务端错误"
	}
	ctx.Error(err)
	ctx.AbortWithStatusJSON(500, res)
}
