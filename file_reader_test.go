package mp3reader_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/zuiwuchang/mp3reader"
)

func TestFileReader(t *testing.T) {
	filename := "bin/a.mp3"
	// filename = "bin/b.mp3"
	// filename = "bin/0.mp3"
	f, e := os.Open(filename)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	r, e := mp3reader.NewFileReader(f, mp3reader.WithFileReaderAllowInvalidFrame(false))
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r.Duration(), r.Frames())
	fmt.Println(r.V1)
	fmt.Println(r.V2)

	frames := 0
	for i := 0; i < 2; i++ {
		frame, e := r.ReadFrame()
		if e != nil {
			if e != io.EOF {
				t.Fatal(e)
			}
			return
		}

		fmt.Println(frames, frame)
		frames++
	}

	for {
		frame, e := r.ReadFrame()
		if e != nil {
			if e != io.EOF {
				t.Fatal(e)
			}
			break
		}

		fmt.Println(frames, frame)
		frames++
	}
	fmt.Println(`frames=`, frames, r.Frames())
}
