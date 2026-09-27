package operatorconfig

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// UnmarshalYAML accepts a scalar env name or a mapping with env and required.
func (s *Secret) UnmarshalYAML(value *yaml.Node) error {
	if value == nil {
		return fmt.Errorf("operatorconfig: secret entry is nil")
	}
	switch value.Kind {
	case yaml.ScalarNode:
		s.Env = strings.TrimSpace(value.Value)
		s.Required = false
		return requireSecretEnv(s.Env)
	case yaml.MappingNode:
		var raw struct {
			Env      string `yaml:"env"`
			Required bool   `yaml:"required"`
		}
		if err := value.Decode(&raw); err != nil {
			return err
		}
		s.Env = strings.TrimSpace(raw.Env)
		s.Required = raw.Required
		return requireSecretEnv(s.Env)
	default:
		return fmt.Errorf("operatorconfig: secret entry must be a string or mapping")
	}
}

// ParseSecretsYAML decodes a top-level secrets key from a YAML document fragment.
func ParseSecretsYAML(data []byte) ([]Secret, error) {
	var doc struct {
		Secrets []Secret `yaml:"secrets"`
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("operatorconfig: parse secrets: %w", err)
	}
	return doc.Secrets, nil
}

func secretsFromViper(v *viper.Viper) ([]Secret, error) {
	if v == nil || !v.IsSet("secrets") {
		return nil, nil
	}
	sub := v.Get("secrets")
	if sub == nil {
		return nil, nil
	}
	b, err := yaml.Marshal(map[string]interface{}{"secrets": sub})
	if err != nil {
		return nil, fmt.Errorf("operatorconfig: marshal secrets: %w", err)
	}
	return ParseSecretsYAML(b)
}

func requireSecretEnv(env string) error {
	if strings.TrimSpace(env) == "" {
		return fmt.Errorf("operatorconfig: secret env name required")
	}
	return nil
}
