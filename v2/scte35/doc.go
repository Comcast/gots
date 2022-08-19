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
	"github.com/Comcast/gots/v2"
	"github.com/Comcast/gots/v2/psi"
)

// SpliceCommandType is a type used to describe the types of splice commands.
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

var SpliceCommandTypeNames = map[SpliceCommandType]string{
	SpliceNull:           "SpliceNull",
	SpliceSchedule:       "SpliceSchedule",
	SpliceInsert:         "SpliceInsert",
	TimeSignal:           "TimeSignal",
	BandwidthReservation: "BandwidthReservation",
	PrivateCommand:       "PrivateCommand",
}

// DeviceRestrictions type is used to specify what group that the segment is restricted to.
type DeviceRestrictions byte

const (
	RestrictGroup0 DeviceRestrictions = 0x00
	RestrictGroup1 DeviceRestrictions = 0x01
	RestrictGroup2 DeviceRestrictions = 0x02
	RestrictNone   DeviceRestrictions = 0x03 // no restrictions
)

var DeviceRestrictionsNames = map[DeviceRestrictions]string{
	RestrictGroup0: "RestrictGroup0",
	RestrictGroup1: "RestrictGroup1",
	RestrictGroup2: "RestrictGroup2",
	RestrictNone:   "RestrictNone",
}

// SegDescType is the Segmentation Descriptor Type - not really needed for processing according
// to method below, but included here for backwards compatibility/porting
type SegDescType uint8

const (
	SegDescNotIndicated                     SegDescType = 0x00
	SegDescContentIdentification                        = 0x01
	SegDescProgramStart                                 = 0x10
	SegDescProgramEnd                                   = 0x11
	SegDescProgramEarlyTermination                      = 0x12
	SegDescProgramBreakaway                             = 0x13
	SegDescProgramResumption                            = 0x14
	SegDescProgramRunoverPlanned                        = 0x15
	SegDescProgramRunoverUnplanned                      = 0x16
	SegDescProgramOverlapStart                          = 0x17
	SegDescProgramBlackoutOverride                      = 0x18
	SegDescProgramStartInProgress                       = 0x19
	SegDescChapterStart                                 = 0x20
	SegDescChapterEnd                                   = 0x21
	SegDescBreakStart                                   = 0x22
	SegDescBreakEnd                                     = 0x23
	SegDescOpeningCreditStart                           = 0x24
	SegDescOpeningCreditEnd                             = 0x25
	SegDescClosingCreditStart                           = 0x26
	SegDescClosingCreditEnd                             = 0x27
	SegDescProviderAdvertisementStart                   = 0x30
	SegDescProviderAdvertisementEnd                     = 0x31
	SegDescDistributorAdvertisementStart                = 0x32
	SegDescDistributorAdvertisementEnd                  = 0x33
	SegDescProviderPOStart                              = 0x34
	SegDescProviderPOEnd                                = 0x35
	SegDescDistributorPOStart                           = 0x36
	SegDescDistributorPOEnd                             = 0x37
	SegDescProviderPromoStart                           = 0x3C
	SegDescProviderPromoEnd                             = 0x3D
	SegDescUnscheduledEventStart                        = 0x40
	SegDescUnscheduledEventEnd                          = 0x41
	SegDescAlternateContentOpportunityStart             = 0x42
	SegDescAlternateContentOpportunityEnd               = 0x43
	SegDescProviderAdBlockStart                         = 0x44
	SegDescProviderAdBlockEnd                           = 0x45
	SegDescNetworkStart                                 = 0x50
	SegDescNetworkEnd                                   = 0x51
)

