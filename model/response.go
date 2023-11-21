package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(body string) ErrorResponse {
	return ErrorResponse{
		Error: body,
	}
}

func NewSuccessResponse(body string) SuccessResponse {
	return SuccessResponse{
		Message: body,
	}
}