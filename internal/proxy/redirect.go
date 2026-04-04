package proxy

import (
	"net/http"

	"github.com/ntquang98/go-proxy/internal/reqctx"
	"github.com/ntquang98/go-proxy/internal/rules"
)

func (h *Handler) handleRedirect(req *http.Request, rule *rules.Rule) (*http.Request, *http.Response) {
	logger := reqctx.GetLogger(req.Context())

	req.URL.Scheme = "https"
	req.URL.Host = rule.Redirect

	logger.Info("redirect applied", "target", rule.Redirect)

	return req, nil
}
