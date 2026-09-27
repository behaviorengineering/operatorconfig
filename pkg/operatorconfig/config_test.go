package operatorconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadExpandsEnvInYAML(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("MY_HOST", "example.com")
	cfgPath := filepath.Join(dir, "app", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "url: https://${MY_HOST}/v1\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Load(Options{App: "app"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Path != cfgPath {
		t.Fatalf("path %q", res.Path)
	}
	if res.Viper.GetString("url") != "https://example.com/v1" {
		t.Fatalf("url %q", res.Viper.GetString("url"))
	}
}
