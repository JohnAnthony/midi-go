package main

import (
	"github.com/JohnAnthony/midi-go"
)

func main() {
	cmajscaleup := []uint8 {60, 62, 64, 65, 67, 69, 71, 72}
	cmajscaledown := []uint8 {72, 71, 69, 67, 65, 64, 62, 60}
	
	tracks := make([]midi.Track, 1)
	tracks[0] = midi.NewTrack()

	for _, v := range cmajscaleup {
		tracks[0].AddEvent(1, midi.NoteOn(0, 48, 0xff))
		tracks[0].AddEvent(1, midi.NoteOn(0, v, 0xff))
	}
	
	for _, v := range cmajscaledown {
		tracks[0].AddEvent(1, midi.NoteOn(0, 48, 0xff))
		tracks[0].AddEvent(1, midi.NoteOn(0, v, 0xff))
	}

	midi.WriteOut("test.midi", midi.Asyncronous, 4, tracks)
}
