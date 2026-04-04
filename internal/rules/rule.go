// Package rules
package rules

type RuleType string

const (
	RuleMock     RuleType = "mock"
	RuleModify   RuleType = "modify"
	RudeRedirect RuleType = "redirect"
)

type Rule struct {
	ID      string   `json:"id"`
	Enabled bool     `json:"enabled"`
	Type    RuleType `json:"type"`

	// Match conditions
	URLContains string `json:"url_contains"`
	Method      string `json:"method"`

	// Actions
	FilePath string            `json:"file_path,omitempty"`
	Redirect string            `json:"redirect,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`

	// JSON modification
	JSONPatch map[string]any `json:"json_patch,omitempty"`
}
