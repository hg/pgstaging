package command

import (
	"github.com/hg/pgstaging/proc"
	"github.com/hg/pgstaging/util"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

var (
	rePort = regexp.MustCompile(`port\s*=\s*(\d+)`)
)

func CreateCluster(name, pass string) error {
	target := pathTarget(name)

	if util.FileExists(target) {
		return fmt.Errorf("path %s already exists", target)
	}

	err := stepCreateCluster(name)
	if err == nil {
		err = stepRemovePgData(name)
	}
	if err == nil {
		err = stepAllowAuth(name)
	}
	if err == nil {
		err = stepCreateSubvolume(name)
	}
	if err == nil {
		err = stepRemoveFiles(name)
	}
	if err == nil {
		err = stepConfigure(name)
	}
	if err == nil {
		err = StartCluster(name)
	}
	if err == nil {
		err = stepSetPasswd(name, pass)
	}
	if err == nil && pass != "" {
		err = stepStorePasswd(name, pass)
	}
	return err
}

func stepStorePasswd(name string, pass string) error {
	p := PathPasswd(name)
	log.Printf("saving password to %s", p)
	return os.WriteFile(p, []byte(pass), 0o600)
}

func stepRemovePgData(name string) error {
	return os.RemoveAll(pathTarget(name))
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

func stepConfigure(name string) error {
	return proc.Run(
		[]string{"pg_conftool", version, name, "set", "listen_addresses", "*"},
		[]string{"systemctl", "daemon-reload"},
	)
}

func stepSetPasswd(name, pass string) error {
	if pass == "" {
		pass = "sc"
	}

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

	sql := fmt.Sprintf("ALTER USER sc PASSWORD '%s'", pass)

	return proc.RunAs("postgres", []string{
		"psql",
		"--port", port,
		"--command", sql,
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
