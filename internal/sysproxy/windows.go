//go:build windows

package sysproxy

import (
	"fmt"
	"os/exec"
)

type WindowsManager struct {
	host   string
	port   int
	backup ProxyState
}

func New(host string, port int) Manager {
	return &WindowsManager{
		host: host,
		port: port,
	}
}

func (w *WindowsManager) Backup() error {
	// TODO: read registry later
	// assume disabled
	w.backup = ProxyState{
		Enabled: false,
	}
	return nil
}

func (w *WindowsManager) Enable() error {
	server := fmt.Sprintf("%s:%d", w.host, w.port)
	cmds := [][]string{
		{
			"reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f",
		},
		{
			"reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyServer", "/t", "REG_SZ", "/d", server, "/f",
		},
	}

	for _, c := range cmds {
		if err := exec.Command(c[0], c[1:]...).Run(); err != nil {
			return err
		}
	}

	return nil
}

func (w *WindowsManager) Disable() error {
	cmd := exec.Command("reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")

	return cmd.Run()
}

func (w *WindowsManager) Restore() error {
	if w.backup.Enabled {
		return w.Enable()
	}
	return w.Disable()
}
