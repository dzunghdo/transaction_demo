// Package middleware provides HTTP middleware functions for the Gin web framework.
// This middleware specifically handles panic recovery to ensure the application doesn't crash
// when unexpected errors occur during request processing.
package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// Recover creates a middleware function that recovers from panics and logs errors.
// It is designed to be used with the Gin web framework to ensure that the application
// continues to run even when an unexpected error occurs during request handling.
// How it works:
// 1. It uses a deferred function to catch any panic that occurs during the request processing.
// 2. If a panic occurs, it checks if the error is related to a broken pipe (i.e., the client has disconnected).
// 3. If the error is a broken pipe, it logs the error and aborts the request without sending a response.
// 4. For all other panics, it logs the error and returns a 500 Internal Server Error response.
//
// Returns a gin.HandlerFunc that can be used as middleware in the Gin router.
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check if the error is due to a broken connection (client disconnected)
				if isBrokenPipe(err) {
					fmt.Println("broken pipe error, stack: ", string(debug.Stack()), "error", err)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error))
					c.Abort()
					return
				}
				// Log all other panics as errors with full stack trace for debugging
				fmt.Println("panic, stack: ", string(debug.Stack()), "error", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// Continue to the next middleware/handler in the chain
		c.Next()
	}
}

// isBrokenPipe checks if the error is a broken pipe or connection reset by peer error.
// These errors typically occur when the client disconnects before the server finishes
// processing the request (e.g., user closes browser, network timeout, etc.).
//
// Returns:
//   - true if the error indicates a broken connection, false otherwise
func isBrokenPipe(err interface{}) bool {
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				return true
			}
		}
	}
	return false
}
