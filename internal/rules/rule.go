// Package rules
package rules

import "regexp"

type RuleType string

const (
	RuleMock     RuleType = "mock"
	RuleModify   RuleType = "modify"
	RudeRedirect RuleType = "redirect"
)

type Rule struct {
	ID       string `json:"id"`
	Enabled  bool   `json:"enabled"`
	Priority int    `json:"priority"`

	Type RuleType `json:"type"`

	// Match conditions
	Method      string `json:"method"`
	URLContains string `json:"url_contains"`
	URLRegex    string `json:"url_regex"`

	// Actions
	FilePath string `json:"file_path,omitempty"`
	Redirect string `json:"redirect,omitempty"`

	// Modify headers
	Headers map[string]string `json:"headers,omitempty"`

	// JSON modification
	JSONPatch map[string]any `json:"json_patch,omitempty"`

	// compiled (runtime only)
	compiledURL *regexp.Regexp
}
