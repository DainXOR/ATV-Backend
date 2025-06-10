package types

import (
	"dainxor/atv/logger"
	"fmt"
	"net/http"
)

type HttpCode int
type httpCode struct{}

type _200 struct{}
type _300 struct{}
type _400 struct{}
type _500 struct{}

var Http httpCode

func (c HttpCode) AsInt() int {
	return int(c)
}
func (c HttpCode) AsString() string {
	return fmt.Sprint(c.AsInt())
}
func (c HttpCode) Name() string {
	return http.StatusText(c.AsInt())
}

func (HttpCode) C200() _200 {
	return _200{}
}
func (HttpCode) C300() _300 {
	return _300{}
}
func (HttpCode) C400() _400 {
	return _400{}
}
func (HttpCode) C500() _500 {
	return _500{}
}

// 2xx Success
func (httpCode) Ok() HttpCode {
	logger.Warning("Using deprecated Ok() method, use C200().Ok() instead")
	return http.StatusOK
}
func (httpCode) Created() HttpCode {
	logger.Warning("Using deprecated Created() method, use C200().Created() instead")
	return http.StatusCreated
}
func (httpCode) Accepted() HttpCode {
	logger.Warning("Using deprecated Accepted() method, use C200().Accepted() instead")
	return http.StatusAccepted
}
func (httpCode) NoContent() HttpCode {
	logger.Warning("Using deprecated NoContent() method, use C200().NoContent() instead")
	return http.StatusNoContent
}
func (httpCode) ResetContent() HttpCode {
	logger.Warning("Using deprecated ResetContent() method, use C200().ResetContent() instead")
	return http.StatusResetContent
}
func (httpCode) PartialContent() HttpCode {
	logger.Warning("Using deprecated PartialContent() method, use C200().PartialContent() instead")
	return http.StatusPartialContent
}
func (httpCode) MultiStatus() HttpCode {
	logger.Warning("Using deprecated MultiStatus() method, use C200().MultiStatus() instead")
	return http.StatusMultiStatus
}
func (httpCode) AlreadyReported() HttpCode {
	logger.Warning("Using deprecated AlreadyReported() method, use C200().AlreadyReported() instead")
	return http.StatusAlreadyReported
}
func (httpCode) IMUsed() HttpCode {
	logger.Warning("Using deprecated IMUsed() method, use C200().IMUsed() instead")
	return http.StatusIMUsed
}

// Future 2xx codes
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
func (httpCode) MultipleChoices() HttpCode {
	logger.Warning("Using deprecated MultipleChoices() method, use C300().MultipleChoices() instead")
	return http.StatusMultipleChoices
}
func (httpCode) MovedPermanently() HttpCode {
	logger.Warning("Using deprecated MovedPermanently() method, use C300().MovedPermanently() instead")
	return http.StatusMovedPermanently
}
func (httpCode) Found() HttpCode {
	logger.Warning("Using deprecated Found() method, use C300().Found() instead")
	return http.StatusFound
}
func (httpCode) SeeOther() HttpCode {
	logger.Warning("Using deprecated SeeOther() method, use C300().SeeOther() instead")
	return http.StatusSeeOther
}
func (httpCode) NotModified() HttpCode {
	logger.Warning("Using deprecated NotModified() method, use C300().NotModified() instead")
	return http.StatusNotModified
}
func (httpCode) UseProxy() HttpCode {
	logger.Warning("Using deprecated UseProxy() method, use C300().UseProxy() instead")
	return http.StatusUseProxy
}
func (httpCode) TemporaryRedirect() HttpCode {
	logger.Warning("Using deprecated TemporaryRedirect() method, use C300().TemporaryRedirect() instead")
	return http.StatusTemporaryRedirect
}
func (httpCode) PermanentRedirect() HttpCode {
	logger.Warning("Using deprecated PermanentRedirect() method, use C300().PermanentRedirect() instead")
	return http.StatusPermanentRedirect
}

