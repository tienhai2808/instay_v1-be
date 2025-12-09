package middleware

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RequestMiddleware struct {
	logger *zap.Logger
}

func NewRequestMiddleware(logger *zap.Logger) *RequestMiddleware {
	return &RequestMiddleware{logger}
}

func (m *RequestMiddleware) Recovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, err any) {
		var recErr error
		switch v := err.(type) {
		case error:
			recErr = v
		default:
			recErr = fmt.Errorf("%v", v)
		}

		stack := string(debug.Stack())
		m.logger.Error("panic recovered",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("stack", stack),
			zap.Error(recErr),
		)

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
			Message: "internal server error",
		})
	})
}

func (m *RequestMiddleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		if apiErr, ok := err.Err.(*common.APIError); ok {
			common.ToAPIResponse(c, apiErr.Status, apiErr.Message, nil)
			return
		}

		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
	}
}
