package command

import (
	"fmt"
	"path"
)

const version = "15"
const mnt = "/opt/dev/mnt"

func pathTarget(name string) string {
	return path.Join(mnt, name)
}

func pathData(name string) string {
	return path.Join(pathTarget(name), "data")
}

func service(name string) string {
	return fmt.Sprintf("postgres@15-%s.service", name)
}
