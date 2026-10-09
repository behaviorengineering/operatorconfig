package operatorconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyEnvDefaultsSetsWhenEmpty(t *testing.T) {
	t.Setenv("TEST_DEFAULT_ENV", "")
	err := ApplyEnvDefaults(Options{
		EnvDefaults: []EnvDefault{{Env: "TEST_DEFAULT_ENV", Value: "from-default"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("TEST_DEFAULT_ENV") != "from-default" {
		t.Fatalf("got %q", os.Getenv("TEST_DEFAULT_ENV"))
	}
}

func TestApplyEnvDefaultsSkipsWhenSet(t *testing.T) {
	t.Setenv("TEST_DEFAULT_ENV", "existing")
	err := ApplyEnvDefaults(Options{
		EnvDefaults: []EnvDefault{{Env: "TEST_DEFAULT_ENV", Value: "from-default"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("TEST_DEFAULT_ENV") != "existing" {
		t.Fatalf("got %q", os.Getenv("TEST_DEFAULT_ENV"))
	}
}

func TestApplyEnvDefaultsDerive(t *testing.T) {
	t.Setenv("TEST_DERIVE_ENV", "")
	err := ApplyEnvDefaults(Options{
		EnvDefaults: []EnvDefault{
			{Env: "TEST_DERIVE_ENV", Derive: func() (string, error) { return "derived", nil }},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("TEST_DERIVE_ENV") != "derived" {
		t.Fatalf("got %q", os.Getenv("TEST_DERIVE_ENV"))
	}
}

func TestLoadAppliesEnvDefaultsBeforeExpand(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("LOOP_HOST", "")
	cfgPath := filepath.Join(dir, "app", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "url: http://${LOOP_HOST}:1320\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Load(Options{
		App: "app",
		EnvDefaults: []EnvDefault{{Env: "LOOP_HOST", Value: "127.0.0.1"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Viper.GetString("url") != "http://127.0.0.1:1320" {
		t.Fatalf("url %q", res.Viper.GetString("url"))
	}
}
