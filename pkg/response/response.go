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

func ErrorParameterInvalid(ctx *gin.Context, err error) {
	if err == nil {
		Success(ctx)
		return
	}

	res := response{
		Code:    codeParamInvalid,
		Message: "参数无效",
	}
	lang := validator.GetTranslateLang(ctx)
	transErr := validator.TranslateError(err, lang)
	res.Data = transErr.Error()

	ctx.Error(err)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
}

// Error 宽松错误处理 - 接受普通error和错误码
//func Error(ctx *gin.Context, code code, err error) {
//	res := response{
//		Code: code,
//	}
//	msg, ok := errCodeMsgMap[code]
//	if ok {
//		res.Message = msg
//	} else {
//		res.Code = codeUnKnowError
//		res.Message = "未知错误"
//	}
//
//	if err != nil {
//		ctx.Error(err)
//	}
//	ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
//}

func Error(ctx *gin.Context, appErr *AppError) {
	if appErr.Err == nil {
		Success(ctx)
		return
	}
	res := response{
		Code: appErr.Code,
	}
	msg, ok := errCodeMsgMap[appErr.Code]
	if ok {
		res.Message = msg
	} else {
		res.Code = codeUnKnowError
		res.Message = "未知错误"
	}

	ctx.Error(appErr.Err)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
}
