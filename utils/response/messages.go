package response

import (
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
)

const (
	// MessageRequestError is default response message for 400 BadRequest error
	MessageRequestError = "Your request is not valid."
	// MessageLockedError is default response message for 423 Locked error
	MessageLockedError = "Specified content is locked and not editable."
	// MessageNotFoundError is default response message for 404 NotFound error
	MessageNotFoundError = "Specified content was not exist."
	// MessageUnauthorizedError is default response message for 401 Unauthorized error
	MessageUnauthorizedError = "Probably your password incorrect."
	// MessageConflictedError is default response message for 409 Conflict error
	MessageConflictedError = "Specified content was already exists."
	// MessagePermissionError is default response message for 403 Forbidden error
	MessagePermissionError = "You don't have enough permission to do it."
	// MessageInternalError is default response message for 500 Internal error
	MessageInternalError = "Unfortunately, the server exploded."
)

// NewRequestError creates 400 BadRequest response
func NewRequestError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusBadRequest,
		Body: gen.GeneralMessageResponse{Message: MessageRequestError},
	}
}

// NewRequestErrorWithMessage creates 400 BadRequest response with using message
func NewRequestErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusBadRequest,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

// NewLockedError creates 423 Locked response
func NewLockedError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusLocked,
		Body: gen.GeneralMessageResponse{Message: MessageLockedError},
	}
}

// NewLockedErrorWithMessage creates 423 Locked response with using message
func NewLockedErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusLocked,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

// NewNotFoundError creates 404 NotFound response
func NewNotFoundError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusNotFound,
		Body: gen.GeneralMessageResponse{Message: MessageNotFoundError},
	}
}

// NewNotFoundErrorWithMessage creates 404 NotFound response with using message
func NewNotFoundErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusNotFound,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

// NewUnauthorizedError creates 401 Unauthorized response
func NewUnauthorizedError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusUnauthorized,
		Body: gen.GeneralMessageResponse{Message: MessagePermissionError},
	}
}

// NewUnauthorizedErrorWithMessage creates 401 Unauthorized response with using message
func NewUnauhorizedErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusUnauthorized,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

// NewPermissionError creates 403 Forbidden response
func NewPermissionError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusForbidden,
		Body: gen.GeneralMessageResponse{Message: MessagePermissionError},
	}
}

// NewPermissionErrorWithMessage creates 403 Forbidden response with using message
func NewPermissionErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusForbidden,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

// NewInternalError creates 500 InternalError response
func NewInternalError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusInternalServerError,
		Body: gen.GeneralMessageResponse{Message: MessageInternalError},
	}
}

// NewInternalErrorWithMessage creates 500 InternalError response with using message
func NewInternalErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusInternalServerError,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

// NewConflictedError creates 409 ConflictError response
func NewConflictedError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusConflict,
		Body: gen.GeneralMessageResponse{Message: MessageInternalError},
	}
}

// NewConflictedErrorWithMessage creates 409 ConflictError response with using message
func NewConflictedErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusConflict,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}
