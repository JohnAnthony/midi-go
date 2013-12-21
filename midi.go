package midi

import (
	"bufio"
	"os"
	"log"
)

//////////////////////////////////////////////////
///// Self-rolled binary trees for track data
//////////////////////////////////////////////////

type tracktree struct {
	left  *tracktree
	right *tracktree
	time  uint32
	event []byte
}

func (tt *tracktree) add(time uint32, data []byte) {
	if time <= tt.time && tt.left != nil {
		tt.left.add(time, data)
		return
	} else if time >= tt.time && tt.right != nil {
		tt.right.add(time, data)
		return
	}

	newt := tracktree{
		left:  nil,
		right: nil,
		time:  time,
		event: data,
	}

	if time <= tt.time {
		tt.left = &newt
	} else {
		tt.right = &newt
	}
}

func (tt *tracktree) dump(w *bufio.Writer) {
	if tt.left != nil {
		tt.left.dump(w)
	}

	w.Write(MakeDelta(tt.time))
	w.Write(tt.event)

	if tt.right != nil {
		tt.right.dump(w)
	}
}

//////////////////////////////////////////////////
///// Track handling
//////////////////////////////////////////////////

type Track struct {
	length uint32
	tree   *tracktree
}

var trackhead = []byte{0x4d, 0x54, 0x72, 0x6b}

func NewTrack() Track {
	return Track{
		length: 8,
		tree:   nil,
	}
}

func (t *Track) AddEvent(time uint32, data []byte) {
	if t.tree == nil {
		t.tree = &tracktree{
			left:  nil,
			right: nil,
			time:  time,
			event: data,
		}
	} else {
		t.tree.add(time, data)
	}

	t.length += uint32(len(data))
}

func (t *Track) dump(w *bufio.Writer) {
	w.Write(trackhead)
	w.WriteByte(byte(t.length >> 24))
	w.WriteByte(byte(t.length >> 16))
	w.WriteByte(byte(t.length >> 8))
	w.WriteByte(byte(t.length))

	t.tree.dump(w)
}

//////////////////////////////////////////////////
///// Correctly formatting delta time
//////////////////////////////////////////////////

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

//////////////////////////////////////////////////
///// Making events
//////////////////////////////////////////////////

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

//////////////////////////////////////////////////
///// Write tracks out!
//////////////////////////////////////////////////

var fileheader = []byte{0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06}

func WriteOut(path string, sync byte, tperq uint16, tracks []Track) {
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
	w.WriteByte(sync)
	ntracks := len(tracks)
	w.WriteByte(byte(ntracks >> 8))
	w.WriteByte(byte(ntracks))
	w.WriteByte(byte(tperq >> 8))
	w.WriteByte(byte(tperq))

	for i := 0; i < len(tracks); i++ {
		tracks[i].dump(w)
	}
}
