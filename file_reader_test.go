package mp3reader_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/zuiwuchang/mp3reader"
)

func TestFileReader(t *testing.T) {
	{
		f, e := os.OpenFile("bin/a.mp3", os.O_RDONLY|os.O_WRONLY, 0)
		if e != nil {
			return
		}
		_, e = f.Seek(-128, io.SeekEnd)
		if e != nil {
			return
		}
		_, e = f.Seek(33, io.SeekCurrent)
		if e != nil {
			return
		}
		_, e = f.Write([]byte("梅豔芳"))
		if e != nil {
			return
		}
	}

	f, e := os.Open("bin/a.mp3")
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	r, e := mp3reader.NewFileReader(f)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r.V1)
	fmt.Println(r.V1.Raw())
	// b, e := ioutil.ReadFile("bin/a.mp3")
	// if e != nil {
	// 	t.Fatal(e)
	// }
	// fmt.Println(string(b[:10]))
	// fmt.Println(string(b[len(b)-128:]))
}
