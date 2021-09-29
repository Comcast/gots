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
	"github.com/Comcast/gots"
)

const (
	PayloadFlag                   AdaptationFieldControlOptions = 1 // 10
	AdaptationFieldFlag           AdaptationFieldControlOptions = 2 // 01
	PayloadAndAdaptationFieldFlag AdaptationFieldControlOptions = 3 // 11
)

// AdaptationField is an optional part of the packet.
type AdaptationField Packet

// AdaptationFieldControlOptions is a set of constants for
// selecting the adaptation field control.
type AdaptationFieldControlOptions byte

// NewPacket creates a new packet with a Null ID, sync byte, and with the adaptation field control set to payload only.
// This function is error free.
func NewAdaptationField() *AdaptationField {
	p := New()
	p.SetAdaptationFieldControl(AdaptationFieldFlag)
	af, _ := p.AdaptationField()
	return af
}

// setBit sets a bit in the adaptation field
// index is the byte in the packet where the bit is located
// mask is the bitmask that has that bit set to one and the rest of the bits set to zero
// value is the value that the bit will be set to
func (af *AdaptationField) setBit(index int, mask byte, value bool) {
	if value {
		af[index] |= mask
	} else {
		af[index] &= ^mask
	}
}

// getBit gets a bit in the adaptation field
// index is the byte in the packet where the bit is located
// mask is the bitmask that has that bit set to one and the rest of the bits set to zero
// value is the value of that bit in the adaptation field
func (af *AdaptationField) getBit(index int, mask byte) bool {
	return af[index]&mask != 0
}

// bitDelta returns the difference of a bit (boolean) and the bit that is
// currently set in the adaptation field.
// 1 is returned if the packet is changed from 0 to 1.
// -1 is returned if the bit is changed from 1 to 0.
// 0 is returned if the bit is unchanged.
// this can be used to find if a field is growing or shrinking.
func (af *AdaptationField) bitDelta(index int, mask byte, value bool) int {
	if value == af.getBit(index, mask) {
		return 0 // same
	}
	if value {
		return 1 // growing
	}
	return -1 // shrinking
}

// valid returns any errors that prevent the AdaptationField from being valid.
func (af *AdaptationField) valid() error {
	if !af.getBit(3, 0x20) {
		return gots.ErrNoAdaptationField
	}
	if af[4] == 0 {
		return gots.ErrAdaptationFieldZeroLength
	}
	return nil
}

// initAdaptationField initializes the adaptation field to have all false flags
// it will also occupy the remainder of the packet.
func initAdaptationField(p *Packet) {
	af := (*AdaptationField)(p)
	af[4] = 183  // adaptation field will take up the rest of the packet by Default
	af[5] = 0x00 // no flags are set by default
	for i := 6; i < len(af); i++ {
		af[i] = 0xFF
	}
}

// returns if the adaptation field has a PCR, this does not check for errors.
func (af *AdaptationField) hasPCR() bool {
	return af.getBit(5, 0x10)
}

// returns if the adaptation field has an OPCR, this does not check for errors.
func (af *AdaptationField) hasOPCR() bool {
	return af.getBit(5, 0x08)
}

// returns the adaptation field's Splicing Point flag, this does not check for errors.
func (af *AdaptationField) hasSplicingPoint() bool {
	return af.getBit(5, 0x04)
}

// returns if the adaptation field has a Transport Private Data, this does not check for errors.
func (af *AdaptationField) hasTransportPrivateData() bool {
	return af.getBit(5, 0x02)
}

// returns if the adaptation field has an adaptation field extension, this does not check for errors.
func (af *AdaptationField) hasAdaptationFieldExtension() bool {
	return af.getBit(5, 0x01)
}

// pcrLength returns the length of the PCR, if there is no PCR then its length is zero
func (af *AdaptationField) pcrLength() int {
	if af.hasPCR() {
		return 6
	}
	return 0
}

const pcrStart = 6 // start of the pcr with respect to the start of the packet

// opcrLength returns the length of the OPCR, if there is no OPCR then its length is zero
func (af *AdaptationField) opcrLength() int {
	if af.hasOPCR() {
		return 6
	}
	return 0
}

