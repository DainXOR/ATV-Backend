package types

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
