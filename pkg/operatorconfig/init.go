package operatorconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitUserConfig creates the user config directory and seeds config files.
// Writes config.yaml with mode 0600 when missing or when force is true.
// Always refreshes config.yaml.example beside the live file.
func InitUserConfig(opts Options, example []byte, force bool) (written bool, err error) {
	path, err := UserConfigPath(opts.App, opts.filename())
	if err != nil {
		return false, err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return false, fmt.Errorf("operatorconfig: mkdir %s: %w", dir, err)
	}
	examplePath := filepath.Join(dir, opts.filename()+".example")
	if err := os.WriteFile(examplePath, example, 0o644); err != nil {
		return false, fmt.Errorf("operatorconfig: write example %s: %w", examplePath, err)
	}
	exists := fileExists(path)
	if exists && !force {
		return false, nil
	}
	if len(example) == 0 {
		return false, fmt.Errorf("operatorconfig: example template is empty")
	}
	if err := os.WriteFile(path, example, 0o600); err != nil {
		return false, fmt.Errorf("operatorconfig: write config %s: %w", path, err)
	}
	return true, nil
}
