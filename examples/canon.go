package main

import (
	"github.com/JohnAnthony/midi-go"
)

func main() {
	// For those of you in the know, this is Pachabel's Canon
	notes := []uint8 {
		62, 64, 66, 67, 69, 71, 73, 74, // D Maj ascending
		57, 59, 61, 62, 64, 66, 67, 69, // A Maj ascending
		59, 61, 62, 64, 66, 67, 69, 71, // B Min ascending
		54, 56, 57, 59, 61, 62, 64, 66, // F# Min ascending
		55, 57, 59, 60, 62, 64, 66, 67, // G Maj ascending
		62, 64, 66, 67, 69, 71, 73, 74, // D Maj ascending
		55, 57, 59, 61, 62, 64, 66, 67, // G Maj ascending
		57, 59, 61, 62, 64, 66, 67, 69, // A Maj ascending
	}

	tracks := make([]midi.Track, 2)
	tracks[0] = midi.NewTrack()
	tracks[1] = midi.NewTrack()

	for i, v := range notes {
		tracks[0].AddEvent(0, midi.NoteOn(0, v, 0x66))
		tracks[0].AddEvent(2, midi.NoteOff(0, v, 0x66))

		if i % 8 == 0 { // Bass note
			tracks[1].AddEvent(0, midi.NoteOn(1, v - 24, 0xff))
			tracks[1].AddEvent(16, midi.NoteOff(1, v - 24, 0xff))
		}
	}

	midi.WriteOut("test.midi", midi.Syncronous, 4, tracks)
}
