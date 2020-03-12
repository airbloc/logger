package loggergrpc

import "google.golang.org/grpc/codes"

var grpcCodeToString = map[codes.Code]string{
	codes.OK:                 "OK",
	codes.Canceled:           "Canceled",
	codes.Unknown:            "Unknown",
	codes.InvalidArgument:    "InvalidArgument",
	codes.DeadlineExceeded:   "DeadlineExceeded",
	codes.NotFound:           "NotFound",
	codes.AlreadyExists:      "AlreadyExists",
	codes.PermissionDenied:   "PermissionDenied",
	codes.ResourceExhausted:  "ResourceExhausted",
	codes.FailedPrecondition: "FailedPrecondition",
	codes.Aborted:            "Aborted",
	codes.OutOfRange:         "OutOfRange",
	codes.Unimplemented:      "Unimplemented",
	codes.Internal:           "Internal",
	codes.Unavailable:        "Unavailable",
	codes.DataLoss:           "DataLoss",
	codes.Unauthenticated:    "Unauthenticated",
}
