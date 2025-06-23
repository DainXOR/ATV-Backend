package types

type body = map[string]any

type JSONResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Extra   any    `json:"extra,omitempty"`
}

func Response(data any, message string, extra ...any) JSONResponse {
	response := JSONResponse{
		Data:    data,
		Message: message,
	}

	if len(extra) > 0 {
		response.Extra = extra
	}

	return response
}

func EmptyResponse(message string, extra ...any) JSONResponse {
	response := JSONResponse{
		Data:    body{},
		Message: message,
	}

	if len(extra) > 0 {
		response.Extra = extra
	}

	return response
}

func ErrorResponse(err error) JSONResponse {
	return Response(
		body{},
		err.Error(),
	)
}
