package response

import (
	"net/http"
	"order-service/constants"
	errorConst "order-service/constants/error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Token   *string     `json:"token,omitempty"`
}

type ParamHTTPResp struct {
	Code    int
	Error   error
	Message *string
	Gin     *gin.Context
	Data    interface{}
	Token   *string
}

func HttpResponse(param ParamHTTPResp) {
	if param.Error == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
			Token:   param.Token,
		})
		return
	}

	message := errorConst.ErrInternalServer.Error()
	if param.Message != nil {
		message = *param.Message
	} else if param.Error != nil {
		if errorConst.ErrMapping(param.Error) {
			message = param.Error.Error()
		}
	}

	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    param.Data,
	})

	return
}
