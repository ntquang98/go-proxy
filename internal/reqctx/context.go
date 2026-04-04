// Package reqctx
package reqctx

type ctxKey string

const (
	traceIDKey ctxKey = "trace_id"
	loggerKey  ctxKey = "logger"
	ruleKey ctxKey = "rule"
)
