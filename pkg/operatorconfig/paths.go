package operatorconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UserConfigPath returns $XDG_CONFIG_HOME/<app>/<filename> or ~/.config/<app>/<filename>.
func UserConfigPath(app, filename string) (string, error) {
	app = strings.TrimSpace(app)
	if app == "" {
		return "", fmt.Errorf("operatorconfig: app name required")
	}
	if strings.TrimSpace(filename) == "" {
		filename = "config.yaml"
	}
	if xdg := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); xdg != "" {
		return filepath.Join(xdg, app, filename), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", app, filename), nil
}

// ResolveConfigPath picks the first config file using Options discovery order.
// Returns ("", nil) when no candidate exists. Override paths must exist or an error is returned.
func ResolveConfigPath(opts Options) (string, error) {
	if p := strings.TrimSpace(opts.ConfigFlagPath); p != "" {
		return requireExisting(p)
	}
	if opts.ConfigEnv != "" {
		if p := strings.TrimSpace(os.Getenv(opts.ConfigEnv)); p != "" {
			return requireExisting(p)
		}
	}
	userPath, err := UserConfigPath(opts.App, opts.filename())
	if err != nil {
		return "", err
	}
	if fileExists(userPath) {
		return userPath, nil
	}
	for _, candidate := range opts.ExtraPaths {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if fileExists(candidate) {
			return candidate, nil
		}
	}
	return "", nil
}

func requireExisting(path string) (string, error) {
	if !fileExists(path) {
		return "", fmt.Errorf("operatorconfig: config file %q not found", path)
	}
	return path, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
