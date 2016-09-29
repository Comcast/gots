package adaptationfield

import "github.com/Comcast/gots/packet"

var emptyByteSlice []byte

// Length returns the length of the adaptation field in bytes
func Length(packet packet.Packet) uint8 {
	return uint8(packet[4])
}

// IsDiscontinuous returns the discontinuity indicator for this adaptation field
func IsDiscontinuous(packet packet.Packet) bool {
	return (packet[5] & 0x80) != 0
}

// IsRandomAccess returns the random access indicator for this adaptation field
func IsRandomAccess(packet packet.Packet) bool {
	return (packet[5] & 0x40) != 0
}

// IsESHigherPriority returns true if this elementary stream is
// high priority. Corresponds to the elementary stream
// priority indicator.
func IsESHigherPriority(packet packet.Packet) bool {
	return (packet[5] & 0x20) != 0
}

// HasPCR returns true when the PCR flag is set
func HasPCR(packet packet.Packet) bool {
	return (packet[5] & 0x10) != 0
}

// HasOPCR returns true when the OPCR flag is set
func HasOPCR(packet packet.Packet) bool {
	return (packet[5] & 0x08) != 0
}

// HasSplicingPoint returns true when the splicing countdown field is present
func HasSplicingPoint(packet packet.Packet) bool {
	return (packet[5] & 0x04) != 0
}

// HasTransportPrivateData returns true when the private data field is present
func HasTransportPrivateData(packet packet.Packet) bool {
	return (packet[5] & 0x02) != 0
}

// HasAdaptationFieldExtension returns true if this adaptation field contains an extension field
func HasAdaptationFieldExtension(packet packet.Packet) bool {
	return (packet[5] & 0x01) != 0
}

// EncoderBoundaryPoint returns the byte array located in the optional TransportPrivateData of the (also optional)
// AdaptationField of the Packet. If either of these optional fields are missing an empty byte array is returned with an error
func EncoderBoundaryPoint(pkt packet.Packet) ([]byte, error) {
	if badLen(pkt) {
		return emptyByteSlice, mpegts.ErrInvalidPacketLength
	}
	hasAdapt, err := packet.ContainsAdaptationField(pkt)
	if err != nil {
		return emptyByteSlice, nil
	}
	if hasAdapt && Length(pkt) > 0 && HasTransportPrivateData(pkt) {
		ebp, err := TransportPrivateData(pkt)
		if err != nil {
			return nil, err
		}
		return ebp, nil
	}
	return nil, mpegts.ErrNoEBP
}

// PCR is the Program Clock Reference.
// First 33 bits are PCR base.
// Next 6 bits are reserved.
// Final 9 bits are PCR extension.
func PCR(packet packet.Packet) ([]byte, error) {
	if !HasPCR(packet) {
		return emptyByteSlice, mpegts.ErrNoPCR
	}
	offset := 6
	return packet[offset : offset+6], nil
}

// OPCR is the Original Program Clock Reference.
// First 33 bits are original PCR base.
// Next 6 bits are reserved.
// Final 9 bits are original PCR extension.
func OPCR(packet packet.Packet) ([]byte, error) {
	if badLen(packet) {
		return emptyByteSlice, mpegts.ErrInvalidPacketLength
	}
	if !HasOPCR(packet) {
		return emptyByteSlice, mpegts.ErrNoOPCR
	}
	offset := 6
	if HasPCR(packet) {
		offset += 6
	}
	return packet[offset : offset+6], nil
}

// SpliceCountdown returns a count of how many packets after this one until
// a splice point occurs or an error if none exist. This function calls
// HasSplicingPoint to check for the existence of a splice countdown.
func SpliceCountdown(packet packet.Packet) (uint8, error) {
	if badLen(packet) {
		return 0, mpegts.ErrInvalidPacketLength
	}
	if !HasSplicingPoint(packet) {
		return 0, mpegts.ErrNoSplicePoint
	}
	offset := 6
	if HasPCR(packet) {
		offset += 6
	}
	if HasOPCR(packet) {
		offset += 6
	}
	return packet[offset], nil
}

// TransportPrivateData returns the private data from this adaptation field
// or an empty array and an error if there is none. This function calls
// HasTransportPrivateData to check for the existence of private data.
func TransportPrivateData(packet packet.Packet) ([]byte, error) {
	if badLen(packet) {
		return emptyByteSlice, mpegts.ErrInvalidPacketLength
	}
	if !HasTransportPrivateData(packet) {
		return emptyByteSlice, mpegts.ErrNoPrivateTransportData
	}
	offset := 6
	if HasPCR(packet) {
		offset += 6
	}
	if HasOPCR(packet) {
		offset += 6
	}
	if HasSplicingPoint(packet) {
		offset++
	}
	dataLength := uint8(packet[offset])
	offset++
	return packet[uint8(offset) : uint8(offset)+dataLength], nil
}

// badLen returns true if the packet has invalid length
func badLen(pkt packet.Packet) bool {
	if len(pkt) != packet.PacketSize {
		return true
	}
	return false
}
