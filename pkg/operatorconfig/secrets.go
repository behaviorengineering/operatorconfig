package operatorconfig

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// ResolveSecrets fills empty process env vars from the keyring, then an optional SOPS secrets file.
// Environment always wins when non-empty. Required secrets fail closed when still empty.
func ResolveSecrets(opts Options, kr Keyring) error {
	return resolveSecrets(opts, kr, nil)
}

func effectiveKeyring(opts Options, kr Keyring) Keyring {
	if opts.Keyring != nil {
		return opts.Keyring
	}
	if kr != nil {
		return kr
	}
	return DefaultKeyring()
}

func resolveSecrets(opts Options, kr Keyring, v *viper.Viper) error {
	kr = effectiveKeyring(opts, kr)
	service := strings.TrimSpace(opts.App)
	if service == "" {
		return fmt.Errorf("operatorconfig: app name required for secrets")
	}
	if len(opts.Secrets) == 0 {
		return nil
	}

	var sopsMap map[string]string
	var sopsLoaded bool

	for _, sec := range opts.Secrets {
		env := strings.TrimSpace(sec.Env)
		if env == "" {
			return fmt.Errorf("operatorconfig: secret env name required")
		}
		raw := os.Getenv(env)
		clean := SanitizeSecret(raw)
		if clean != "" {
			if raw != clean {
				if err := setResolvedSecret(env, clean, v); err != nil {
					return err
				}
			}
			continue
		}
		got, err := kr.Get(service, env)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return fmt.Errorf("operatorconfig: keyring get %q: %w", env, err)
			}
		} else {
			got = SanitizeSecret(got)
			if got != "" {
				if err := setResolvedSecret(env, got, v); err != nil {
					return err
				}
				continue
			}
		}

		if !sopsLoaded {
			sopsMap, err = loadSOPSSecrets(opts, kr)
			if err != nil {
				return err
			}
			sopsLoaded = true
		}
		if sopsMap != nil {
			if val := SanitizeSecret(sopsMap[env]); val != "" {
				if err := setResolvedSecret(env, val, v); err != nil {
					return err
				}
				continue
			}
		}

		if sec.Required {
			return fmt.Errorf("operatorconfig: required secret %q not in environment, keyring, or secrets file", env)
		}
	}
	return nil
}

func setResolvedSecret(env, val string, v *viper.Viper) error {
	val = SanitizeSecret(val)
	if err := os.Setenv(env, val); err != nil {
		return fmt.Errorf("operatorconfig: setenv %q: %w", env, err)
	}
	if v != nil {
		v.Set(env, val)
	}
	return nil
}

func loadSOPSSecrets(opts Options, kr Keyring) (map[string]string, error) {
	path, err := SecretsEncPath(opts)
	if err != nil {
		return nil, err
	}
	if !fileExists(path) {
		return nil, nil
	}
	if err := ensureSOPSAgeKeyFromKeyring(opts, kr); err != nil {
		return nil, err
	}
	sf := effectiveSecretFile(opts)
	m, err := sf.Load(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}
