package common

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse(status int, message string, data any) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func Success(message string, data any) Response {
	return NewResponse(200, message, data)
}

func Error(status int, message string) Response {
	return NewResponse(status, message, nil)
}
