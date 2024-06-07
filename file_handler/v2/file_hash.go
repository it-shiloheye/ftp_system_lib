package filehandler

import (
	"fmt"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

type FileHash struct {
	*FileBasic
	Hash    string `json:"hash"`
	ModTime string `json:"last_mod_time"`
}

func HashFile(Fo *FileBasic, bs *BytesStore) (hash string, err error) {
	loc := " NewFileHash(Fo *FileBasic, bs *BytesStore)(hash string, err error)"
	if Fo == nil || Fo.File == nil {
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			Message:  "FileBasic or os.File pointer is nil",
		}
		return
	}
	if bs == nil || bs.h == nil {
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			Message:  "ByteStore pointer provided is invalid",
		}
		return
	}

	var err1, err2 error
	bs.Reset()

	_, err1 = bs.ReadFrom(Fo.File)
	if err1 != nil {
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			After:    `_, err1 = bs.ReadFrom(Fo.File)`,
			Message:  err1.Error(),
			Err:      true, CallStack: []error{err1},
		}
		return
	}

	hash, err2 = bs.Hash()
	if err2 != nil {
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			After:    `hash, err2 = bs.Hash()`,
			Message:  err2.Error(),
			Err:      true, CallStack: []error{err2},
		}
		return
	}
	return
}

func NewFileHashOpen(file_path string) (Fh *FileHash, err error) {
	loc := "NewFileHashOpen(file_path string) (Fh *FileHash, err error)"
	Fh = &FileHash{}
	var err1 error
	Fh.FileBasic, err1 = Open(file_path)
	if err1 != nil {
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			After:    fmt.Sprintf(`Fh.FileBasic, err1 = Open("%s")`, file_path),
			Message:  err1.Error(),
			Err:      true, CallStack: []error{err1},
		}
		return
	}

	Fh.ModTime = fmt.Sprint(Fh.fs.ModTime())

	return
}

func NewFileHashCreate(file_path string) (Fh *FileHash, err error) {
	loc := "NewFileHashOpen(file_path string) (Fh *FileHash, err error)"
	Fh = &FileHash{}
	var err1 error
	Fh.FileBasic, err1 = Create(file_path)
	if err1 != nil {
		err = &ftp_context.LogItem{
			Location: loc,
			Time:     time.Now(),
			After:    fmt.Sprintf(`Fh.FileBasic, err1 = Create("%s")`, file_path),
			Message:  err1.Error(),
			Err:      true, CallStack: []error{err1},
		}
		return
	}

	return
}
