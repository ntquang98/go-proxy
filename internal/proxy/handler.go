package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/ntquang98/go-proxy/internal/reqctx"
	"github.com/ntquang98/go-proxy/internal/rules"
)

type Handler struct {
	engine *rules.Engine
}

func NewHandler(engine *rules.Engine) *Handler {
	return &Handler{engine: engine}
}

func (h *Handler) enrichContext(req *http.Request) *http.Request {
	traceID := reqctx.NewTraceID()
	ctx := reqctx.WithTraceID(req.Context(), traceID)
	logger := slog.With(
		"trace_id", traceID,
		"method", req.Method,
		"url", req.URL.String(),
	)
	ctx = reqctx.WithLogger(ctx, logger)
	return req.WithContext(ctx)
}

func (h *Handler) attachRule(req *http.Request) (*http.Request, *rules.Rule) {
	rule := h.engine.Match(req)
	if rule == nil {
		return req, nil
	}

	ctx := reqctx.WithRule(req.Context(), rule)
	return req.WithContext(ctx), rule
}

// handle redirect and mock
func (h *Handler) HandleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	req = h.enrichContext(req)

	logger := reqctx.GetLogger(req.Context())
	start := time.Now()

	logger.Debug("incoming request")

	req, rule := h.attachRule(req)
	if rule == nil {
		return req, nil
	}

	logger.Info("rule matched",
		"rule_id", rule.ID,
		"type", rule.Type,
	)

	defer func() {
		logger.Info("request done",
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}()

	switch rule.Type {
	case rules.RudeRedirect:
		return h.handleRedirect(req, rule)

	case rules.RuleMock:
		return h.handleMock(req, rule)
	}

	return req, nil
}

// handle modify JSON, replace file
func (h *Handler) HandleResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil || resp.Request == nil {
		return resp
	}

	ctxReq := resp.Request.Context()
	logger := reqctx.GetLogger(ctxReq)
	rule := reqctx.GetRule(ctxReq)
	if rule == nil {
		return resp
	}

	logger.Info("processing response", "rule_id", rule.ID)

	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	switch rule.Type {
	case rules.RuleModify:
		logger.Info("applying json patch")
		var jsonData map[string]any
		json.Unmarshal(bodyBytes, &jsonData)

		for k, v := range rule.JSONPatch {
			jsonData[k] = v
		}
		bodyBytes, _ = json.Marshal(jsonData)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	resp.ContentLength = int64(len(bodyBytes))
	resp.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
	resp.Header.Del("Content-Encoding")

	logger.Info("response modified", "size", len(bodyBytes))

	return resp
}
