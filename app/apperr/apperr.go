// Package apperr provides a way to handle application errors with specific types and HTTP status codes.
package apperr

type ErrorType string

// ErrorType represents the type of error
const (
	ErrTypeBadRequest     ErrorType = "bad_request"     // 400
	ErrTypeNotFound       ErrorType = "not_found"       // 404
	ErrTypeAlreadyExists  ErrorType = "already_exists"  // 409
	ErrTypeInternalServer ErrorType = "internal_server" // 500
)

// mapErrTypeStatus maps the error type to the corresponding HTTP status code
var mapErrTypeStatus = map[ErrorType]int{
	ErrTypeBadRequest:     400, // Bad Request
	ErrTypeNotFound:       404, // Not Found
	ErrTypeAlreadyExists:  409, // Conflict
	ErrTypeInternalServer: 500, // Internal Server Error
}

type AppError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func NewAppError(code string, errType ErrorType) *AppError {
	statusCode, ok := mapErrTypeStatus[errType]
	if !ok {
		statusCode = 500 // Default to Internal Server Error if type is unknown
	}
	return &AppError{
		Status: statusCode,
		Code:   code,
	}
}

func (e AppError) Error() string {
	return e.Err.Error()
}

func (e AppError) WithMessage(message string) AppError {
	e.Message = message
	return e
}

func (e AppError) WithError(err error) AppError {
	e.Err = err
	return e
}
