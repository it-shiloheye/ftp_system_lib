package filehandler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
)

type BytesStore struct {
	h hash.Hash
	bytes.Buffer
}

func (bs *BytesStore) Hash() (hash string, err error) {
	loc := "func (bs *BytesStore) Hash() (hash string, err error)"
	bs.h.Reset()
	_, err1 := bs.WriteTo(bs.h)
	if err1 != nil {
		err = ftp_context.NewLogItem(loc, true).
			SetAfter("_, err = bs.CopyTo(bs.h)").
			SetMessage(err1.Error()).
			AppendParentError(err1)
		return
	}
	hash = fmt.Sprintf("%x", bs.h.Sum(nil))

	return
}

func NewBytesStore() (bs *BytesStore) {
	bs = &BytesStore{
		h:      sha256.New(),
		Buffer: bytes.Buffer{},
	}
	bs.Grow(100_000)
	return
}
