package packet

import (
	"encoding/hex"
	"testing"
)

const (
	testPacket = "4700001" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000000000000000000" +
		"000000000000000000000000000000000"
)

func createPacketEmptyBody(t *testing.T, header string) (p Packet) {
	headerBytes, _ := hex.DecodeString(header)
	bodyBytes := make([]byte, 188-len(headerBytes))
	packetByets := Packet(append(headerBytes, bodyBytes...))

	p, err := FromBytes(packetByets)
	if err != nil {
		t.Error("packet Error checking failed.")
	}
	return
}

func assertPacket(t *testing.T, target Packet, generated Packet) {
	if err := target.CheckErrors(); err != nil {
		t.Error("error in target packet")
	}

	if err := generated.CheckErrors(); err != nil {
		t.Error("error in generated packet during modification.")
	}

	if !target.Equals(generated) {
		t.Errorf("crafted packet:\n%X \ndoes not match expected packet:\n%X.", generated, target)
	}
}

func TestFromBytes(t *testing.T) {
	bytes, _ := hex.DecodeString(testPacket)
	p, err := FromBytes(bytes)
	if err != nil {
		t.Error("Packet Error checking failed.")
	}
	if (len(bytes) != 188) && (len(p) != 188) {
		t.Error("Packets are not 188 bytes.")
		return
	}
	for i := range bytes {
		if bytes[i] != p[i] {
			t.Error("Packet generated with FromBytes did not copy bytes correctly.")
			return
		}
	}
}

func TestNewPacket(t *testing.T) {
	target := createPacketEmptyBody(t, "471FFF10")
	generated := NewPacket()
	if err := generated.CheckErrors(); err != nil {
		t.Error("Default packet has errors.")
	}
	assertPacket(t, target, generated)
}

func TestSetTransportErrorIndicator(t *testing.T) {
	generated := NewPacket()

	target := createPacketEmptyBody(t, "479FFF10")
	generated.SetTransportErrorIndicator(true)
	assertPacket(t, target, generated)

	target = createPacketEmptyBody(t, "471FFF10")
	generated.SetTransportErrorIndicator(false)
	assertPacket(t, target, generated)
}

func TestTransportErrorIndicator(t *testing.T) {
	pkt := createPacketEmptyBody(t, "479FFF10")
	if pkt.TransportErrorIndicator() != true {
		t.Error("Failed to set read set TransportErrorIndicator flag.")
	}
}

func TestSetPayloadUnitStartIndicator(t *testing.T) {
	generated := NewPacket()

	target := createPacketEmptyBody(t, "475FFF10")
	generated.SetPayloadUnitStartIndicator(true)
	assertPacket(t, target, generated)

	target = createPacketEmptyBody(t, "471FFF10")
	generated.SetPayloadUnitStartIndicator(false)
	assertPacket(t, target, generated)
}

func TestPayloadUnitStartIndicator(t *testing.T) {
	pkt := createPacketEmptyBody(t, "475FFF10")
	if pkt.PayloadUnitStartIndicator() != true {
		t.Error("Failed to set read set PayloadUnitStartIndicator flag.")
	}
}

func TestSetTP(t *testing.T) {
	generated := NewPacket()

	target := createPacketEmptyBody(t, "473FFF10")
	generated.SetTransportPriority(true)
	assertPacket(t, target, generated)

	target = createPacketEmptyBody(t, "471FFF10")
	generated.SetTransportPriority(false)
	assertPacket(t, target, generated)
}

func TestTP(t *testing.T) {
	pkt := createPacketEmptyBody(t, "473FFF10")
	if pkt.TransportPriority() != true {
		t.Error("Failed to set read set TP flag.")
	}
}

func TestSetPID(t *testing.T) {
	generated := createPacketEmptyBody(t, "47e09010")

	target := createPacketEmptyBody(t, "47f76A10")
	generated.SetPID(0x176A)
	assertPacket(t, target, generated)

	generated = NewPacket()

	target = createPacketEmptyBody(t, "47000010")
	generated.SetPID(0x0000)
	assertPacket(t, target, generated)

	target = createPacketEmptyBody(t, "471fec10")
	generated.SetPID(0x1fec)
	assertPacket(t, target, generated)
}