var SegDescTypeNames = map[SegDescType]string{
	SegDescNotIndicated:                     "SegDescNotIndicated",
	SegDescContentIdentification:            "SegDescContentIdentification",
	SegDescProgramStart:                     "SegDescProgramStart",
	SegDescProgramEnd:                       "SegDescProgramEnd",
	SegDescProgramEarlyTermination:          "SegDescProgramEarlyTermination",
	SegDescProgramBreakaway:                 "SegDescProgramBreakaway",
	SegDescProgramResumption:                "SegDescProgramResumption",
	SegDescProgramRunoverPlanned:            "SegDescProgramRunoverPlanned",
	SegDescProgramRunoverUnplanned:          "SegDescProgramRunoverUnplanned",
	SegDescProgramOverlapStart:              "SegDescProgramOverlapStart",
	SegDescProgramBlackoutOverride:          "SegDescProgramBlackoutOverride",
	SegDescProgramStartInProgress:           "SegDescProgramStartInProgress",
	SegDescChapterStart:                     "SegDescChapterStart",
	SegDescChapterEnd:                       "SegDescChapterEnd",
	SegDescBreakStart:                       "SegDescBreakStart",
	SegDescBreakEnd:                         "SegDescBreakEnd",
	SegDescProviderAdvertisementStart:       "SegDescProviderAdvertisementStart",
	SegDescProviderAdvertisementEnd:         "SegDescProviderAdvertisementEnd",
	SegDescOpeningCreditStart:               "SegDescOpeningCreditStart",
	SegDescOpeningCreditEnd:                 "SegDescOpeningCreditEnd",
	SegDescClosingCreditStart:               "SegDescClosingCreditStart",
	SegDescClosingCreditEnd:                 "SegDescClosingCreditEnd",
	SegDescDistributorAdvertisementStart:    "SegDescDistributorAdvertisementStart",
	SegDescDistributorAdvertisementEnd:      "SegDescDistributorAdvertisementEnd",
	SegDescProviderPOStart:                  "SegDescProviderPOStart",
	SegDescProviderPOEnd:                    "SegDescProviderPOEnd",
	SegDescDistributorPOStart:               "SegDescDistributorPOStart",
	SegDescDistributorPOEnd:                 "SegDescDistributorPOEnd",
	SegDescProviderPromoStart:               "SegDescProviderPromoStart",
	SegDescProviderPromoEnd:                 "SegDescProviderPromoEnd",
	SegDescUnscheduledEventStart:            "SegDescUnscheduledEventStart",
	SegDescUnscheduledEventEnd:              "SegDescUnscheduledEventEnd",
	SegDescAlternateContentOpportunityStart: "SegDescAlternateContentOpportunityStart",
	SegDescAlternateContentOpportunityEnd:   "SegDescAlternateContentOpportunityEnd",
	SegDescProviderAdBlockStart:             "SegDescProviderAdBlockStart",
	SegDescProviderAdBlockEnd:               "SegDescProviderAdBlockEnd",
	SegDescNetworkStart:                     "SegDescNetworkStart",
	SegDescNetworkEnd:                       "SegDescNetworkEnd",
}

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
	SegUPADSINFO                   = 0x0e
	SegUPIDURN                     = 0x0f
)

var SegUPIDTypeNames = map[SegUPIDType]string{
	SegUPIDNotUsed:     "SegUPIDNotUsed",
	SegUPIDUserDefined: "SegUPIDUserDefined",
	SegUPIDISCI:        "SegUPIDISCI",
	SegUPIDAdID:        "SegUPIDAdID",
	SegUPIDUMID:        "SegUPIDUMID",
	SegUPIDISAN:        "SegUPIDISAN",
	SegUPIDVISAN:       "SegUPIDVISAN",
	SegUPIDTID:         "SegUPIDTID",
	SegUPIDTI:          "SegUPIDTI",
	SegUPIDADI:         "SegUPIDADI",
	SegUPIDEIDR:        "SegUPIDEIDR",
	SegUPIDATSCID:      "SegUPIDATSCID",
	SegUPIDMPU:         "SegUPIDMPU",
	SegUPIDMID:         "SegUPIDMID",
	SegUPADSINFO:       "SegUPADSINFO",
	SegUPIDURN:         "SegUPIDURN",
}

