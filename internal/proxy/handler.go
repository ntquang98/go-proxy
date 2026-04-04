package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/elazarl/goproxy"
	"github.com/ntquang98/go-proxy/internal/modifier"
	"github.com/ntquang98/go-proxy/internal/rules"
)

type Handler struct {
	engine *rules.Engine
}

func NewHandler(engine *rules.Engine) *Handler {
	return &Handler{engine: engine}
}

// handle redirect and mock
func (h *Handler) HandleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	rule := h.engine.Match(req)
	log.Println("Incoming:", req.URL.String())
	if rule == nil {
		return req, nil
	}
	log.Println("Matched rule:", rule.ID)

	switch rule.Type {
	case rules.RudeRedirect:
		req.URL.Scheme = "https"
		req.URL.Host = rule.Redirect
	case rules.RuleMock:
		log.Println("MOCK HIT:", req.URL.String())
		log.Println("MOCK TO FILE:", rule.FilePath)
		data, err := os.ReadFile(rule.FilePath)
		if err != nil {
			log.Println("read file error: ", err, rule.FilePath)
			return req, nil
		}

		contentType := modifier.DetectContentType(rule.FilePath, data)
		resp := goproxy.NewResponse(
			req,
			contentType,
			http.StatusOK,
			"",
		)

		resp.Body = io.NopCloser(bytes.NewBuffer(data))
		resp.ContentLength = int64(len(data))

		resp.Header.Set("Content-Length", strconv.Itoa(len(data)))
		resp.Header.Set("Content-Type", contentType)
		resp.Header.Set("Cache-Control", "no-store")
		resp.Header.Set("X-Proxy", "mapmapsq-goproxy")

		resp.Header.Del("Content-Encoding")

		return req, resp
	}

	return req, nil
}

// handle modify JSON, replace file
func (h *Handler) HandleResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if ctx == nil || ctx.Req == nil {
		return resp
	}

	if resp == nil || resp.Request == nil {
		return resp
	}

	rule := h.engine.Match(resp.Request)
	if rule == nil {
		return resp
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	switch rule.Type {
	case rules.RuleModify:
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

	return resp
}
