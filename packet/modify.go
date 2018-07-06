package packet

import (
	"bytes"
	"fmt"
	"github.com/Comcast/gots"
)

// flags that are reserved and should not be used.
const (
	invalidTransportScramblingControlFlag TransportScramblingControlOptions = 1 // 01
	invalidAdaptationFieldControlFlag     AdaptationFieldControlOptions     = 0 // 00
)

// NewPacket creates a new packet with a Null ID, sync byte, and with the adaptation field control set to payload only.
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
func (p Packet) SetTransportErrorIndicator(value bool) {
	p.setBit(1, 0x80, value)
}

// TransportErrorIndicator returns the Transport Error Indicator
func (p Packet) TransportErrorIndicator() bool {
	if p.invalid() {
		return false // defalt value for nil packet
	}
	return p.getBit(1, 0x80)
}

// SetPayloadUnitStartIndicator sets the Payload unit start indicator (PUSI) flag
// PUSI is a flag that indicates the start of PES data or PSI
// (Program-Specific Information) such as AT, CAT, PMT or NIT.  The PUSI
// flag is contained in the second bit of the second byte of the Packet.
func (p Packet) SetPayloadUnitStartIndicator(value bool) {
	p.setBit(1, 0x40, value)
}

// PayloadUnitStartIndicator returns the Payload unit start indicator (PUSI) flag
// PUSI is a flag that indicates the start of PES data or PSI
// (Program-Specific Information) such as AT, CAT, PMT or NIT.  The PUSI
// flag is contained in the second bit of the second byte of the Packet.
func (p Packet) PayloadUnitStartIndicator() bool {
	if p.invalid() {
		return false // defalt value for nil packet
	}
	return p.getBit(1, 0x40)
}

// SetTransportPriority sets the Transport Priority flag
func (p Packet) SetTransportPriority(value bool) {
	p.setBit(1, 0x20, value)
}

// TP returns the Transport Priority flag
func (p Packet) TransportPriority() bool {
	if p.invalid() {
		return false // defalt value for nil packet
	}
	return p.getBit(1, 0x20)
}

// SetPID sets the program ID
func (p Packet) SetPID(pid int) {
	if p.invalid() {
		return
	}
	p[1] = p[1]&^byte(0x1f) | byte(pid>>8)&byte(0x1f)
	p[2] = byte(pid)
}

// PID Returns the program ID
func (p Packet) PID() int {
	if p.invalid() {
		return NullPacketPid // defalt value for nil packet
	}
	return int(p[1]&0x1f)<<8 | int(p[2])
}

// SetTransportScramblingControl sets the Transport Scrambling Control flag.
func (p Packet) SetTransportScramblingControl(value TransportScramblingControlOptions) {
	if p.invalid() {
		return
	}
	p[3] = p[3]&^byte(0xC0) | byte(value)<<6
}

// TransportScramblingControl returns the Transport Scrambling Control flag.
func (p Packet) TransportScramblingControl() TransportScramblingControlOptions {
	if p.invalid() {
		return NoScrambleFlag // defalt value for nil packet
	}
	return TransportScramblingControlOptions((p[3] & 0xC0) >> 6)
}

// SetAdaptationFieldControl sets the Adaptation Field Control flag.
func (p Packet) SetAdaptationFieldControl(value AdaptationFieldControlOptions) {
	if p.invalid() {
		return
	}
	p[3] = p[3]&^byte(0x30) | byte(value)<<4
	if p.HasAdaptationField() {
		initAdaptationField(p)
	}
}

// AdaptationFieldControl returns the Adaptation Field Control.
func (p Packet) AdaptationFieldControl() AdaptationFieldControlOptions {
	if p.invalid() {
		return PayloadFlag // defalt value for nil packet
	}
	return AdaptationFieldControlOptions((p[3] & 0x30) >> 4)
}

// HasPayload returns true if the adaptation field control specifies that there is a payload.
func (p Packet) HasPayload() bool {
	if p.invalid() {
		return true // defalt value for nil packet
	}
	return p.getBit(3, 0x10)
}

// HasPayload returns true if the adaptation field control specifies that there is an adaptation field.
func (p Packet) HasAdaptationField() bool {
	if p.invalid() {
		return false // defalt value for nil packet
	}
	return p.getBit(3, 0x20)
}

// SetContinuityCounter sets the continuity counter.
// The continuity counter should be an integer between 0 and 15.
// If the number is out of this range then it will discard the extra bits.
// The effect is the same as modulus by 16.
func (p Packet) SetContinuityCounter(value int) {
	if p.invalid() {
		return
	}
	// if value is greater than 15 it will overflow and start at 0 again.
	p[3] = p[3]&^byte(0x0F) | byte(value&0x0F)
}

// ContinuityCounter returns the continuity counter.
// The continuity counter is an integer between 0 and 15.
func (p Packet) ContinuityCounter() int {
	if p.invalid() {
		return 0 // defalt value for nil packet
	}
	return int(p[3] & 0x0F)
}

// SetContinuityCounter sets the continuity counter to 0.
func (p Packet) ZeroContinuityCounter() {
	p.SetContinuityCounter(0)
}

// IncContinuityCounter increments the continuity counter.
// The continuity counter is an integer between 0 and 15.
// If the number is out of this range (overflow)
// after incrementing then it will discard the extra bits.
// The effect is the same as modulus by 16.
func (p Packet) IncContinuityCounter() {
	p.SetContinuityCounter(p.ContinuityCounter() + 1)
}

// IsNull returns true if the packet PID is equal to 8191, the null packet pid.
func (p Packet) IsNull() bool {
	return p.PID() == NullPacketPid
}

// IsNull returns true if the packet PID is equal to 0, the PAT packet pid.
func (p Packet) IsPAT() bool {
	return p.PID() == 0
}

// AdaptationField returns the AdaptationField of the packet.
// If the packet does not have an adaptation field then a nil
// AdaptationField is returned.
func (p Packet) AdaptationField() AdaptationField {
	if p.invalid() {
		return nil // defalt value for nil packet
	}
	return parseAdaptationField(p)
}

// Payload returns a slice to a copy of the payload bytes in the packet.
// TODO: write tests
func (p Packet) Payload(packet Packet) []byte {
	offset := 4 // packet header bytes
	if a := p.AdaptationField(); a != nil {
		offset += 1 + a.Length()
	}
	payload := make([]byte, PacketSize-offset)
	copy(payload, p[offset:])
	return payload
}

// Equal returns true if the bytes of the two packets are equal
func (p Packet) Equals(r Packet) bool {
	return bytes.Equal(p, r)
}

// CheckErrors checks the packet for errors
func (p Packet) CheckErrors() error {
	if len(p) != PacketSize {
		return gots.ErrInvalidPacketLength
	}
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
func (p Packet) syncByte() byte {
	return p[0]
}

// invalid returns true if the length of the packet slice is
// anything but PacketSize (188)
func (p Packet) invalid() bool {
	return len(p) != PacketSize
}

// setBit sets a bit in a packet. If packet is nil, or slice has a bad
// length, it does nothing.
func (p Packet) setBit(index int, mask byte, value bool) {
	if p.invalid() {
		return
	}
	if value {
		p[index] |= mask
	} else {
		p[index] &= ^mask
	}
}

// getBit returns true if a bit in a packet is set to 1.
func (p Packet) getBit(index int, mask byte) bool {
	return p[index]&mask != 0
}

func (p Packet) String() string {
	if p.invalid() {
		return "Null"
	}
	return fmt.Sprintf("%X", []byte(p))
}
