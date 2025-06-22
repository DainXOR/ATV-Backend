package types

import (
	"errors"
	"fmt"
)

type HttpError struct {
	Code   HttpCode `json:"code"`
	Err    error    `json:"error"`
	Detail string   `json:"detail"`
}

func Error(code HttpCode, information ...string) HttpError {
	message := ""
	detail := ""

	if len(information) > 0 {
		message = information[0]

		if len(information) > 1 {
			detail = information[1]
		}
	}

	return HttpError{
		Code:   code,
		Err:    errors.New(message),
		Detail: detail,
	}
}

func (m *HttpError) Error() string {
	return fmt.Sprintf("%s (%d): %s - %s", Http.Name(m.Code), m.Code, m.Err, m.Detail)
}

func ErrorNotFound(information ...string) HttpError {
	return Error(Http.NotFound(), information...)
}

func ErrorInternal(information ...string) HttpError {
	return Error(Http.InternalServerError(), information...)
}