// SCTE35 represent operations available on a SCTE35 message.
type SCTE35 interface {
	// HasPTS returns true if there is a pts time.
	HasPTS() bool
	// SetHasPTS sets if this SCTE35 message has a PTS.
	SetHasPTS(flag bool)
	// PTS returns the PTS time of the signal if it exists. Includes adjustment.
	PTS() gots.PTS
	// SetPTS sets the PTS time of the signal's command. There will be no PTS adjustment using this function.
	// If HasPTS is false, then it will have no effect until it is set to true. Also this command has no
	// effect with a null splice command.
	SetPTS(pts gots.PTS)
	// SetAdjustPTS will modify the pts adjustment field. The desired PTS value
	// after adjustment should be passed, The adjustment value will be calculated
	// during the call to Data().
	SetAdjustPTS(pts gots.PTS)
	// Tier returns which authorization tier this message was assigned to.
	// The tier value of 0XFFF is the default and will ignored. When using tier values,
	// the SCTE35 message must fit entirely into the ts payload without being split up.
	// The tier is a 12 bit unsigned integer.
	Tier() uint16
	// SetTier sets which authorization tier this message was assigned to.
	// The tier value of 0XFFF is the default and will ignored. When using tiers,
	// the SCTE35 message must fit entirely into the ts payload without being split up.
	// The tier is a 12 bit unsigned integer.
	SetTier(tier uint16)
	// Command returns the signal's splice command.
	Command() SpliceCommandType
	// CommandInfo returns an object describing fields of the signal's splice
	// command structure.
	CommandInfo() SpliceCommand
	// SetCommandInfo sets the object describing fields of the signal's splice
	// command structure
	SetCommandInfo(commandInfo SpliceCommand)
	// Descriptors returns a slice of the signals SegmentationDescriptors sorted
	// by descriptor weight (least important signals first)
	Descriptors() []SegmentationDescriptor
	// SetDescriptors sets a slice of the signals SegmentationDescriptors they
	// should be sorted by descriptor weight (least important signals first)
	SetDescriptors(descriptors []SegmentationDescriptor)
	// AlignmentStuffing returns how many stuffing bytes will be added to the SCTE35
	// message at the end.
	AlignmentStuffing() uint
	// SetAlignmentStuffing sets how many stuffing bytes will be added to the SCTE35
	// message at the end.
	SetAlignmentStuffing(alignmentStuffing uint)
	// UpdateData will encode the SCTE35 information to bytes and return it.
	// UpdateData will make the next call to Data() return these new bytes.
	UpdateData() []byte
	// Data returns the raw data bytes of the scte signal. It will not change
	// unless a call to UpdateData() is issued before this.
	Data() []byte
	// String returns a string representation of the SCTE35 message.
	// String function is for debugging and testing.
	String() string
}

// SpliceCommand represent operations available on all SpliceCommands.
type SpliceCommand interface {
	// CommandType returns the signal's splice command type value.
	CommandType() SpliceCommandType
	// HasPTS returns true if there is a pts time on the command.
	HasPTS() bool
	// PTS returns the PTS time of the command, not including adjustment.
	PTS() gots.PTS
	// SetHasPTS sets the flag that indicates if there is a pts time on the command.
	SetHasPTS(value bool)
	// SetPTS sets the PTS.
	SetPTS(value gots.PTS)
	// Data returns the bytes of this splice command.
	Data() []byte
}

// TimeSignalCommand is a type of SpliceCommand.
type TimeSignalCommand interface {
	SpliceCommand
}

// Component is an interface for components, a structure in SpliceInsertCommand.
type Component interface {
	// ComponentTag returns the tag of the component.
	ComponentTag() byte
	// HasPTS returns if the component has a PTS.
	HasPTS() bool
	// PTS returns the PTS of the component.
	PTS() gots.PTS
	// SetComponentTag sets the component tag, which is used for the identification of the component.
	SetComponentTag(value byte)
	// SetHasPTS sets a flag that determines if the component has a PTS.
	SetHasPTS(value bool)
	// SetPTS sets the PTS of the component.
	SetPTS(value gots.PTS)
}

// ComponentOffset is an interface for componentOffset, a structure in SegmentationDescriptor.
type ComponentOffset interface {
	// ComponentTag returns the tag of the component.
	ComponentTag() byte
	// PTSOffset returns the PTS offset of the component.
	PTSOffset() gots.PTS
	// SetComponentTag sets the component tag, which is used for the identification of the component.
	SetComponentTag(value byte)
	// SetPTSOffset sets the PTS offset of the component.
	SetPTSOffset(value gots.PTS)
}

