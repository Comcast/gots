package packet

import (
	"bytes"
	"github.com/Comcast/gots"
)

const (
	invalidTransportScramblingControlFlag TransportScramblingControlOptions = 1 // 01
	invalidAdaptationFieldControlFlag     AdaptationFieldControlOptions     = 0 // 00
)

func NewPacket() (pkt Packet) {
	pkt = make(Packet, PacketSize)
	pkt[0] = 0x47
	return
}

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
	return p.getBit(1, 0x40)
}

func (p Packet) SetTransportPriority(value bool) {
	p.setBit(1, 0x20, value)
}

func (p Packet) TransportPriority() bool {
	return p.getBit(1, 0x20)
}

func (p Packet) SetPID(pid int) {
	p[1] = p[1]&^byte(0x1f) | byte(pid>>8)&byte(0x1f)
	p[2] = byte(pid)
}

func (p Packet) PID() int {
	return int(p[1]&0x1f)<<8 | int(p[2])
}

func (p Packet) SetTransportScramblingControl(value TransportScramblingControlOptions) {
	p[3] = p[3]&^byte(0xC0) | byte(value)<<6
}

func (p Packet) TransportScramblingControl() TransportScramblingControlOptions {
	return TransportScramblingControlOptions((p[3] & 0xC0) >> 6)
}

func (p Packet) SetAdaptationFieldControl(value AdaptationFieldControlOptions) {
	p[3] = p[3]&^byte(0x30) | byte(value)<<4
	// TODO: adaptation field class
}

func (p Packet) AdaptationFieldControl() AdaptationFieldControlOptions {
	return AdaptationFieldControlOptions((p[3] & 0x30) >> 4)
}

func (p Packet) HasPayload() bool {
	return p.getBit(3, 0x10)
}

func (p Packet) HasAdaptationField() bool {
	return p.getBit(3, 0x20)
}

// overflows after 15 and starts againat 0
func (p Packet) SetContinuityCounter(value int) {
	p[3] = p[3]&^byte(0x0F) | byte(value&0x0F)
}

func (p Packet) ContinuityCounter() int {
	return int(p[3] & 0x0F)
}

func (p Packet) ZeroContinuityCounter() {
	p.SetContinuityCounter(0)
}

func (p Packet) IncContinuityCounter() {
	p.SetContinuityCounter(p.ContinuityCounter() + 1)
}

func (p Packet) IsNull() bool {
	return p.PID() == NullPacketPid
}

func (p Packet) IsPAT() bool {
	return p.PID() == 0
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

// setBit sets a bit in a packet. If packet is nil, or slice has a bad
// length, it does nothing.
func (p Packet) setBit(index int, mask byte, value bool) {
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
