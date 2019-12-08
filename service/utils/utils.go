package utils

import (
	"context"
	"encoding/json"
	"github.com/go-bongo/bongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"
)

// RecoveryInterceptor ...
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	// [TODO] add mailer to recover function
	defer func() {
		if erro := recover(); erro != nil {
			res, err = nil, status.Errorf(codes.Internal, "%v", erro)
			grpclog.SetLogger(log.New(os.Stdin, "[UserService::Error] ", 0))
			grpclog.Errorf("Error occured during RPC method=%s; Error=%v", info.FullMethod, erro)
		}
	}()

	start := time.Now()

	resp, erro := handler(ctx, req)

	grpclog.SetLogger(log.New(os.Stdin, "[UserService::Log] ", 0))
	grpclog.Printf("Handled RPC method=%s; Duration=%s; Error=%v", info.FullMethod, time.Since(start), erro)
	return resp, erro
}

// ErrorHandler ...
func ErrorHandler(err error) error {
	switch v := err.(type) {
	case *bongo.ValidationError:
		return status.Error(codes.InvalidArgument, v.Errors[0].Error())
	case *bongo.DocumentNotFoundError:
		return status.Error(codes.NotFound, v.Error())
	default:
		return status.Error(codes.Internal, v.Error())
	}
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

// Marshal ...
func Marshal(source interface{}) ([]byte, error) {
	load, err := json.Marshal(source)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return load, nil
}
