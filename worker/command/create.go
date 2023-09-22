package command

import (
	"errors"
	"fmt"
	"github.com/hg/pgstaging/proc"
	"github.com/hg/pgstaging/util"
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

func CreateCluster(name, pass string, force bool) Result {
	target := pathTarget(name)

	if util.FileExists(target) {
		if !force {
			return Result{Err: fmt.Errorf("path %s already exists", target)}
		}
		if err := DropCluster(name); err != nil {
			return Result{Err: fmt.Errorf("could not drop cluster")}
		}
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
	var port uint16
	if err == nil {
		port, err = stepGetPort(name)
	}
	if err == nil {
		err = stepSetPasswd(pass, port)
	}
	if err == nil && pass != "" {
		err = stepStorePasswd(name, pass)
	}
	if err == nil {
		return Result{Data: port}
	}
	return Result{Err: err}
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

func stepGetPort(name string) (uint16, error) {
	cmd := exec.Command("pg_conftool", version, name, "get", "port")
	out, err := cmd.Output()

	if err != nil {
		return 0, fmt.Errorf("could not get port: %v", err)
	}

	m := rePort.FindSubmatch(out)
	if m == nil {
		return 0, errors.New("port not found")
	}
	raw := string(m[1])

	port, err := strconv.ParseUint(raw, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("could not parse port '%s': %v", raw, err)
	}

	return uint16(port), nil
}

func stepSetPasswd(pass string, port uint16) error {
	if pass == "" {
		pass = "sc"
	}

	sql := fmt.Sprintf("ALTER USER sc PASSWORD '%s'", pass)

	return proc.RunAs("postgres", []string{
		"psql",
		"--port", strconv.Itoa(int(port)),
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
