package packet

import (
	"fmt"
	"github.com/Comcast/gots"
)

// valid returns true if the length of the packet slice is
// anything but PacketSize (188)
func (af AdaptationField) valid() error {
	if len(af) != PacketSize {
		return gots.ErrInvalidPacketLength
	}
	if hasAF, _ := Packet(af).HasAdaptationField(); !hasAF {
		return gots.ErrNoAdaptationField
	}
	return nil
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
	return AdaptationField(p)
}

func (af AdaptationField) hasPCR() bool {
	return af.getBit(5, 0x10)
}
func (af AdaptationField) hasOPCR() bool {
	return af.getBit(5, 0x08)
}

func (af AdaptationField) hasSplicingPoint() bool {
	return af.getBit(5, 0x04)
}

func (af AdaptationField) hasTransportPrivateData() bool {
	return af.getBit(5, 0x02)
}

func (af AdaptationField) hasAdaptationFieldExtension() bool {
	return af.getBit(5, 0x01)
}

// pcrLength returns the length of the PCR, if there is no PCR then its length is zero
func (af AdaptationField) pcrLength() int {
	if af.hasPCR() {
		return 6
	}
	return 0
}

const pcrStart = 6

// opcrLength returns the length of the OPCR, if there is no OPCR then its length is zero
func (af AdaptationField) opcrLength() int {
	if af.hasOPCR() {
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
	if af.hasSplicingPoint() {
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
	if af.hasTransportPrivateData() {
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
	if af.hasAdaptationFieldExtension() {
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
func (af AdaptationField) stuffingStart() int {
	return pcrStart +
		af.pcrLength() + af.opcrLength() + af.spliceCountdownLength() +
		af.transportPrivateDataLength() + af.adaptationExtensionLength()
}

// length returns the length of the adaptation field
func (af AdaptationField) stuffingEnd() int {
	return int(af[4]) + 5
}

func (af AdaptationField) setBit(index int, mask byte, value bool) {
	if value {
		af[index] |= mask
	} else {
		af[index] &= ^mask
	}
}

func (af AdaptationField) getBit(index int, mask byte) bool {
	return af[index]&mask != 0
}

func (af AdaptationField) bitDelta(index int, mask byte, value bool) int {
	if value != af.getBit(index, mask) {
		if value {
			return 1 // growing
		}
		return -1 // shrinking
	}
	return 0 // same
}

func (af AdaptationField) setLength(length int) {
	af[4] = byte(length)
}

func (af AdaptationField) stuffAF() {
	for i := af.stuffingStart(); i < af.stuffingEnd(); i++ {
		af[i] = 0xFF // stuffing byte must be 0xFF
	}
}

func (af AdaptationField) resizeAF(start int, delta int) {
	if delta > 0 { // shifting for growing
		end := af.stuffingStart()
		startRight := start + delta
		endRight := af.stuffingStart() + delta
		src := af[start:end]
		dst := af[startRight:endRight]
		copy(dst, src)
		if af.stuffingEnd() < endRight {
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
		endRight := af.stuffingStart()
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
func (af AdaptationField) Length() (int, error) {
	if err := af.valid(); err != nil {
		return 0, err
	}
	return int(af[4]), nil
}

// SetDiscontinuity sets the Discontinuity field of the packet
func (af AdaptationField) SetDiscontinuity(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	af.setBit(5, 0x80, value)
	return nil
}

// Discontinuity returns the value of the discontinuity field in the packet
func (af AdaptationField) Discontinuity() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.getBit(5, 0x80), nil
}

func (af AdaptationField) SetRandomAccess(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	af.setBit(5, 0x40, value)
	return nil
}

func (af AdaptationField) RandomAccess() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.getBit(5, 0x40), nil
}

func (af AdaptationField) SetESPriority(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	af.setBit(5, 0x20, value)
	return nil
}

func (af AdaptationField) ESPriority() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.getBit(5, 0x20), nil
}

func (af AdaptationField) SetHasPCR(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 6 * af.bitDelta(5, 0x10, value)
	af.resizeAF(pcrStart, delta)
	af.setBit(5, 0x10, value)
	return nil
}

func (af AdaptationField) HasPCR() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasPCR(), nil
}

func (af AdaptationField) SetPCR(PCR uint64) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasPCR() {
		return gots.ErrNoPCR
	}
	gots.InsertPCR(af[pcrStart:af.opcrStart()], PCR)
	return nil
}

func (af AdaptationField) PCR() (uint64, error) {
	if err := af.valid(); err != nil {
		return 0, err
	}
	if !af.hasPCR() {
		return 0, gots.ErrNoPCR
	}
	return gots.ExtractPCR(af[pcrStart:af.opcrStart()]), nil
}

func (af AdaptationField) SetHasOPCR(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 6 * af.bitDelta(5, 0x08, value)
	af.resizeAF(af.opcrStart(), delta)
	af.setBit(5, 0x08, value)
	return nil
}

func (af AdaptationField) HasOPCR() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasOPCR(), nil
}

func (af AdaptationField) SetOPCR(PCR uint64) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasOPCR() {
		return gots.ErrNoOPCR
	}
	gots.InsertPCR(af[af.opcrStart():af.spliceCountdownStart()], PCR)
	return nil
}

func (af AdaptationField) OPCR() (uint64, error) {
	if err := af.valid(); err != nil {
		return 0, err
	}
	if !af.hasOPCR() {
		return 0, gots.ErrNoOPCR
	}
	return gots.ExtractPCR(af[af.opcrStart():af.spliceCountdownStart()]), nil
}

func (af AdaptationField) SetHasSplicingPoint(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 1 * af.bitDelta(5, 0x04, value)
	af.resizeAF(af.spliceCountdownStart(), delta)
	af.setBit(5, 0x04, value)
	return nil
}

func (af AdaptationField) HasSplicingPoint() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasSplicingPoint(), nil
}

