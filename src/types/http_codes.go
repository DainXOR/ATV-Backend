package types

import (
	"net/http"
)

type HttpCode = int
type internalHttpCode struct{}

type _200 struct{}
type _300 struct{}
type _400 struct{}
type _500 struct{}

var Http internalHttpCode

func (internalHttpCode) Name(c HttpCode) string {
	if c == 407 {
		return "Proxy Authentication Required"
	}

	return http.StatusText(c)
}

func (internalHttpCode) C200() _200 {
	return _200{}
}
func (internalHttpCode) C300() _300 {
	return _300{}
}
func (internalHttpCode) C400() _400 {
	return _400{}
}
func (internalHttpCode) C500() _500 {
	return _500{}
}

// 2xx Success
func (_200) Ok() HttpCode {
	return http.StatusOK
}
func (_200) Created() HttpCode {
	return http.StatusCreated
}
func (_200) Accepted() HttpCode {
	return http.StatusAccepted
}
func (_200) NoContent() HttpCode {
	return http.StatusNoContent
}
func (_200) ResetContent() HttpCode {
	return http.StatusResetContent
}
func (_200) PartialContent() HttpCode {
	return http.StatusPartialContent
}
func (_200) MultiStatus() HttpCode {
	return http.StatusMultiStatus
}
func (_200) AlreadyReported() HttpCode {
	return http.StatusAlreadyReported
}
func (_200) IMUsed() HttpCode {
	return http.StatusIMUsed
}

// 3xx Redirection
func (_300) MultipleChoices() HttpCode {
	return http.StatusMultipleChoices
}
func (_300) MovedPermanently() HttpCode {
	return http.StatusMovedPermanently
}
func (_300) Found() HttpCode {
	return http.StatusFound
}
func (_300) SeeOther() HttpCode {
	return http.StatusSeeOther
}
func (_300) NotModified() HttpCode {
	return http.StatusNotModified
}
func (_300) UseProxy() HttpCode {
	return http.StatusUseProxy
}
func (_300) TemporaryRedirect() HttpCode {
	return http.StatusTemporaryRedirect
}
func (_300) PermanentRedirect() HttpCode {
	return http.StatusPermanentRedirect
}

// 4xx Client Error
func (_400) BadRequest() HttpCode {
	return http.StatusBadRequest
}
func (_400) Unauthorized() HttpCode {
	return http.StatusUnauthorized
}
func (_400) Forbidden() HttpCode {
	return http.StatusForbidden
}
func (_400) NotFound() HttpCode {
	return http.StatusNotFound
}
func (_400) MethodNotAllowed() HttpCode {
	return http.StatusMethodNotAllowed
}
func (_400) NotAcceptable() HttpCode {
	return http.StatusNotAcceptable
}
func (_400) ProxyAuthenticationRequired() HttpCode {
	return 407
}
func (_400) RequestTimeout() HttpCode {
	return http.StatusRequestTimeout
}
func (_400) Conflict() HttpCode {
	return http.StatusConflict
}
func (_400) Gone() HttpCode {
	return http.StatusGone
}
func (_400) LengthRequired() HttpCode {
	return http.StatusLengthRequired
}
func (_400) PreconditionFailed() HttpCode {
	return http.StatusPreconditionFailed
}
func (_400) ContentTooLarge() HttpCode {
	return http.StatusRequestEntityTooLarge
}
func (_400) RequestURITooLong() HttpCode {
	return http.StatusRequestURITooLong
}
func (_400) UnsupportedMediaType() HttpCode {
	return http.StatusUnsupportedMediaType
}
func (_400) RequestedRangeNotSatisfiable() HttpCode {
	return http.StatusRequestedRangeNotSatisfiable
}
func (_400) ExpectationFailed() HttpCode {
	return http.StatusExpectationFailed
}
func (_400) Teapot() HttpCode {
	return http.StatusTeapot
}
func (_400) MisdirectedRequest() HttpCode {
	return http.StatusMisdirectedRequest
}
func (_400) UnprocessableEntity() HttpCode {
	return http.StatusUnprocessableEntity
}
func (_400) Locked() HttpCode {
	return http.StatusLocked
}
func (_400) FailedDependency() HttpCode {
	return http.StatusFailedDependency
}
func (_400) UpgradeRequired() HttpCode {
	return http.StatusUpgradeRequired
}
func (_400) PreconditionRequired() HttpCode {
	return http.StatusPreconditionRequired
}
func (_400) TooManyRequests() HttpCode {
	return http.StatusTooManyRequests
}
func (_400) RequestHeaderFieldsTooLarge() HttpCode {
	return http.StatusRequestHeaderFieldsTooLarge
}
func (_400) UnavailableForLegalReasons() HttpCode {
	return http.StatusUnavailableForLegalReasons
}

// 5xx Server Error
func (_500) InternalServerError() HttpCode {
	return http.StatusInternalServerError
}
func (_500) NotImplemented() HttpCode {
	return http.StatusNotImplemented
}
func (_500) BadGateway() HttpCode {
	return http.StatusBadGateway
}
func (_500) ServiceUnavailable() HttpCode {
	return http.StatusServiceUnavailable
}
func (_500) GatewayTimeout() HttpCode {
	return http.StatusGatewayTimeout
}
func (_500) HTTPVersionNotSupported() HttpCode {
	return http.StatusHTTPVersionNotSupported
}
func (_500) VariantAlsoNegotiates() HttpCode {
	return http.StatusVariantAlsoNegotiates
}
func (_500) InsufficientStorage() HttpCode {
	return http.StatusInsufficientStorage
}
func (_500) LoopDetected() HttpCode {
	return http.StatusLoopDetected
}
func (_500) NotExtended() HttpCode {
	return http.StatusNotExtended
}
func (_500) NetworkAuthenticationRequired() HttpCode {
	return http.StatusNetworkAuthenticationRequired
}
