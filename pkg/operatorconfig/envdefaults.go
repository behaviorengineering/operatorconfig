package operatorconfig

import (
	"fmt"
	"os"
	"strings"
)

// EnvDefault sets a process environment variable when it is still empty after secret resolution.
type EnvDefault struct {
	Env    string
	Value  string
	Derive func() (string, error)
}

// ApplyEnvDefaults sets env vars from opts.EnvDefaults when the process env is empty.
// Non-empty env (including values from ResolveSecrets) is never overwritten.
func ApplyEnvDefaults(opts Options) error {
	if len(opts.EnvDefaults) == 0 {
		return nil
	}
	for _, d := range opts.EnvDefaults {
		env := strings.TrimSpace(d.Env)
		if env == "" {
			return fmt.Errorf("operatorconfig: env default name required")
		}
		if strings.TrimSpace(os.Getenv(env)) != "" {
			continue
		}
		val := strings.TrimSpace(d.Value)
		if val == "" && d.Derive != nil {
			derived, err := d.Derive()
			if err != nil {
				return fmt.Errorf("operatorconfig: derive %q: %w", env, err)
			}
			val = strings.TrimSpace(derived)
		}
		if val == "" {
			continue
		}
		if err := os.Setenv(env, val); err != nil {
			return fmt.Errorf("operatorconfig: setenv %q: %w", env, err)
		}
	}
	return nil
}
