package command

import (
	"github.com/hg/pgstaging/proc"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

var (
	rePort = regexp.MustCompile(`port\s*=\s*(\d+)`)
)

func CreateCluster(name string) error {
	target := pathTarget(name)

	if fileExists(target) {
		return fmt.Errorf("path %s already exists", target)
	}

	err := stepCreateCluster(name)
	if err == nil {
		err = stepAllowAuth(name)
	}
	if err == nil {
		err = stepRemoveFiles(name)
	}
	if err == nil {
		err = stepCreateSubvolume(name)
	}
	if err == nil {
		err = stepStart(name)
	}
	if err == nil {
		err = stepSetPasswd(name)
	}
	return err
}

func stepRemoveFiles(name string) error {
	data := pathData(name)

	for _, file := range []string{
		"postgresql.auto.conf",
		"standby.signal",
		"postmaster.pid",
		"postmaster.opts",
	} {
		err := os.Remove(path.Join(data, file))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("could not remove old file: %v", err)
		}
	}

	return nil
}

func stepStart(name string) error {
	return proc.Run(
		[]string{"pg_conftool", version, name, "set", "listen_addresses", "*"},
		[]string{"systemctl", "daemon-reload"},
		[]string{"systemctl", "start", service(name)},
	)
}

func stepSetPasswd(name string) error {
	cmd := exec.Command("pg_conftool", version, name, "get", "port")
	out, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("could not get port: %v", err)
	}

	m := rePort.FindSubmatch(out)
	if m == nil {
		return errors.New("port not found")
	}
	port := string(m[1])

	_, err = strconv.ParseUint(port, 10, 16)
	if err != nil {
		return fmt.Errorf("could not parse port '%s': %v", port, err)
	}

	return proc.RunAs("postgres", []string{
		"psql",
		"--port", port,
		"--command", "ALTER USER sc PASSWORD 'sc'",
	})
}

func stepCreateSubvolume(name string) error {
	source := pathTarget("base")
	target := pathTarget(name)

	return proc.Run([]string{
		"btrfs", "subvolume", "snapshot", source, target,
	})
}

func stepCreateCluster(name string) error {
	return proc.RunAs("postgres", []string{
		"pg_createcluster",
		"-d", pathData(name),
		version, name,
		"--", "--no-sync",
	})
}

func stepAllowAuth(name string) error {
	hba := pathPgHba(name)
	err := appendText(hba, "\nhost all all 0.0.0.0/0 scram-sha-256\n")

	if err != nil {
		return fmt.Errorf("failed to update pg_hba.conf: %v", err)
	}
	return nil
}
