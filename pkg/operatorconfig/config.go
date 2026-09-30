package operatorconfig

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Result holds the resolved config path and loaded Viper instance.
type Result struct {
	Path  string
	Viper *viper.Viper
}

// Load resolves the config path, reads YAML when present, resolves secrets and env defaults,
// then expands ${VAR} in string values.
func Load(opts Options) (Result, error) {
	path, err := ResolveConfigPath(opts)
	if err != nil {
		return Result{}, err
	}
	resolveOpts := opts
	if path == "" {
		if err := resolveSecrets(resolveOpts, nil, nil); err != nil {
			return Result{}, err
		}
		if err := ApplyEnvDefaults(opts.EnvDefaults); err != nil {
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
	if err := ApplyEnvDefaults(opts.EnvDefaults); err != nil {
		return Result{}, err
	}
	if err := expandViperStrings(v); err != nil {
		return Result{}, err
	}
	return Result{Path: path, Viper: v}, nil
}
