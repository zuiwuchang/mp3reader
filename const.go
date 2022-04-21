package mp3reader

import "strconv"

type Version uint8

const (
	Version2_5 Version = iota
	VersionUndefined
	Version2
	Version1
)

func (v Version) String() string {
	switch v {
	case Version2_5:
		return `MPEG 2.5`
	case VersionUndefined:
		return `MPEG Undefined`
	case Version2:
		return `MPEG 2`
	case Version1:
		return `MPEG 1`
	}
	return `Unkonw MPEG ` + strconv.Itoa(int(v))
}

type Layer uint8

const (
	LayerUndefined Layer = iota
	Layer3
	Layer2
	Layer1
)

func (v Layer) String() string {
	switch v {
	case LayerUndefined:
		return `Layer Undefined`
	case Layer3:
		return `Layer 3`
	case Layer2:
		return `Layer 2`
	case Layer1:
		return `Layer 1`
	}
	return `Unkonw Layer ` + strconv.Itoa(int(v))
}

type Mode uint8

const (
	Stereo Mode = iota
	JointStereo
	DualChannel
	SingleChannel
)

func (v Mode) String() string {
	switch v {
	case Stereo:
		return `Stereo`
	case JointStereo:
		return `Joint Stereo`
	case DualChannel:
		return `Dual Channel`
	case SingleChannel:
		return `Single Channel`
	}
	return `Unkonw Mode ` + strconv.Itoa(int(v))
}