func TestPID(t *testing.T) {
	pkt := createPacketEmptyBody(t, "47344410")
	if pkt.PID() != 0x1444 {
		t.Errorf("Reads different PID than expected. Expected: %d", 0x1444)
	}
}

func TestSetTransportScramblingControl(t *testing.T) {
	generated := NewPacket()

	target := createPacketEmptyBody(t, "471FFFD0")
	generated.SetTransportScramblingControl(ScrambleOddKeyFlag)
	assertPacket(t, target, generated)

	target = createPacketEmptyBody(t, "471FFF90")
	generated.SetTransportScramblingControl(ScrambleEvenKeyFlag)
	assertPacket(t, target, generated)

	target = createPacketEmptyBody(t, "471FFF10")
	generated.SetTransportScramblingControl(NoScrambleFlag)
	assertPacket(t, target, generated)
}

func TestTransportScramblingControl(t *testing.T) {
	pkt := createPacketEmptyBody(t, "471FFF90")
	if pkt.TransportScramblingControl() != ScrambleEvenKeyFlag {
		t.Error("Failed to set TransportScramblingControl to ScrambleEvenKeyFlag.")
	}
	pkt = createPacketEmptyBody(t, "471FFFD0")
	if pkt.TransportScramblingControl() != ScrambleOddKeyFlag {
		t.Error("Failed to set TransportScramblingControl to ScrambleOddKeyFlag.")
	}
	pkt = createPacketEmptyBody(t, "471FFF10")
	if pkt.TransportScramblingControl() != NoScrambleFlag {
		t.Error("Failed to set TransportScramblingControl to NoScrambleFlag.")
	}
}

func TestSetAdaptationFieldControl(t *testing.T) {
	generated := NewPacket()

	target := createPacketEmptyBody(t, "471FFF10")
	generated.SetAdaptationFieldControl(PayloadFlag)
	assertPacket(t, target, generated)

	target = createPacketEmptyBody(t, "471FFF30")
	generated.SetAdaptationFieldControl(PayloadAndAdaptationFieldFlag)
	assertPacket(t, target, generated)

	target = createPacketEmptyBody(t, "471FFF20")
	generated.SetAdaptationFieldControl(AdaptationFieldFlag)
	assertPacket(t, target, generated)
}

func TestAdaptationFieldControl(t *testing.T) {
	pkt := createPacketEmptyBody(t, "471FFF90")
	if pkt.AdaptationFieldControl() != PayloadFlag {
		t.Error("Failed to set AdaptationFieldControl to PayloadFlag.")
	}
	pkt = createPacketEmptyBody(t, "471FFFB0")
	if pkt.AdaptationFieldControl() != PayloadAndAdaptationFieldFlag {
		t.Error("Failed to set AdaptationFieldControl to PayloadAndAdaptationFieldFlag.")
	}
	pkt = createPacketEmptyBody(t, "471FFFA0")
	if pkt.AdaptationFieldControl() != AdaptationFieldFlag {
		t.Error("Failed to set AdaptationFieldControl to AdaptationFieldFlag.")
	}
}

func TestSetContinuityCounter(t *testing.T) {
	target := createPacketEmptyBody(t, "471FFF1f")
	generated := NewPacket()

	generated.SetContinuityCounter(15)

	assertPacket(t, target, generated)
}

func TestContinuityCounterModify(t *testing.T) {
	pkt := createPacketEmptyBody(t, "471FFF19")
	if pkt.ContinuityCounter() != 9 {
		t.Errorf("Reads different ContinuityCounter than expected.")
	}
}

func TestIncContinuityCounter(t *testing.T) {
	target := createPacketEmptyBody(t, "471FFF10")
	generated := NewPacket()

	generated.SetContinuityCounter(0xFE) // cc = 14
	generated.IncContinuityCounter()     // cc = 15
	generated.IncContinuityCounter()     // cc = 0, overflow

	assertPacket(t, target, generated)
}

func TestIsPAT(t *testing.T) {
	pkt := createPacketEmptyBody(t, "47000010")
	if !pkt.IsPAT() {
		t.Errorf("Packet should be a PAT.")
	}
}

func TestIsNull(t *testing.T) {
	pkt := createPacketEmptyBody(t, "471FFF10")
	if !pkt.IsNull() {
		t.Errorf("Packet should be a Null.")
	}
}

