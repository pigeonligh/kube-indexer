package server

import "github.com/gin-gonic/gin"

const (
	errorNotInit        = "data not init"
	errorInvalidBody    = "invalid body"
	errorKindNotFound   = "kind not found"
	errorObjectNotFound = "object not found"
)

func (s *restfulServer) responseError(ctx *gin.Context, code int, err string) {
	ctx.JSON(code, gin.H{
		"error": err,
	})
}
