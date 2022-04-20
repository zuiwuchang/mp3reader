package mp3reader

import (
	"errors"
	"io"
)

type FileReader struct {
	f    io.ReadSeeker
	size int64
	V2   ID3V2
	V1   ID3V1
}

func NewFileReader(f io.ReadSeeker) (r *FileReader, e error) {
	size, e := f.Seek(0, io.SeekEnd)
	if e != nil {
		return
	} else if size < 10 {
		e = errors.New("mp3: ReadSeeker too short")
		return
	}
	var (
		v2 ID3V2
		v1 ID3V1
	)
	// v1
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

	// v2
	_, e = f.Seek(0, io.SeekStart)
	if e != nil {
		return
	}
	rawHeader := make([]byte, 10)
	_, e = io.ReadAtLeast(f, rawHeader, 10)
	if e != nil {
		return
	}
	if string(rawHeader[:3]) == "ID3" {
		v2.rawHeader = rawHeader
		
	}

	r = &FileReader{size: size,
		f:  f,
		V2: v2,
		V1: v1,
	}
	return
}
