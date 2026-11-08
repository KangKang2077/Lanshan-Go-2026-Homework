package response

import (
	"errors"
	"net/http"
	"week_06_web/app_err"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: "OK",
		Message: "success",
		Data:    data})
	return
}

func Error(c *gin.Context, err error) {
	var appErr *app_err.Error
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, Response{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
	})
}
