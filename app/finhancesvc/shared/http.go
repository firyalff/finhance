package shared

type PagedResponseMeta struct {
	TotalData int64 `json:"total_data"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

type PagedResponse struct {
	Meta PagedResponseMeta `json:"meta"`
	Data []interface{}     `json:"data"`
}

func GeneratePagedResponse(total, limit, offset int, data []interface{}) PagedResponse {
	return PagedResponse{
		Meta: PagedResponseMeta{
			TotalData: int64(total),
			Limit:     int32(limit),
			Offset:    int32(offset),
		},
		Data: data,
	}
}

type ErrorResponse struct {
	ErrorCode string      `json:"error"`
	Details   interface{} `json:"details"`
}

func GenerateErrorResponse(errorCode string, errorDetails interface{}) ErrorResponse {

	return ErrorResponse{
		ErrorCode: errorCode,
		Details:   errorDetails,
	}

}
