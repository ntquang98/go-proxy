// Package sysproxy
package sysproxy

type ProxyState struct {
	Enabled bool
	Server  string
}

type Manager interface {
	Enable() error
	Disable() error
	Backup() error
	Restore() error
}
