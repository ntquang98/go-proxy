//go:build darwin

package sysproxy

import (
	"fmt"
	"os/exec"
)

type MacManager struct {
	service string
	host    string
	port    int
}

func New(host string, port int) Manager {
	return &MacManager{
		service: "Wi-Fi",
		host:    host,
		port:    port,
	}
}

func (m *MacManager) Backup() error {
	// TODO: read current proxy via networksetup -getwebproxy
	return nil
}

func (m *MacManager) Enable() error {
	portStr := fmt.Sprintf("%d", m.port)
	cmds := [][]string{
		{"networksetup", "-setwebproxy", m.service, m.host, portStr},
		{"networksetup", "-setsecurewebproxy", m.service, m.host, portStr},
	}

	for _, c := range cmds {
		if err := exec.Command(c[0], c[1:]...).Run(); err != nil {
			return err
		}
	}

	return nil
}

func (m *MacManager) Disable() error {
	cmds := [][]string{
		{"networksetup", "-setwebproxystate", m.service, "off"},
		{"networksetup", "-setsecurewebproxystate", m.service, "off"},
	}

	for _, c := range cmds {
		if err := exec.Command(c[0], c[1:]...).Run(); err != nil {
			return err
		}
	}

	return nil
}

func (m *MacManager) Restore() error {
	return m.Disable()
}
