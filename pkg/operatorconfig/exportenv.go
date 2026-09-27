package operatorconfig

import (
	"fmt"
	"os"
	"strings"
)

// WriteDotenv writes resolved secret env vars to path with mode 0600.
// Keys must already be present in the process environment after ResolveSecrets.
func WriteDotenv(path string, secrets []Secret) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("operatorconfig: dotenv path required")
	}
	var b strings.Builder
	for _, sec := range secrets {
		env := strings.TrimSpace(sec.Env)
		if env == "" {
			continue
		}
		val := os.Getenv(env)
		if val == "" {
			continue
		}
		b.WriteString(env)
		b.WriteByte('=')
		b.WriteString(quoteDotenvValue(val))
		b.WriteByte('\n')
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		return fmt.Errorf("operatorconfig: write dotenv %q: %w", path, err)
	}
	return nil
}

func quoteDotenvValue(val string) string {
	if strings.ContainsAny(val, " \t\n\"'#") {
		return fmt.Sprintf("%q", val)
	}
	return val
}