// opcrStart returns the start index of where the OPCR field should
// be with respect to the start of the packet.
func (af *AdaptationField) opcrStart() int {
	return pcrStart + af.pcrLength()
}

// spliceCountdownLength returns the length of the splice countdown,
// if there is no splice countdown then its length is zero.
func (af *AdaptationField) spliceCountdownLength() int {
	if af.hasSplicingPoint() {
		return 1
	}
	return 0
}

// spliceCountdownStart returns the start index of where the splice
// countdown field should be with respect to the start of the packet.
func (af *AdaptationField) spliceCountdownStart() int {
	return pcrStart + af.pcrLength() + af.opcrLength()
}

// transportPrivateDataLength returns the length of the transport private data,
// if there is no transport private data then its length is zero.
func (af *AdaptationField) transportPrivateDataLength() int {
	if !af.hasTransportPrivateData() {
		return 0
	}

	if af.transportPrivateDataStart() >= PacketSize {
		return 0
	}

	// cannot extend beyond adaptation field, number of bytes
	// for field stored in transportPrivateDataLength
	return 1 + int(af[af.transportPrivateDataStart()])
}

// transportPrivateDataStart returns the start index of where the
// transport private data should be with respect to the start of the packet.
func (af *AdaptationField) transportPrivateDataStart() int {
	return pcrStart + af.pcrLength() + af.opcrLength() + af.spliceCountdownLength()
}

// adaptationExtensionLength returns the length of the adaptation field extension,
// if there is no adaptation field extension then its length is zero
func (af *AdaptationField) adaptationExtensionLength() int {
	if !af.hasAdaptationFieldExtension() {
		return 0
	}

	if af.adaptationExtensionStart() >= PacketSize {
		return 0
	}

	return 1 + int(af[af.adaptationExtensionStart()])
}

// adaptationExtensionStart returns the start index of where the
// adaptation extension start should be with respect to the start of the packet.
func (af *AdaptationField) adaptationExtensionStart() int {
	return pcrStart + af.pcrLength() + af.opcrLength() +
		af.spliceCountdownLength() + af.transportPrivateDataLength()
}

// stuffingStart returns the start index of where the
// stuffing bytes should be with respect to the start of the packet.
func (af *AdaptationField) stuffingStart() int {
	return pcrStart +
		af.pcrLength() + af.opcrLength() + af.spliceCountdownLength() +
		af.transportPrivateDataLength() + af.adaptationExtensionLength()
}

// stuffingEnd returns the index where the stuffing bytes end
// (first index without stuffing bytes) with respect to the start of the packet.
func (af *AdaptationField) stuffingEnd() int {
	stuffingEnd := int(af[4]) + 5

	if stuffingEnd >= PacketSize {
		return PacketSize - 1
	}

	return stuffingEnd
}

// setLength sets the length field of the adaptation field.
func (af *AdaptationField) setLength(length int) {
	af[4] = byte(length)
}

// stuffAF will ensure that the stuffing bytes are all 0xFF.
func (af *AdaptationField) stuffAF() {
	for i := af.stuffingStart(); i < af.stuffingEnd(); i++ {
		af[i] = 0xFF // stuffing byte must be 0xFF
	}
}

// resizeAF will resize the adaptation field to insert a new field into it.
// start is the start of the field being manipulated (smallest index)
// delta is how much shifting needs to be done.
// this function must be called before the field is marked as present.
func (af *AdaptationField) resizeAF(start int, delta int) error {
	if delta > 0 { // shifting for growing
		end := af.stuffingStart()
		startRight := start + delta
		endRight := af.stuffingStart() + delta
		if af.stuffingEnd() < endRight {
			// cannot grow in size, payload will be corrupted.
			return gots.ErrAdaptationFieldCannotGrow
		}
		src := af[start:end]
		dst := af[startRight:endRight]
		copy(dst, src)
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
	return nil
}

// Length returns the length of the adaptation field
func (af *AdaptationField) Length() int {
	return int(af[4])
}

// SetDiscontinuity sets the Discontinuity field of the packet.
func (af *AdaptationField) SetDiscontinuity(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	af.setBit(5, 0x80, value)
	return nil
}

// Discontinuity returns the value of the discontinuity field in the packet.
func (af *AdaptationField) Discontinuity() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.getBit(5, 0x80), nil
}

