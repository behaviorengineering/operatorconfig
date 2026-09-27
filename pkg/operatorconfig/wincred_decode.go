package operatorconfig

import (
	"encoding/binary"
	"unicode/utf16"
)

// DecodeWinCredBlob normalizes a Windows CredentialBlob to UTF-8 text.
// Generic credentials from the Control Panel are often UTF-16 LE; go-keyring
// returns the raw bytes as a string without decoding.
func DecodeWinCredBlob(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if u16, ok := decodeUTF16LEASCII(b); ok {
		return u16
	}
	return trimNULs(string(b))
}

func decodeUTF16LEASCII(b []byte) (string, bool) {
	if len(b) < 2 || len(b)%2 != 0 {
		return "", false
	}
	// UTF-16 LE ASCII: odd-index bytes are 0x00.
	for i := 1; i < len(b); i += 2 {
		if b[i] != 0 {
			return "", false
		}
	}
	u16 := make([]uint16, len(b)/2)
	for i := range u16 {
		u16[i] = binary.LittleEndian.Uint16(b[i*2 : i*2+2])
	}
	runes := utf16.Decode(u16)
	end := len(runes)
	for i, r := range runes {
		if r == 0 {
			end = i
			break
		}
		if r < 0x20 || r == 0x7f {
			return "", false
		}
	}
	return string(runes[:end]), true
}

func trimNULs(s string) string {
	for len(s) > 0 && s[len(s)-1] == 0 {
		s = s[:len(s)-1]
	}
	return s
}
