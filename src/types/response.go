package types

type jsonResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Extra   any    `json:"extra,omitempty"`
}

func Response(data any, message string, extra ...any) jsonResponse {
	response := jsonResponse{
		Data:    data,
		Message: message,
	}

	if len(extra) > 0 {
		response.Extra = extra
	}

	return response
}