// SetRandomAccess sets the value of the random access field in the packet.
func (af *AdaptationField) SetRandomAccess(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	af.setBit(5, 0x40, value)
	return nil
}

// RandomAccess returns the value of the random access field in the packet.
func (af *AdaptationField) RandomAccess() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.getBit(5, 0x40), nil
}

// SetElementaryStreamPriority sets the Elementary Stream Priority Flag.
func (af *AdaptationField) SetElementaryStreamPriority(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	af.setBit(5, 0x20, value)
	return nil
}

// ElementaryStreamPriority returns the Elementary Stream Priority Flag.
func (af *AdaptationField) ElementaryStreamPriority() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.getBit(5, 0x20), nil
}

// SetHasPCR sets HasPCR
// HasPCR determines if the packet has a PCR
func (af *AdaptationField) SetHasPCR(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 6 * af.bitDelta(5, 0x10, value)
	err := af.resizeAF(pcrStart, delta)
	if err != nil {
		return err
	}
	af.setBit(5, 0x10, value)
	return nil
}

// HasPCR returns if the packet has a PCR
func (af *AdaptationField) HasPCR() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasPCR(), nil
}

// SetPCR sets the PCR of the adaptation field.
// If impossible an error is returned.
func (af *AdaptationField) SetPCR(PCR uint64) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasPCR() {
		return gots.ErrNoPCR
	}
	gots.InsertPCR(af[pcrStart:af.opcrStart()], PCR)
	return nil
}

// PCR returns the PCR from the adaptation field if possible.
// if it is not possible an error is returned.
func (af *AdaptationField) PCR() (uint64, error) {
	if err := af.valid(); err != nil {
		return 0, err
	}
	if !af.hasPCR() {
		return 0, gots.ErrNoPCR
	}
	return gots.ExtractPCR(af[pcrStart:af.opcrStart()]), nil
}

// SetHasOPCR sets HasOPCR
// HasOPCR determines if the packet has a OPCR
func (af *AdaptationField) SetHasOPCR(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 6 * af.bitDelta(5, 0x08, value)
	err := af.resizeAF(af.opcrStart(), delta)
	if err != nil {
		return err
	}
	af.setBit(5, 0x08, value)
	return nil
}

// HasOPCR returns if the packet has an OPCR
func (af *AdaptationField) HasOPCR() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasOPCR(), nil
}

// SetOPCR sets the OPCR of the adaptation field.
// If impossible an error is returned.
func (af *AdaptationField) SetOPCR(PCR uint64) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasOPCR() {
		return gots.ErrNoOPCR
	}
	gots.InsertPCR(af[af.opcrStart():af.spliceCountdownStart()], PCR)
	return nil
}

// OPCR returns the OPCR from the adaptation field if possible.
// if it is not possible an error is returned.
func (af *AdaptationField) OPCR() (uint64, error) {
	if err := af.valid(); err != nil {
		return 0, err
	}
	if !af.hasOPCR() {
		return 0, gots.ErrNoOPCR
	}
	return gots.ExtractPCR(af[af.opcrStart():af.spliceCountdownStart()]), nil
}

// SetHasSplicingPoint sets HasSplicingPoint
// HasSplicingPoint determines if the packet has a Splice Countdown
func (af *AdaptationField) SetHasSplicingPoint(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 1 * af.bitDelta(5, 0x04, value)
	err := af.resizeAF(af.spliceCountdownStart(), delta)
	if err != nil {
		return err
	}
	af.setBit(5, 0x04, value)
	return nil
}

// HasSplicingPoint returns if the packet has a Splice Countdown
func (af *AdaptationField) HasSplicingPoint() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasSplicingPoint(), nil
}

// SetSpliceCountdown sets the Splice Countdown of the adaptation field.
// If impossible an error is returned.
func (af *AdaptationField) SetSpliceCountdown(value byte) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasSplicingPoint() {
		return gots.ErrNoSplicePoint
	}
	af[af.spliceCountdownStart()] = value
	return nil
}

