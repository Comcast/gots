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
	"fmt"
	"github.com/Comcast/gots"
)

// flags that are reserved and should not be used.
const (
	invalidTransportScramblingControlFlag TransportScramblingControlOptions = 1 // 01
	invalidAdaptationFieldControlFlag     AdaptationFieldControlOptions     = 0 // 00
)

// NewPacket creates a new packet with a Null ID, sync byte, and with the adaptation field control set to payload only.
// This function is error free.
func NewPacket() (pkt Packet) {
	pkt = make(Packet, PacketSize)
	//Default packet is the Null packet
	pkt[0] = 0x47 // sets the sync byte
	pkt[1] = 0x1F // equivalent to pkt.SetPID(NullPacketPid)
	pkt[2] = 0xFF // equivalent to pkt.SetPID(NullPacketPid)
	pkt[3] = 0x10 // equivalent to pkt.SetAdaptationFieldControl(PayloadFlag)
	return
}

// FromBytes creates a ts packet from a slice of bytes 188 in length.
// If the bytes provided have errors or the slice is not 188 in length,
// then an error vill be returned along with a nill slice.
func FromBytes(bytes []byte) (pkt Packet, err error) {
	pkt = Packet(bytes)
	err = pkt.CheckErrors()
	return
}

// SetTransportErrorIndicator sets the Transport Error Indicator flag.
func (p Packet) SetTransportErrorIndicator(value bool) error {
	return p.setBit(1, 0x80, value)
}

// TransportErrorIndicator returns the Transport Error Indicator
func (p Packet) TransportErrorIndicator() (bool, error) {
	return p.getBit(1, 0x80)
}

// SetPayloadUnitStartIndicator sets the Payload unit start indicator (PUSI) flag
// PUSI is a flag that indicates the start of PES data or PSI
// (Program-Specific Information) such as AT, CAT, PMT or NIT.  The PUSI
// flag is contained in the second bit of the second byte of the Packet.
func (p Packet) SetPayloadUnitStartIndicator(value bool) error {
	return p.setBit(1, 0x40, value)
}

// PayloadUnitStartIndicator returns the Payload unit start indicator (PUSI) flag
// PUSI is a flag that indicates the start of PES data or PSI
// (Program-Specific Information) such as AT, CAT, PMT or NIT.  The PUSI
// flag is contained in the second bit of the second byte of the Packet.
func (p Packet) PayloadUnitStartIndicator() (bool, error) {
	return p.getBit(1, 0x40)
}

// SetTransportPriority sets the Transport Priority flag
func (p Packet) SetTransportPriority(value bool) error {
	return p.setBit(1, 0x20, value)
}

// TransportPriority returns the Transport Priority flag
func (p Packet) TransportPriority() (bool, error) {
	return p.getBit(1, 0x20)
}

// SetPID sets the program ID
func (p Packet) SetPID(pid int) error {
	if err := p.valid(); err != nil {
		return err
	}
	p[1] = p[1]&^byte(0x1f) | byte(pid>>8)&byte(0x1f)
	p[2] = byte(pid)
	return nil
}

// PID Returns the program ID
func (p Packet) PID() (int, error) {
	if err := p.valid(); err != nil {
		return 0, err
	}
	return int(p[1]&0x1f)<<8 | int(p[2]), nil
}

// SetTransportScramblingControl sets the Transport Scrambling Control flag.
func (p Packet) SetTransportScramblingControl(value TransportScramblingControlOptions) error {
	if err := p.valid(); err != nil {
		return err
	}
	p[3] = p[3]&^byte(0xC0) | byte(value)<<6
	return nil
}

// TransportScramblingControl returns the Transport Scrambling Control flag.
func (p Packet) TransportScramblingControl() (TransportScramblingControlOptions, error) {
	if err := p.valid(); err != nil {
		return NoScrambleFlag, err
	}
	return TransportScramblingControlOptions((p[3] & 0xC0) >> 6), nil
}

// SetAdaptationFieldControl sets the Adaptation Field Control flag.
func (p Packet) SetAdaptationFieldControl(value AdaptationFieldControlOptions) error {
	if err := p.valid(); err != nil {
		return err
	}
	p[3] = p[3]&^byte(0x30) | byte(value)<<4
	if b, _ := p.HasAdaptationField(); b {
		initAdaptationField(p)
	}
	return nil
}

// AdaptationFieldControl returns the Adaptation Field Control.
func (p Packet) AdaptationFieldControl() (AdaptationFieldControlOptions, error) {
	if err := p.valid(); err != nil {
		return PayloadFlag, err
	}
	return AdaptationFieldControlOptions((p[3] & 0x30) >> 4), nil
}

// HasPayload returns true if the adaptation field control specifies that there is a payload.
func (p Packet) HasPayload() (bool, error) {
	return p.getBit(3, 0x10)
}

