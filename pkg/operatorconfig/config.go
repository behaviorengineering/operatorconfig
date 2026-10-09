package operatorconfig

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Result holds the resolved config path and loaded Viper instance.
type Result struct {
	Path  string
	Viper *viper.Viper
}

// Load resolves the config path, reads YAML when present, expands ${VAR} in string values,
// and resolves registered secrets (env then keyring).
func Load(opts Options) (Result, error) {
	path, err := ResolveConfigPath(opts)
	if err != nil {
		return Result{}, err
	}
	if path == "" {
		if len(opts.Secrets) > 0 {
			if err := resolveSecrets(opts, nil, nil); err != nil {
				return Result{}, err
			}
		}
		if err := ApplyEnvDefaults(opts); err != nil {
			return Result{}, err
		}
		return Result{Path: "", Viper: nil}, nil
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)
	if prefix := strings.TrimSpace(opts.EnvPrefix); prefix != "" {
		v.SetEnvPrefix(prefix)
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	}
	if err := v.ReadInConfig(); err != nil {
		return Result{}, fmt.Errorf("operatorconfig: read config %s: %w", path, err)
	}
	resolveOpts := opts
	if len(resolveOpts.Secrets) == 0 {
		secrets, err := secretsFromViper(v)
		if err != nil {
			return Result{}, err
		}
		resolveOpts.Secrets = secrets
	}
	if len(resolveOpts.Secrets) > 0 {
		if err := resolveSecrets(resolveOpts, nil, v); err != nil {
			return Result{}, err
		}
	}
	if err := ApplyEnvDefaults(opts); err != nil {
		return Result{}, err
	}
	if err := expandViperStrings(v); err != nil {
		return Result{}, err
	}
	return Result{Path: path, Viper: v}, nil
}

func expandViperStrings(v *viper.Viper) error {
	for _, key := range v.AllKeys() {
		switch val := v.Get(key).(type) {
		case nil:
			continue
		case string:
			v.Set(key, os.Expand(val, os.Getenv))
		case []interface{}:
			out := make([]interface{}, len(val))
			for i, item := range val {
				if s, ok := item.(string); ok {
					out[i] = os.Expand(s, os.Getenv)
				} else {
					out[i] = item
				}
			}
			v.Set(key, out)
		case map[string]interface{}:
			expanded := make(map[string]interface{}, len(val))
			for k, item := range val {
				if s, ok := item.(string); ok {
					expanded[k] = os.Expand(s, os.Getenv)
				} else {
					expanded[k] = item
				}
			}
			v.Set(key, expanded)
		}
	}
	return nil
}