// SpliceInsertCommand is a type of SpliceCommand.
type SpliceInsertCommand interface {
	SpliceCommand
	// EventID returns the event id
	EventID() uint32
	// SetEventID sets the event id
	SetEventID(value uint32)
	// IsEventCanceled returns the event cancel indicator
	IsEventCanceled() bool
	// SetIsEventCanceled sets the the event cancel indicator
	SetIsEventCanceled(value bool)
	// IsOut returns the value of the out of network indicator
	IsOut() bool
	// SetIsOut sets the out of network indicator
	SetIsOut(value bool)
	// IsProgramSplice returns if the program_splice_flag is set
	IsProgramSplice() bool
	// SetIsProgramSplice sets the program splice flag
	SetIsProgramSplice(value bool)
	// HasDuration returns true if there is a duration
	HasDuration() bool
	// SetHasDuration sets the duration flag
	SetHasDuration(value bool)
	// SpliceImmediate returns if the splice_immediate_flag is set
	SpliceImmediate() bool
	// SetSpliceImmediate sets the splice immediate flag
	SetSpliceImmediate(value bool)
	// Components returns the components of the splice command
	Components() []Component
	// IsAutoReturn returns the boolean value of the auto return field
	IsAutoReturn() bool
	// SetIsAutoReturn sets the auto_return flag
	SetIsAutoReturn(value bool)
	// Duration returns the PTS duration of the command
	Duration() gots.PTS
	// SetDuration sets the PTS duration of the command
	SetDuration(value gots.PTS)
	// UniqueProgramId returns the unique_program_id field
	UniqueProgramId() uint16
	// SetUniqueProgramId sets the unique program Id
	SetUniqueProgramId(value uint16)
	// AvailNum returns the avail_num field, index of this avail or zero if unused
	AvailNum() uint8
	// SetAvailNum sets the avail_num field, zero if unused. otherwise index of the avail
	SetAvailNum(value uint8)
	// AvailsExpected returns avails_expected field, number of avails for program
	AvailsExpected() uint8
	// SetAvailsExpected sets the expected number of avails
	SetAvailsExpected(value uint8)
}

// UPID describes the UPID, this is only used for MID.
type UPID interface {
	// UPIDType returns the type of UPID stored
	UPIDType() SegUPIDType
	// UPID returns the actual UPID
	UPID() []byte
	// SetUPIDType will set the type of the UPID
	SetUPIDType(value SegUPIDType)
	// SetUPID set the actual UPID
	SetUPID(value []byte)
}

