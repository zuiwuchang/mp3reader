package mp3reader

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type FileReader struct {
	frameReader *FrameReader
	size        int64
	V2          ID3V2
	V1          ID3V1

	first *Frame
	cache *Frame

	duration time.Duration
	frames   int64
	index    int64
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
	frameReader := NewFrameReader(
		reader,
		WithReaderAllowInvalidFrame(opts.allowInvalidFrame),
	)
	first, e := frameReader.ReadFrame()
	if e != nil {
		return
	}
	dataSize := size - int64(len(v1.raw)) - limit
	frames, duration := getDuration(first, dataSize)

	r = &FileReader{
		frameReader: frameReader,
		size:        dataSize,
		V2:          v2,
		V1:          v1,
		first:       first,
		cache:       first,
		duration:    duration,
		frames:      frames,
	}
	return
}

func (r *FileReader) ReadFrame() (frame *Frame, e error) {
	if r.cache != nil {
		frame = r.cache
		r.cache = nil
		return
	}
	if r.index >= r.frames {
		e = io.EOF
		return
	}
	frame, e = r.frameReader.ReadFrame()
	r.index++
	return
}
func (f *FileReader) Duration() time.Duration {
	return f.duration
}
func (f *FileReader) Frames() int64 {
	return f.frames
}
func getDuration(frame *Frame, dataSize int64) (int64, time.Duration) {
	raw := frame.raw
	if frame.header.CRC() {
		raw = raw[2:]
	}
	version := frame.header.Version()
	if version == Version1 {
		if frame.header.Channels() == 1 {
			raw = raw[17:]
		} else {
			raw = raw[32:]
		}
	} else {
		if frame.header.Channels() == 1 {
			raw = raw[9:]
		} else {
			raw = raw[17:]
		}
	}

	tag := string(raw[:4])
	if tag == "Info" || tag == "Xing" {
		raw = raw[4:]
		flag := binary.BigEndian.Uint32(raw)
		raw = raw[4:]
		if flag&0x1 != 0 {
			frames := int64(binary.BigEndian.Uint32(raw))
			return frames, time.Duration(frames) * frame.header.Duration()
		}
		if flag&0x2 != 0 {
			size := binary.BigEndian.Uint32(raw)
			dataSize = int64(size)
		}
	}
	l := len(frame.header.raw) + len(frame.raw)
	return -1, time.Duration(dataSize/int64(l)) * frame.header.Duration()
}

func Duration(frame *Frame, dataSize int64) time.Duration {
	_, val := getDuration(frame, dataSize)
	return val
}
