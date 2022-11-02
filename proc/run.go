package proc

import (
	"log"
	"os/exec"
)

func Run(commands ...[]string) error {
	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		out, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("command (%v) failed: %v", args, string(out))
			return err
		}
	}

	return nil
}
