package operatorconfig

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/getsops/sops/v3/decrypt"
	"gopkg.in/yaml.v3"
)

// SecretFile loads plaintext key/value pairs from an encrypted secrets file.
type SecretFile interface {
	Load(path string) (map[string]string, error)
}

// SOPSSecretFile decrypts a file with the SOPS CLI-compatible format (yaml, dotenv, etc.).
type SOPSSecretFile struct{}

// Load decrypts path when it exists. Missing file returns (nil, os.ErrNotExist).
func (SOPSSecretFile) Load(path string) (map[string]string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, os.ErrNotExist
	}
	if !fileExists(path) {
		return nil, os.ErrNotExist
	}
	format := sopsFormatForPath(path)
	plain, err := decrypt.File(path, format)
	if err != nil {
		return nil, fmt.Errorf("operatorconfig: sops decrypt %q: %w", path, err)
	}
	return parseSecretsCleartext(plain, format)
}

// MemSecretFile is a test double mapping paths to plaintext env maps.
type MemSecretFile struct {
	Files map[string]map[string]string
}

func (m *MemSecretFile) Load(path string) (map[string]string, error) {
	if m == nil || m.Files == nil {
		return nil, os.ErrNotExist
	}
	if got, ok := m.Files[path]; ok {
		return got, nil
	}
	return nil, os.ErrNotExist
}

func effectiveSecretFile(opts Options) SecretFile {
	if opts.SecretFile != nil {
		return opts.SecretFile
	}
	return SOPSSecretFile{}
}

func sopsFormatForPath(path string) string {
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".env"), strings.HasSuffix(lower, ".enc.env"):
		return "dotenv"
	default:
		return "yaml"
	}
}

func parseSecretsCleartext(plain []byte, format string) (map[string]string, error) {
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "dotenv" {
		return parseDotenvSecrets(plain)
	}
	return parseYAMLSecrets(plain)
}

func parseYAMLSecrets(plain []byte) (map[string]string, error) {
	var doc map[string]interface{}
	if err := yaml.Unmarshal(plain, &doc); err != nil {
		return nil, fmt.Errorf("operatorconfig: parse secrets yaml: %w", err)
	}
	out := make(map[string]string)
	for k, v := range doc {
		if k == "secrets" {
			continue
		}
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out, nil
}

func parseDotenvSecrets(plain []byte) (map[string]string, error) {
	out := make(map[string]string)
	for _, line := range bytes.Split(plain, []byte{'\n'}) {
		s := strings.TrimSpace(string(line))
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		before, after, ok := strings.Cut(s, "=")
		if !ok {
			continue
		}
		key := strings.TrimSpace(before)
		val := strings.TrimSpace(after)
		val = strings.Trim(val, `"'`)
		if key != "" {
			out[key] = val
		}
	}
	return out, nil
}

// SecretsEncPath returns the encrypted secrets file path for opts.
func SecretsEncPath(opts Options) (string, error) {
	if p := strings.TrimSpace(opts.SecretsEncPath); p != "" {
		return p, nil
	}
	return UserConfigPath(opts.App, "secrets.enc.yaml")
}

func ensureSOPSAgeKeyFromKeyring(opts Options, kr Keyring) error {
	if strings.TrimSpace(os.Getenv("SOPS_AGE_KEY")) != "" {
		return nil
	}
	if p := strings.TrimSpace(os.Getenv("SOPS_AGE_KEY_FILE")); p != "" {
		if fileExists(p) {
			return nil
		}
	}
	service := strings.TrimSpace(opts.App)
	if service == "" {
		return nil
	}
	kr = effectiveKeyring(opts, kr)
	got, err := kr.Get(service, "SOPS_AGE_KEY")
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		return fmt.Errorf("operatorconfig: keyring get SOPS_AGE_KEY: %w", err)
	}
	got = strings.TrimSpace(got)
	if got == "" {
		return nil
	}
	if err := os.Setenv("SOPS_AGE_KEY", got); err != nil {
		return fmt.Errorf("operatorconfig: setenv SOPS_AGE_KEY: %w", err)
	}
	return nil
}
