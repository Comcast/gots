/*
MIT License

Copyright 2016 Comcast Cable Communications Management, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package packet

import (
	"github.com/Comcast/gots/"
)

// flags that are reserved and should not be used.
const (
	invalidTransportScramblingControlFlag TransportScramblingControlOptions = 1 // 01
	invalidAdaptationFieldControlFlag     AdaptationFieldControlOptions     = 0 // 00
)

// New creates a new packet with a Null ID, sync byte, and with the adaptation
// field control set to payload only.
func New() *Packet {
	//Default packet is the Null packet
	return &Packet{
		0: 0x47,          // Sync byte
		1: 0x1f, 2: 0xff, // pkt.SetPID(NullPacketPid)
		3: 0x10, // pkt.SetAdaptationFieldControl(PayloadFlag)
	}
}

// FromBytes creates a ts packet from a slice of bytes 188 in length.
// If the bytes provided have errors or the slice is not 188 in length,
// then an error vill be returned along with a nill slice.
func FromBytes(bytes []byte) (*Packet, error) {
	if len(bytes) != PacketSize {
		return nil, gots.ErrInvalidPacketLength
	}
	var pkt Packet
	copy(pkt[:], bytes)
	err := pkt.CheckErrors()
	return &pkt, err
}

// SetTransportErrorIndicator sets the Transport Error Indicator flag.
func (p *Packet) SetTransportErrorIndicator(value bool) {
	p.setBit(1, 0x80, value)
}

// TransportErrorIndicator returns the Transport Error Indicator
func (p *Packet) TransportErrorIndicator() bool {
	return p.getBit(1, 0x80)
}

// SetPayloadUnitStartIndicator sets the Payload unit start indicator (PUSI) flag
// PUSI is a flag that indicates the start of PES data or PSI
// (Program-Specific Information) such as AT, CAT, PMT or NIT.  The PUSI
// flag is contained in the second bit of the second byte of the Packet.
func (p *Packet) SetPayloadUnitStartIndicator(value bool) {
	p.setBit(1, 0x40, value)
}

// PayloadUnitStartIndicator returns the Payload unit start indicator (PUSI) flag
// PUSI is a flag that indicates the start of PES data or PSI
// (Program-Specific Information) such as AT, CAT, PMT or NIT.  The PUSI
// flag is contained in the second bit of the second byte of the Packet.
func (p *Packet) PayloadUnitStartIndicator() bool {
	return p.getBit(1, 0x40)
}

// SetTransportPriority sets the Transport Priority flag
func (p *Packet) SetTransportPriority(value bool) {
	p.setBit(1, 0x20, value)
}

// TransportPriority returns the Transport Priority flag
func (p *Packet) TransportPriority() bool {
	return p.getBit(1, 0x20)
}

// SetPID sets the program ID
func (p *Packet) SetPID(pid int) {
	p[1] = p[1]&^byte(0x1f) | byte(pid>>8)&byte(0x1f)
	p[2] = byte(pid)
}

// PID Returns the program ID
func (p *Packet) PID() int {
	return int(p[1]&0x1f)<<8 | int(p[2])
}

// SetTransportScramblingControl sets the Transport Scrambling Control flag.
func (p *Packet) SetTransportScramblingControl(value TransportScramblingControlOptions) {
	p[3] = p[3]&^byte(0xC0) | byte(value)<<6
}

// TransportScramblingControl returns the Transport Scrambling Control flag.
func (p *Packet) TransportScramblingControl() TransportScramblingControlOptions {
	return TransportScramblingControlOptions((p[3] & 0xC0) >> 6)
}

// SetAdaptationFieldControl sets the Adaptation Field Control flag.
func (p *Packet) SetAdaptationFieldControl(value AdaptationFieldControlOptions) error {
	hasAFBefore := p.HasAdaptationField()
	p[3] = p[3]&^byte(0x30) | byte(value)<<4
	hasAFAfter := p.HasAdaptationField()
	// if it didnt have an AF but now has one. Init the AF
	if !hasAFBefore && hasAFAfter {
		initAdaptationField(p)
	}

	if value == PayloadAndAdaptationFieldFlag {
		af, _ := p.AdaptationField()
		if af.Length() == 183 {
			if af.stuffingStart() < PacketSize {
				af.setLength(182)
				af.stuffAF()
			} else {
				return gots.ErrAdaptationFieldTooLarge
			}
		}
	}

	return nil
}

// AdaptationFieldControl returns the Adaptation Field Control.
func (p *Packet) AdaptationFieldControl() AdaptationFieldControlOptions {
	return AdaptationFieldControlOptions((p[3] & 0x30) >> 4)
}

// HasPayload returns true if the adaptation field control specifies that there is a payload.
func (p *Packet) HasPayload() bool {
	return p.getBit(3, 0x10)
}

// HasPayload returns true if the adaptation field control specifies that there is an adaptation field.
func (p *Packet) HasAdaptationField() bool {
	return p.getBit(3, 0x20)
}

// SetContinuityCounter sets the continuity counter.
// The continuity counter should be an integer between 0 and 15.
// If the number is out of this range then it will discard the extra bits.
// The effect is the same as modulus by 16.
func (p *Packet) SetContinuityCounter(value int) {
	// if value is greater than 15 it will overflow and start at 0 again.
	p[3] = p[3]&^byte(0x0F) | byte(value&0x0F)
}

// ContinuityCounter returns the continuity counter.
// The continuity counter is an integer between 0 and 15.
func (p *Packet) ContinuityCounter() int {
	return int(p[3] & 0x0F)
}

// SetContinuityCounter sets the continuity counter to 0.
func (p *Packet) ZeroContinuityCounter() {
	p.SetContinuityCounter(0)
}

// IncContinuityCounter increments the continuity counter.
// The continuity counter is an integer between 0 and 15.
// If the number is out of this range (overflow)
// after incrementing then it will discard the extra bits.
// The effect is the same as modulus by 16.
func (p *Packet) IncContinuityCounter() {
	cc := p.ContinuityCounter()
	cc += 1
	p.SetContinuityCounter(cc)
}

// IsNull returns true if the packet PID is equal to 8191, the null packet pid.
func (p *Packet) IsNull() bool {
	return p.PID() == NullPacketPid
}

// IsNull returns true if the packet PID is equal to 0, the PAT packet pid.
func (p *Packet) IsPAT() bool {
	return p.PID() == 0
}

// AdaptationField returns the AdaptationField of the packet.
// If the packet does not have an adaptation field then a nil
// AdaptationField is returned.
func (p *Packet) AdaptationField() (*AdaptationField, error) {
	hasAF := p.HasAdaptationField()
	if hasAF {
		return (*AdaptationField)(p), nil
	}
	return nil, gots.ErrNoAdaptationField
}

// SetAdaptationField copies the AdaptationField into a packet.
// If the packet does not have an adaptation field then an error is returned
// AdaptationField must fit in the same size as the existing adaptation field
// and its stuffing bytes.
func (p *Packet) SetAdaptationField(af *AdaptationField) error {
	if !p.HasAdaptationField() {
		return gots.ErrNoAdaptationField
	}
	oldAF, _ := p.AdaptationField()
	if oldAF.stuffingEnd() < af.stuffingStart() {
		return gots.ErrAdaptationFieldTooLarge
	}
	copy(oldAF[5:oldAF.stuffingEnd()], af[5:af.stuffingStart()]) // Copy without length.
	oldAF.stuffAF()
	return nil
}

func (p *Packet) payloadStart() int {
	if p.HasAdaptationField() {
		return 4 + 1 + (*AdaptationField)(p).Length()
	}
	return 4 // packet header bytes
}

// stuffingStart returns where the stuffing begins, this is also the first byte where the payload can begin.
// if there is no payload then it is stuffed until the very end
func (p *Packet) stuffingStart() int {
	af, err := p.AdaptationField()
	if err != nil {
		return 4
	}
	if af.Length() == 0 {
		return 5
	}
	return af.stuffingStart()
}

func (p *Packet) freeSpace() int {
	return PacketSize - p.stuffingStart()
}

// Payload returns a slice to a copy of the payload bytes in the packet.
func (p *Packet) Payload() ([]byte, error) {
	afc := p.AdaptationFieldControl()
	if afc == AdaptationFieldFlag {
		return nil, gots.ErrNoPayload
	}
	offset := p.payloadStart()
	payload := make([]byte, PacketSize-offset)
	copy(payload, p[offset:])
	return payload, nil
}

// SetPayload sets the payload of the packet. If the payload cannot fit in the
// packet an integer will be returned that is the number of bytes that were
// able to fit in the packet.
func (p *Packet) SetPayload(data []byte) (int, error) {
	afc := p.AdaptationFieldControl()
	if afc == AdaptationFieldFlag {
		return 0, gots.ErrNoPayload
	}
	freeSpace := p.freeSpace()
	if freeSpace > len(data) {
		p.SetAdaptationFieldControl(PayloadAndAdaptationFieldFlag)
		af, _ := p.AdaptationField()

		af.setLength(PacketSize - (len(data) + 4 + 1)) // header length + adaptation field length
		af.stuffAF()
	} else {
		af, _ := p.AdaptationField()
		if af != nil {
			af.setLength(PacketSize - (freeSpace + 4 + 1)) // header length + adaptation field length
		}
	}

	offset := p.payloadStart()
	return copy(p[offset:], data), nil
}

// Equal returns true if the bytes of the two packets are equal
func (p *Packet) Equals(r *Packet) bool {
	return Equal(p, r)
}

// CheckErrors checks the packet for errors
func (p *Packet) CheckErrors() error {
	if p.syncByte() != SyncByte {
		return gots.ErrBadSyncByte
	}
	if p.TransportScramblingControl() == invalidTransportScramblingControlFlag {
		return gots.ErrInvalidTSCFlag
	}
	if p.AdaptationFieldControl() == invalidAdaptationFieldControlFlag {
		return gots.ErrInvalidAFCFlag
	}
	return nil
}

// syncByte returns the Sync byte.
func (p *Packet) syncByte() byte {
	return p[0]
}

// setBit sets a bit in a packet. If packet is nil, or slice has a bad
// length, it does nothing.
func (p *Packet) setBit(index int, mask byte, value bool) {
	if value {
		p[index] |= mask
	} else {
		p[index] &= ^mask
	}
}

// getBit returns true if a bit in a packet is set to 1.
func (p *Packet) getBit(index int, mask byte) bool {
	return p[index]&mask != 0
}
