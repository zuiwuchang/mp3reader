package mp3reader

import (
	"errors"
	"fmt"
	"io"
)

type FrameReader struct {
	r    io.Reader
	opts readerOptions
}

func NewFrameReader(r io.Reader, opt ...ReaderOption) *FrameReader {
	opts := defaultReaderOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	return &FrameReader{
		r:    r,
		opts: opts,
	}
}

func (f *FrameReader) ReadFrame() (frame *Frame, e error) {
	header, frameSize, e := f.readHeader()
	if e != nil {
		return
	}
	raw := make([]byte, frameSize-4)
	_, e = io.ReadAtLeast(f.r, raw, len(raw))
	if e != nil {
		return
	}
	frame = &Frame{
		header: header,
		raw:    raw,
	}
	return
}
func (f *FrameReader) readHeader() (header FrameHeader, frameSize int, e error) {
	raw := make([]byte, 4)
	_, e = io.ReadAtLeast(f.r, raw, 4)
	if e != nil {
		return
	}
	header = FrameHeader{
		raw: raw,
	}
	frameSize = header.frameSize()
	if frameSize > -1 {
		return
	} else if !f.opts.allowInvalidFrame {
		e = errors.New(`mp3: FrameHeader invalid`)
		return
	}
	for {
		copy(raw, raw[1:])
		_, e = io.ReadAtLeast(f.r, raw[3:], 1)
		if e != nil {
			return
		}
		header = FrameHeader{
			raw: raw,
		}
		frameSize = header.frameSize()
		if frameSize > -1 {
			break
		}
	}
	return
}

type Frame struct {
	header FrameHeader
	raw    []byte
}

func (f *Frame) String() string {
	return fmt.Sprintf(`Frame{
	%s
	raw=%d duration=%s
}`,
		f.header,
		len(f.raw), f.header.Duration(),
	)
}
func (f *Frame) Header() FrameHeader {
	return f.header
}

func (f *Frame) Raw() []byte {
	return f.raw
}
