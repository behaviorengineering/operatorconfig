package operatorconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSecretSOPSFillsAfterKeyringMiss(t *testing.T) {
	dir := t.TempDir()
	encPath := filepath.Join(dir, "secrets.enc.yaml")
	if err := os.WriteFile(encPath, []byte("placeholder"), 0o600); err != nil {
		t.Fatal(err)
	}
	mem := NewMemKeyring()
	t.Setenv("TOKEN", "")
	opts := Options{
		App:            "app",
		SecretsEncPath: encPath,
		Secrets:        []Secret{{Env: "TOKEN", Required: true}},
		SecretFile: &MemSecretFile{
			Files: map[string]map[string]string{
				encPath: {"TOKEN": "from-sops"},
			},
		},
	}
	if err := ResolveSecrets(opts, mem); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("TOKEN"); got != "from-sops" {
		t.Fatalf("got %q want from-sops", got)
	}
}

func TestSecretSOPSSkipsWhenFileMissing(t *testing.T) {
	dir := t.TempDir()
	encPath := filepath.Join(dir, "missing.enc.yaml")
	mem := NewMemKeyring()
	t.Setenv("TOKEN", "")
	opts := Options{
		App:            "app",
		SecretsEncPath: encPath,
		Secrets:        []Secret{{Env: "TOKEN", Required: false}},
		SecretFile:     &MemSecretFile{Files: map[string]map[string]string{}},
	}
	if err := ResolveSecrets(opts, mem); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("TOKEN"); got != "" {
		t.Fatalf("got %q want empty", got)
	}
}

func TestSecretOrderEnvKeyringSOPS(t *testing.T) {
	dir := t.TempDir()
	encPath := filepath.Join(dir, "secrets.enc.yaml")
	if err := os.WriteFile(encPath, []byte("placeholder"), 0o600); err != nil {
		t.Fatal(err)
	}
	mem := NewMemKeyring()
	if err := mem.Set("app", "TOKEN", "from-keyring"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TOKEN", "")
	opts := Options{
		App:            "app",
		SecretsEncPath: encPath,
		Secrets:        []Secret{{Env: "TOKEN", Required: true}},
		SecretFile: &MemSecretFile{
			Files: map[string]map[string]string{
				encPath: {"TOKEN": "from-sops"},
			},
		},
	}
	if err := ResolveSecrets(opts, mem); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("TOKEN"); got != "from-keyring" {
		t.Fatalf("got %q want from-keyring", got)
	}
}