func TestHasAdaptationField(t *testing.T) {
	pkt := createPacketEmptyBody(t, "471FFF10")
	if pkt.HasAdaptationField() {
		t.Errorf("Packet should not have Adaptation Field (AdaptationFieldControl = 01).")
	}
	pkt = createPacketEmptyBody(t, "471FFF20")
	if !pkt.HasAdaptationField() {
		t.Errorf("Packet should have Adaptation Field (AdaptationFieldControl = 10).")
	}
	pkt = createPacketEmptyBody(t, "471FFF30")
	if !pkt.HasAdaptationField() {
		t.Errorf("Packet should have Adaptation Field (AdaptationFieldControl = 11).")
	}
}

func TestHasPayload(t *testing.T) {
	pkt := createPacketEmptyBody(t, "471FFF10")
	if !pkt.HasPayload() {
		t.Errorf("Packet should have Payload (AdaptationFieldControl = 01).")
	}
	pkt = createPacketEmptyBody(t, "471FFF20")
	if pkt.HasPayload() {
		t.Errorf("Packet should not have Payload (AdaptationFieldControl = 10).")
	}
	pkt = createPacketEmptyBody(t, "471FFF30")
	if !pkt.HasPayload() {
		t.Errorf("Packet should have Payload (AdaptationFieldControl = 11).")
	}
}

func TestHeaderComboBasic(t *testing.T) {
	target := createPacketEmptyBody(t, "47EFA098")
	generated := NewPacket()

	generated.SetContinuityCounter(7)
	generated.IncContinuityCounter()
	generated.SetPID(4000)
	generated.SetTransportErrorIndicator(true)
	generated.SetPayloadUnitStartIndicator(true)
	generated.SetTransportPriority(true)
	generated.SetAdaptationFieldControl(PayloadFlag)
	generated.SetTransportScramblingControl(ScrambleEvenKeyFlag)

	assertPacket(t, target, generated)

	if generated.ContinuityCounter() != 8 {
		t.Errorf("Reads different ContinuityCounter than expected.")
	}
	if generated.PID() != 4000 {
		t.Errorf("Reads different PID than expected.")
	}
	if generated.TransportErrorIndicator() != true {
		t.Errorf("Reads different TransportErrorIndicator than expected.")
	}
	if generated.PayloadUnitStartIndicator() != true {
		t.Errorf("Reads different PayloadUnitStartIndicator than expected.")
	}
	if generated.TransportPriority() != true {
		t.Errorf("Reads different TP than expected.")
	}
	if generated.AdaptationFieldControl() != PayloadFlag {
		t.Errorf("Reads different AdaptationFieldControl than expected.")
	}
	if generated.TransportScramblingControl() != ScrambleEvenKeyFlag {
		t.Errorf("Reads different TransportScramblingControl than expected.")
	}
}

func BenchmarkNewStyleCreate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		pkt := NewPacket()
		pkt.SetPID(13)
		pkt.SetAdaptationFieldControl(PayloadFlag)
		pkt.SetPayloadUnitStartIndicator(true)
		pkt.SetContinuityCounter(7)
	}
}

func BenchmarkOldStyleCreate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SetCC(
			Create(
				uint16(13),
				WithHasPayloadFlag,
				WithPUSI),
			7)
	}
}

func BenchmarkNewStyleRead(b *testing.B) {
	pkt := NewPacket()
	pkt.SetPID(13)
	pkt.SetAdaptationFieldControl(PayloadFlag)
	pkt.SetPayloadUnitStartIndicator(true)
	pkt.CheckErrors() // no errors possible

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pkt.ContinuityCounter()
		pkt.PayloadUnitStartIndicator()
		pkt.AdaptationFieldControl()
		pkt.PID()
	}
}

func BenchmarkOldStyleRead(b *testing.B) {
	pkt := NewPacket()
	pkt.SetPID(13)
	pkt.SetAdaptationFieldControl(PayloadFlag)
	pkt.SetPayloadUnitStartIndicator(true)
	pkt.CheckErrors() // no errors possible

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ContinuityCounter(pkt)
		PayloadUnitStartIndicator(pkt)
		ContainsPayload(pkt)
		Pid(pkt)
	}
}
