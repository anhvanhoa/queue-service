package logger

import (
	"context"
	"time"
)

type Log interface {
	Info(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	LogGRPC(ctx context.Context, method string, req any, resp any, err error, duration time.Duration)
	LogGRPCRequest(ctx context.Context, method string, req any)
	LogGRPCResponse(ctx context.Context, method string, resp any, err error, duration time.Duration)
}
