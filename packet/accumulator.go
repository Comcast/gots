package packet

import (
	"bytes"

	"github.com/comcast/gots"
)

var (
	emptyByteArray []byte
)

type doneFunc func([]byte) (bool, error)

type accumulator struct {
	f       doneFunc
	packets []Packet
}

// NewAccumulator creates a new packet accumulator
// that is done when the provided doneFunc returns true.
// PacketAccumulator is not thread safe
func NewAccumulator(f doneFunc) Accumulator {
	return &accumulator{f: f}
}

// Add a packet to the accumulator. If the added packet completes
// the accumulation, based on the provided doneFunc, true is returned.
// Returns an error if the packet is not valid.
func (a *accumulator) Add(pkt Packet) (bool, error) {
	if badLen(pkt) {
		return false, mpegts.ErrInvalidPacketLength
	}
	if payloadUnitStartIndicator(pkt) {
		a.packets = make([]Packet, 0)
	} else if len(a.packets) == 0 {
		// First packet must have payload unit start indicator
		return false, mpegts.ErrNoPayloadUnitStartIndicator
	}
	a.packets = append(a.packets, pkt)
	b, err := a.Parse()
	if err != nil {
		return false, err
	}
	done, err := a.f(b)
	if err != nil {
		return false, err
	}
	return done, nil
}

// Parses the accumulated packets and returns the
// concatenated payloads or any error that occurred, not both
func (a *accumulator) Parse() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, pkt := range a.packets {
		pay, err := Payload(pkt)
		if err != nil {
			return emptyByteArray, err
		}
		buf.Write(pay)
	}
	return buf.Bytes(), nil
}

func (a *accumulator) Packets() []Packet {
	return a.packets
}
