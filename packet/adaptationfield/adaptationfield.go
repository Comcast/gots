//package adaptationfield

// Length returns the length of the adaptation field in bytes
// func Length(pkt *packet.Packet) uint8 {
// 	return uint8(pkt[4])
// }

// IsDiscontinuous returns the discontinuity indicator for this adaptation field
// func IsDiscontinuous(pkt *packet.Packet) bool {
// 	return pkt[5]&0x80 != 0
// }

// IsRandomAccess returns the random access indicator for this adaptation field
// func IsRandomAccess(pkt *packet.Packet) bool {
// 	return pkt[5]&0x40 != 0
// }

// IsESHigherPriority returns true if this elementary stream is
// high priority. Corresponds to the elementary stream
// priority indicator.
// func IsESHigherPriority(pkt *packet.Packet) bool {
// 	return pkt[5]&0x20 != 0
// }

// HasPCR returns true when the PCR flag is set
// func HasPCR(pkt *packet.Packet) bool {
// 	return pkt[5]&0x10 != 0
// }

// HasOPCR returns true when the OPCR flag is set
// func HasOPCR(pkt *packet.Packet) bool {
// 	return pkt[5]&0x08 != 0
// }

// HasSplicingPoint returns true when the splicing countdown field is present
// func HasSplicingPoint(pkt *packet.Packet) bool {
// 	return pkt[5]&0x04 != 0
// }

// HasTransportPrivateData returns true when the private data field is present
// func HasTransportPrivateData(pkt *packet.Packet) bool {
// 	return pkt[5]&0x02 != 0
// }

// HasAdaptationFieldExtension returns true if this adaptation field contains an extension field
// func HasAdaptationFieldExtension(pkt *packet.Packet) bool {
// 	return pkt[5]&0x01 != 0
// }

// // EncoderBoundaryPoint returns the byte array located in the optional TransportPrivateData of the (also optional)
// // AdaptationField of the Packet. If either of these optional fields are missing an empty byte array is returned with an error
// func EncoderBoundaryPoint(pkt *packet.Packet) ([]byte, error) {
// 	hasAdapt, err := packet.ContainsAdaptationField(pkt)
// 	if err != nil {
// 		return nil, nil
// 	}
// 	if hasAdapt && Length(pkt) > 0 && HasTransportPrivateData(pkt) {
// 		ebp, err := TransportPrivateData(pkt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return ebp, nil
// 	}
// 	return nil, gots.ErrNoEBP
// }

// PCR is the Program Clock Reference.
// First 33 bits are PCR base.
// Next 6 bits are reserved.
// Final 9 bits are PCR extension.
// func PCR(pkt *packet.Packet) ([]byte, error) {
// 	if !HasPCR(pkt) {
// 		return nil, gots.ErrNoPCR
// 	}
// 	offset := 6
// 	return pkt[offset : offset+6], nil
// }

// OPCR is the Original Program Clock Reference.
// First 33 bits are original PCR base.
// Next 6 bits are reserved.
// Final 9 bits are original PCR extension.
// func OPCR(pkt *packet.Packet) ([]byte, error) {
// 	if !HasOPCR(pkt) {
// 		return nil, gots.ErrNoOPCR
// 	}
// 	offset := 6
// 	if HasPCR(pkt) {
// 		offset += 6
// 	}
// 	return pkt[offset : offset+6], nil
// }

// SpliceCountdown returns a count of how many packets after this one until
// a splice point occurs or an error if none exist. This function calls
// HasSplicingPoint to check for the existence of a splice countdown.
// func SpliceCountdown(pkt *packet.Packet) (uint8, error) {
// 	if !HasSplicingPoint(pkt) {
// 		return 0, gots.ErrNoSplicePoint
// 	}
// 	offset := 6
// 	if HasPCR(pkt) {
// 		offset += 6
// 	}
// 	if HasOPCR(pkt) {
// 		offset += 6
// 	}
// 	return pkt[offset], nil
// }

// TransportPrivateData returns the private data from this adaptation field
// or an empty array and an error if there is none. This function calls
// HasTransportPrivateData to check for the existence of private data.
// func TransportPrivateData(pkt *packet.Packet) ([]byte, error) {
// 	if !HasTransportPrivateData(pkt) {
// 		return nil, gots.ErrNoPrivateTransportData
// 	}
// 	offset := 6
// 	if HasPCR(pkt) {
// 		offset += 6
// 	}
// 	if HasOPCR(pkt) {
// 		offset += 6
// 	}
// 	if HasSplicingPoint(pkt) {
// 		offset++
// 	}
// 	dataLength := uint8(pkt[offset])
// 	offset++
// 	return pkt[uint8(offset) : uint8(offset)+dataLength], nil
// }
