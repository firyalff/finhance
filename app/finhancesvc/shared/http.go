package shared

import (
	"github.com/gin-gonic/gin"
)

type PagedRequest struct {
	Limit  int32  `form:"limit" validate:"required,gt=0"`
	Offset int32  `form:"offset"`
	Filter string `form:"filter"`
}

type DetailByIDRequest struct {
	ID string `uri:"id"`
}

type PagedResponseMeta struct {
	TotalData int64 `json:"total_data"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

type PagedResponse struct {
	Meta PagedResponseMeta `json:"meta"`
	Data interface{}       `json:"data"`
}

func GeneratePagedResponse(total, limit, offset int, data interface{}) gin.H {
	if data == nil {
		data = make([]interface{}, 0)
	}

	return gin.H{
		"meta": PagedResponseMeta{
			TotalData: int64(total),
			Limit:     int32(limit),
			Offset:    int32(offset),
		},
		"data": data,
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