// Future 3xx codes
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
func (httpCode) BadRequest() HttpCode {
	logger.Warning("Using deprecated BadRequest() method, use C400().BadRequest() instead")
	return http.StatusBadRequest
}
func (httpCode) Unauthorized() HttpCode {
	logger.Warning("Using deprecated Unauthorized() method, use C400().Unauthorized() instead")
	return http.StatusUnauthorized
}
func (httpCode) Forbidden() HttpCode {
	logger.Warning("Using deprecated Forbidden() method, use C400().Forbidden() instead")
	return http.StatusForbidden
}
func (httpCode) NotFound() HttpCode {
	logger.Warning("Using deprecated NotFound() method, use C400().NotFound() instead")
	return http.StatusNotFound
}
func (httpCode) MethodNotAllowed() HttpCode {
	logger.Warning("Using deprecated MethodNotAllowed() method, use C400().MethodNotAllowed() instead")
	return http.StatusMethodNotAllowed
}
func (httpCode) NotAcceptable() HttpCode {
	logger.Warning("Using deprecated NotAcceptable() method, use C400().NotAcceptable() instead")
	return http.StatusNotAcceptable
}
func (httpCode) ProxyAuthenticationRequired() HttpCode {
	logger.Warning("Using deprecated ProxyAuthenticationRequired() method, use C400().ProxyAuthenticationRequired() instead")
	return 407
}
func (httpCode) RequestTimeout() HttpCode {
	logger.Warning("Using deprecated RequestTimeout() method, use C400().RequestTimeout() instead")
	return http.StatusRequestTimeout
}
func (httpCode) Conflict() HttpCode {
	logger.Warning("Using deprecated Conflict() method, use C400().Conflict() instead")
	return http.StatusConflict
}
func (httpCode) Gone() HttpCode {
	logger.Warning("Using deprecated Gone() method, use C400().Gone() instead")
	return http.StatusGone
}
func (httpCode) LengthRequired() HttpCode {
	logger.Warning("Using deprecated LengthRequired() method, use C400().LengthRequired() instead")
	return http.StatusLengthRequired
}
func (httpCode) PreconditionFailed() HttpCode {
	logger.Warning("Using deprecated PreconditionFailed() method, use C400().PreconditionFailed() instead")
	return http.StatusPreconditionFailed
}
func (httpCode) ContentTooLarge() HttpCode {
	logger.Warning("Using deprecated ContentTooLarge() method, use C400().ContentTooLarge() instead")
	return http.StatusRequestEntityTooLarge
}
func (httpCode) RequestURITooLong() HttpCode {
	logger.Warning("Using deprecated RequestURITooLong() method, use C400().RequestURITooLong() instead")
	return http.StatusRequestURITooLong
}
func (httpCode) UnsupportedMediaType() HttpCode {
	logger.Warning("Using deprecated UnsupportedMediaType() method, use C400().UnsupportedMediaType() instead")
	return http.StatusUnsupportedMediaType
}
func (httpCode) RequestedRangeNotSatisfiable() HttpCode {
	logger.Warning("Using deprecated RequestedRangeNotSatisfiable() method, use C400().RequestedRangeNotSatisfiable() instead")
	return http.StatusRequestedRangeNotSatisfiable
}
func (httpCode) ExpectationFailed() HttpCode {
	logger.Warning("Using deprecated ExpectationFailed() method, use C400().ExpectationFailed() instead")
	return http.StatusExpectationFailed
}
func (httpCode) Teapot() HttpCode {
	logger.Warning("Using deprecated Teapot() method, use C400().Teapot() instead")
	return http.StatusTeapot
}
func (httpCode) MisdirectedRequest() HttpCode {
	logger.Warning("Using deprecated MisdirectedRequest() method, use C400().MisdirectedRequest() instead")
	return http.StatusMisdirectedRequest
}
func (httpCode) UnprocessableEntity() HttpCode {
	logger.Warning("Using deprecated UnprocessableEntity() method, use C400().UnprocessableEntity() instead")
	return http.StatusUnprocessableEntity
}
func (httpCode) Locked() HttpCode {
	logger.Warning("Using deprecated Locked() method, use C400().Locked() instead")
	return http.StatusLocked
}
func (httpCode) FailedDependency() HttpCode {
	logger.Warning("Using deprecated FailedDependency() method, use C400().FailedDependency() instead")
	return http.StatusFailedDependency
}
func (httpCode) UpgradeRequired() HttpCode {
	logger.Warning("Using deprecated UpgradeRequired() method, use C400().UpgradeRequired() instead")
	return http.StatusUpgradeRequired
}
func (httpCode) PreconditionRequired() HttpCode {
	logger.Warning("Using deprecated PreconditionRequired() method, use C400().PreconditionRequired() instead")
	return http.StatusPreconditionRequired
}
func (httpCode) TooManyRequests() HttpCode {
	logger.Warning("Using deprecated TooManyRequests() method, use C400().TooManyRequests() instead")
	return http.StatusTooManyRequests
}
func (httpCode) RequestHeaderFieldsTooLarge() HttpCode {
	logger.Warning("Using deprecated RequestHeaderFieldsTooLarge() method, use C400().RequestHeaderFieldsTooLarge() instead")
	return http.StatusRequestHeaderFieldsTooLarge
}
func (httpCode) UnavailableForLegalReasons() HttpCode {
	logger.Warning("Using deprecated UnavailableForLegalReasons() method, use C400().UnavailableForLegalReasons() instead")
	return http.StatusUnavailableForLegalReasons
}

