package mp3reader

import "fmt"

type ID3V2 struct {
	rawHeader []byte
	raw       []byte
}

func (v2 ID3V2) String() string {
	if v2.rawHeader != nil {
		f0, f1, f2 := v2.Flag()
		return fmt.Sprintf(`ID3V2{
	%s, version=%d, revision=%d, flag=%v %v %v, size=%d
	len(frame)=%d
}`,
			v2.TAG(),
			v2.Version(), v2.Revision(),
			f0, f1, f2,
			v2.Size(),
			len(v2.raw),
		)
	}
	return `ID3V2{}`
}

// return raw of header
func (v2 ID3V2) RawHeader() []byte {
	return v2.rawHeader
}

// return raw of frame
func (v2 ID3V2) Raw() []byte {
	return v2.raw
}
func (v2 ID3V2) TAG() string {
	if v2.rawHeader != nil {
		return string(v2.rawHeader[:3])
	}
	return ``
}
func (v2 ID3V2) Version() uint8 {
	if v2.rawHeader != nil {
		return v2.rawHeader[3]
	}
	return 0
}
func (v2 ID3V2) Revision() uint8 {
	if v2.rawHeader != nil {
		return v2.rawHeader[4]
	}
	return 0
}
func (v2 ID3V2) Flag() (bool, bool, bool) {
	if v2.rawHeader != nil {
		flag := v2.rawHeader[5]
		return flag&0x80 != 0, flag&0x40 != 0, flag&0x20 != 0
	}
	return false, false, false
}
func (v2 ID3V2) Size() uint32 {
	if v2.rawHeader != nil {
		buf := v2.rawHeader[6:10]
		return (uint32(buf[0]&0x7F) << 21) | (uint32(buf[1]&0x7F) << 14) |
			(uint32(buf[2]&0x7F) << 7) | uint32(buf[3]&0x7F)
	}
	return 0
}
