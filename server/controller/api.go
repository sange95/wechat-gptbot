package controller

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	State  string      `json:"state"`
	Reason string      `json:"reason,omitempty"`
	Data   interface{} `json:"data"`
}

func (cr *Response) ResponseError(reason string, err error) {}

func (cr *Response) ResponseOk(data interface{}, instanceId, verb string, c *gin.Context) {}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		State: "Success",
		Data:  data,
	}
}

func NewFailureResponse(reason string) *Response {
	return &Response{
		State:  "Failure",
		Reason: reason,
	}
}
