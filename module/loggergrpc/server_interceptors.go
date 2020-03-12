package loggergrpc

import (
	"context"
	"github.com/airbloc/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// UnaryServerLogger returns a new unary server interceptors that logs unary requests.
func UnaryServerLogger(log logger.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	opt := createOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		resp, err = handler(ctx, req)
		elapsed := time.Now().Sub(startTime).String()

		if panicErr, ok := err.(*logger.PanicError); ok {
			err = status.Error(codes.Internal, panicErr.Error())
			log.Log(opt.errorLogLevel, "Request {}({}) – {}", []interface{}{info.FullMethod, elapsed, panicErr.Pretty()})

		} else if s := status.Convert(err); s != nil {
			code := grpcCodeToString[s.Code()]
			log.Log(opt.errorLogLevel, "Request {}({}) – {}: {}", []interface{}{info.FullMethod, elapsed, code, s.Message()})
		} else {
			log.Log(opt.requestLogLevel, "Request {}({}) – OK", []interface{}{info.FullMethod, elapsed})
		}
		return
	}
}

// UnaryServerRecover returns a new unary server interceptors
// that recovers panic and wraps with logger.WrapRecover.
func UnaryServerRecover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if panicErr := logger.WrapRecover(recover()); panicErr != nil {
				err = panicErr
			}
		}()
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that logs streams.
func StreamServerLogger(log logger.Logger, opts ...Option) grpc.StreamServerInterceptor {
	opt := createOptions(opts)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		err := handler(srv, ss)
		elapsed := time.Now().Sub(startTime).String()

		if panicErr, ok := err.(*logger.PanicError); ok {
			err = status.Error(codes.Internal, panicErr.Error())
			log.Log(opt.errorLogLevel, "Streaming {}({}) – {}", []interface{}{info.FullMethod, elapsed, panicErr.Pretty()})

		} else if s := status.Convert(err); s != nil {
			code := grpcCodeToString[s.Code()]
			log.Log(opt.errorLogLevel, "Streaming {}({}) – {}: {}", []interface{}{info.FullMethod, elapsed, code, s.Message()})
		} else {
			log.Log(opt.requestLogLevel, "Streaming {}({}) – OK", []interface{}{info.FullMethod, elapsed})
		}
		return err
	}
}

// StreamServerRecover returns a new streaming server interceptors
// that recovers panic and wraps with logger.WrapRecover.
func StreamServerRecover() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if panicErr := logger.WrapRecover(recover()); panicErr != nil {
				err = panicErr
			}
		}()
		return handler(srv, ss)
	}
}
