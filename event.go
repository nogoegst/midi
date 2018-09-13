// event.go - MIDI event definition
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of midi, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package midi

import "time"

// Event represents MIDI event.
type Event struct {
	Status    byte
	Data1     int8
	Data2     int8
	Timestamp time.Time
}
