// midi.go - midi interface for OpenBSD midi(4)
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of midi, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package midi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Device represents a MIDI device.
type Device struct {
	f *os.File
}

// OpenDevice opens a MIDI device by its name (e.g. "rmidi0").
func OpenDevice(name string) (*Device, error) {
	p := filepath.Join("/dev", name)
	f, err := os.OpenFile(p, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("midi: %v", err)
	}
	device := &Device{
		f: f,
	}
	return device, nil
}

// Event represents MIDI event.
type Event struct {
	Status    byte
	Data1     int8
	Data2     int8
	Timestamp time.Time
}

func parseEvents(p []byte) (events []Event, err error) {
	if len(p)%3 != 0 {
		panic("buffer length is not divisible by 3")
	}
	timestamp := time.Now()
	for i := 0; i < len(p); i += 3 {
		e := Event{
			Status:    p[i],
			Data1:     int8(p[i+1]),
			Data2:     int8(p[i+2]),
			Timestamp: timestamp,
		}
		events = append(events, e)
	}
	return events, nil
}

// Read reads out bunch of events from device d.
func (d *Device) Read() ([]Event, error) {
	buf := make([]byte, 1024)
	n, err := d.f.Read(buf)
	if err != nil {
		return []Event{}, err
	}
	return parseEvents(buf[:n])
}

// WriteShort writes an event to device d.
func (d *Device) WriteShort(ev Event) error {
	p := make([]byte, 3)
	p[0] = ev.Status
	p[1] = byte(ev.Data1)
	p[2] = byte(ev.Data2)

	n, err := d.f.Write(p)
	if err != nil {
		return err
	}
	if n != 3 {
		return errors.New("midi: partial write")
	}
	return nil
}

// Close closes the device.
func (d *Device) Close() error {
	d.f.SetDeadline(time.Now())
	return d.f.Close()
}
