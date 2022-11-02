package command

import (
	"github.com/hg/pgstaging/proc"
)

func StartCluster(name string) error {
	return proc.Run(
		[]string{"systemctl", "start", service(name)},
	)
}
