package filehandler

import (
	"encoding/json"

	"io"
	"io/fs"
	"strings"
	"time"

	"os"

	"github.com/it-shiloheye/ftp_system_lib/base"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

type FileType string

type FileBasic struct {
	Path string   `json:"path"`
	Type FileType `json:"type"`
	Size int64    `json:"size"`
	// directly read write to the file
	*os.File
	Fs os.FileInfo
	d  fs.DirEntry
}

func init() {
	var _ io.ReadWriteCloser = &FileBasic{}
}

func (Fo *FileBasic) IsOpen() bool {
	return Fo.File != nil
}

func (Fo *FileBasic) ReadAll() (data []byte, err error) {
	data, err = io.ReadAll(Fo.File)
	if err != nil {
		err = ftp_context.NewLogItem("FileBasic.ReadAll", true).
			SetAfter("data, err = io.ReadAll(Fo.File)").
			Set("path", Fo.Path).AppendParentError(err).
			SetMessage(err.Error())

	}
	return
}

func (Fo *FileBasic) ModTime() string {
	return Fo.Fs.ModTime().Format(time.RFC822Z)
}

func Open(file_path string) (Fo *FileBasic, err error) {
	loc := "Fo.File, err = base.OpenFile(Fo.Path, os.O_RDWR|os.O_SYNC)"
	var err1, err2 error
	Fo = &FileBasic{
		Path: file_path,
	}
	Fo.File, err1 = base.OpenFile(file_path, os.O_RDWR|os.O_SYNC)
	if err1 != nil {
		err = ftp_context.NewLogItem(loc, true).
			SetAfter("base.OpenFile").
			Set("path", Fo.Path).
			SetMessagef("path: %s \nerror:\n%s", file_path, err1.Error()).
			AppendParentError(err1)

		return

	}

	Fo.Fs, err2 = Fo.File.Stat()
	if err2 != nil {
		err = ftp_context.NewLogItem(loc, true).
			SetAfter("Fo.Fs, err2 = Fo.File.Stat()").
			SetMessagef("path: %s \nerror:\n%s", file_path, err2.Error()).AppendParentError(err2)
		return
	}

	Fo.Type = Ext(Fo)
	return
}

func Ext(Fo *FileBasic) FileType {
	if len(Fo.Type) > 0 {
		return Fo.Type
	}
	if Fo.d.IsDir() {
		Fo.Type = "dir"
		return Fo.Type
	}
	stp_1 := strings.Split(Fo.Path, ".")
	stp_2 := len(stp_1)
	stp_3 := stp_1[stp_2-1]
	if len(stp_3) > 4 {
		Fo.Type = "unknown"
		return Fo.Type
	}

	Fo.Type = FileType(stp_3)
	return Fo.Type
}

func Create(file_path string) (Fo *FileBasic, err error) {
	loc := "Create(file_path string) (Fo *FileBasic,err error)"
	var err1, err2 error
	Fo = &FileBasic{
		Path: file_path,
	}
	Fo.File, err1 = base.OpenFile(file_path, os.O_RDWR|os.O_SYNC|os.O_CREATE)
	if err != nil {
		err = ftp_context.NewLogItem(loc, true).
			SetAfter("Fo.File, err1 = base.OpenFile(file_path, os.O_RDWR|os.O_SYNC|os.O_CREATE)").
			Set("path", Fo.Path).
			SetMessagef("path: %s \nerror:\n%s", file_path, err1.Error()).
			AppendParentError(err1)
		return

	}

	Fo.Fs, err2 = Fo.File.Stat()
	if err2 != nil {
		err = ftp_context.NewLogItem(loc, true).
			SetAfter("Fo.Fs, err2 = Fo.File.Stat()").
			SetMessagef("path: %s \nerror:\n%s", file_path, err2.Error()).AppendParentError(err2)
		return
	}

	Fo.Type = Ext(Fo)
	return
}

func NewFileBasic(path string) (Fo *FileBasic) {
	Fo = &FileBasic{
		Path: path,
	}
	return
}

func WriteJson(Fo *FileBasic, v any) (err error) {
	loc := "WriteJson(Fo *FileBasic, v any)(err error)"
	var err1, err2 error
	t, err1 := json.MarshalIndent(v, " ", "\t")
	if err1 != nil {
		err = ftp_context.NewLogItem(loc, true).
			SetAfter(`t, err1 := json.MarshalIndent(v, " ", "\t")`).
			Set("path", Fo.Path).
			SetMessagef("path: %s \nerror:\n%s", Fo.Path, err1.Error()).
			AppendParentError(err1)
		return
	}

	_, err2 = Fo.Write(t)
	if err2 != nil {
		err = ftp_context.NewLogItem(loc, true).
			SetMessagef("_, err2 = Fo.Write(t)").
			SetMessagef("path: %s \nerror:\n%s", Fo.Path, err2.Error()).
			AppendParentError(err2)
		return
	}
	return nil
}
