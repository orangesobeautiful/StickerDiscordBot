package utils

import (
	"os"
	"path/filepath"

	"backend/app/pkg/log"
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
