package util

import (
	"os/user"
	"strconv"
)

type UserID struct {
	UID uint32
	GID uint32
}

var root UserID

func GetUserId(name string) (UserID, error) {
	usr, err := user.Lookup(name)
	if err != nil {
		return root, err
	}
	uid, err := strconv.ParseInt(usr.Uid, 10, 32)
	if err != nil {
		return root, err
	}
	gid, err := strconv.ParseInt(usr.Gid, 10, 32)
	if err != nil {
		return root, err
	}
	return UserID{uint32(uid), uint32(gid)}, nil
}
