package helpers

import "testing"

func TestFakeFeedback(t *testing.T) {
	if ReadFakeFeedback() == true {
		t.Errorf("Feedback should be false initially")
	}
	SetFakeTrigger(true)
	if ReadFakeFeedback() == false {
		t.Errorf("Feedback did not update")
	}
}
