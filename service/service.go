package service

import (
	"github.com/hg/pgstaging/proc"
	"fmt"
	"io"
	"os"
)

const appname = "stagingdb"
const service = appname + ".service"

const svcPath = "/etc/systemd/system/" + service
const binPath = "/usr/local/bin/" + appname

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
	err = out.Sync()
	return
}

func Install() (err error) {
	if err = installBinary(); err == nil {
		err = installService()
	}
	return
}

func Enable() error {
	return proc.Run(
		[]string{"systemctl", "enable", service},
		[]string{"systemctl", "daemon-reload"},
	)
}
