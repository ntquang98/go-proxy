package rules

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type Engine struct {
	Rules []Rule
}

func NewEngine(rules []Rule) *Engine {
	return &Engine{Rules: rules}
}

func BuildEngine(rules []Rule) (*Engine, error) {
	for i := range rules {
		r := &rules[i]

		if r.URLRegex != "" {
			re, err := regexp.Compile(r.URLRegex)
			if err != nil {
				return nil, fmt.Errorf("rule %s invalid regex: %w", r.ID, err)
			}
			r.compiledURL = re
		}
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	return &Engine{Rules: rules}, nil
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

		if r.compiledURL != nil {
			if !r.compiledURL.MatchString(req.URL.String()) {
				continue
			}
		} else if r.URLContains != "" {
			if !strings.Contains(req.URL.String(), r.URLContains) {
				continue
			}
		}

		return r
	}

	return nil
}
