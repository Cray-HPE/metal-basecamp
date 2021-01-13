package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Adapted from: https://stackoverflow.com/questions/48224908/better-error-handling

// AppError for app specifics.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf(e.Message)
}

// AppErrorReporter reports errors.
func AppErrorReporter() gin.HandlerFunc {
	return appErrorReporterT(gin.ErrorTypeAny)
}

func appErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)

		if len(detectedErrors) > 0 {
			log.Println("Handle APP error")

			for i := range detectedErrors {
				log.Println(i)
			}

			err := detectedErrors[0].Err
			log.Println(err)
			var parsedError *AppError
			switch err.(type) {
			case *AppError:
				parsedError = err.(*AppError)
			default:
				parsedError = &AppError{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				}
			}
			c.JSON(parsedError.Code, parsedError)
			c.Abort()
			return
		}

	}
}
