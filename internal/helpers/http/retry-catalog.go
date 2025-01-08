package http

import "net/http"

var (
	code2xx = map[int]bool{
		http.StatusOK:                   true,
		http.StatusCreated:              true,
		http.StatusAccepted:             true,
		http.StatusNonAuthoritativeInfo: true,
		http.StatusNoContent:            true,
		http.StatusResetContent:         true,
		http.StatusPartialContent:       true,
		http.StatusMultiStatus:          true,
		http.StatusAlreadyReported:      true,
		http.StatusIMUsed:               true,
	}

	code4xx = map[int]bool{
		http.StatusBadRequest:                   true,
		http.StatusUnauthorized:                 true,
		http.StatusPaymentRequired:              true,
		http.StatusForbidden:                    true,
		http.StatusNotFound:                     true,
		http.StatusMethodNotAllowed:             true,
		http.StatusNotAcceptable:                true,
		http.StatusProxyAuthRequired:            true,
		http.StatusRequestTimeout:               true,
		http.StatusConflict:                     true,
		http.StatusGone:                         true,
		http.StatusLengthRequired:               true,
		http.StatusPreconditionFailed:           true,
		http.StatusRequestEntityTooLarge:        true,
		http.StatusRequestURITooLong:            true,
		http.StatusUnsupportedMediaType:         true,
		http.StatusRequestedRangeNotSatisfiable: true,
		http.StatusExpectationFailed:            true,
		http.StatusTeapot:                       true,
		http.StatusMisdirectedRequest:           true,
		http.StatusUnprocessableEntity:          true,
		http.StatusLocked:                       true,
		http.StatusFailedDependency:             true,
		http.StatusTooEarly:                     true,
		http.StatusUpgradeRequired:              true,
		http.StatusPreconditionRequired:         true,
		http.StatusTooManyRequests:              true,
		http.StatusRequestHeaderFieldsTooLarge:  true,
		http.StatusUnavailableForLegalReasons:   true,
	}

	Is5XXCode = map[int]bool{
		http.StatusInternalServerError:           true,
		http.StatusNotImplemented:                true,
		http.StatusBadGateway:                    true,
		http.StatusServiceUnavailable:            true,
		http.StatusGatewayTimeout:                true,
		http.StatusHTTPVersionNotSupported:       true,
		http.StatusVariantAlsoNegotiates:         true,
		http.StatusInsufficientStorage:           true,
		http.StatusLoopDetected:                  true,
		http.StatusNotExtended:                   true,
		http.StatusNetworkAuthenticationRequired: true,
	}
)

func Is2XXCode(code int) bool {
	return code2xx[code]
}

func Is4XXCode(code int) bool {
	return code4xx[code]
}
