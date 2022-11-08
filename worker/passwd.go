package worker

import (
	"github.com/hg/pgstaging/worker/command"
	"crypto/subtle"
	"errors"
	"log"
	"os"
)

func validPassword(name string, pass string, admin string) bool {
	if subtle.ConstantTimeCompare([]byte(pass), []byte(admin)) != 0 {
		log.Print("admin password matched")
		return true
	}

	p := command.PathPasswd(name)

	ref, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return pass == ""
		}
		log.Printf("could not read password file: %v", err)
		return false
	}

	return subtle.ConstantTimeCompare(ref, []byte(pass)) != 0
}
