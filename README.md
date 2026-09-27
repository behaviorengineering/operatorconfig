# operatorconfig

Shared operator configuration for Go CLIs and daemons:

- XDG user config discovery (`~/.config/<app>/config.yaml`)
- Optional Viper load with `${VAR}` expansion in string fields
- Secret resolution: process environment, then [zalando/go-keyring](https://github.com/zalando/go-keyring), then an optional SOPS file (`secrets.enc.yaml` under the app config dir by default). Values are single-line; `SanitizeSecret` trims whitespace on keyring store and rewrites dirty process env on resolve.

Containers and CI should inject secrets via the environment on the host before `docker compose`; there is no host keyring inside typical Docker images.

## Module

```text
github.com/behaviorengineering/operatorconfig/pkg/operatorconfig
```

## Quick start

```go
opts := operatorconfig.Options{
    App:       "myapp",
    ConfigEnv: "MYAPP_CONFIG",
    Secrets: []operatorconfig.Secret{
        {Env: "API_TOKEN", Required: true},
    },
}
path, err := operatorconfig.ResolveConfigPath(opts)
_ = operatorconfig.ResolveSecrets(opts, operatorconfig.DefaultKeyring())
```

## export-env CLI

```bash
go run ./cmd/export-env --app myapp --secrets API_TOKEN --dry-run
go run ./cmd/export-env --app myapp --secrets API_TOKEN --out .env.deploy
```

Bare `export-env` prints the agent operating guide.

See `ai-copilots/skills/operatorconfig/SKILL.md` for operator guidance.
