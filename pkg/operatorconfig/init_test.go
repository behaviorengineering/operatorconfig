package operatorconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitUserConfigCreatesFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	example := []byte("hello: world\n")
	written, err := InitUserConfig(Options{App: "demo"}, example, false)
	if err != nil {
		t.Fatal(err)
	}
	if !written {
		t.Fatal("expected written true")
	}
	live := filepath.Join(dir, "demo", "config.yaml")
	info, err := os.Stat(live)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode %o want 0600", info.Mode().Perm())
	}
	ex := filepath.Join(dir, "demo", "config.yaml.example")
	if _, err := os.Stat(ex); err != nil {
		t.Fatal(err)
	}
}

func TestInitUserConfigNoClobber(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	appDir := filepath.Join(dir, "demo")
	if err := os.MkdirAll(appDir, 0o700); err != nil {
		t.Fatal(err)
	}
	live := filepath.Join(appDir, "config.yaml")
	if err := os.WriteFile(live, []byte("keep: me\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	written, err := InitUserConfig(Options{App: "demo"}, []byte("new: data\n"), false)
	if err != nil {
		t.Fatal(err)
	}
	if written {
		t.Fatal("expected not written")
	}
	raw, err := os.ReadFile(live)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "keep: me\n" {
		t.Fatalf("config clobbered: %q", raw)
	}
}

func TestInitUserConfigForce(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	appDir := filepath.Join(dir, "demo")
	if err := os.MkdirAll(appDir, 0o700); err != nil {
		t.Fatal(err)
	}
	live := filepath.Join(appDir, "config.yaml")
	if err := os.WriteFile(live, []byte("old\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	written, err := InitUserConfig(Options{App: "demo"}, []byte("new\n"), true)
	if err != nil {
		t.Fatal(err)
	}
	if !written {
		t.Fatal("expected written")
	}
	raw, _ := os.ReadFile(live)
	if string(raw) != "new\n" {
		t.Fatalf("got %q", raw)
	}
}
