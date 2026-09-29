package operatorconfig

import (
	"os"
	"testing"
)

func TestApplyEnvDefaults_staticAndDerive(t *testing.T) {
	t.Setenv("STATIC_DEF", "")
	t.Setenv("DERIVED_DEF", "")
	err := ApplyEnvDefaults([]EnvDefault{
		{Env: "STATIC_DEF", Value: "alpha"},
		{Env: "DERIVED_DEF", Derive: func() string { return "beta" }},
	})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("STATIC_DEF") != "alpha" {
		t.Fatalf("STATIC_DEF %q", os.Getenv("STATIC_DEF"))
	}
	if os.Getenv("DERIVED_DEF") != "beta" {
		t.Fatalf("DERIVED_DEF %q", os.Getenv("DERIVED_DEF"))
	}
}

func TestApplyEnvDefaults_doesNotOverwriteSetEnv(t *testing.T) {
	t.Setenv("KEEP_ME", "from-shell")
	err := ApplyEnvDefaults([]EnvDefault{{Env: "KEEP_ME", Value: "default"}})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("KEEP_ME") != "from-shell" {
		t.Fatalf("got %q", os.Getenv("KEEP_ME"))
	}
}
