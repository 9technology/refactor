package confirm

import (
	"testing"
)

func TestConfirmation(t *testing.T) {
	var confirmation Confirmation
	if confirmation.state != unconfirmed {
		t.Errorf("should start unconfirmed")
	}
	confirmation.ConfirmOnce()
	if !confirmation.Next() {
		t.Errorf("ConfirmOnce() should cause Next() to return true")
	}
	if confirmation.Next() {
		t.Errorf("ConfirmOnce() shouldn't return true multiple times")
	}
	confirmation.ConfirmAll()
	if !confirmation.Next() {
		t.Errorf("ConfirmAll() should cause Next() to always return true")
	}
	if !confirmation.Next() {
		t.Errorf("ConfirmAll() should cause Next() to always return true")
	}
	if !confirmation.Next() {
		t.Errorf("ConfirmAll() should cause Next() to always return true")
	}
}
