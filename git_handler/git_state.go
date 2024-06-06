package githandler

import (
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

type GitStateStruct struct {
	F filehandler.FileHash `json:"git_state"`

	buffer_store chan *filehandler.BytesStore

	CommitMessage string `json:"commit_message"`
}

func (g *GitStateStruct) populate_buffer_store(n int) {
	if n < 1 {
		return
	}
	if g.buffer_store == nil {
		g.buffer_store = make(chan *filehandler.BytesStore)
	}

	for n > 0 {
		g.buffer_store <- filehandler.NewBytesStore()
		n -= 1
	}
}
