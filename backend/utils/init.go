package utils

import (
	"backend/pkg/log"
	"os"
	"path/filepath"
)

func Init() error {
	var err error
	exePath, err = os.Executable()
	if err != nil {
		log.Errorf("os.Executable failed, err=%s", err)
		return err
	}
	exeDir = filepath.Dir(exePath)

	return nil
}
