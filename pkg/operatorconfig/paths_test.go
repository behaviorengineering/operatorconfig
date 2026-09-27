package operatorconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigPathOverrideEnv(t *testing.T) {
	explicit := filepath.Join(t.TempDir(), "explicit.yaml")
	if err := os.WriteFile(explicit, []byte("x: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TESTAPP_CONFIG", explicit)
	got, err := ResolveConfigPath(Options{
		App:       "testapp",
		ConfigEnv: "TESTAPP_CONFIG",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != explicit {
		t.Fatalf("got %q want %q", got, explicit)
	}
}

func TestResolveConfigPathXDG(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "myapp", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfg), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfg, []byte("x: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)
	got, err := ResolveConfigPath(Options{App: "myapp"})
	if err != nil {
		t.Fatal(err)
	}
	if got != cfg {
		t.Fatalf("got %q want %q", got, cfg)
	}
}

func TestResolveConfigPathOverrideMissing(t *testing.T) {
	t.Setenv("TESTAPP_CONFIG", filepath.Join(t.TempDir(), "missing.yaml"))
	_, err := ResolveConfigPath(Options{
		App:       "testapp",
		ConfigEnv: "TESTAPP_CONFIG",
	})
	if err == nil {
		t.Fatal("expected error for missing override")
	}
}

func TestResolveConfigPathNoFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	got, err := ResolveConfigPath(Options{App: "empty"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("got %q want empty", got)
	}
}

func TestResolveConfigPathExtraPaths(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cwdCfg := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(cwdCfg, []byte("x: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ResolveConfigPath(Options{
		App:        "nouser",
		ExtraPaths: []string{cwdCfg},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != cwdCfg {
		t.Fatalf("got %q want %q", got, cwdCfg)
	}
}

func TestResolveConfigPathFlagWins(t *testing.T) {
	flagPath := filepath.Join(t.TempDir(), "flag.yaml")
	if err := os.WriteFile(flagPath, []byte("x: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(t.TempDir(), "other.yaml")
	if err := os.WriteFile(other, []byte("x: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TESTAPP_CONFIG", other)
	got, err := ResolveConfigPath(Options{
		App:            "testapp",
		ConfigEnv:      "TESTAPP_CONFIG",
		ConfigFlagPath: flagPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != flagPath {
		t.Fatalf("got %q want %q", got, flagPath)
	}
}
