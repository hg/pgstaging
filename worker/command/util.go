package command

import (
	"errors"
	"io"
	"os"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func appendText(path, text string) error {
	conf, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}

	_, err = conf.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	_, err = conf.WriteString(text)
	closeErr := conf.Close()

	if err == nil {
		return closeErr
	}
	return err
}
