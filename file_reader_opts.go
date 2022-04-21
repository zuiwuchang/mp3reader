package mp3reader

const defaultBufferSize = 1024 * 32

var defaultFileReaderOptions = fileReaderOptions{
	bufferSize:        defaultBufferSize,
	allowInvalidFrame: false,
}

type fileReaderOptions struct {
	bufferSize        int
	allowInvalidFrame bool
}

type FileReaderOption interface {
	apply(*fileReaderOptions)
}

type funcFileReaderOption struct {
	f func(*fileReaderOptions)
}

func (fdo *funcFileReaderOption) apply(do *fileReaderOptions) {
	fdo.f(do)
}
func newFileReaderFuncOption(f func(*fileReaderOptions)) *funcFileReaderOption {
	return &funcFileReaderOption{
		f: f,
	}
}

func WithFileReaderBufferSize(bufferSize int) FileReaderOption {
	return newFileReaderFuncOption(func(o *fileReaderOptions) {
		if bufferSize < 0 {
			bufferSize = defaultBufferSize
		} else {
			o.bufferSize = bufferSize
		}
	})
}
func WithFileReaderAllowInvalidFrame(allowInvalidFrame bool) FileReaderOption {
	return newFileReaderFuncOption(func(o *fileReaderOptions) {
		o.allowInvalidFrame = allowInvalidFrame
	})
}
