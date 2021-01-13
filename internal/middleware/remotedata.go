package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RemoteMetaData middleware to inject all data found for MAC within local files
func RemoteMetaData() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Error(&AppError{Code: http.StatusNotImplemented, Message: fmt.Sprintf("TODO")})
		c.Abort()
	}
}
