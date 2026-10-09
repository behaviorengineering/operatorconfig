package operatorconfig

// Options configures config discovery, Viper env binding, and secret resolution.
type Options struct {
	App            string   // XDG dir name and keyring service, e.g. "polypus"
	ConfigEnv      string   // override path env, e.g. "POLYPUS_CONFIG"
	ConfigFlagPath string   // from --config; wins over ConfigEnv if non-empty
	Filename       string   // default "config.yaml"
	ExtraPaths     []string // checked after XDG if file missing
	EnvPrefix      string   // Viper SetEnvPrefix; empty = no prefix
	Secrets        []Secret
	SecretsEncPath string     // optional SOPS file; default ~/.config/<app>/secrets.enc.yaml
	Keyring        Keyring      // when nil, ResolveSecrets uses DefaultKeyring()
	SecretFile     SecretFile   // when nil, ResolveSecrets uses SOPSSecretFile for hop 3
	EnvDefaults    []EnvDefault // non-secret env fallbacks after ResolveSecrets, before ${VAR} expand
}

// Secret names a process env var and optional keyring account (same name).
type Secret struct {
	Env      string
	Required bool
}

func (o Options) filename() string {
	if o.Filename != "" {
		return o.Filename
	}
	return "config.yaml"
}
