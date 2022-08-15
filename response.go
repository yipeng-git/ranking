package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type baseResponse struct {
	Result    interface{} `json:"result"`
	ErrorMsg  string      `json:"errorMsg"`
	ErrorCode int         `json:"errorCode"`
}

func responseOk(ctx *gin.Context, ret interface{}) {
	ctx.JSON(http.StatusOK, baseResponse{
		ErrorCode: 0,
		ErrorMsg:  "ok",
		Result:    ret,
	})
}

func responseErr(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, baseResponse{
		ErrorCode: 1,
		ErrorMsg:  msg,
	})
}
