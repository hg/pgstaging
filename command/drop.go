package command

import (
	"github.com/hg/pgstaging/proc"
)

func DropCluster(name string) error {
	return proc.Run(
		[]string{"pg_dropcluster", "--stop", version, name},
		[]string{"btrfs", "subvolume", "delete", pathTarget(name)},
		[]string{"systemctl", "daemon-reload"},
	)
}
