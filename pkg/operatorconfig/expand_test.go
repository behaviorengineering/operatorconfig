package operatorconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_secretsBeforeExpand(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("DB_PASS", "")
	cfgPath := filepath.Join(dir, "app", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "secrets:\n  - DB_PASS\npassword: \"${DB_PASS}\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	mem := NewMemKeyring()
	if err := mem.Set("app", "DB_PASS", "ring-secret"); err != nil {
		t.Fatal(err)
	}
	res, err := Load(Options{
		App:     "app",
		Keyring: mem,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Viper.GetString("password") != "ring-secret" {
		t.Fatalf("password %q", res.Viper.GetString("password"))
	}
}

func TestLoad_envDefaultsBeforeExpand(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("GATEWAY_HOST", "")
	cfgPath := filepath.Join(dir, "app", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "url: \"https://${GATEWAY_HOST}/v1\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Load(Options{
		App: "app",
		EnvDefaults: []EnvDefault{
			{Env: "GATEWAY_HOST", Value: "gateway.example"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Viper.GetString("url") != "https://gateway.example/v1" {
		t.Fatalf("url %q", res.Viper.GetString("url"))
	}
}

func TestLoad_unresolvedPlaceholderFails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgPath := filepath.Join(dir, "app", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "url: \"https://example.com/${\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(Options{App: "app"})
	if err == nil {
		t.Fatal("expected error for unresolved placeholder")
	}
}
