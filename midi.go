package midi

import (
	"bufio"
	"bytes"
	"container/list"
	"log"
	"os"
)

////////////////////////////////////////////////////////////////////////////////
///// Track handling
////////////////////////////////////////////////////////////////////////////////

type Track struct {
	length uint32
	events *list.List
}

var trackhead = []byte{0x4d, 0x54, 0x72, 0x6b}

func NewTrack() Track {
	return Track{
		length: 8, // 8 bytes of header
		events: list.New(),
	}
}

func (t *Track) AddEvent(time uint32, data []byte) {
	tslice := MakeDelta(time)
	fullevent := bytes.Join([][]byte{tslice, data}, nil)
	t.events.PushBack(fullevent)
	t.length += uint32(len(fullevent))
}

func (t *Track) dump(w *bufio.Writer) {
	t.AddEvent(0, EndOfTrack())

	w.Write(trackhead)
	w.WriteByte(byte(t.length >> 24))
	w.WriteByte(byte(t.length >> 16))
	w.WriteByte(byte(t.length >> 8))
	w.WriteByte(byte(t.length))

	for e := t.events.Front(); e != nil; e = e.Next() {
		w.Write(e.Value.([]byte))
	}
}

////////////////////////////////////////////////////////////////////////////////
///// Correctly formatting delta time
////////////////////////////////////////////////////////////////////////////////

func MakeDelta(t uint32) []byte {
	if t>>7 == 0 {
		return []byte{
			byte(t),
		}
	} else if t>>14 == 0 {
		return []byte{
			byte(256 & t >> 7),
			byte(127 & t),
		}
	} else if t>>21 == 0 {
		return []byte{
			byte(256 & t >> 14),
			byte(256 & t >> 7),
			byte(127 & t),
		}
	} else {
		return []byte{
			byte(256 & t >> 21),
			byte(256 & t >> 14),
			byte(256 & t >> 7),
			byte(127 & t),
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
///// Making events
////////////////////////////////////////////////////////////////////////////////

func NoteOff(channel byte, nn byte, vv byte) []byte {
	return []byte{
		0x80 | channel,
		nn,
		vv,
	}
}

func NoteOn(channel uint8, nn uint8, vv uint8) []byte {
	return []byte{
		0x90 | channel,
		nn,
		vv,
	}
}

func KeyAfterTouch(channel uint8, nn uint8, vv uint8) []byte {
	return []byte{
		0xA0 | channel,
		nn,
		vv,
	}
}

func ControlChange(channel uint8, cc uint8, vv uint8) []byte {
	return []byte{
		0xB0 | channel,
		cc,
		vv,
	}
}

func ProgramChange(channel uint8, pp uint8) []byte {
	return []byte{
		0xC0 | channel,
		pp,
	}
}

func ChannelAfterTouch(channel uint8, cc uint8) []byte {
	return []byte{
		0xD0 | channel,
		cc,
	}
}

func PitchWheelChange(channel uint8, bb uint8, tt uint8) []byte {
	return []byte{
		0xE0 | channel,
		bb,
		tt,
	}
}

////////////////////////////////////////////////////////////////////////////////
///// Meta Events
////////////////////////////////////////////////////////////////////////////////

// WIP: There are actually 13 of these

func EndOfTrack() []byte {
	return []byte{
		0xFF,
		0x2F,
		0x00,
	}
}

////////////////////////////////////////////////////////////////////////////////
///// Write tracks out!
////////////////////////////////////////////////////////////////////////////////

type SyncType byte

const (
	OneTrack SyncType = 0
	Syncronous = 1
	Asyncronous = 2
)

var fileheader = []byte{0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06}

func WriteOut(path string, sync SyncType, tperq uint16, tracks []Track) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	// header format is 4D 54 68 64 00 00 00 06 ff ff nn nn dd dd
	// "ff ff" is the sync format (sync)
	// "nn nn" is the number of tracks
	// "dd dd" is the number of delta ticks per quarter note (tperq)
	w.Write(fileheader)
	w.WriteByte(0)
	w.WriteByte(byte(sync))
	ntracks := len(tracks)
	w.WriteByte(byte(ntracks >> 8))
	w.WriteByte(byte(ntracks))
	w.WriteByte(byte(tperq >> 8))
	w.WriteByte(byte(tperq))

	for i := 0; i < len(tracks); i++ {
		tracks[i].dump(w)
	}
}
