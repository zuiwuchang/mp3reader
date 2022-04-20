package mp3reader

import (
	"bytes"
	"fmt"
	"strconv"
)

type ID3V1 struct {
	raw []byte
}

func (v1 ID3V1) Raw() []byte {
	return v1.raw
}
func (v1 ID3V1) String() string {
	if v1.raw == nil {
		return `ID3V1{}`
	}
	return fmt.Sprintf(`ID3V1{ %s, title=%s, artist=%s, album=%s, year=%s, notes=%s, type=%s }`,
		v1.TAG(),
		v1.Title(),
		v1.Artist(),
		v1.Album(),
		v1.Year(),
		v1.Notes(),
		v1.Type(),
	)
}
func (v1 ID3V1) TAG() string {
	if v1.raw != nil {
		return string(v1.raw[:3])
	}
	return ``
}
func (v1 ID3V1) str(offset, size int) string {
	b := v1.raw[offset : offset+size]
	i := bytes.Index(b, []byte{0})
	if i == 0 {
		return ``
	} else if i < 0 {
		return string(b)
	}
	return string(b[:i])
}
func (v1 ID3V1) Title() string {
	if v1.raw != nil {
		return v1.str(3, 30)
	}
	return ``
}
func (v1 ID3V1) Artist() string {
	if v1.raw != nil {
		return v1.str(33, 30)
	}
	return ``
}
func (v1 ID3V1) Album() string {
	if v1.raw != nil {
		return v1.str(63, 30)
	}
	return ``
}
func (v1 ID3V1) Year() string {
	if v1.raw != nil {
		return v1.str(93, 4)
	}
	return ``
}
func (v1 ID3V1) Notes() string {
	if v1.raw != nil {
		return v1.str(97, 30)
	}
	return ``
}

func (v1 ID3V1) Type() string {
	if v1.raw != nil {
		i := int(v1.raw[127])
		if i < len(types) {
			return types[i]
		}
		return strconv.Itoa(i)
	}
	return ``
}

var types = []string{"Blues",
	"ClassicRock",
	"Country",
	"Dance",
	"Disco",
	"Funk",
	"Grunge",
	"Hip-Hop",
	"Jazz",
	"Metal",
	"NewAge",
	"Oldies",
	"Other",
	"Pop",
	"R&B",
	"Rap",
	"Reggae",
	"Rock",
	"Techno",
	"Industrial",
	"Alternative",
	"Ska",
	"Deathl",
	"Pranks",
	"Soundtrack",
	"Euro-Techno",
	"Ambient",
	"Trip-Hop",
	"Vocal",
	"Jazz+Funk",
	"Fusion",
	"Trance",
	"Classical",
	"Instrumental",
	"Acid",
	"House",
	"Game",
	"SoundClip",
	"Gospel",
	"Noise",
	"AlternRock",
	"Bass",
	"Soul",
	"Punk",
	"Space",
	"Meditative",
	"InstrumentalPop",
	"InstrumentalRock",
	"Ethnic",
	"Gothic",
	"Darkwave",
	"Techno-Industrial",
	"Electronic",
	"Pop-Folk",
	"Eurodance",
	"Dream",
	"SouthernRock",
	"Comedy",
	"Cult",
	"Gangsta",
	"Top40",
	"ChristianRap",
	"Pop/Funk",
	"Jungle",
	"NativeAmerican",
	"Cabaret",
	"NewWave",
	"Psychadelic",
	"Rave",
	"Showtunes",
	"Trailer",
	"Lo-Fi",
	"Tribal",
	"AcidPunk",
	"AcidJazz",
	"Polka",
	"Retro",
	"Musical",
	"Rock&Roll",
	"HardRock",
	"Folk",
	"Folk-Rock",
	"NationalFolk",
	"Swing",
	"FastFusion",
	"Bebob",
	"Latin",
	"Revival",
	"Celtic",
	"Bluegrass",
	"Avantgarde",
	"GothicRock",
	"ProgessiveRock",
	"PsychedelicRock",
	"SymphonicRock",
	"SlowRock",
	"BigBand",
	"Chorus",
	"EasyListening",
	"Acoustic",
	"Humour",
	"Speech",
	"Chanson",
	"Opera",
	"ChamberMusic",
	"Sonata",
	"Symphony",
	"BootyBass",
	"Primus",
	"PornGroove",
	"Satire",
	"SlowJam",
	"Club",
	"Tango",
	"Samba",
	"Folklore",
	"Ballad",
	"PowerBallad",
	"RhythmicSoul",
	"Freestyle",
	"Duet",
	"PunkRock",
	"DrumSolo",
	"Acapella",
	"Euro-House",
	"DanceHall",
	"Goa",
	"Drum&Bass",
	"Club-House",
	"Hardcore",
	"Terror",
	"Indie",
	"BritPop",
	"Negerpunk",
	"PolskPunk",
	"Beat",
	"ChristianGangstaRap",
	"Heavyl",
	"Blackl",
	"Crossover",
	"ContemporaryChristian",
	"ChristianRock",
	"Merengue",
	"Salsa",
	"Trashl",
	"Anime",
	"JPop",
	"Synthpop",
}
