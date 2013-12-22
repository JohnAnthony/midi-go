package main

import (
	"github.com/JohnAnthony/midi-go"
)

func main() {
	cmajscaleup := []uint8 {60, 62, 64, 65, 67, 69, 71, 72}
	cmajscaledown := []uint8 {72, 71, 69, 67, 65, 64, 62, 60}
	aminscaleup := []uint8 {57, 59, 60, 62, 64, 65, 67, 69}
	aminscaledown := []uint8 {69, 67, 65, 64, 62, 60, 59, 57}
	fmajscaleup := []uint8 {53, 55, 57, 59, 60, 62, 64, 65}
	fmajscaledown := []uint8 {65, 64, 62, 60, 59, 57, 55, 53}
	
	tracks := make([]midi.Track, 1)
	tracks[0] = midi.NewTrack()

	for _, v := range cmajscaleup {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0xff))
		tracks[0].AddEvent(1, midi.NoteOff(0, v, 0xff))
	}
	
	for _, v := range cmajscaledown {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0xff))
		tracks[0].AddEvent(1, midi.NoteOff(0, v, 0xff))
	}

	for _, v := range aminscaleup {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0xff))
		tracks[0].AddEvent(1, midi.NoteOff(0, v, 0xff))
	}
	
	for _, v := range aminscaledown {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0xff))
		tracks[0].AddEvent(1, midi.NoteOff(0, v, 0xff))
	}

	for _, v := range fmajscaleup {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0xff))
		tracks[0].AddEvent(1, midi.NoteOff(0, v, 0xff))
	}
	
	for _, v := range fmajscaledown {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0xff))
		tracks[0].AddEvent(1, midi.NoteOff(0, v, 0xff))
	}

	for _, v := range aminscaleup {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0xff))
		tracks[0].AddEvent(1, midi.NoteOff(0, v, 0xff))
	}
	
	for _, v := range aminscaledown {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0xff))
		tracks[0].AddEvent(1, midi.NoteOff(0, v, 0xff))
	}

	midi.WriteOut("test.midi", midi.OneTrack, 4, tracks)
}
