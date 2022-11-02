package command

import (
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func CreateCluster(name string) error {
	target := pathTarget(name)

	if fileExists(target) {
		return fmt.Errorf("path %s already exists", target)
	}

	user, err := user.Lookup("postgres")
	if err != nil {
		return err
	}

	uid, err := strconv.ParseInt(user.Uid, 10, 32)
	if err != nil {
		return err
	}
	gid, err := strconv.ParseInt(user.Gid, 10, 32)
	if err != nil {
		return err
	}

	cmd := exec.Command("pg_createcluster", "-d", pathData(name), version, name, "--", "--no-sync")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("could not create cluster: %v", err)
	}

	configPath := fmt.Sprintf("/etc/postgresql/$version/%s/pg_hba.conf", name)
	err = appendText(configPath, "\nhost all all 0.0.0.0/0 scram-sha-256\n")

	if err != nil {
		return fmt.Errorf("failed to update pg_hba.conf: %v", err)
	}

	return nil
}
