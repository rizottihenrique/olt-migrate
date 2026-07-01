package crypto

import (
	"testing"
)

func TestDecryptPassword(t *testing.T) {
	tests := []struct {
		encrypted string
		expected  string
	}{
		{"=<;:", "abcd"},
		{"9876", "efgh"},
		{"543210", "ijklmn"},
		{"/.-,+*", "opqrst"},
		{")('&%$", "uvwxyz"},
		{"nmlkjihgfe", "0123456789"},
		{"m2k4", "1l3j"},
	}

	for _, tt := range tests {
		result := DecryptPassword(tt.encrypted)
		if result != tt.expected {
			t.Errorf("DecryptPassword(%q) = %q; esperava %q", tt.encrypted, result, tt.expected)
		}
	}
}

func TestEncryptPassword(t *testing.T) {
	tests := []struct {
		plaintext string
		expected  string
	}{
		{"abcd", "=<;:"},
		{"efgh", "9876"},
		{"0123456789", "nmlkjihgfe"},
		{"1l3j", "m2k4"},
	}

	for _, tt := range tests {
		result := EncryptPassword(tt.plaintext)
		if result != tt.expected {
			t.Errorf("EncryptPassword(%q) = %q; esperava %q", tt.plaintext, result, tt.expected)
		}
	}
}
