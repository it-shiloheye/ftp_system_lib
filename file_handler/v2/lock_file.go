package filehandler

import (
	"errors"
	"fmt"

	"os"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

type LockFile struct {
	name string
}

func Lock(file_path string) (lf *LockFile, err error) {
	err1 := os.MkdirAll(file_path, os.FileMode(os.ModeExclusive))
	if err1 != nil {
		if errors.Is(err1, os.ErrExist) {
			err = err1
			return
		}
		err = &ftp_context.LogItem{
			Location: fmt.Sprintf(`Lock("%s" string) (lf *LockFile, err error)`, file_path),
			Time:     time.Now(),
			After:    fmt.Sprintf(`err1 := os.MkdirAll("%s", os.FileMode(os.ModeExclusive))`, file_path),
			Message:  err1.Error(),
		}
		return
	}

	lf = &LockFile{
		name: file_path,
	}

	return
}

func (lf *LockFile) Unlock() error {
	err1 := os.Remove(lf.name)
	if err1 == nil {
		return nil
	}
	err2 := &ftp_context.LogItem{
		Location: `func (lf *LockFile) Unlock() error`,
		Time:     time.Now(),
		After:    fmt.Sprintf(`err1 := os.Remove("%s")`, lf.name),
		Message:  err1.Error(),
	}
	return err2
}
