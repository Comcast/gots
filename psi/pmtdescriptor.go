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

import (
	"fmt"
	"strconv"
)

type pmtDescriptor struct {
	tag  uint8
	data []byte
}

// NewPmtDescriptor creates a new PMTDescriptor with the provided tag and byte contents.
func NewPmtDescriptor(tag uint8, data []byte) PmtDescriptor {
	descriptor := &pmtDescriptor{}
	descriptor.tag = tag
	descriptor.data = data
	return descriptor
}

func (descriptor *pmtDescriptor) Tag() uint8 {
	return descriptor.tag
}

func (descriptor *pmtDescriptor) String() string {
	return descriptor.decode()
}

func (descriptor *pmtDescriptor) Format() string {
	return fmt.Sprintf("[tag=%b, decoded=%s]\n", descriptor.tag, descriptor.decode())
}

func (descriptor *pmtDescriptor) decode() string {
	switch descriptor.tag {
	case LANGUAGE:
		return fmt.Sprintf("ISO 639 Language (code=%s, audioType=%b)",
			descriptor.DecodeIso639LanguageCode(), descriptor.DecodeIso639AudioType())
	case MAXIMUM_BITRATE:
		return fmt.Sprintf("Maximum Bit-Rate (%d)", descriptor.DecodeMaximumBitRate())
	case AUDIO_STREAM:
		return fmt.Sprintf("Audio Stream (%d)", descriptor.tag)
	case REGISTRATION:
		return fmt.Sprintf("Registration (%d)", descriptor.tag)
	case CONDITIONAL_ACCESS:
		return fmt.Sprintf("Conditional Access (%d)", descriptor.tag)
	case SYSTEM_CLOCK:
		return fmt.Sprintf("System Clock (%d)", descriptor.tag)
	case COPYRIGHT:
		return fmt.Sprintf("Copyright (%d)", descriptor.tag)
	case AVC_VIDEO:
		return fmt.Sprintf("AVC Video (%d)", descriptor.tag)
	case DOLBY_DIGITAL:
		return fmt.Sprintf("Dolby Digital (%d)", descriptor.tag)
	case SCTE_ADAPTATION:
		return fmt.Sprintf("SCTE Adaptation (%d)", descriptor.tag)
	case EBP:
		return fmt.Sprintf("EBP (%d)", descriptor.tag)
	}
	return "unknown tag (" + strconv.Itoa(int(descriptor.tag)) + ")"
}

func (descriptor *pmtDescriptor) IsIso639LanguageDescriptor() bool {
	return descriptor.tag == LANGUAGE
}

func (descriptor *pmtDescriptor) IsMaximumBitrateDescriptor() bool {
	return descriptor.tag == MAXIMUM_BITRATE
}

func (descriptor *pmtDescriptor) IsEBPDescriptor() bool {
	return descriptor.tag == EBP
}

// Return the decoded Maximum_bitrate in units of 50 bytes per second
func (descriptor *pmtDescriptor) DecodeMaximumBitRate() uint32 {
	if descriptor.IsMaximumBitrateDescriptor() {
		return uint32(descriptor.data[0]&0x1f)<<16 | uint32(descriptor.data[1])<<8 | uint32(descriptor.data[2])
	}
	return 0
}

func (descriptor *pmtDescriptor) DecodeIso639LanguageCode() string {
	if LANGUAGE == descriptor.tag {
		return string(descriptor.data[0:3])
	}
	return ""
}

func (descriptor *pmtDescriptor) DecodeIso639AudioType() byte {
	return descriptor.data[3]
}

// IsIFrameProfile determines from the PMT if the profile is an I-Frame only track
// or not. An I-Frame only track is defined to be true if and only if the
// 'EBP_distance' is equal to '1'. The 'EBP_distance' is found in the PMT EBP
// descriptor as defined on page 16-17 of OC-SP-EBP-I01-130018.pdf.
// https://www.teamccp.com/confluence/download/attachments/59024185/OC-SP-EBP-I01-130118.pdf?version=1&modificationDate=1378666671000&api=v2
func (descriptor *pmtDescriptor) IsIFrameProfile() bool {
	if EBP == descriptor.tag && 0 < len(descriptor.data) {

		offset := 0
		num_partitions := uint8((descriptor.data[offset] & 0xF8) >> 3)
		timescale_flag := 1 == uint8((descriptor.data[offset]&0x04)>>2)
		offset++

		var EBP_distance_width_minus_1 uint8
		if timescale_flag {
			return false
		}

		indx := uint8(0)
		for indx < num_partitions {
			indx++
			EBP_data_explicit_flag := 1 == uint8((descriptor.data[offset]&0x80)>>7)
			representation_id_flag := 1 == uint8((descriptor.data[offset]&0x04)>>6)

			if EBP_data_explicit_flag {
				offset++
				if 0 == EBP_distance_width_minus_1 {
					EBP_distance := uint8(descriptor.data[offset])
					return 1 == EBP_distance
				} else {
					return false
				}
			} else {
				offset += 2
			}

			if representation_id_flag {
				offset += 8
			}
		}
		return false
	}
	return false
}
