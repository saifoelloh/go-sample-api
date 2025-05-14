package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saifoelloh/ranger/pkg/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			if extErr, ok := err.Err.(*errors.Extension); ok {
				// Log full error with location/scope for engineer
				fmt.Printf("[ERROR] %s/%s - %s\n", extErr.Scope, extErr.Location, extErr.Error())
				if extErr.Detail != nil {
					fmt.Printf("Detail: %+v\n", extErr.Detail)
				}

				// Send structured error to client
				c.AbortWithStatusJSON(extErr.StatusCode, gin.H{
					"error": gin.H{
						"message":        extErr.Message,
						"locale_message": extErr.LocaleMessage,
						"scope":          extErr.Scope,
						"location":       extErr.Location,
						"error_code":     extErr.ErrorCode,
						"status_code":    extErr.StatusCode,
						"detail":         extErr.Detail,
					},
				})
				return
			}

			// Fallback for non-IExtension errors
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message":     "An unexpected error occurred",
					"error_code":  "internal/server-error",
					"status_code": http.StatusInternalServerError,
				},
			})
		}
	}
}
