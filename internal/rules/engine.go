package rules

import (
	"net/http"
	"strings"
)

type Engine struct {
	Rules []Rule
}

func NewEngine(rules []Rule) *Engine {
	return &Engine{Rules: rules}
}

func (e *Engine) Match(req *http.Request) *Rule {
	for i := range e.Rules {
		r := &e.Rules[i]
		if !r.Enabled {
			continue
		}

		if r.Method != "" && r.Method != req.Method {
			continue
		}

		if r.URLContains != "" &&
			!strings.Contains(req.URL.String(), r.URLContains) {
			continue
		}

		return r
	}

	return nil
}
