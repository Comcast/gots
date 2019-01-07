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

// Package psi provides mechanisms for collecting and querying program specific information in an MPEG transport stream.
package psi

// Stream type constants
const (
	PmtStreamTypeMpeg2VideoH262 uint8 = 2  // H262
	PmtStreamTypeMpeg4Video     uint8 = 27 // H264
	PmtStreamTypeMpeg4VideoH264 uint8 = 27 // H264
	PmtStreamTypeMpeg4VideoH265 uint8 = 36 // H265

	PmtStreamTypeAac uint8 = 15  // AAC
	PmtStreamTypeAc3 uint8 = 129 // DD
	PmtStreamTypeEc3 uint8 = 135 // DD+

	PmtStreamTypeScte35 uint8 = 134 // SCTE-35
)

// Program Element Stream Descriptor Type.
const (
	VIDEO_STREAM       uint8 = 2   // 0000 0010 (0x02)
	AUDIO_STREAM       uint8 = 3   // 0000 0011 (0x03)
	REGISTRATION       uint8 = 5   // 0000 1000 (0x05)
	CONDITIONAL_ACCESS uint8 = 9   // 0000 1001 (0x09)
	LANGUAGE           uint8 = 10  // 0000 1010 (0x0A)
	SYSTEM_CLOCK       uint8 = 11  // 0000 1011 (0x0B)
	DOLBY_DIGITAL      uint8 = 12  // 0000 1100 (0x0C)
	COPYRIGHT          uint8 = 13  // 0000 1101 (0x0D)
	MAXIMUM_BITRATE    uint8 = 14  // 0000 1110 (0x0E)
	AVC_VIDEO          uint8 = 40  // 0010 1000 (0x28)
	STREAM_IDENTIFIER  uint8 = 82  // 0101 0010 (0x52)
	SCTE_ADAPTATION    uint8 = 151 // 1001 0111 (0x97)
	EBP                uint8 = 233 // 1110 1001 (0xE9)
	EC3                uint8 = 204 // 1100 1100 (0xCC)
)

// Unaccounted bytes before the end of the SectionLength field
const (
	// Pointerfield(1) + table id(1) + flags(.5) + section length (2.5)
	PSIHeaderLen uint16 = 4
	CrcLen       uint16 = 4
)

// TableHeader struct represents operations available on all PSI
type TableHeader struct {
	TableID                uint8
	SectionSyntaxIndicator bool
	PrivateIndicator       bool
	SectionLength          uint16
}

// PmtStreamType is used to represent elementary steam type inside a PMT
type PmtStreamType interface {
	StreamType() uint8
	StreamTypeDescription() string
	IsStreamWherePresentationLagsEbp() bool
	IsAudioContent() bool
	IsVideoContent() bool
	IsSCTE35Content() bool
}

// PAT interface represents operations on a Program Association Table. Currently only single program transport streams (SPTS)are supported
type PAT interface {
	NumPrograms() int
	ProgramMap() map[uint16]uint16
	SPTSpmtPID() (uint16, error)
}

// PMT is a Program Map Table.
type PMT interface {
	Pids() []uint16
	IsPidForStreamWherePresentationLagsEbp(pid uint16) bool
	ElementaryStreams() []PmtElementaryStream
	RemoveElementaryStreams(pids []uint16)
	String() string
}

// PmtElementaryStream represents an elementary stream inside a PMT
type PmtElementaryStream interface {
	PmtStreamType
	ElementaryPid() uint16
	Descriptors() []PmtDescriptor
	MaxBitRate() uint64
}

// PmtDescriptor represents operations currently necessary on descriptors found in the PMT
type PmtDescriptor interface {
	Tag() uint8
	Format() string
	IsIso639LanguageDescriptor() bool
	IsMaximumBitrateDescriptor() bool
	IsIFrameProfile() bool
	IsEBPDescriptor() bool
	DecodeMaximumBitRate() uint32
	DecodeIso639LanguageCode() string
	DecodeIso639AudioType() byte
	IsDolbyATMOS() bool
}
