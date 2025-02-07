package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"template/internal/common"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func Log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := common.Log

		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()

		fields := map[string]any{
			"method":    ctx.Request.Method,
			"host":      ctx.Request.Host,
			"uri":       ctx.Request.RequestURI,
			"status":    ctx.Writer.Status(),
			"client_ip": ctx.ClientIP(),
			"latency":   endTime.Sub(startTime).String(),
		}

		lastErr := ctx.Errors.Last()
		var ve validator.ValidationErrors
		var he common.ClientError

		if lastErr != nil && !errors.As(lastErr, &he) && !errors.As(lastErr, &ve) {
			log.WithFields(fields).Error(lastErr)
			return
		}

		log.WithFields(fields).Infof("REQUEST %s %s NO SERVER ERROR", ctx.Request.Method, ctx.Request.RequestURI)
	}
}

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if gin.Mode() == gin.TestMode {
			ctx.Set("user_id", 1)
			ctx.Next()
			return
		}

		header := strings.Split(ctx.Request.Header.Get("Authorization"), " ")
		if len(header) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "login to access this resource"})
			return
		}

		user_id, err := common.JwtGetID(header[1])

		if errors.As(err, &common.AuthenticationError{}) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "re-login to access this resource"})
			return
		}

		if err != nil {
			ctx.Error(err)
		}

		ctx.Set("user_id", user_id)
		ctx.Next()
	}
}

func Errors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			var je *json.UnmarshalTypeError
			var ve validator.ValidationErrors
			var he common.ClientError

			switch {
			case errors.As(ctx.Errors[0], &je):
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					gin.H{"message": "invalid JSON format"})

			case errors.As(ctx.Errors[0], &ve):
				errs := make(map[string]string, 0)
				for _, err := range ve {
					errs[err.StructField()] = fmt.Sprintf("cannot satistfy %v tag", err.Tag())
				}
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					gin.H{"message": "invalid body format", "errors": errs})

			case errors.As(ctx.Errors[0], &he):
				ctx.AbortWithStatusJSON(he.HTTPStatus(),
					gin.H{"message": fmt.Sprintf("failed to process request: %s", he.Error())})

			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError,
					gin.H{"message": "internal server error. Please contact admin"})
			}
		}
	}
}
