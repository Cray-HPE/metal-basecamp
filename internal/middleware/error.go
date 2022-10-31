/*
 MIT License
 
 (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 
 Permission is hereby granted, free of charge, to any person obtaining a
 copy of this software and associated documentation files (the "Software"),
 to deal in the Software without restriction, including without limitation
 the rights to use, copy, modify, merge, publish, distribute, sublicense,
 and/or sell copies of the Software, and to permit persons to whom the
 Software is furnished to do so, subject to the following conditions:
 
 The above copyright notice and this permission notice shall be included
 in all copies or substantial portions of the Software.
 
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 OTHER DEALINGS IN THE SOFTWARE.
*/

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
