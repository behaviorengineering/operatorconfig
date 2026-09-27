package operatorconfig

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestSecretEnvWins(t *testing.T) {
	mem := NewMemKeyring()
	if err := mem.Set("app", "TOKEN", "from-keyring"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TOKEN", "from-env")
	opts := Options{
		App:     "app",
		Secrets: []Secret{{Env: "TOKEN", Required: true}},
	}
	if err := ResolveSecrets(opts, mem); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("TOKEN"); got != "from-env" {
		t.Fatalf("got %q want from-env", got)
	}
}

func TestSecretKeyringFillsEnv(t *testing.T) {
	mem := NewMemKeyring()
	if err := mem.Set("app", "TOKEN", "ring-token"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TOKEN", "")
	opts := Options{
		App:     "app",
		Secrets: []Secret{{Env: "TOKEN", Required: true}},
	}
	if err := ResolveSecrets(opts, mem); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("TOKEN"); got != "ring-token" {
		t.Fatalf("got %q want ring-token", got)
	}
}

func TestSecretRequiredMiss(t *testing.T) {
	mem := NewMemKeyring()
	t.Setenv("MISSING_SECRET", "")
	opts := Options{
		App:     "app",
		Secrets: []Secret{{Env: "MISSING_SECRET", Required: true}},
	}
	err := ResolveSecrets(opts, mem)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "MISSING_SECRET") {
		t.Fatalf("error %v should mention env name", err)
	}
}

func TestSecretEmptyNameInOptionsFailsClosed(t *testing.T) {
	opts := Options{
		App:     "app",
		Secrets: []Secret{{Env: "  ", Required: false}},
	}
	err := ResolveSecrets(opts, NewMemKeyring())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "secret env name required") {
		t.Fatalf("error %v", err)
	}
}

func TestSecretKeyringBackendErrorFailsClosed(t *testing.T) {
	t.Setenv("TOKEN", "")
	opts := Options{
		App:     "app",
		Secrets: []Secret{{Env: "TOKEN", Required: false}},
	}
	err := ResolveSecrets(opts, errKeyring{err: fmt.Errorf("dbus unavailable")})
	if err == nil {
		t.Fatal("expected keyring backend error")
	}
	if !strings.Contains(err.Error(), "TOKEN") {
		t.Fatalf("error %v should mention env name", err)
	}
	if !strings.Contains(err.Error(), "dbus unavailable") {
		t.Fatalf("error %v should wrap backend cause", err)
	}
}

type errKeyring struct {
	err error
}

func (e errKeyring) Get(service, account string) (string, error) { return "", e.err }
func (e errKeyring) Set(service, account, secret string) error   { return nil }
func (e errKeyring) Delete(service, account string) error        { return nil }
