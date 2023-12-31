package utils

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type ErrorResponses struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors"`
}

type PaginationResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Meta       Meta        `json:"meta"`
}

type Meta struct {
	CurrentPage int         `json:"current_page" example:"1"`
	PrevPage    int         `json:"prev_page" example:"1"`
	NextPage    interface{} `json:"next_page"`
	Total       int         `json:"total" example:"1"`
}

func NewResponse(statusCode int, message string, data interface{}) Response {
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func NewErrorResponse(statusCode int, message string, errors interface{}) ErrorResponses {
	return ErrorResponses{
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}
}

func NewPaginationResponse(statusCode int, message string, data interface{}, page int, limit int, total int) PaginationResponse {
	var (
		nextPage interface{}
		prevPage int
	)
	if page*limit >= total {
		nextPage = nil
	} else {
		nextPage = page + 1
	}

	if page == 1 {
		prevPage = page
	} else {
		prevPage = page - 1
	}

	return PaginationResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta: Meta{
			CurrentPage: page,
			NextPage:    nextPage,
			PrevPage:    prevPage,
			Total:       total,
		},
	}
}
