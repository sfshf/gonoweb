package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceID middleware
const (
	HeaderKey_TraceID  = "traceID"
	GinContext_TraceID = "traceID"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(HeaderKey_TraceID)
		if traceID == "" {
			uuid, _ := uuid.NewRandom()
			traceID = fmt.Sprintf("%s", uuid)
		}
		// 配置到请求上下文里
		c.Set(GinContext_TraceID, traceID)
		c.Next()
		c.Writer.Header().Add(HeaderKey_TraceID, traceID)
	}
}
