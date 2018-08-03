package adaptationfield

import "github.com/Comcast/gots/packet"

func SetPrivateData(pkt *packet.Packet, af []byte) {
	offset := 6
	if HasPCR(pkt) {
		offset += 6
	}
	if HasOPCR(pkt) {
		offset += 6
	}
	if HasSplicingPoint(pkt) {
		offset++
	}
	pkt[offset] = byte(0x04) // data length
	offset++
	// FIXME(kortschak): Handle len(af) != 4.
	copy(pkt[offset:offset+4], af)
}
