// Package handler provides common functionality for HTTP handlers in the application.
package handler

import (
	"transaction_demo/app/apperr"

	"github.com/gin-gonic/gin"
)

// BaseHandler provides common functionality for HTTP handlers in the application.
type BaseHandler struct{}

// RenderResponse renders a successful HTTP response with the provided status code and data.
// It standardizes the JSON response format across the application.
//
// Parameters:
//   - ctx: The Gin context for the HTTP request
//   - status: HTTP status code to return
//   - data: The response payload to be serialized as JSON
//   - meta: Additional metadata (currently unused but available for future extension)
func (h *BaseHandler) RenderResponse(
	ctx *gin.Context,
	status int,
	data interface{},
	meta interface{},
) {
	ctx.JSON(status, data)
}

// RenderError handles error responses by converting errors to a standardized format.
// It ensures that all errors are properly formatted as AppError instances with
// appropriate HTTP status codes and error details.
//
// If the provided error is not already an AppError, it wraps it in a generic
// internal server error. Otherwise, it uses the error as-is.
//
// Parameters:
//   - ctx: The Gin context for the HTTP request
//   - err: The error to be rendered in the response
func (h *BaseHandler) RenderError(
	ctx *gin.Context,
	err error,
) {
	var appErr apperr.AppError
	if _, ok := err.(apperr.AppError); !ok {
		appErr = apperr.NewAppError("undefined_error", apperr.ErrTypeInternalServer).WithError(err)
	} else {
		appErr = err.(apperr.AppError)
	}

	ctx.JSON(appErr.Status, appErr)
}
