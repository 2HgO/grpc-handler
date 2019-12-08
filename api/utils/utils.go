package utils

import (
	"fmt"
	"encoding/json"
	pb "api/config/api"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseData unmarshalls the response data from the rpc call
func ResponseData(model interface{}, res *pb.Res) (statusCode int, errMsg string, errType *string) {
	if res.Success {
		err := json.Unmarshal(res.GetData(), model)
		if err != nil {
			statusCode, errMsg, errType  = ErrorHandler(err)
			return
		}
	}
	return 200, "", nil
}

// ResponsePayload constructs the response to be sent back for the api call
func ResponsePayload(model interface{}, res *pb.Res, err error) (statusCode int, payload gin.H) {
	var (
		code int
		msg string
		erro *string
	)
	switch {
		case err != nil:
			code, msg, erro = ErrorHandler(err)
		case res.Success:
			json.Unmarshal(res.GetData(), &model)
			fallthrough
		case !res.Success:
			code, msg = int(res.GetCode()), res.GetMessage()
			if res.GetError() != "" {
				erro = new(string)
				*erro = res.GetError()
			}
	}
	return code, gin.H{"success": res.GetSuccess(), "message": msg, "data": model, "error": erro}
}

func handleStatusCode(code codes.Code) int {
	switch code {
		case codes.OK:
			return 200
		case codes.InvalidArgument:
			return 400
		case codes.NotFound:
			return 404
		case codes.Unauthenticated:
			return 401
		case codes.Unknown, codes.Internal:
			return 500
		default:
			return 500
	}
}

// ErrorHandler handles error messages
func ErrorHandler(err error) (statusCode int, errorMsg string, _error *string) {
	if msg, ok := status.FromError(err); ok {
		_error = new(string)
		statusCode, errorMsg, *_error = handleStatusCode(msg.Code()), msg.Message(), msg.Code().String()
		return
	}
	var (
		code int
		msg  string
		erro = new(string)
	)
	switch v := err.(type) {
		case Error:
			code, msg, erro = v.Code(), v.Error(), v.Type()
		case validator.ValidationErrors:
			msg = fmt.Sprint("Validation failed on field { ", v[0].Field(), " }, Condition: ", v[0].ActualTag())
			if v[0].Param() != "" {
				msg += fmt.Sprint(" { ", v[0].Param(), " }")
			}
			if v[0].Value() != nil {
				msg += fmt.Sprint(", Value Recieved: ", v[0].Value())
			}
			code, *erro = 400, "Validation Error"
		default:
			code, msg, *erro = 500, v.Error(), "Validation Error"
	}

	return code, msg, erro
}

// Error is the default structure of an error
type Error struct {
	statusCode int
	message   string
	errType *string
}

// Code returns the error code associated with an error
func (e Error) Code() int {
	return e.statusCode
}

// Error satisfies the interface for an error type and returns the error message associated with an error
func (e Error) Error() string {
	return e.message
}

// Type returns the error type associated with an error
func (e Error) Type() *string {
	return e.errType
}

// NewErr returns a new Error
func NewErr(code int, msg string, err *string) Error {
	return Error{code, msg, err}
}

// Unmarshal ...
func Unmarshal(source, dest interface{}) error {
	load, err := json.Marshal(source)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if err := json.Unmarshal(load, dest); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

// PanicHandler ...
func PanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"message": "Looks like something bad occured on our end, we're on our way to fix it",
					"error":   err,
				})
			}
		}()
		c.Next()
	}
}