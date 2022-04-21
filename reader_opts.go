package mp3reader

var defaultReaderOptions = readerOptions{
	allowInvalidFrame: true,
}

type readerOptions struct {
	allowInvalidFrame bool
}

type ReaderOption interface {
	apply(*readerOptions)
}

type funcReaderOption struct {
	f func(*readerOptions)
}

func (fdo *funcReaderOption) apply(do *readerOptions) {
	fdo.f(do)
}
func newReaderFuncOption(f func(*readerOptions)) *funcReaderOption {
	return &funcReaderOption{
		f: f,
	}
}

func WithReaderAllowInvalidFrame(allowInvalidFrame bool) ReaderOption {
	return newReaderFuncOption(func(o *readerOptions) {
		o.allowInvalidFrame = allowInvalidFrame
	})
}
