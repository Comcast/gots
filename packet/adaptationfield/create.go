package adaptationfield

import "github.comcast.com/viper-cog/libmpegts/packet"

func SetPrivateData(pkt *packet.Packet, af []byte) {
	offset := 6
	if HasPCR(*pkt) {
		offset += 6
	}
	if HasOPCR(*pkt) {
		offset += 6
	}
	if HasSplicingPoint(*pkt) {
		offset++
	}
	(*pkt)[offset] = byte(0x04) // data length
	offset++
	for i, b := range af {
		(*pkt)[offset+i] = b
	}
}