func (af AdaptationField) SetSpliceCountdown(value int) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasSplicingPoint() {
		return gots.ErrNoSplicePoint
	}
	af[af.spliceCountdownStart()] = byte(value)
	return nil
}

func (af AdaptationField) SpliceCountdown() (int, error) {
	if err := af.valid(); err != nil {
		return 0, err
	}
	if !af.hasSplicingPoint() {
		return 0, gots.ErrNoSplicePoint
	}
	return int(int8(af[af.spliceCountdownStart()])), nil // int8 cast is for 2s complement numbers
}

func (af AdaptationField) SetHasTransportPrivateData(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 1 * af.bitDelta(5, 0x02, value)
	af.resizeAF(af.transportPrivateDataStart(), delta)
	af[af.transportPrivateDataStart()] = 0 // zero length by default
	af.setBit(5, 0x02, value)
	return nil
}

func (af AdaptationField) HasTransportPrivateData() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasTransportPrivateData(), nil
}

func (af AdaptationField) SetTransportPrivateData(data []byte) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasTransportPrivateData() {
		return gots.ErrNoPrivateTransportData
	}
	delta := len(data) - (af.transportPrivateDataLength() - 1)
	start := af.transportPrivateDataStart() + 1
	end := start + len(data)
	af.resizeAF(start, delta)
	copy(af[start:end], data)
	af[start-1] = byte(len(data))
	return nil
}

func (af AdaptationField) TransportPrivateData() []byte {
	return af[af.transportPrivateDataStart():af.adaptationExtensionStart()]
}

func (af AdaptationField) SetHasAdaptationFieldExtension(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 1 * af.bitDelta(5, 0x01, value)
	af.resizeAF(af.adaptationExtensionStart(), delta)
	af[af.adaptationExtensionStart()] = 0
	af.setBit(5, 0x01, value)
	return nil
}

func (af AdaptationField) HasAdaptationFieldExtension() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasAdaptationFieldExtension(), nil
}

func (af AdaptationField) SetAdaptationFieldExtension(data []byte) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasAdaptationFieldExtension() {
		return gots.ErrNoAdaptationFieldExtension
	}
	delta := len(data) - (af.adaptationExtensionLength() - 1)
	start := af.adaptationExtensionStart() + 1
	end := start + len(data)
	af.resizeAF(start, delta)
	copy(af[start:end], data)
	af[start-1] = byte(len(data))
	return nil
}

func (af AdaptationField) AdaptationFieldExtension() []byte {
	return af[af.adaptationExtensionStart():af.stuffingStart()]
}

func (af AdaptationField) String() string {
	if af.valid() != nil {
		return "Null"
	}
	return fmt.Sprintf("%X", []byte(af[4:]))
}
