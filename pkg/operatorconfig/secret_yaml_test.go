package operatorconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSecretUnmarshalRejectsEmptyName(t *testing.T) {
	_, err := ParseSecretsYAML([]byte("secrets:\n  - \"\"\n"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "secret env name required") {
		t.Fatalf("error %v", err)
	}
}

func TestSecretUnmarshalScalarAndMapping(t *testing.T) {
	secrets, err := ParseSecretsYAML([]byte(`
secrets:
  - CF_AI_API_KEY
  - env: CF_ACCOUNT_ID
    required: true
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(secrets) != 2 {
		t.Fatalf("len %d", len(secrets))
	}
	if secrets[0].Env != "CF_AI_API_KEY" || secrets[0].Required {
		t.Fatalf("first %+v", secrets[0])
	}
	if secrets[1].Env != "CF_ACCOUNT_ID" || !secrets[1].Required {
		t.Fatalf("second %+v", secrets[1])
	}
}

func TestLoadResolvesSecretsFromYAMLFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	mem := NewMemKeyring()
	if err := mem.Set("app", "TOKEN", "from-ring"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TOKEN", "")
	cfgPath := filepath.Join(dir, "app", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "secrets:\n  - TOKEN\nfoo: bar\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(Options{App: "app", Keyring: mem})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("TOKEN") != "from-ring" {
		t.Fatalf("TOKEN %q", os.Getenv("TOKEN"))
	}
}

func TestLoadSkipsKeyringWhenNoSecretsInFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgPath := filepath.Join(dir, "app", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, []byte("foo: bar\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	boom := &boomKeyring{}
	_, err := Load(Options{App: "app", Keyring: boom})
	if err != nil {
		t.Fatal(err)
	}
	if boom.called {
		t.Fatal("expected no keyring Get")
	}
}

type boomKeyring struct {
	called bool
}

func (b *boomKeyring) Get(service, account string) (string, error) {
	b.called = true
	return "", ErrNotFound
}

func (b *boomKeyring) Set(service, account, secret string) error { return nil }
func (b *boomKeyring) Delete(service, account string) error      { return nil }
