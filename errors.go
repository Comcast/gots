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

package gots

import "errors"

var (
	// ErrUnrecognizedEbpType is returned if the EBP cannot be parsed
	ErrUnrecognizedEbpType = errors.New("unrecognized EBP")
	// ErrNoEBP is returned when an attempt is made to extract an EBP from a packet that does not contain one
	ErrNoEBP = errors.New("packet does not contain EBP")
	// ErrInvalidPacketLength denotes an packet length that is not packet.PacketSize bytes in length
	ErrInvalidPacketLength = errors.New("invalid packet length")
	// ErrNoPayload denotes that the attempted operation is not valid on a packet with no payload
	ErrNoPayload = errors.New("packet does not contain payload")
	// ErrNoPrivateTransportData is returned when an attempt is made to access private transport data when none exists
	ErrNoPrivateTransportData = errors.New("adaptation field has no private transport data")
	// ErrNoSplicePoint is returned when an attempt to access a splice countdown and no splice point exists
	ErrNoSplicePoint = errors.New("adaptation field has no splice point")
	// ErrNoPCR is returned when an attempt is made to access adaptation field PRC that does not exist
	ErrNoPCR = errors.New("adaptation field has no Program Clock Reference")
	// ErrNoOPCR is returned when an attempt is made to access an adaptation field OPCR that does not exist
	ErrNoOPCR = errors.New("adaptation field has no Original Program Clock Reference")
	// ErrPATNotFound is returned when expected PAT packet is not found when
	// reading TS packets.
	ErrPATNotFound = errors.New("No PAT was found while reading TS")
	// ErrPMTNotFound is returned when expected PMT packet(s) are not found when
	// reading TS packets.
	ErrPMTNotFound = errors.New("No PMT was found while reading TS")
	// ErrParsePMTDescriptor is returned when a PMT descriptor cannot be parsed
	ErrParsePMTDescriptor = errors.New("unable to parse PMT descriptor")
	// ErrInvalidPATLength is returned when a PAT cannot be parsed because there are not enough bytes
	ErrInvalidPATLength = errors.New("too few bytes to parse PAT")
	// ErrNoPayloadUnitStartIndicator should be returned when a packet is expected to have a PUSI and does not.
	ErrNoPayloadUnitStartIndicator = errors.New("packet does not have payload unit start indicator")
	// ErrUnknownTableID is returned when PSI is parsed with an unknown table id
	ErrUnknownTableID = errors.New("Unknown table id received")
	// ErrInvalidSCTE35Length is returned when a SCTE35 cue cannot be parsed because there are not enough bytes
	ErrInvalidSCTE35Length = errors.New("too few bytes to parse SCTE35")
	// ErrSCTE35EncryptionUnsupported is returned when a scte35 cue cannot be parsed because it is encrypted
	ErrSCTE35EncryptionUnsupported = errors.New("SCTE35 is encrypted, which is not supported")
	// ErrSCTE35UnsupportedSpliceCommand is returned when a SCTE35 cue
	// cannot be parsed because the command type is not supported
	ErrSCTE35UnsupportedSpliceCommand = errors.New("SCTE35 cue can't be parsed because only time_signal with a pts value and splice_null commands are supported")
	// ErrSCTE35InvalidDescriptorID is returned when a segmentation descriptor is found with an id that is not CUEI
	ErrSCTE35InvalidDescriptorID = errors.New("SCTE35 segmentation descriptor has a id that is not \"CUEI\"")
	// ErrSCTE35DuplicateSignal is returned when a duplicate or equivalent descriptor is received by state
	ErrSCTE35DuplicateDescriptor = errors.New("Duplicate or equivalent descriptor received by scte35.State")
	// ErrSCTE35InvalidDescriptor is returned when a descriptor is invalid given the current state (i.e. a ProgramResumption received when no it breakaway)
	ErrSCTE35InvalidDescriptor = errors.New("Invalid descriptor given the current state")
	// ErrSCTE35MissingOut is returned when an in descriptor is received by state with no matching out
	ErrSCTE35MissingOut = errors.New("In descriptor received with no matching out")
	// ErrSCTE35DescriptorNotFound is returned when a descriptor is closed that's not in the open list
	ErrSCTE35DescriptorNotFound = errors.New("Cannot close descriptor that's not in the open list")
	// ErrNilPAT is returned when a PAT is passed into a function for which it cannot be nil.
	ErrNilPAT = errors.New("Nil PAT not allowed here.")
	// ErrSyncByteNotFound is returned when a packet sync byte could not be found
	// when reading.
	ErrSyncByteNotFound = errors.New("Sync-byte not found.")
)
