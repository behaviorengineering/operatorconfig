package operatorconfig

import (
	"encoding/binary"
	"testing"
)

func TestDecodeWinCredBlobUTF16AccountID(t *testing.T) {
	const want = "0123456789abcdef0123456789abcdef"
	b := utf16LEASCII(want)
	got := DecodeWinCredBlob(b)
	if got != want {
		t.Fatalf("got len=%d want len=%d", len(got), len(want))
	}
	if len(got) != 32 {
		t.Fatalf("got %q", got)
	}
}

func TestDecodeWinCredBlobUTF8Unchanged(t *testing.T) {
	const want = "secret-token-value"
	got := DecodeWinCredBlob([]byte(want))
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestDecodeWinCredBlobOddLengthUTF8(t *testing.T) {
	b := []byte("abc")
	got := DecodeWinCredBlob(b)
	if got != "abc" {
		t.Fatalf("got %q", got)
	}
}

func utf16LEASCII(s string) []byte {
	out := make([]byte, len(s)*2)
	for i, r := range s {
		binary.LittleEndian.PutUint16(out[i*2:], uint16(r))
	}
	return out
}