// SegmentationDescriptor describes the segmentation descriptor interface.
type SegmentationDescriptor interface {
	// SCTE35 returns the SCTE35 signal this segmentation descriptor was found in.
	SCTE35() SCTE35
	// EventID returns the event id
	EventID() uint32
	// SetEventID sets the event id
	SetEventID(value uint32)
	// IsEventCanceled returns the event cancel indicator
	IsEventCanceled() bool
	// SetIsEventCanceled sets the the event cancel indicator
	SetIsEventCanceled(value bool)
	// HasProgramSegmentation returns if the descriptor has program segmentation
	HasProgramSegmentation() bool
	// SetHasProgramSegmentation if the descriptor has program segmentation
	SetHasProgramSegmentation(value bool)
	// HasDuration returns true if there is a duration associated with the descriptor
	HasDuration() bool
	// SetHasDuration determines if a duration is present in the descriptor
	SetHasDuration(value bool)
	// Duration returns the duration of the descriptor, 40 bit unsigned integer.
	Duration() gots.PTS
	// SetDuration sets the duration of the descriptor, 40 bit unsigned integer. extra bits will be truncated.
	SetDuration(value gots.PTS)
	// IsDeliveryNotRestricted returns if the delivery is not restricted
	IsDeliveryNotRestricted() bool
	// SetIsDeliveryNotRestricted sets a flag that determines if the delivery is not restricted
	SetIsDeliveryNotRestricted(bool)
	// IsWebDeliveryAllowed returns if web delivery is allowed, this field has no meaning if delivery is not restricted.
	IsWebDeliveryAllowed() bool
	// SetIsWebDeliveryAllowed sets a flag that determines if web delivery is allowed, this field has no meaning if delivery is not restricted.
	SetIsWebDeliveryAllowed(bool)
	// HasNoRegionalBlackout returns true if there is no regional blackout, this field has no meaning if delivery is not restricted.
	HasNoRegionalBlackout() bool
	// SetHasNoRegionalBlackout sets a flag that determines if there is no regional blackout, this field has no meaning if delivery is not restricted.
	SetHasNoRegionalBlackout(bool)
	// IsArchiveAllowed returns true if there are restrictions to storing/recording this segment, this field has no meaning if delivery is not restricted.
	IsArchiveAllowed() bool
	// SetIsArchiveAllowed sets a flag that determines if there are restrictions to storing/recording this segment, this field has no meaning if delivery is not restricted.
	SetIsArchiveAllowed(bool)
	// DeviceRestrictions returns which device group the segment is restriced to, this field has no meaning if delivery is not restricted.
	DeviceRestrictions() DeviceRestrictions
	// SetDeviceRestrictions sets which device group the segment is restriced to, this field has no meaning if delivery is not restricted.
	SetDeviceRestrictions(DeviceRestrictions)
	// Components will return components' offsets
	Components() []ComponentOffset
	// SetComponents will set components' offsets
	SetComponents([]ComponentOffset)
	// UPIDType returns the type of the upid
	UPIDType() SegUPIDType
	// SetUPIDType sets the type of upid, only works if UPIDType is not SegUPIDMID
	SetUPIDType(value SegUPIDType)
	// UPID returns the upid of the descriptor, if the UPIDType is not SegUPIDMID
	UPID() []byte
	// SetUPID sets the upid of the descriptor
	SetUPID(value []byte)
	// MID returns multiple UPIDs, if UPIDType is SegUPIDMID
	MID() []UPID
	// SetMID sets multiple UPIDs, only works if UPIDType is SegUPIDMID
	SetMID(value []UPID)
	// TypeID returns the segmentation type for descriptor
	TypeID() SegDescType
	// SetTypeID sets the type id
	SetTypeID(value SegDescType)
	// SegmentNumber returns the segment number for this descriptor.
	SegmentNumber() uint8
	// SetSegmentNumber sets the segment number for this descriptor.
	SetSegmentNumber(value uint8)
	// SegmentsExpected returns the number of expected segments for this descriptor.
	SegmentsExpected() uint8
	// SetSegmentsExpected sets the number of expected segments for this descriptor.
	SetSegmentsExpected(value uint8)
	// HasSubSegments returns true if this segmentation descriptor has subsegment fields.
	HasSubSegments() bool
	// SetHasSubSegments sets the field that determines if this segmentation descriptor has sub subsegments.
	SetHasSubSegments(bool)
	// SubSegmentNumber returns the sub-segment number for this descriptor.
	SubSegmentNumber() uint8
	// SetSubSegmentNumber sets the sub-segment number for this descriptor.
	SetSubSegmentNumber(value uint8)
	// SubSegmentsExpected returns the number of expected sub-segments for this descriptor.
	SubSegmentsExpected() uint8
	// SetSubSegmentsExpected sets the number of expected sub-segments for this descriptor.
	SetSubSegmentsExpected(value uint8)
	// StreamSwitchSignalID returns the signalID of streamswitch signal if
	// present in the descriptor
	StreamSwitchSignalId() (string, error)
	// IsOut returns true if a signal is an out
	IsOut() bool
	// IsIn returns true if a signal is an in
	IsIn() bool
	// CanClose returns true if this descriptor can close the passed in descriptor
	CanClose(out SegmentationDescriptor) bool
	// Equal returns true/false if segmentation descriptor is functionally
	// equal (i.e. a duplicate)
	Equal(sd SegmentationDescriptor) bool
	// Data returns the raw data bytes of the SegmentationDescriptor
	Data() []byte
	// SegmentNum is deprecated, use SegmentNumber instead.
	// SegmentNum returns the segment_num descriptor field.
	SegmentNum() uint8
}

// State maintains current state for all signals and descriptors.  The intended
// usage is to call ParseSCTE35() on raw data to create a signal, and then call
// ProcessSignal with that signal.  This returns the list of descriptors closed
// by that signal. If signals have a duration and need to be closed implicitly
// after some timer has passed, then Close() can be used for that.  Some
// example code is below.
// s := scte35.NewState()
// scte,_ := scte.ParseSCTE35(bytes)
// for _,d := range(scte.Descriptors()) {
//   closed = s.ProcessDescriptor(d)
//   ...handle closed signals appropriately here
//   if d.HasDuration() {
//     time.AfterFunc(d.Duration() + someFudgeDelta,
//                    func() { closed = s.Close(d) })
//   }
// }
type State interface {
	// Open returns a list of open signals
	Open() []SegmentationDescriptor
	// ProcessDescriptor takes a scte35 descriptor and returns a list of descriptors closed by it
	ProcessDescriptor(desc SegmentationDescriptor) ([]SegmentationDescriptor, error)
	// Close acts like Process and acts as if an appropriate close has been
	// received for this given descriptor.
	Close(desc SegmentationDescriptor) ([]SegmentationDescriptor, error)
}

// SCTE done func is the same as the PMT because they're both psi
func SCTE35AccumulatorDoneFunc(b []byte) (bool, error) {
	return psi.PmtAccumulatorDoneFunc(b)
}
