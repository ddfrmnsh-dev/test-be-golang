package helper

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  bool        `json:"status"`
}

func APIResponse(message string, data interface{}) Response {
	return Response{
		Message: message,
		Data:    data,
		Status:  true,
	}
}
func APIErrorResponse(message string) Response {
	return Response{
		Message: message,
		Status:  false,
	}
}
