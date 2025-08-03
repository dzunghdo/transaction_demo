package apperr

// Define the error constants for the application
var (
	ErrInvalidInput      = NewAppError("INVALID_INPUT", ErrTypeBadRequest)
	ErrInsufficientFunds = NewAppError("INSUFFICIENT_FUNDS", ErrTypeBadRequest)
	ErrNotFound          = NewAppError("NOT_FOUND", ErrTypeNotFound)
	ErrAlreadyExists     = NewAppError("ALREADY_EXISTS", ErrTypeAlreadyExists)
	ErrResourceBusy      = NewAppError("RESOURCE_BUSY", ErrTypeBadRequest)
	ErrInternalServer    = NewAppError("INTERNAL_SERVER_ERROR", ErrTypeInternalServer)
)
