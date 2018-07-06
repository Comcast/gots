package packet

import (
	"fmt"
)

type AdaptationField []byte

func (af AdaptationField) invalid() bool {
	// order matters, short circuit
	return len(af) != PacketSize || !Packet(af).HasAdaptationField()
}

// all false flags
func initAdaptationField(p Packet) {
	af := AdaptationField(p)
	af[4] = 1
	af[5] = 0x00 // no flags are set by default
}

func parseAdaptationField(p Packet) AdaptationField {
	if p.HasAdaptationField() { // nil packet does not have an adaptation field
		return AdaptationField(p)
	}
	return nil
}

func (af AdaptationField) pcrLength() int {
	if af.HasPCR() {
		return 6
	}
	return 0
}

func (af AdaptationField) opcrLength() int {
	if af.HasOPCR() {
		return 6
	}
	return 0
}

func (af AdaptationField) spliceCountdownLength() int {
	if af.HasSplicingPoint() {
		return 1
	}
	return 0
}

func (af AdaptationField) transportPrivateDataLength() int {
	if af.HasTransportPrivateData() {
		indexLength := 6 +
			af.pcrLength() +
			af.opcrLength() +
			af.spliceCountdownLength()
		return 1 + int(af[indexLength])
	}
	return 0
}

func (af AdaptationField) adaptationExtensionLength() int {
	if af.HasAdaptationFieldExtension() {
		// TODO: make new type and call its methods to find the length
		return 1
	}
	return 0
}

func (af AdaptationField) calculateMinLength() int {
	if af.invalid() {
		return 0
	}
	return 6 +
		af.pcrLength() + af.opcrLength() + af.spliceCountdownLength() +
		af.transportPrivateDataLength() + af.adaptationExtensionLength()
}

func (af AdaptationField) setBit(index int, mask byte, value bool) {
	if af.invalid() {
		return
	}
	if value {
		af[index] = af[index] | mask
	} else {
		af[index] = af[index] & ^mask
	}
}

func (af AdaptationField) getBit(index int, mask byte) bool {
	if af.invalid() {
		return false
	}
	return af[index]&mask != 0
}

func (af AdaptationField) setBitReturnDelta(index int, mask byte, value bool) int {
	if af.invalid() {
		return 0
	}
	oldValue := af.getBit(index, mask)
	af.setBit(index, mask, value)
	if value != oldValue {
		if value {
			return 1 // growing
		}
		return -1 // shrinking
	}
	return 0 // same
}

func (af AdaptationField) setLength(length int) {
	if af.invalid() {
		return
	}
	af[4] = byte(length)
}

func (af AdaptationField) resizeAF(start int, delta int) {
	if delta > 0 { // growing
		// move existing bytes to new location
		payloadStart := af.calculateMinLength()
		endNew := payloadStart - 1
		end := endNew - delta
		for start <= end {
			af[endNew] = af[end]
			endNew--
			end--
		}
		packetEnd := len(af) - 1
		if af.Length() < payloadStart-5 {
			// erase payload
			for payloadStart <= packetEnd {
				af[payloadStart] = 0xFF
				payloadStart++
			}
			af.setLength(183)
		}
	}
	if delta < 0 { // shrinking
		end := af.calculateMinLength() - delta
		startNew := start
		start := startNew - delta
		for start <= end {
			af[startNew] = af[start]
			startNew++
			start++
		}
		//preserve payload size and stuff it
		for startNew <= end {
			af[startNew] = 0xFF
			startNew++
		}
	}
}

// Length returns the length of the adaptation field
func (af AdaptationField) Length() int {
	if af.invalid() {
		return 0
	}
	return int(af[4])
}

// SetDiscontinuity sets the Discontinuity field of the packet
func (af AdaptationField) SetDiscontinuity(value bool) {
	af.setBit(5, 0x80, value)
}

// Discontinuity returns the value of the discontinuity field in the packet
func (af AdaptationField) Discontinuity() bool {
	return af.getBit(5, 0x80)
}

func (af AdaptationField) SetRandomAccess(value bool) {
	af.setBit(5, 0x40, value)
}

func (af AdaptationField) RandomAccess() bool {
	return af.getBit(5, 0x40)
}

func (af AdaptationField) SetESPriority(value bool) {
	af.setBit(5, 0x20, value)
}

func (af AdaptationField) ESPriority() bool {
	return af.getBit(5, 0x20)
}

func (af AdaptationField) SetHasPCR(value bool) {
	delta := 6 * af.setBitReturnDelta(5, 0x10, value)
	af.resizeAF(6, delta)
}

func (af AdaptationField) HasPCR() bool {
	return af.getBit(5, 0x10)
}

func (af AdaptationField) SetHasOPCR(value bool) {
	delta := 6 * af.setBitReturnDelta(5, 0x08, value)
	af.resizeAF(6+
		af.pcrLength(),
		delta,
	)
}

func (af AdaptationField) HasOPCR() bool {
	return af.getBit(5, 0x08)
}

func (af AdaptationField) SetHasSplicingPoint(value bool) {
	delta := 1 * af.setBitReturnDelta(5, 0x04, value)
	af.resizeAF(6+
		af.pcrLength()+
		af.opcrLength(),
		delta,
	)
}

func (af AdaptationField) HasSplicingPoint() bool {
	return af.getBit(5, 0x04)
}

func (af AdaptationField) SetHasTransportPrivateData(value bool) {
	delta := 1 * af.setBitReturnDelta(5, 0x02, value)
	// TODO: craft a TP
	af.resizeAF(6+
		af.pcrLength()+
		af.opcrLength()+
		af.spliceCountdownLength(),
		delta,
	)
}

func (af AdaptationField) HasTransportPrivateData() bool {
	return af.getBit(5, 0x02)
}

func (af AdaptationField) SetHasAdaptationFieldExtension(value bool) {
	delta := af.setBitReturnDelta(5, 0x01, value) * 1
	af.resizeAF(6+
		af.pcrLength()+
		af.opcrLength()+
		af.spliceCountdownLength()+
		af.transportPrivateDataLength(),
		delta,
	)
}

func (af AdaptationField) HasAdaptationFieldExtension() bool {
	return af.getBit(5, 0x01)
}

func (af AdaptationField) String() string {
	if af.invalid() {
		return "Null"
	}
	return fmt.Sprintf("%X", []byte(af[4:]))
}
