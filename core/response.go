package core

import (
	"simplest_script/core/tool"
	"time"
)

type Response struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
}

func Success(requestId string, message string, data any) *Response {
	if len(requestId) == 0 {
		requestId = time.Now().Format("20060102150405") + tool.Uuid()
	}

	return &Response{
		Code:      CodeOK,
		Data:      data,
		Message:   message,
		RequestId: requestId,
	}
}

func Fail(requestId string, message string, data any) *Response {
	if len(requestId) == 0 {
		requestId = time.Now().Format("20060102150405") + tool.Uuid()
	}

	return &Response{
		Code:      CodeFail,
		Data:      data,
		Message:   message,
		RequestId: requestId,
	}
}

func LoginFail(requestId string, message string, data any) *Response {
	if len(requestId) == 0 {
		requestId = time.Now().Format("20060102150405") + tool.Uuid()
	}

	return &Response{
		Code:      CodeLoginFail,
		Data:      data,
		Message:   message,
		RequestId: requestId,
	}
}
