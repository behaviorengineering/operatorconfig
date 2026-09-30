package operatorconfig

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

// Matches ${VAR} or $VAR (uppercase env-style names).
var envPlaceholderRE = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

func expandEnvString(s string) string {
	return strings.TrimSpace(envPlaceholderRE.ReplaceAllStringFunc(s, func(match string) string {
		if strings.HasPrefix(match, "${") {
			return os.Getenv(match[2 : len(match)-1])
		}
		varName := match[1:]
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match
	}))
}

func expandViperStrings(v *viper.Viper) error {
	for _, key := range v.AllKeys() {
		val := v.Get(key)
		expanded, err := expandValue(val)
		if err != nil {
			return fmt.Errorf("operatorconfig: expand %s: %w", key, err)
		}
		v.Set(key, expanded)
	}
	return assertNoUnresolvedPlaceholders(v)
}

func expandValue(val interface{}) (interface{}, error) {
	if val == nil {
		return nil, nil
	}
	switch v := val.(type) {
	case string:
		return expandEnvString(v), nil
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, item := range v {
			e, err := expandValue(item)
			if err != nil {
				return nil, err
			}
			out[i] = e
		}
		return out, nil
	case map[string]interface{}:
		out := make(map[string]interface{}, len(v))
		for k, item := range v {
			e, err := expandValue(item)
			if err != nil {
				return nil, err
			}
			out[k] = e
		}
		return out, nil
	default:
		return val, nil
	}
}

func assertNoUnresolvedPlaceholders(v *viper.Viper) error {
	var bad []string
	for _, key := range v.AllKeys() {
		if err := collectUnresolved(key, v.Get(key), &bad); err != nil {
			return err
		}
	}
	if len(bad) > 0 {
		return fmt.Errorf("operatorconfig: unresolved env placeholders: %s", strings.Join(bad, ", "))
	}
	return nil
}

func collectUnresolved(path string, val interface{}, bad *[]string) error {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case string:
		if strings.Contains(v, "${") {
			*bad = append(*bad, fmt.Sprintf("%s=%q", path, v))
		}
	case []interface{}:
		for i, item := range v {
			if err := collectUnresolved(fmt.Sprintf("%s[%d]", path, i), item, bad); err != nil {
				return err
			}
		}
	case map[string]interface{}:
		for k, item := range v {
			if err := collectUnresolved(path+"."+k, item, bad); err != nil {
				return err
			}
		}
	}
	return nil
}
