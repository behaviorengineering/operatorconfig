package operatorconfig

import "testing"

func TestMemKeyringSetTrimsWhitespace(t *testing.T) {
	mem := NewMemKeyring()
	if err := mem.Set("app", "TOKEN", "  tok\r\n"); err != nil {
		t.Fatal(err)
	}
	got, err := mem.Get("app", "TOKEN")
	if err != nil {
		t.Fatal(err)
	}
	if got != "tok" {
		t.Fatalf("got %q want tok", got)
	}
}

func TestSanitizeSecret(t *testing.T) {
	if SanitizeSecret("  a\n") != "a" {
		t.Fatalf("got %q", SanitizeSecret("  a\n"))
	}
}
