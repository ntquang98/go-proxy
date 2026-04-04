package reqctx

import (
	"context"
	"strconv"
	"time"
)

func NewTraceID() string {
	// TODO: use UUID
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey).(string); ok {
		return v
	}

	return ""
}