// Future 4xx codes
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
func (httpCode) InternalServerError() HttpCode {
	logger.Warning("Using deprecated InternalServerError() method, use C500().InternalServerError() instead")
	return http.StatusInternalServerError
}
func (httpCode) NotImplemented() HttpCode {
	logger.Warning("Using deprecated NotImplemented() method, use C500().NotImplemented() instead")
	return http.StatusNotImplemented
}
func (httpCode) BadGateway() HttpCode {
	logger.Warning("Using deprecated BadGateway() method, use C500().BadGateway() instead")
	return http.StatusBadGateway
}
func (httpCode) ServiceUnavailable() HttpCode {
	logger.Warning("Using deprecated ServiceUnavailable() method, use C500().ServiceUnavailable() instead")
	return http.StatusServiceUnavailable
}
func (httpCode) GatewayTimeout() HttpCode {
	logger.Warning("Using deprecated GatewayTimeout() method, use C500().GatewayTimeout() instead")
	return http.StatusGatewayTimeout
}
func (httpCode) HTTPVersionNotSupported() HttpCode {
	logger.Warning("Using deprecated HTTPVersionNotSupported() method, use C500().HTTPVersionNotSupported() instead")
	return http.StatusHTTPVersionNotSupported
}
func (httpCode) VariantAlsoNegotiates() HttpCode {
	logger.Warning("Using deprecated VariantAlsoNegotiates() method, use C500().VariantAlsoNegotiates() instead")
	return http.StatusVariantAlsoNegotiates
}
func (httpCode) InsufficientStorage() HttpCode {
	logger.Warning("Using deprecated InsufficientStorage() method, use C500().InsufficientStorage() instead")
	return http.StatusInsufficientStorage
}
func (httpCode) LoopDetected() HttpCode {
	logger.Warning("Using deprecated LoopDetected() method, use C500().LoopDetected() instead")
	return http.StatusLoopDetected
}
func (httpCode) NotExtended() HttpCode {
	logger.Warning("Using deprecated NotExtended() method, use C500().NotExtended() instead")
	return http.StatusNotExtended
}
func (httpCode) NetworkAuthenticationRequired() HttpCode {
	logger.Warning("Using deprecated NetworkAuthenticationRequired() method, use C500().NetworkAuthenticationRequired() instead")
	return http.StatusNetworkAuthenticationRequired
}

// Future 5xx codes
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
