package reqctx

import (
	"context"

	"github.com/ntquang98/go-proxy/internal/rules"
)

func WithRule(ctx context.Context, rule *rules.Rule) context.Context {
	return context.WithValue(ctx, ruleKey, rule)
}

func GetRule(ctx context.Context) *rules.Rule {
	if r, ok := ctx.Value(ruleKey).(*rules.Rule); ok {
		return r
	}

	return nil
}
