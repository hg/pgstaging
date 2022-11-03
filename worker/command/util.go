package command

import (
	"errors"
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
	conf, err := os.OpenFile(path, os.O_WRONLY, os.ModeAppend)
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
