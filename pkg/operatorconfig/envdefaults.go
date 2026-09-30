package operatorconfig

import (
	"fmt"
	"os"
	"strings"
)

// EnvDefault fills a process env var when it is empty, after secrets resolve and before YAML expand.
type EnvDefault struct {
	Env    string
	Value  string
	Derive func() string
}

// ApplyEnvDefaults sets each Env when the current process value is empty.
// Derive runs when Value is empty; static Value is used when Derive is nil.
func ApplyEnvDefaults(defaults []EnvDefault) error {
	for _, d := range defaults {
		env := strings.TrimSpace(d.Env)
		if env == "" {
			return fmt.Errorf("operatorconfig: env default name required")
		}
		if strings.TrimSpace(os.Getenv(env)) != "" {
			continue
		}
		val := strings.TrimSpace(d.Value)
		if d.Derive != nil {
			val = strings.TrimSpace(d.Derive())
		}
		if val == "" {
			continue
		}
		if err := os.Setenv(env, val); err != nil {
			return fmt.Errorf("operatorconfig: set env %q: %w", env, err)
		}
	}
	return nil
}
