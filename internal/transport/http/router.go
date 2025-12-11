package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		logger.Info("health check called")
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	return r
}
