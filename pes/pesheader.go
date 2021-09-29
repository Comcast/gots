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

package pes

import (
	"errors"
	"fmt"

	"github.com/Comcast/gots"
)

// stream_id possibilities
const (
	STREAM_ID_ALL_AUDIO_STREAMS           uint8 = 184
	STREAM_ID_ALL_VIDEO_STREAMS                 = 185
	STREAM_ID_PROGRAM_STREAM_MAP                = 188
	STREAM_ID_PRIVATE_STREAM_1                  = 189
	STREAM_ID_PADDNG_STREAM                     = 190
	STREAM_ID_PRIVATE_STREAM_2                  = 191
	STREAM_ID_ECM_STREAM                        = 240
	STREAM_ID_EMM_STREAM                        = 241
	STREAM_ID_DSM_CC_STREAM                     = 242
	STREAM_ID_ISO_IEC_13552_STREAM              = 243
	STREAM_ID_ITU_T_H222_1_TYPE_A               = 244
	STREAM_ID_ITU_T_H222_1_TYPE_B               = 245
	STREAM_ID_ITU_T_H222_1_TYPE_C               = 246
	STREAM_ID_ITU_T_H222_1_TYPE_D               = 247
	STREAM_ID_ITU_T_H222_1_TYPE_E               = 248
	STREAM_ID_ANCILLARY_STREAM                  = 249
	STREAM_ID_MPEG_4_SL_PACKETIZED_STREAM       = 250
	STREAM_ID_MPEG_4_FLEXMUX_STREAM             = 251
	STREAM_ID_METADATA_STREAM                   = 252
	STREAM_ID_EXTENDED_STREAM_ID                = 253
	STREAM_ID_RESERVED                          = 254
	STREAM_ID_PROGRAM_STREAM_DIRECTORY          = 255
)

// PESHeader represents operations available on a packetized elementary stream header.
type PESHeader interface {
	// HasPTS returns true if the header has a PTS time
	HasPTS() bool
	//PTS return the PTS time in the header
	PTS() uint64
	// HasDTS returns true if the header has a DTS time
	HasDTS() bool
	//DTS return the DTS time in the header
	DTS() uint64
	// Data returns the PES data
	Data() []byte
	// StreamId returns the stream id
	StreamId() uint8
	// DataAligned returns true if the data_alignment_indicator is set
	DataAligned() bool
	// PacketStartCodePrefix returns the packet_start_code_prefix. Note that this is a 24 bit value.
	PacketStartCodePrefix() uint32
}

/*
 * ============================================================
 *  NAME                                              | # BITS
 * ============================================================
 *  PACKET START CODE PREFIX (0X000001)               |   24
 *  STREAM ID                                         |    8
 *       ex: AUDIO STREAMS (0XC0-0XDF)                |
 *       ex: VIDEO STREAMS (0XE0-0XEF)                |
 *  PES PACKET LENGTH                                 |   16
 *      (CAN BE 0; NOT SPECIFIED FOR VIDEO IN MPEGTS) |
 *  OPTIONAL PES HEADER                               |  VAR
 *  STUFFING BYTES                                    |  VAR
 *  DATA                                              |
 *                                                    |
 *  +--- OPTIONAL PES HEADER                          |
 *  | MARKER BITS (BINARY(10) OR HEX(0X2)             |    2
 *  | SCRAMBLING CONTROL (00 IMPLIES NOT SCRAMBLED)   |    2
 *  | PRIORITY                                        |    1
 *  | DATA ALIGNMENT INDICATOR                        |    1
 *  |    1 INDICATES THAT THE PES PACKET HEADER IS    |
 *  |    IMMEDIATELY FOLLOWED BY THE VIDEO START CODE |
 *  |    OR AUDIO SYNCWORD                            |
 *  | COPYRIGHT (1 IMPLIES COPYRIGHTED)               |    1
 *  | ORIGINAL OR COPY (1 IMPLIES ORIGINAL)           |    1
 *  | PTS DTS INDICATOR (11=BOTH, 10=ONLY PTS)        |    2
 *  | ESCR FLAG                                       |    1
 *  | ES RATE FLAG                                    |    1
 *  | DSM TRICK MODE FLAG                             |    1
 *  | ADDITIONAL COPY INFO FLAG                       |    1
 *  | CRC FLAG                                        |    1
 *  | EXTENSION FLAG                                  |    1
 *  | PES HEADER LENGTH (REMAINDER LENGTH OF HEADER)  |    8
 *  |                                                 |
 *  | +--- OPTIONAL FIELDS (IF PTS_DTS_FLAG = 10)     |
 *  | | MARKER BITS '0010'                            |    4
 *  | | PTS 32-30                                     |    3
 *  | | MARKER BIT '1'                                |    1
 *  | | PTS 29-15                                     |   15
 *  | | MARKER BIT '1'                                |    1
 *  | | PTS 14-0                                      |   15
 *  | | MARKER BIT '1'                                |    1
 *  | +--- OPTIONAL FIELDS (IF PTS_DTS_FLAG = 10)     |
 *  |                                                 |
 *  | +--- OPTIONAL FIELDS (IF PTS_DTS_FLAG = 11)     |
 *  | | MARKER BITS '0011'                            |    4
 *  | | PTS 32-30                                     |    3
 *  | | MARKER BIT '1'                                |    1
 *  | | PTS 29-15                                     |   15
 *  | | MARKER BIT '1'                                |    1
 *  | | PTS 14-0                                      |   15
 *  | | MARKER BIT '1'                                |    1
 *  | | MARKER BITS '0001'                            |    4
 *  | | DTS 32-30                                     |    3
 *  | | MARKER BIT '1'                                |    1
 *  | | DTS 29-15                                     |   15
 *  | | MARKER BIT '1'                                |    1
 *  | | DTS 14-0                                      |   15
 *  | | MARKER_BIT '1'                                |    1
 *  | +--- OPTIONAL FIELDS (IF PTS_DTS_FLAG = 10)     |
 *  |                                                 |
 *  | OPTIONAL FIELDS (PRESENCE IS SET BY FLAG ABOVE) |  VAR
 *  | STUFFING BYTES (0xFF)                           |  VAR
 *  +--- OPTIONAL PES HEADER                          |
 * ============================================================
 */
