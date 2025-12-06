package server

import "github.com/gin-gonic/gin"

type ServerManager interface {
	Start() error
	wsHandler(c *gin.Context)
}
