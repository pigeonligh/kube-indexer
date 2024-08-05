package server

import "github.com/gin-gonic/gin"

const (
	errorNotInit        = "data not init"
	errorInvalidBody    = "invalid body"
	errorKindNotFound   = "kind not found"
	errorObjectNotFound = "object not found"
)

func (rr *restfulRegisterer) responseError(ctx *gin.Context, code int, err string) {
	ctx.JSON(code, gin.H{
		"error": err,
	})
}
