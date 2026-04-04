package proxy

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/elazarl/goproxy"
	"github.com/ntquang98/go-proxy/internal/modifier"
	"github.com/ntquang98/go-proxy/internal/reqctx"
	"github.com/ntquang98/go-proxy/internal/rules"
)

func (h *Handler) handleMock(req *http.Request, rule *rules.Rule) (*http.Request, *http.Response) {
	logger := reqctx.GetLogger(req.Context())

	data, err := os.ReadFile(rule.FilePath)
	if err != nil {
		logger.Error("read file failed", "error", err)
		return req, nil
	}

	contentType := modifier.DetectContentType(rule.FilePath, data)
	resp := goproxy.NewResponse(req, contentType, http.StatusOK, "")

	resp.Body = io.NopCloser(bytes.NewBuffer(data))
	resp.ContentLength = int64(len(data))

	resp.Header.Set("Content-Type", contentType)
	resp.Header.Set("Content-Length", strconv.Itoa(len(data)))
	resp.Header.Set("Cache-Control", "no-store")
	resp.Header.Set("X-Proxy", "mapmapsq-goproxy")

	resp.Header.Del("Content-Encoding")

	logger.Info("mock response served", "file", rule.FilePath)

	return req, resp
}
