package mp3reader

import (
	"encoding/binary"
	"fmt"
	"time"
)

type FrameHeader struct {
	raw []byte
}

func (f FrameHeader) String() string {
	flag := f.flag()
	sync := flag & 0xFFE0
	low := uint8(flag)
	version := Version((low & 0x18) >> 3)
	layer := Layer((low & 0x6) >> 1)
	crc := low&0x1 != 0
	low = f.raw[3]
	mode := Mode(low >> 6)
	return fmt.Sprintf(`FrameHeader{ sync=0x%X, version=%s, layer=%s, crc=%v, samples=%d, bitrate=%d, samplerate=%d, padding=%v, frameSize=%d, duration=%s, mode=%s, channels=%d }`,
		sync,
		version,
		layer,
		crc,
		f.Samples(), f.Bitrate(), f.Samplerate(),
		f.Padding(), f.FrameSize(), f.Duration(),
		mode, f.Channels(),
	)
}
func (f FrameHeader) Raw() []byte {
	return f.raw
}
func (f FrameHeader) frameSize() int {
	flag := f.flag()
	sync := flag & 0xFFE0
	if sync != 0xFFE0 {
		return -1
	}
	low := uint8(flag)
	layer := Layer((low & 0x6) >> 1)
	if layer == LayerUndefined {
		return -1
	}
	size := f.FrameSize()
	if size == -1 {
		return -1
	}
	return size
}
func (f FrameHeader) IsValid() bool {
	return f.frameSize() > -1
}
func (f FrameHeader) flag() uint16 {
	return binary.BigEndian.Uint16(f.raw)
}
func (f FrameHeader) Sync() uint16 {
	return f.flag() & 0xFFE0
}
func (f FrameHeader) Version() Version {
	low := uint8(f.flag())
	return Version((low & 0x18) >> 3)
}
func (f FrameHeader) Layer() Layer {
	low := uint8(f.flag())
	return Layer((low & 0x6) >> 1)
}

func (f FrameHeader) CRC() bool {
	low := uint8(f.flag())
	return low&0x1 == 0
}

var bitrates = [2][3][16]int{
	{
		// MPEG 1 Layer 3
		{0, 32000, 40000, 48000, 56000, 64000, 80000, 96000,
			112000, 128000, 160000, 192000, 224000, 256000, 320000},

		// MPEG 1 Layer 2
		{0, 32000, 48000, 56000, 64000, 80000, 96000, 112000,
			128000, 160000, 192000, 224000, 256000, 320000, 384000},

		// MPEG 1 Layer 1
		{0, 32000, 64000, 96000, 128000, 160000, 192000, 224000,
			256000, 288000, 320000, 352000, 384000, 416000, 448000},
	},
	{
		// MPEG2 2 Layer 3
		{0, 8000, 16000, 24000, 32000, 40000, 48000, 56000,
			64000, 80000, 96000, 112000, 128000, 144000, 160000},

		// MPEG 2 Layer 2
		{0, 8000, 16000, 24000, 32000, 40000, 48000, 56000,
			64000, 80000, 96000, 112000, 128000, 144000, 160000},

		// MPEG 2 Layer 1
		{0, 32000, 48000, 56000, 64000, 80000, 96000, 112000,
			128000, 144000, 160000, 176000, 192000, 224000, 256000},
	},
}

func (f FrameHeader) Bitrate() int {
	low := uint8(f.flag())
	version := Version((low & 0x18) >> 3)
	if version == VersionUndefined {
		return -1
	}
	layer := Layer((low & 0x6) >> 1)
	if layer == LayerUndefined {
		return -1
	}
	low = f.raw[2]
	index := low >> 4
	if index == 0xF {
		return -1
	}
	switch version {
	case Version1:
		return bitrates[0][layer-1][index]
	case Version2:
		fallthrough
	case Version2_5:
		return bitrates[1][layer-1][index]
	}
	return -1
}

var samplerate = [3][3]int{
	// MPEG 2.5
	{11025, 12000, 8000},
	// MPEG 2
	{22050, 24000, 16000},
	// MPEG 1
	{44100, 48000, 32000},
}

func (f FrameHeader) Samplerate() int {
	low := uint8(f.flag())
	version := Version((low & 0x18) >> 3)
	if version == VersionUndefined {
		return -1
	} else if version != Version2_5 {
		version--
	}
	index := f.raw[2] >> 2 & 0x3
	if index == 0x3 {
		return -1
	}
	return samplerate[version][index]
}
func (f FrameHeader) Padding() bool {
	return f.raw[2]>>1&0x1 != 0
}
func (f FrameHeader) Mode() Mode {
	return Mode(f.raw[3] >> 6)
}
func (f FrameHeader) Channels() int {
	if f.Mode() == SingleChannel {
		return 1
	}
	return 2
}
func (f FrameHeader) Samples() int {
	switch f.Layer() {
	case Layer1:
		return 384
	case Layer2:
		return 1152
	case Layer3:
		switch f.Version() {
		case Version1:
			return 1152
		case Version2:
			fallthrough
		case Version2_5:
			return 576
		}
	}
	return -1
}
func (f FrameHeader) FrameSize() int {
	bitrate := f.Bitrate()
	if bitrate == -1 {
		return -1
	}
	samples := f.Samples()
	if samples == -1 {
		return -1
	}
	samplerate := f.Samplerate()
	if samplerate == -1 {
		return -1
	}

	var padding int
	if f.Padding() {
		padding = 1
		layer := f.Layer()
		if layer == Layer1 {
			padding *= 4
		} else if layer == LayerUndefined {
			return -1
		}
	}
	val := samples/8*bitrate/samplerate + padding
	min := 4 + 32
	if f.CRC() {
		min += 2
	}
	if val < min {
		return -1
	}
	return val
}

func (f FrameHeader) Duration() time.Duration {
	samples := f.Samples()
	if samples == -1 {
		return -1
	}
	samplerate := f.Samplerate()
	if samplerate == -1 {
		return -1
	}
	return time.Duration(samples) * 1000 * 1000 * 1000 / time.Duration(samplerate) * time.Nanosecond
}
