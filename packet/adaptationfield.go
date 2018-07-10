package packet

import (
	"fmt"
	"github.com/Comcast/gots"
)

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
	af[4] = 183  // adaptation field will take up the rest of the packet by Default
	af[5] = 0x00 // no flags are set by default
	for i := 6; i < len(af); i++ {
		af[i] = 0xFF
	}
}

// parseAdaptationField parses the adaptation field that is present in a packet.
// no need to report errors since this is handled during packet creation.
func parseAdaptationField(p Packet) AdaptationField {
	if p.HasAdaptationField() { // nil packet does not have an adaptation field
		return AdaptationField(p)
	}
	return nil
}

// pcrLength returns the length of the PCR, if there is no PCR then its length is zero
func (af AdaptationField) pcrLength() int {
	if af.HasPCR() {
		return 6
	}
	return 0
}

const pcrStart = 6

// opcrLength returns the length of the OPCR, if there is no OPCR then its length is zero
func (af AdaptationField) opcrLength() int {
	if af.HasOPCR() {
		return 6
	}
	return 0
}

// opcrStart returns the start index of where the OPCR field should be
func (af AdaptationField) opcrStart() int {
	return pcrStart + af.pcrLength()
}

// spliceCountdownLength returns the length of the splice countdown, if there is no splice countdown then its length is zero
func (af AdaptationField) spliceCountdownLength() int {
	if af.HasSplicingPoint() {
		return 1
	}
	return 0
}

// spliceCountdownStart returns the start index of where the splice countdown field should be
func (af AdaptationField) spliceCountdownStart() int {
	return pcrStart + af.pcrLength() + af.opcrLength()
}

// transportPrivateDataLength returns the length of the transport private data,
// if there is no transport private data then its length is zero
func (af AdaptationField) transportPrivateDataLength() int {
	if af.HasTransportPrivateData() {
		// cannot extend beyond adaptation field, number of bytes
		// for field stored in transportPrivateDataLength
		return 1 + int(af[af.transportPrivateDataStart()])
	}
	return 0
}

// transportPrivateDataStart returns the start index of where the transport private data should be
func (af AdaptationField) transportPrivateDataStart() int {
	return pcrStart + af.pcrLength() + af.opcrLength() + af.spliceCountdownLength()
}

// adaptationExtensionLength returns the length of the adaptation field extension,
//  if there is no adaptation field extension then its length is zero
func (af AdaptationField) adaptationExtensionLength() int {
	if af.HasAdaptationFieldExtension() {
		return 1 + int(af[af.adaptationExtensionStart()])
	}
	return 0
}

// adaptationExtensionStartreturns the length of the adaptation field extension,
// if there is no adaptation field extension then its length is zero
func (af AdaptationField) adaptationExtensionStart() int {
	return pcrStart + af.pcrLength() + af.opcrLength() +
		af.spliceCountdownLength() + af.transportPrivateDataLength()
}

// calculates the length of the Adaptation Field
// (with respect to the start of the packet) excluding stuffing
func (af AdaptationField) calculateMinLength() int {
	return pcrStart +
		af.pcrLength() + af.opcrLength() + af.spliceCountdownLength() +
		af.transportPrivateDataLength() + af.adaptationExtensionLength()
}

func (af AdaptationField) setBit(index int, mask byte, value bool) {
	if af.invalid() {
		return
	}
	if value {
		af[index] |= mask
	} else {
		af[index] &= ^mask
	}
}

func (af AdaptationField) getBit(index int, mask byte) bool {
	if af.invalid() {
		return false
	}
	return af[index]&mask != 0
}

