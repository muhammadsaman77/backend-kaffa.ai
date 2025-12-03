package middlewares

import (
	"time"

	"backend-kaffa.ai/configs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	configs.Log.Info("request",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", c.Writer.Status()),
		zap.Duration("duration", time.Since(start)),
		zap.String("ip", c.ClientIP()),
	)
}
