package service

import (
	"github.com/hg/pgstaging/proc"
	"fmt"
	"io"
	"os"
)

const (
	appname = "pgstaging"
	service = appname + ".service"
	svcPath = "/etc/systemd/system/" + service
	binPath = "/usr/local/bin/" + appname
)

var content = fmt.Sprintf(`
[Unit]
Description = staging database server
After = network-online.target
Wants = network-online.target

[Service]
Type = simple
ExecStart = %s run
Restart = on-failure
RestartSec = 5

[Install]
WantedBy = multi-user.target
`, binPath)

func installService() (err error) {
	f, err := os.OpenFile(svcPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not create systemd service: %v", err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			err = e
		}
	}()
	_, err = f.WriteString(content)
	return
}

func installBinary() (err error) {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	in, err := os.Open(exe)
	if err != nil {
		return
	}
	out, err := os.OpenFile(binPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	if err = out.Sync(); err != nil {
		return
	}
	err = os.Chown(binPath, 0, 0)
	return
}

func Install() error {
	if err := installBinary(); err != nil {
		return fmt.Errorf("could not install binary: %v", err)
	}
	if err := installService(); err != nil {
		return fmt.Errorf("could not install service: %v", err)
	}
	return nil
}

func Enable() error {
	return proc.Run(
		[]string{"systemctl", "daemon-reload"},
		[]string{"systemctl", "enable", service},
	)
}