func (af AdaptationField) bitDelta(index int, mask byte, value bool) int {
	oldValue := af.getBit(index, mask)
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
	if delta > 0 { // shifting for growing
		end := af.calculateMinLength()
		startRight := start + delta
		endRight := af.calculateMinLength() + delta
		src := af[start:end]
		dst := af[startRight:endRight]
		copy(dst, src)
		if af.Length() < endRight-5 {
			// erase payload, it is corrupt anyways
			for i := endRight; i < len(af); i++ {
				af[i] = 0xFF // packet is stuffed until the very end.
			}
			// payload must be at least one byte in size.
			// this will remind the user of the library
			// that the payload was destroyed
			af.setLength(183)
		}
	}
	if delta < 0 {
		startRight := start - delta
		endRight := af.calculateMinLength()
		end := endRight + delta
		src := []byte(af[startRight:endRight])
		dst := []byte(af[start:end])
		copy(dst, src)
		for i := end; i < endRight; i++ {
			af[i] = 0xFF // fill in the gap with stuffing
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
	delta := 6 * af.bitDelta(5, 0x10, value)
	af.resizeAF(pcrStart, delta)
	af.setBit(5, 0x10, value)
}

func (af AdaptationField) HasPCR() bool {
	return af.getBit(5, 0x10)
}

func (af AdaptationField) SetPCR(PCR uint64) {
	if !af.HasPCR() {
		return
	}
	gots.InsertPCR(af[pcrStart:af.opcrStart()], PCR)
}

func (af AdaptationField) PCR() uint64 {
	if !af.HasPCR() {
		return 0
	}
	return gots.ExtractPCR(af[pcrStart:af.opcrStart()])
}

func (af AdaptationField) SetHasOPCR(value bool) {
	delta := 6 * af.bitDelta(5, 0x08, value)
	af.resizeAF(af.opcrStart(), delta)
	af.setBit(5, 0x08, value)
}

func (af AdaptationField) HasOPCR() bool {
	return af.getBit(5, 0x08)
}

func (af AdaptationField) SetOPCR(PCR uint64) {
	if !af.HasOPCR() {
		return
	}
	gots.InsertPCR(af[af.opcrStart():af.spliceCountdownStart()], PCR)
}

func (af AdaptationField) OPCR() uint64 {
	if !af.HasOPCR() {
		return 0
	}
	return gots.ExtractPCR(af[af.opcrStart():af.spliceCountdownStart()])
}

func (af AdaptationField) SetHasSplicingPoint(value bool) {
	delta := 1 * af.bitDelta(5, 0x04, value)
	af.resizeAF(af.spliceCountdownStart(), delta)
	af.setBit(5, 0x04, value)
}

func (af AdaptationField) HasSplicingPoint() bool {
	return af.getBit(5, 0x04)
}

func (af AdaptationField) SetSpliceCountdown(value int) {
	if !af.HasSplicingPoint() {
		return
	}
	af[af.spliceCountdownStart()] = byte(value)
}

func (af AdaptationField) SpliceCountdown() int {
	if !af.HasSplicingPoint() {
		return 0
	}
	return int(int8(af[af.spliceCountdownStart()])) // int8 cast is for 2s complement numbers
}

func (af AdaptationField) SetHasTransportPrivateData(value bool) {
	delta := 1 * af.bitDelta(5, 0x02, value)
	af.resizeAF(af.transportPrivateDataStart(), delta)
	af[af.transportPrivateDataStart()] = 0 // zero length by default
	af.setBit(5, 0x02, value)
}

func (af AdaptationField) HasTransportPrivateData() bool {
	return af.getBit(5, 0x02)
}

func (af AdaptationField) SetTransportPrivateData(data []byte) {
	delta := len(data) - (af.transportPrivateDataLength() - 1)
	start := af.transportPrivateDataStart() + 1
	end := start + len(data)
	af.resizeAF(start, delta)
	copy(af[start:end], data)
	af[start-1] = byte(len(data))
}

func (af AdaptationField) TransportPrivateData() []byte {
	return af[af.transportPrivateDataStart():af.adaptationExtensionStart()]
}

func (af AdaptationField) SetHasAdaptationFieldExtension(value bool) {
	delta := 1 * af.bitDelta(5, 0x01, value)
	af.resizeAF(af.adaptationExtensionStart(), delta)
	af[af.adaptationExtensionStart()] = 0
	af.setBit(5, 0x01, value)
}

func (af AdaptationField) HasAdaptationFieldExtension() bool {
	return af.getBit(5, 0x01)
}

func (af AdaptationField) SetAdaptationFieldExtension(data []byte) {
	delta := len(data) - (af.adaptationExtensionLength() - 1)
	start := af.adaptationExtensionStart() + 1
	end := start + len(data)
	af.resizeAF(start, delta)
	copy(af[start:end], data)
	af[start-1] = byte(len(data))
}

func (af AdaptationField) AdaptationFieldExtension() []byte {
	return af[af.adaptationExtensionStart():af.calculateMinLength()]
}

func (af AdaptationField) String() string {
	if af.invalid() {
		return "Null"
	}
	return fmt.Sprintf("%X", []byte(af[4:]))
}
