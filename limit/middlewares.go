package limit

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"strconv"
)

func RequestLimitMiddleware(limiter *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestNum := limiter.GetRequestNum(c.ClientIP())
		timeDuration := limiter.GetPipelineTTL(c.ClientIP())
		limit := limiter.GetLimit()
		period := limiter.GetPeriod()

		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(limit - requestNum, 10))
		c.Writer.Header().Set("X-RateLimit-Reset", timeDuration.String())
		if !limiter.Allow(c.ClientIP(), limit, period * time.Second) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{})
			return
		}
		c.Next()
	}
}