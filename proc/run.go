package proc

import (
	"github.com/hg/pgstaging/util"
	"log"
	"os/exec"
	"syscall"
)

func run(uid, gid uint32, commands [][]string) error {
	for _, args := range commands {
		log.Printf("running %v as uid %d", args, uid)

		cmd := exec.Command(args[0], args[1:]...)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uid,
				Gid: gid,
			},
		}

		out, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("command (%v) failed: %v", args, string(out))
			return err
		}
	}

	return nil
}

func RunAs(user string, commands ...[]string) error {
	id, err := util.GetUserId(user)
	if err != nil {
		return err
	}
	return run(id.UID, id.GID, commands)
}

func Run(commands ...[]string) error {
	return run(0, 0, commands)
}
