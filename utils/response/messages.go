package response

import (
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
)

const (
	MessageRequestError      = "Your request is not valid."
	MessageLockedError       = "Specified content is locked and not editable."
	MessageNotFoundError     = "Specified content was not exist."
	MessageUnauthorizedError = "Probably your password incorrect."
	MessageConflictedError   = "Specified content was already exists."
	MessagePermissionError   = "You don't have enough permission to do it."
	MessageInternalError     = "Unfortunately, the server exploded."
)

func NewRequestError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusBadRequest,
		Body: gen.GeneralMessageResponse{Message: MessageRequestError},
	}
}

func NewRequestErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusBadRequest,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

func NewLockedError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusLocked,
		Body: gen.GeneralMessageResponse{Message: MessageLockedError},
	}
}

func NewLockedErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusLocked,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

func NewNotFoundError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusNotFound,
		Body: gen.GeneralMessageResponse{Message: MessageNotFoundError},
	}
}

func NewNotFoundErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusNotFound,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

func NewUnauthorizedError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusUnauthorized,
		Body: gen.GeneralMessageResponse{Message: MessagePermissionError},
	}
}

func NewUnauhorizedErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusUnauthorized,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

func NewPermissionError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusForbidden,
		Body: gen.GeneralMessageResponse{Message: MessagePermissionError},
	}
}

func NewPermissionErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusForbidden,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

func NewInternalError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusInternalServerError,
		Body: gen.GeneralMessageResponse{Message: MessageInternalError},
	}
}

func NewInternalErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusInternalServerError,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}

func NewConflictedError() gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusConflict,
		Body: gen.GeneralMessageResponse{Message: MessageInternalError},
	}
}

func NewConflictedErrorWithMessage(message string) gen.ImplResponse {
	return gen.ImplResponse{
		Code: http.StatusConflict,
		Body: gen.GeneralMessageResponse{Message: message},
	}
}
