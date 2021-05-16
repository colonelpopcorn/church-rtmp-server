package main

import (
	"testing"
)

func TestGenerateGuid(t *testing.T) {
	guid := generateGUID()
	if len(guid) != 22 {
		t.Error("GUID not long enough")
	}
}

func TestGeneratePassword(t *testing.T) {
	pwd := generatePassword(12)
	if len(pwd) != (12 * 2) {
		t.Errorf("Expected length of password was %d, got %d", (12 * 2), len(pwd))
	}
}
