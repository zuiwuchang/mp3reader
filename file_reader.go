package mp3reader

import (
	"io"
)

type FileReader struct {
	f    io.ReadSeeker
	size int64
	V1   ID3V1
}

func NewFileReader(f io.ReadSeeker) (r *FileReader, e error) {
	size, e := f.Seek(0, io.SeekEnd)
	if e != nil {
		return
	}
	var (
		v1 ID3V1
	)
	if size > 128 {
		_, e = f.Seek(-128, io.SeekEnd)
		if e != nil {
			return
		}
		raw := make([]byte, 128)
		_, e = io.ReadAtLeast(f, raw, len(raw))
		if e != nil {
			return
		}
		if string(raw[:3]) == `TAG` {
			v1 = ID3V1{
				raw: raw,
			}
		}
	}
	r = &FileReader{size: size,
		f:  f,
		V1: v1,
	}
	return
}