// HasPayload returns true if the adaptation field control specifies that there is an adaptation field.
func (p Packet) HasAdaptationField() (bool, error) {
	return p.getBit(3, 0x20)
}

// SetContinuityCounter sets the continuity counter.
// The continuity counter should be an integer between 0 and 15.
// If the number is out of this range then it will discard the extra bits.
// The effect is the same as modulus by 16.
func (p Packet) SetContinuityCounter(value int) error {
	if err := p.valid(); err != nil {
		return err
	}
	// if value is greater than 15 it will overflow and start at 0 again.
	p[3] = p[3]&^byte(0x0F) | byte(value&0x0F)
	return nil
}

// ContinuityCounter returns the continuity counter.
// The continuity counter is an integer between 0 and 15.
func (p Packet) ContinuityCounter() (int, error) {
	if err := p.valid(); err != nil {
		return 0, err
	}
	return int(p[3] & 0x0F), nil
}

// SetContinuityCounter sets the continuity counter to 0.
func (p Packet) ZeroContinuityCounter() error {
	return p.SetContinuityCounter(0)
}

// IncContinuityCounter increments the continuity counter.
// The continuity counter is an integer between 0 and 15.
// If the number is out of this range (overflow)
// after incrementing then it will discard the extra bits.
// The effect is the same as modulus by 16.
func (p Packet) IncContinuityCounter() error {
	cc, err := p.ContinuityCounter()
	cc += 1
	if err != nil {
		return err
	}
	return p.SetContinuityCounter(cc)
}

// IsNull returns true if the packet PID is equal to 8191, the null packet pid.
func (p Packet) IsNull() (bool, error) {
	pid, err := p.PID()
	return pid == NullPacketPid, err
}

// IsNull returns true if the packet PID is equal to 0, the PAT packet pid.
func (p Packet) IsPAT() (bool, error) {
	pid, err := p.PID()
	return pid == 0, err
}

// AdaptationField returns the AdaptationField of the packet.
// If the packet does not have an adaptation field then a nil
// AdaptationField is returned.
func (p Packet) AdaptationField() (AdaptationField, error) {
	hasAF, err := p.HasAdaptationField()
	if err != nil {
		return nil, err
	}
	if hasAF {
		return parseAdaptationField(p), nil
	}
	return nil, nil
}

// Payload returns a slice to a copy of the payload bytes in the packet.
// TODO: write tests
func (p Packet) Payload(packet Packet) ([]byte, error) {
	offset := 4 // packet header bytes
	if hasAF, err := p.HasAdaptationField(); err == nil && hasAF {
		offset += 1 + int(p[4])
	}
	payload := make([]byte, PacketSize-offset)
	copy(payload, p[offset:])
	return payload, nil
}

// Equal returns true if the bytes of the two packets are equal
// func (p Packet) Equals(r Packet) bool {
// 	return bytes.Equal(p, r)
// }

// CheckErrors checks the packet for errors
func (p Packet) sizeErrors() error {
	if len(p) != PacketSize {
		return gots.ErrInvalidPacketLength
	}
	if p.syncByte() != SyncByte {
		return gots.ErrBadSyncByte
	}
	return nil
}

// CheckErrors checks the packet for errors
func (p Packet) CheckErrors() error {
	if len(p) != PacketSize {
		return gots.ErrInvalidPacketLength
	}
	if p.syncByte() != SyncByte {
		return gots.ErrBadSyncByte
	}
	if flag, _ := p.TransportScramblingControl(); flag == invalidTransportScramblingControlFlag {
		return gots.ErrInvalidTSCFlag
	}
	if flag, _ := p.AdaptationFieldControl(); flag == invalidAdaptationFieldControlFlag {
		return gots.ErrInvalidAFCFlag
	}
	return nil
}

// syncByte returns the Sync byte.
func (p Packet) syncByte() byte {
	return p[0]
}

// invalid returns true if the length of the packet slice is
// anything but PacketSize (188)
func (p Packet) valid() error {
	if len(p) != PacketSize {
		return gots.ErrInvalidPacketLength
	}
	return nil
}

// setBit sets a bit in a packet. If packet is nil, or slice has a bad
// length, it does nothing.
func (p Packet) setBit(index int, mask byte, value bool) error {
	if err := p.valid(); err != nil {
		return err
	}
	if value {
		p[index] |= mask
	} else {
		p[index] &= ^mask
	}
	return nil
}

// getBit returns true if a bit in a packet is set to 1.
func (p Packet) getBit(index int, mask byte) (bool, error) {
	if err := p.valid(); err != nil {
		return false, err
	}
	return p[index]&mask != 0, nil
}

func (p Packet) String() string {
	if p.valid() != nil {
		return "Null"
	}
	return fmt.Sprintf("%X", []byte(p))
}