type pESHeader struct {
	packetStartCodePrefix uint32
	dataAlignment         bool
	streamId              uint8
	pesPacketLength       uint16
	ptsDtsIndicator       uint8
	pts                   uint64
	dts                   uint64
	data                  []byte
}

// ExtractTime extracts a PTS time
func ExtractTime(bytes []byte) uint64 {
	var a, b, c, d, e uint64
	a = uint64((bytes[0] >> 1) & 0x07)
	b = uint64(bytes[1])
	c = uint64((bytes[2] >> 1) & 0x7f)
	d = uint64(bytes[3])
	e = uint64((bytes[4] >> 1) & 0x7f)
	return (a << 30) | (b << 22) | (c << 15) | (d << 7) | e
}

// NewPESHeader creates a new PES header with the provided bytes.
// pesBytes is the packet payload that contains the PES data
func NewPESHeader(pesBytes []byte) (PESHeader, error) {
	pes := new(pESHeader)
	var err error

	if CheckLength(pesBytes, "PES", 6) {

		pes.packetStartCodePrefix = uint32(pesBytes[0])<<16 | uint32(pesBytes[1])<<8 | uint32(pesBytes[2])

		pes.streamId = uint8(pesBytes[3])

		pes.pesPacketLength = uint16(pesBytes[4])<<8 | uint16(pesBytes[5])
		pes.dataAlignment = pesBytes[6]&0x04 != 0
		dataStartIndex := 6

		if pes.optionalFieldsExist() && CheckLength(pesBytes, "Optional Fields", 9) {

			ptsDtsIndicator := (uint8(pesBytes[7]) & 0xc0 >> 6)

			pesHeaderLength := pesBytes[8]
			dataStartIndex = 9 + int(pesHeaderLength)

			pes.ptsDtsIndicator = ptsDtsIndicator

			if ptsDtsIndicator != gots.PTS_DTS_INDICATOR_NONE &&
				CheckLength(pesBytes, "PTS", 14) {

				pes.pts = gots.ExtractTime(pesBytes[9:14])

				if pes.ptsDtsIndicator == gots.PTS_DTS_INDICATOR_BOTH &&
					CheckLength(pesBytes, "DTS", 19) {

					pes.dts = gots.ExtractTime(pesBytes[14:19])
				}
			}

		}

		if len(pesBytes) > dataStartIndex {
			pes.data = pesBytes[dataStartIndex:]
		}
	} else {
		err = errors.New("Invalid length for PES header. Too short to parse")
	}

	return pes, err
}

func (pes *pESHeader) optionalFieldsExist() bool {
	if pes.streamId == STREAM_ID_PADDNG_STREAM ||
		pes.streamId == STREAM_ID_PRIVATE_STREAM_2 ||
		pes.streamId == STREAM_ID_ECM_STREAM ||
		pes.streamId == STREAM_ID_EMM_STREAM ||
		pes.streamId == STREAM_ID_DSM_CC_STREAM ||
		pes.streamId == STREAM_ID_ITU_T_H222_1_TYPE_E ||
		pes.streamId == STREAM_ID_PROGRAM_STREAM_DIRECTORY {
		return false
	}
	return true
}

func (pes *pESHeader) PacketStartCodePrefix() uint32 {
	return pes.packetStartCodePrefix
}

func (pes *pESHeader) StreamId() uint8 {
	return pes.streamId
}

func (pes *pESHeader) PTS() uint64 {
	return pes.pts
}

func (pes *pESHeader) DTS() uint64 {
	return pes.dts
}

func (pes *pESHeader) Data() []byte {
	return pes.data
}

func (pes *pESHeader) HasPTS() bool {
	return (pes.ptsDtsIndicator & gots.PTS_DTS_INDICATOR_ONLY_PTS) != 0
}

func (pes *pESHeader) HasDTS() bool {
	return pes.ptsDtsIndicator == gots.PTS_DTS_INDICATOR_BOTH
}

func (pes *pESHeader) Format() string {
	var f = fmt.Sprintf(
		"PES\n"+
			"---\n"+
			"Packet Start Code Prefix: %X \n"+
			"Stream Id: %X \n"+
			"PES Packet Length: %d\n",
		pes.packetStartCodePrefix,
		pes.streamId,
		pes.pesPacketLength)

	if pes.optionalFieldsExist() {
		ptsDtsIndicator := pes.ptsDtsIndicator
		f += fmt.Sprintf("PTS DTS Indicator: %b\n", pes.ptsDtsIndicator)
		if ptsDtsIndicator == gots.PTS_DTS_INDICATOR_BOTH || ptsDtsIndicator == gots.PTS_DTS_INDICATOR_ONLY_PTS {
			f += fmt.Sprintf("PTS: %d\n", pes.pts)
			if ptsDtsIndicator == gots.PTS_DTS_INDICATOR_BOTH {
				f += fmt.Sprintf("DTS: %d\n", pes.dts)
			}
		}

	}

	return f
}

func (pes *pESHeader) DataAligned() bool {
	return pes.dataAlignment
}

// CheckLength the length of the byte array to avoid index out of bound panics
func CheckLength(byteArray []byte, name string, min int) bool {
	if len(byteArray) < min {
		return false
	}
	return true
}
