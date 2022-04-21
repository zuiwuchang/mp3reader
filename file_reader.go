package mp3reader

import (
	"bufio"
	"errors"
	"io"
)

type FileReader struct {
	*FrameReader
	size int64
	V2   ID3V2
	V1   ID3V1
}

func NewFileReader(f io.ReadSeeker, opt ...FileReaderOption) (r *FileReader, e error) {
	size, e := f.Seek(0, io.SeekEnd)
	if e != nil {
		return
	} else if size < 20 {
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
			v1.raw = raw
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
	var limit int64
	if string(rawHeader[:3]) == "ID3" {
		v2.rawHeader = rawHeader
		frameSize := v2.Size()
		limit = 10 + int64(frameSize)
		if limit+int64(len(v1.raw)) > size {
			e = errors.New("mp3: ID3 too large")
			return
		}
		v2.raw = make([]byte, frameSize)
		_, e = io.ReadAtLeast(f, v2.raw, len(v2.raw))
		if e != nil {
			return
		}
	} else {
		_, e = f.Seek(0, io.SeekStart)
		if e != nil {
			return
		}
	}
	var reader io.Reader
	if limit == 0 {
		reader = f
	} else {
		reader = io.LimitReader(f, size-limit-int64(len(v1.raw)))
	}
	opts := defaultFileReaderOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	if opts.bufferSize > 0 {
		reader = bufio.NewReaderSize(reader, opts.bufferSize)
	}
	r = &FileReader{
		FrameReader: NewFrameReader(
			reader,
			WithReaderAllowInvalidFrame(opts.allowInvalidFrame),
		),
		size: size,
		V2:   v2,
		V1:   v1,
	}
	return
}
