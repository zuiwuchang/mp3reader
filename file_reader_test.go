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
	filename = "bin/b.mp3"
	filename = "bin/0.mp3"
	f, e := os.Open(filename)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	r, e := mp3reader.NewFileReader(f, mp3reader.WithFileReaderAllowInvalidFrame(false))
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r.V1)
	fmt.Println(r.V2)
	total := 0
	for i := 0; i < 2; i++ {
		frame, e := r.ReadFrame()
		if e != nil {
			if e != io.EOF {
				t.Fatal(e)
			}
			return
		}

		fmt.Println(total, frame)
		total++
	}
	// return
	for {
		frame, e := r.ReadFrame()
		if e != nil {
			if e != io.EOF {
				t.Fatal(e)
			}
			break
		}

		fmt.Println(total, frame)
		total++
	}

	// dst, e := os.Create("bin/b.mp3")
	// if e != nil {
	// 	t.Fatal(e)
	// }
	// defer dst.Close()
	// n, e := io.Copy(dst, r.R())
	// if e != nil {
	// 	t.Fatal(e)
	// }
	// fmt.Println(n)

	// dst, e := os.Create("bin/ok.mp3")
	// if e != nil {
	// 	t.Fatal(e)
	// }
	// defer dst.Close()
	// for {
	// 	frame, e := r.ReadFrame()
	// 	if e != nil {
	// 		if e != io.EOF  {
	// 			t.Fatal(e)
	// 		}
	// 		break
	// 	}
	// 	_, e = dst.Write(frame.Header().Raw())
	// 	if e != nil {
	// 		t.Fatal(e)
	// 	}
	// 	_, e = dst.Write(frame.CRC())
	// 	if e != nil {
	// 		t.Fatal(e)
	// 	}
	// 	_, e = dst.Write(frame.Raw())
	// 	if e != nil {
	// 		t.Fatal(e)
	// 	}
	// }
}
