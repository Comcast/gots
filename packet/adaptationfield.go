package packet

import (
	"fmt"
)

type AdaptationField []byte

// invalid returns true if the size of the slice (AdaptationField) is invalid,
// or the adaptation field doesnt exist
func (af AdaptationField) invalid() bool {
	// exact size of an adaptation field slice to work correctly
	// this is not related to the Field "Length" in the adaptation field
	// it is the same as the length of a packet. the adaptation must also
	// exist in order to be modified.

	// order matters, short circuit
	return len(af) != PacketSize || !Packet(af).HasAdaptationField()
}

// initAdaptationField initializes the adaptation field to have all false flags
func initAdaptationField(p Packet) {
	af := AdaptationField(p)
	af[4] = 1    // size of empty adaptation field is at least 1, zero size is used only for stuffing
	af[5] = 0x00 // no flags are set by default
}

// parseAdaptationField parses the adaptation field that is present in a packet.
// no need to report errors since this is handled during packet creation.
func parseAdaptationField(p Packet) AdaptationField {
	if p.HasAdaptationField() { // nil packet does not have an adaptation field
		return AdaptationField(p)
	}
	return nil
}

// returns the length of the PCR, if there is no PCR then its length is zero
func (af AdaptationField) pcrLength() int {
	if af.HasPCR() {
		return 6
	}
	return 0
}

// returns the length of the OPCR, if there is no PCR then its length is zero
func (af AdaptationField) opcrLength() int {
	if af.HasOPCR() {
		return 6
	}
	return 0
}

// returns the length of the splice countdown, if there is no splice countdown then its length is zero
func (af AdaptationField) spliceCountdownLength() int {
	if af.HasSplicingPoint() {
		return 1
	}
	return 0
}

// returns the length of the transport private data, if there is no transport private data then its length is zero
func (af AdaptationField) transportPrivateDataLength() int {
	if af.HasTransportPrivateData() {
		// cannot extend beyond adaptation field, number of bytes
		// for field stored in transportPrivateDataLength
		indexLength := 6 +
			af.pcrLength() +
			af.opcrLength() +
			af.spliceCountdownLength()
		return 1 + int(af[indexLength])
	}
	return 0
}

// returns the length of the adaptation field extension, if there is no adaptation field extension then its length is zero
func (af AdaptationField) adaptationExtensionLength() int {
	if af.HasAdaptationFieldExtension() {
		// TODO: make new type and call its methods to find the length
		return 1
	}
	return 0
}

// calculates the length of the Adaptation Field
// (with respect to the start of the packet) excluding stuffing
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

// resizeAF grows an adaptation field and erases the payload.
// Alternatley, with a negative delta, it can shrink a packet
// and keep the payload and stuff the empty space with stuffing bytes.
// this function is called automatically, no need for the library user
// to call it.
func (af AdaptationField) resizeAF(start int, delta int) {
	if delta > 0 { // shifting for growing
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
		// check if payload was corrupted/overwritten by growing
		if af.Length() < payloadStart-5 {
			// erase payload, it is corrupt anyways
			for payloadStart <= packetEnd {
				af[payloadStart] = 0xFF
				payloadStart++
			}
			// packet is stuffed until the very end.
			// this is an invalid packet since payload
			// must be at least one byte in size.
			// this will remind the user of the library
			// that the payload was destroyed
			af.setLength(183)
		}
	}
	if delta < 0 { // shifting for shrinking
		end := af.calculateMinLength() - 1 - delta
		startNew := start
		start := startNew - delta
		for start <= end {
			af[startNew] = af[start]
			startNew++
			start++
		}
		//stuff remaining bytes to preserve payload size
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
		delta, // default len of Transport Private Data
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
