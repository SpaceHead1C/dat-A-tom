package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxLogger struct{}

func ContextWithLogger(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	if l, ok := ctx.Value(ctxLogger{}).(*zap.SugaredLogger); ok {
		return l
	}
	return zap.L().Sugar()
}
