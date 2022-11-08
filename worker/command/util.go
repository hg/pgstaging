package command

import (
	"io"
	"os"
)

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
