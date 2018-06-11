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

// Package pes contains interfaces and operations for packetized elementary stream headers.
package pes

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
