package packet

import (
	"testing"
)

func generatePacketAF(t *testing.T, AFString string) AdaptationField {
	p := createPacketEmptyBody(t, "47000030"+AFString)
	a := p.AdaptationField()
	if a == nil {
		t.Errorf("failed to get adaptation field")
	}
	return a
}

func TestDiscontinuity(t *testing.T) {
	a := generatePacketAF(t, "0180")
	if !a.Discontinuity() {
		t.Errorf("failed to read discontinuity correctly.")
	}
	a = generatePacketAF(t, "0190")
	if !a.Discontinuity() {
		t.Errorf("failed to read discontinuity correctly.")
	}
	a = generatePacketAF(t, "0170")
	if a.Discontinuity() {
		t.Errorf("failed to read discontinuity correctly.")
	}
}

func TestAdaptationField(t *testing.T) {
	p := createPacketEmptyBody(t, "470000300102")
	a := p.AdaptationField()
	if a == nil {
		t.Errorf("Error getting adaptation field.")
	}

	p = createPacketEmptyBody(t, "470000100002")
	a = p.AdaptationField()
	if a != nil {
		t.Errorf("Adaptation field does not exist but something was returned.")
	}
}