// SpliceCountdown returns the Splice Countdown from the adaptation field if possible.
// if it is not possible an error is returned.
func (af *AdaptationField) SpliceCountdown() (int, error) {
	if err := af.valid(); err != nil {
		return 0, err
	}
	if !af.hasSplicingPoint() {
		return 0, gots.ErrNoSplicePoint
	}
	return int(int8(af[af.spliceCountdownStart()])), nil // int8 cast is for 2s complement numbers
}

// SetHasTransportPrivateData sets HasTransportPrivateData
// HasTransportPrivateData determines if the packet has Transport Private Data
func (af *AdaptationField) SetHasTransportPrivateData(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 1 * af.bitDelta(5, 0x02, value)
	err := af.resizeAF(af.transportPrivateDataStart(), delta)
	if err != nil {
		return err
	}
	af[af.transportPrivateDataStart()] = 0 // zero length by default
	af.setBit(5, 0x02, value)
	return nil
}

// HasTransportPrivateData returns if the packet has an Transport Private Data
func (af *AdaptationField) HasTransportPrivateData() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasTransportPrivateData(), nil
}

// SetTransportPrivateData sets the Transport Private Data of the adaptation field.
// If impossible an error is returned.
func (af *AdaptationField) SetTransportPrivateData(data []byte) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasTransportPrivateData() {
		return gots.ErrNoPrivateTransportData
	}
	delta := len(data) - (af.transportPrivateDataLength() - 1)
	start := af.transportPrivateDataStart() + 1
	end := start + len(data)
	err := af.resizeAF(start, delta)
	if err != nil {
		return err
	}
	copy(af[start:end], data)
	af[start-1] = byte(len(data))
	return nil
}

// TransportPrivateData returns the Transport Private Data from the adaptation field if possible.
// if it is not possible an error is returned.
func (af *AdaptationField) TransportPrivateData() ([]byte, error) {
	hasTPD, err := af.HasTransportPrivateData()
	if err != nil {
		return nil, err
	}
	if !hasTPD {
		return nil, gots.ErrNoPrivateTransportData
	}
	return af[af.transportPrivateDataStart():af.adaptationExtensionStart()], nil
}

// SetHasAdaptationFieldExtension sets HasAdaptationFieldExtension
// HasAdaptationFieldExtension determines if the packet has an Adaptation Field Extension
func (af *AdaptationField) SetHasAdaptationFieldExtension(value bool) error {
	if err := af.valid(); err != nil {
		return err
	}
	delta := 1 * af.bitDelta(5, 0x01, value)
	err := af.resizeAF(af.adaptationExtensionStart(), delta)
	if err != nil {
		return err
	}
	af[af.adaptationExtensionStart()] = 0
	af.setBit(5, 0x01, value)
	return nil
}

// HasAdaptationFieldExtension returns if the packet has an Adaptation Field Extension
func (af *AdaptationField) HasAdaptationFieldExtension() (bool, error) {
	if err := af.valid(); err != nil {
		return false, err
	}
	return af.hasAdaptationFieldExtension(), nil
}

// SetAdaptationFieldExtension sets the Adaptation Field Extension of the adaptation field.
// If impossible an error is returned.
func (af *AdaptationField) SetAdaptationFieldExtension(data []byte) error {
	if err := af.valid(); err != nil {
		return err
	}
	if !af.hasAdaptationFieldExtension() {
		return gots.ErrNoAdaptationFieldExtension
	}
	delta := len(data) - (af.adaptationExtensionLength() - 1)
	start := af.adaptationExtensionStart() + 1
	end := start + len(data)
	err := af.resizeAF(start, delta)
	if err != nil {
		return err
	}
	copy(af[start:end], data)
	af[start-1] = byte(len(data))
	return nil
}

// AdaptationFieldExtension returns the Adaptation Field Extension from the adaptation field if possible.
// if it is not possible an error is returned.
func (af *AdaptationField) AdaptationFieldExtension() ([]byte, error) {
	hasAFC, err := af.HasAdaptationFieldExtension()
	if err != nil {
		return nil, err
	}
	if !hasAFC {
		return nil, gots.ErrNoAdaptationFieldExtension
	}
	return af[af.adaptationExtensionStart():af.stuffingStart()], nil
}
