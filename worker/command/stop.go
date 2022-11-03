package command

import (
	"github.com/hg/pgstaging/proc"
)

func StopCluster(name string) error {
	return proc.Run(
		[]string{"systemctl", "stop", service(name)},
	)
}
