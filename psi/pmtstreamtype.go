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

package psi

import "fmt"

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

// PmtStreamType is used to represent elementary steam type inside a PMT
type PmtStreamType interface {
	StreamType() uint8
	StreamTypeDescription() string
	IsStreamWherePresentationLagsEbp() bool
	IsAudioContent() bool
	IsVideoContent() bool
	IsSCTE35Content() bool
}
type pmtStreamType struct {
	code                uint8
	description         string
	presentationLagsEbp bool
}

func (st pmtStreamType) StreamType() uint8 {
	return st.code
}

func (st pmtStreamType) StreamTypeDescription() string {
	return st.description
}

func (st pmtStreamType) IsStreamWherePresentationLagsEbp() bool {
	return st.presentationLagsEbp
}

func (st pmtStreamType) IsAudioContent() bool {
	return st.code == PmtStreamTypeAac ||
		st.code == PmtStreamTypeAc3 ||
		st.code == PmtStreamTypeEc3
}

func (st pmtStreamType) IsVideoContent() bool {
	return st.code == PmtStreamTypeMpeg4VideoH264 ||
		st.code == PmtStreamTypeMpeg4VideoH265 ||
		st.code == PmtStreamTypeMpeg4Video ||
		st.code == PmtStreamTypeMpeg2VideoH262
}

func (st pmtStreamType) IsSCTE35Content() bool {
	return st.code == PmtStreamTypeScte35
}

func (st pmtStreamType) String() string {
	return fmt.Sprintf("streamType=%d", st.code)
}

// LookupPmtStreamType returns the associated PmtStreamType of the provided code. If the code is not recognized, a PmtSteamType of "unknown" is returned.
func LookupPmtStreamType(code uint8) PmtStreamType {

	for _, t := range atscPmtStreamTypes {

		if code >= t.firstCode && code <= t.lastCode {

			return *newPmtStreamType(code, t.description, presentationLagsEbp(code))
		}
	}

	return *newPmtStreamType(code, "unknown", presentationLagsEbp(code))
}

func newPmtStreamType(code uint8, description string, presentationLagsEbp bool) *pmtStreamType {
	return &pmtStreamType{code, description, presentationLagsEbp}
}

func presentationLagsEbp(code uint8) bool {
	switch code {
	case 3, 4, 15, 17, 129, 135, 136:
		return true
	}
	return false
}

type atscPmtStreamType struct {
	firstCode   uint8
	lastCode    uint8
	description string
}

// Code/Descriptions transcribed from ATSC Code Point Registry at
// http://www.atsc.org/cms/index.php/standards/other-technical-documents/78-atsc-code-point-registry
// As of 4/8/2014
var atscPmtStreamTypes = []atscPmtStreamType{
	{0, 0, "ITU-T | ISO/IEC Reserved"},
	{1, 1, "ISO/IEC 11172 Video	"},
	{2, 2, "ITU-T Rec. H.262 | ISO/IEC 13818-2 Video"},
	{3, 3, "ISO/IEC 11172 Audio"},
	{4, 4, "ISO/IEC 13818-3 Audio"},
	{5, 5, "ITU-T Rec. H.222.0 | ISO/IEC 13818-1 private sections"},
	{6, 6, "ITU-T Rec. H.222.0 | ISO/IEC 13818-1 PES packets containing private data"},
	{7, 7, "ISO/IEC 13522 MHEG"},
	{8, 8, "ITU-T Rec. H.222.0 | ISO/IEC 13818-1 DSM-CC"},
	{9, 9, "ITU-T Rec. H.222.0 | ISO/IEC 13818-1/11172-1 auxiliary"},
	{10, 10, "ISO/IEC 13818-6 Multi-protocol Encapsulation"},
	{11, 11, "ISO/IEC 13818-6 DSM-CC U-N Messages"},
	{12, 12, "ISO/IEC 13818-6 Stream Descriptors"},
	{13, 13, "ISO/IEC 13818-6 Sections (any type, including private data)"},
	{14, 14, "ISO/IEC 13818-1 auxiliary"},
	{15, 15, "ISO/IEC 13818-7 Audio (AAC) with ADTS transport"},
	{16, 16, "ISO/IEC 14496-2 Visual"},
	{17, 17, "ISO/IEC 14496-3 Audio with the LATM transport syntax as defined in ISO/IEC 14496-3"},
	{18, 18, "ISO/IEC 14496-1 SL-packetized stream or FlexMux stream carried in PES packets"},
	{19, 19, "ISO/IEC 14496-1 SL-packetized stream or FlexMux stream carried in ISO/IEC 14496_sections"},
	{20, 20, "ISO/IEC 13818-6 DSM-CC Synchronized Download Protocol"},
	{21, 21, "Metadata carried in PES packets"},
	{22, 22, "Metadata carried in metadata_sections	"},
	{23, 23, "Metadata carried in ISO/IEC 13818-6 Data Carousel"},
	{24, 24, "Metadata carried in ISO/IEC 13818-6 Object Carousel"},
	{25, 25, "Metadata carried in ISO/IEC 13818-6 Synchronized Download Protocol"},
	{26, 26, "IPMP stream (defined in ISO/IEC 13818-11, MPEG-2 IPMP)"},
	{27, 27, "AVC video stream as defined in ITU-T Rec. H.264 | ISO/IEC 14496-10 Video"},
	{28, 127, "ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Reserved"},
	{36, 36, "HEVC video stream as defined in ITU-T Rec. H.265 | ISO/IEC 23008-2 Video"},
	{128, 128, "DigiCipher® II video | Identical to ITU-T Rec. H.262 | ISO/IEC 13818-2 Video"},
	{129, 129, "ATSC A/53 audio [2] | AC-3 audio"},
	{130, 130, "SCTE Standard Subtitle"},
	{131, 131, "SCTE Isochronous Data | Reserved"},
	{132, 132, "ATSC/SCTE reserved"},
	{133, 133, "ATSC Program Identifier , SCTE Reserved"},
	{134, 134, "SCTE 35 splice_information_table | [Cueing]"},
	{135, 135, "E-AC-3"},
	//{	135,	159	, "	SCTE Reserved	", false },
	{136, 136, "DTS HD Audio"},
	{137, 137, "ATSC Reserved"},
	{138, 143, "ATSC Reserved"},
	{144, 144, "DVB stream_type value for Time Slicing / MPE-FEC"},
	{145, 145, "IETF Unidirectional Link Encapsulation (ULE)"},
	{146, 148, "ATSC Reserved"},
	{149, 149, "ATSC Data Service Table, Network Resources Table"},
	{150, 159, "ATSC Reserved"},
	{160, 160, "SCTE [IP Data] | ATSC Reserved"},
	{161, 191, "ATSC Reserved"},
	{192, 192, "DCII (DigiCipher®) Text"},
	{193, 193, "ATSC Reserved"},
	{194, 194, "ATSC synchronous data stream | [Isochronous Data]"},
	{195, 195, "SCTE Asynchronous Data"},
	{196, 233, "ATSC User Private Program Elements"},
	{234, 234, "VC-1 Elementary Stream per RP227"},
	{235, 255, "ATSC User Private Program Elements"},
}
