package usermanagementsvc

import (
	"testing"
)

func TestNewService(t *testing.T) {
	s := NewService()
	if s == nil {
		t.Error("Expecting NewService to return non-nil")
	}
}
