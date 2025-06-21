package types

import (
	"dainxor/atv/logger"
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
func (internalHttpCode) Ok() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Ok() method, use C200().Ok() instead")
	return http.StatusOK
}
func (internalHttpCode) Created() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Created() method, use C200().Created() instead")
	return http.StatusCreated
}
func (internalHttpCode) Accepted() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Accepted() method, use C200().Accepted() instead")
	return http.StatusAccepted
}
func (internalHttpCode) NoContent() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated NoContent() method, use C200().NoContent() instead")
	return http.StatusNoContent
}
func (internalHttpCode) ResetContent() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated ResetContent() method, use C200().ResetContent() instead")
	return http.StatusResetContent
}
func (internalHttpCode) PartialContent() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated PartialContent() method, use C200().PartialContent() instead")
	return http.StatusPartialContent
}
func (internalHttpCode) MultiStatus() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated MultiStatus() method, use C200().MultiStatus() instead")
	return http.StatusMultiStatus
}
func (internalHttpCode) AlreadyReported() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated AlreadyReported() method, use C200().AlreadyReported() instead")
	return http.StatusAlreadyReported
}
func (internalHttpCode) IMUsed() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated IMUsed() method, use C200().IMUsed() instead")
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
func (internalHttpCode) MultipleChoices() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated MultipleChoices() method, use C300().MultipleChoices() instead")
	return http.StatusMultipleChoices
}
func (internalHttpCode) MovedPermanently() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated MovedPermanently() method, use C300().MovedPermanently() instead")
	return http.StatusMovedPermanently
}
func (internalHttpCode) Found() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Found() method, use C300().Found() instead")
	return http.StatusFound
}
func (internalHttpCode) SeeOther() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated SeeOther() method, use C300().SeeOther() instead")
	return http.StatusSeeOther
}
func (internalHttpCode) NotModified() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated NotModified() method, use C300().NotModified() instead")
	return http.StatusNotModified
}
func (internalHttpCode) UseProxy() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated UseProxy() method, use C300().UseProxy() instead")
	return http.StatusUseProxy
}
func (internalHttpCode) TemporaryRedirect() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated TemporaryRedirect() method, use C300().TemporaryRedirect() instead")
	return http.StatusTemporaryRedirect
}
func (internalHttpCode) PermanentRedirect() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated PermanentRedirect() method, use C300().PermanentRedirect() instead")
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
func (internalHttpCode) BadRequest() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated BadRequest() method, use C400().BadRequest() instead")
	return http.StatusBadRequest
}
func (internalHttpCode) Unauthorized() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Unauthorized() method, use C400().Unauthorized() instead")
	return http.StatusUnauthorized
}
func (internalHttpCode) Forbidden() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Forbidden() method, use C400().Forbidden() instead")
	return http.StatusForbidden
}
func (internalHttpCode) NotFound() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated NotFound() method, use C400().NotFound() instead")
	return http.StatusNotFound
}
func (internalHttpCode) MethodNotAllowed() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated MethodNotAllowed() method, use C400().MethodNotAllowed() instead")
	return http.StatusMethodNotAllowed
}
func (internalHttpCode) NotAcceptable() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated NotAcceptable() method, use C400().NotAcceptable() instead")
	return http.StatusNotAcceptable
}
func (internalHttpCode) ProxyAuthenticationRequired() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated ProxyAuthenticationRequired() method, use C400().ProxyAuthenticationRequired() instead")
	return 407
}
func (internalHttpCode) RequestTimeout() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated RequestTimeout() method, use C400().RequestTimeout() instead")
	return http.StatusRequestTimeout
}
func (internalHttpCode) Conflict() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Conflict() method, use C400().Conflict() instead")
	return http.StatusConflict
}
func (internalHttpCode) Gone() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Gone() method, use C400().Gone() instead")
	return http.StatusGone
}
func (internalHttpCode) LengthRequired() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated LengthRequired() method, use C400().LengthRequired() instead")
	return http.StatusLengthRequired
}
func (internalHttpCode) PreconditionFailed() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated PreconditionFailed() method, use C400().PreconditionFailed() instead")
	return http.StatusPreconditionFailed
}
func (internalHttpCode) ContentTooLarge() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated ContentTooLarge() method, use C400().ContentTooLarge() instead")
	return http.StatusRequestEntityTooLarge
}
func (internalHttpCode) RequestURITooLong() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated RequestURITooLong() method, use C400().RequestURITooLong() instead")
	return http.StatusRequestURITooLong
}
func (internalHttpCode) UnsupportedMediaType() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated UnsupportedMediaType() method, use C400().UnsupportedMediaType() instead")
	return http.StatusUnsupportedMediaType
}
func (internalHttpCode) RequestedRangeNotSatisfiable() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated RequestedRangeNotSatisfiable() method, use C400().RequestedRangeNotSatisfiable() instead")
	return http.StatusRequestedRangeNotSatisfiable
}
func (internalHttpCode) ExpectationFailed() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated ExpectationFailed() method, use C400().ExpectationFailed() instead")
	return http.StatusExpectationFailed
}
func (internalHttpCode) Teapot() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Teapot() method, use C400().Teapot() instead")
	return http.StatusTeapot
}
func (internalHttpCode) MisdirectedRequest() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated MisdirectedRequest() method, use C400().MisdirectedRequest() instead")
	return http.StatusMisdirectedRequest
}
func (internalHttpCode) UnprocessableEntity() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated UnprocessableEntity() method, use C400().UnprocessableEntity() instead")
	return http.StatusUnprocessableEntity
}
func (internalHttpCode) Locked() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated Locked() method, use C400().Locked() instead")
	return http.StatusLocked
}
func (internalHttpCode) FailedDependency() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated FailedDependency() method, use C400().FailedDependency() instead")
	return http.StatusFailedDependency
}
func (internalHttpCode) UpgradeRequired() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated UpgradeRequired() method, use C400().UpgradeRequired() instead")
	return http.StatusUpgradeRequired
}
func (internalHttpCode) PreconditionRequired() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated PreconditionRequired() method, use C400().PreconditionRequired() instead")
	return http.StatusPreconditionRequired
}
func (internalHttpCode) TooManyRequests() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated TooManyRequests() method, use C400().TooManyRequests() instead")
	return http.StatusTooManyRequests
}
func (internalHttpCode) RequestHeaderFieldsTooLarge() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated RequestHeaderFieldsTooLarge() method, use C400().RequestHeaderFieldsTooLarge() instead")
	return http.StatusRequestHeaderFieldsTooLarge
}
func (internalHttpCode) UnavailableForLegalReasons() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated UnavailableForLegalReasons() method, use C400().UnavailableForLegalReasons() instead")
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
func (internalHttpCode) InternalServerError() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated InternalServerError() method, use C500().InternalServerError() instead")
	return http.StatusInternalServerError
}
func (internalHttpCode) NotImplemented() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated NotImplemented() method, use C500().NotImplemented() instead")
	return http.StatusNotImplemented
}
func (internalHttpCode) BadGateway() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated BadGateway() method, use C500().BadGateway() instead")
	return http.StatusBadGateway
}
func (internalHttpCode) ServiceUnavailable() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated ServiceUnavailable() method, use C500().ServiceUnavailable() instead")
	return http.StatusServiceUnavailable
}
func (internalHttpCode) GatewayTimeout() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated GatewayTimeout() method, use C500().GatewayTimeout() instead")
	return http.StatusGatewayTimeout
}
func (internalHttpCode) HTTPVersionNotSupported() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated HTTPVersionNotSupported() method, use C500().HTTPVersionNotSupported() instead")
	return http.StatusHTTPVersionNotSupported
}
func (internalHttpCode) VariantAlsoNegotiates() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated VariantAlsoNegotiates() method, use C500().VariantAlsoNegotiates() instead")
	return http.StatusVariantAlsoNegotiates
}
func (internalHttpCode) InsufficientStorage() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated InsufficientStorage() method, use C500().InsufficientStorage() instead")
	return http.StatusInsufficientStorage
}
func (internalHttpCode) LoopDetected() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated LoopDetected() method, use C500().LoopDetected() instead")
	return http.StatusLoopDetected
}
func (internalHttpCode) NotExtended() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated NotExtended() method, use C500().NotExtended() instead")
	return http.StatusNotExtended
}
func (internalHttpCode) NetworkAuthenticationRequired() HttpCode {
	logger.Deprecate(1, 2, "Using deprecated NetworkAuthenticationRequired() method, use C500().NetworkAuthenticationRequired() instead")
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
