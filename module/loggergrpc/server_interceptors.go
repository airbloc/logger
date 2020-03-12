package loggergrpc

import (
	"context"
	"github.com/airbloc/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


// UnaryServerLogger returns a new unary server interceptors that logs unary requests.
func UnaryServerLogger(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		timer := log.Timer()
		resp, err = handler(ctx, req)

		if panicErr, ok := err.(*logger.PanicError); ok {
			err = status.Error(codes.Internal, panicErr.Error())
			timer.End("{}Request {} – {}{}", logger.Red, info.FullMethod, panicErr.Pretty(), logger.Reset)

		} else if _, ok := status.FromError(err); ok {
			code := grpcCodeToString[status.Code(err)]
			timer.End("{}Request {} – {}: {}{}", logger.Red, info.FullMethod, code, err.Error(), logger.Reset)
		} else {
			timer.End("Request {} – OK", info.FullMethod)
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
func StreamServerLogger(log logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		timer := log.Timer()
		err := handler(srv, ss)

		if panicErr, ok := err.(*logger.PanicError); ok {
			err = status.Error(codes.Internal, panicErr.Error())
			timer.End("{}Streaming {} – {}{}", logger.Red, info.FullMethod, panicErr.Pretty(), logger.Reset)

		} else if _, ok := status.FromError(err); ok {
			code := grpcCodeToString[status.Code(err)]
			timer.End("{}Streaming {} – {}: {}{}", logger.Red, info.FullMethod, code, err.Error(), logger.Reset)
		} else {
			timer.End("Streaming {} – OK", info.FullMethod)
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
