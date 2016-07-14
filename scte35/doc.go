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

// Package scte35 is for handling scte35 splice signals
package scte35

import (
	"github.com/comcast/gots/pes"
)

// SpliceCommandType - not really needed for processing but included for
// backwards compatibility/porting
type SpliceCommandType uint16

const (
	// SpliceNull is a Null Splice command type
	SpliceNull SpliceCommandType = 0x00
	// SpliceSchedule is a splice schedule command type
	SpliceSchedule = 0x04
	// SpliceInsert is a splice insert command type
	SpliceInsert = 0x05
	// TimeSignal is a splice signal command type
	TimeSignal = 0x06
	// BandwidthReservation is a command type that represents a reservation of bandwidth
	BandwidthReservation = 0x07
	// PrivateCommand is a command type that represents private command data
	PrivateCommand = 0xFF
)

// SegDescType is the Segmentation Descriptor Type - not really needed for processing according
// to method below, but included here for backwards compatibility/porting
type SegDescType uint8

const (
	SegDescNotIndicated                  SegDescType = 0x00
	SegDescContentIdentification                     = 0x01
	SegDescProgramStart                              = 0x10
	SegDescProgramEnd                                = 0x11
	SegDescProgramEarlyTermination                   = 0x12
	SegDescProgramBreakaway                          = 0x13
	SegDescProgramResumption                         = 0x14
	SegDescProgramRunoverPlanned                     = 0x15
	SegDescProgramRunoverUnplanned                   = 0x16
	SegDescProgramOverlapStart                       = 0x17
	SegDescProgramBlackoutOverride                   = 0x18
	SegDescChapterStart                              = 0x20
	SegDescChapterEnd                                = 0x21
	SegDescProviderAdvertisementStart                = 0x30
	SegDescProviderAdvertisementEnd                  = 0x31
	SegDescDistributorAdvertisementStart             = 0x32
	SegDescDistributorAdvertisementEnd               = 0x33
	SegDescPlacementOpportunityStart                 = 0x34
	SegDescPlacementOpportunityEnd                   = 0x35
	SegDescDistributorPoStart                        = 0x36
	SegDescDistributorPoEnd                          = 0x37
	SegDescUnscheduledEventStart                     = 0x40
	SegDescUnscheduledEventEnd                       = 0x41
	SegDescNetworkStart                              = 0x50
	SegDescNetworkEnd                                = 0x51
)

// SegUPIDType is the Segmentation UPID Types - Only type that really needs to be checked is
// SegUPIDURN for CSP
type SegUPIDType uint8

const (
	SegUPIDNotUsed     SegUPIDType = 0x00
	SegUPIDUserDefined             = 0x01
	SegUPIDISCI                    = 0x02
	SegUPIDAdID                    = 0x03
	SegUPIDUMID                    = 0x04
	SegUPIDISAN                    = 0x05
	SegUPIDVISAN                   = 0x06
	SegUPIDTID                     = 0x07
	SegUPIDTI                      = 0x08
	SegUPIDADI                     = 0x09
	SegUPIDEIDR                    = 0x0a
	SegUPIDATSCID                  = 0x0b
	SegUPIDMPU                     = 0x0c
	SegUPIDMID                     = 0x0d
	SegUPIDURN                     = 0x0e
)

// SCTE35 represent operations available on a SCTE35 message.
type SCTE35 interface {
	// HasPTS returns true if there is a pts time
	HasPTS() bool
	// PTS returns the PTS time of the signal if it exists
	PTS() pes.PTS
	// Command returns the signal's splice command
	Command() SpliceCommandType
	// Descriptors returns a slice of the signals SegmentationDescriptors
	Descriptors() []SegmentationDescriptor
	// Data returns the raw data bytes of the scte signal
	Data() []byte
}

// SegmentationDescriptor describes the segmentation descriptor interface.  The intended usage is
// to maintain a sorted list of descriptors.  When a new signal is received for
// every descriptor returned from the signal, walk the list, using Compare() to
// find the place in the list.  If Compare()==0, check if is Equal() for
// duplicates and then CanClose() to see if can close the signal and remove
// from the list.  If Compare() goes from < to > and signal is an out, insert
// it into the list.  Some pseudo code is below (additional edge
// cases/bookkeeping not included):
// scte,_ := ParseSCTE35(bytes)
// for _,d := range(scte.Descriptors()) {
//   for i,o := range(sortedDescs) {
//     if d.Equal(o) { break } // ignore duplicates
//     if d.Compare(o)>0 { continue } // d trumps current
//     if d.Compare(o)==0 && d.CanClose(o) {
//       sortedDescs=sortedDescs[i:] && break
//     } else if d.Compare(o)<0 && d.IsOut() {
//       sortedDescs.InsertAt(d,i) && break
//     }
//   }
// }
type SegmentationDescriptor interface {
	// SCTE35 returns the SCTE35 signal this segmentation descriptor was found in.
	SCTE35() SCTE35
	// TypeID returns the segmentation type for descriptor
	TypeID() SegDescType
	// IsOut returns true if a signal is an out
	IsOut() bool
	// IsIn returns true if a signal is an in
	IsIn() bool
	// HasDuration returns true if there is a duration associated with the descriptor
	HasDuration() bool
	// Duration returns the duration of the descriptor
	Duration() pes.PTS
	// UPIDType returns the type of the upid
	UPIDType() SegUPIDType
	// UPID returns the upid of the descriptor
	UPID() []byte
	// Compare returns results in terms of trumping rules, <0 if sd is less than, 0
	// if equal, and >0 if greater
	Compare(sd SegmentationDescriptor) int
	// CanClose returns true if this descriptor can close the passed in descriptor
	CanClose(out SegmentationDescriptor) bool
	// Equal returns true/false if segmentation descriptor is functionally
	// equal (i.e. a duplicate)
	Equal(sd SegmentationDescriptor) bool
}
