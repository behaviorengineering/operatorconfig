---
name: operatorconfig
description: >-
  Operate github.com/behaviorengineering/operatorconfig: XDG config discovery,
  Viper YAML load, env then keyring then optional SOPS secrets, YAML secrets
  list, Options.Keyring and Options.SecretFile injection. Use when wiring Load,
  ResolveConfigPath, ResolveSecrets, InitUserConfig, or deploy export-env.
---

# operatorconfig library

**Moral:** User config lives under `~/.config/<app>/`. Path discovery is separate from secret resolution. Only env names under `secrets:` are filled from the keyring or SOPS file. Environment wins. Containers and CI inject env only.

## Public API

Package: `github.com/behaviorengineering/operatorconfig/pkg/operatorconfig`

- `ResolveConfigPath(Options)` — override flag/env, then XDG file, then `ExtraPaths`
- `Load(Options)` — read YAML when found, expand `${VAR}`, parse top-level `secrets:` when `Options.Secrets` is empty, then `ResolveSecrets`
- `ResolveSecrets(Options, Keyring)` — env, then keyring, then optional SOPS file; `Required` fails closed
- `SecretsEncPath(Options)` — explicit `SecretsEncPath` or default `~/.config/<app>/secrets.enc.yaml`
- `WriteDotenv(path, secrets)` — after resolve, write mode `0600` dotenv for Compose
- `InitUserConfig(Options, example []byte, force)` — `0700` dir, `0600` live file, refresh `.example`
- `ParseSecretsYAML` / `Secret.UnmarshalYAML` — scalar (`- TOKEN`) or mapping (`env` + `required`)
- `Options.Keyring` — injectable store; nil uses `DefaultKeyring()` (`OSKeyring`)
- `Options.SecretFile` — injectable SOPS loader; nil uses `SOPSSecretFile`; tests use `MemSecretFile`
- `NewMemKeyring()` for tests; `SetTestKeyring` is last-resort process-global swap

CLI: `cmd/export-env` — see [export-env-operator/SKILL.md](../export-env-operator/SKILL.md).

## Secret sanitize (single-line)

**CONSTRAINT:** Secret values MUST be single-line. `SanitizeSecret` (trim space) runs on keyring `Set`, on values loaded from keyring/SOPS, and on process env during `ResolveSecrets` (dirty env is rewritten to the trimmed value). Multiline secrets are not supported.

## Secret resolve order (locked)

**CONSTRAINT:** `ResolveSecrets` MUST try, in order: process env, keyring for declared names, optional SOPS file when present. MUST NOT add SOPS to `ResolveConfigPath`.

- Enforcement: Read `resolveSecrets` in `secrets.go`
- Violation: STOP, keep SOPS on hop 3 only

CORRECT:
```text
env set → skip keyring and SOPS for that name
env empty → keyring Get(service=App, account=env)
Get ErrNotFound → decrypt SecretsEncPath when file exists
Get any other error → fail (even if required: false)
still empty and required → error
```

PROHIBITED:
```text
ResolveConfigPath tries secrets.enc.yaml before config.yaml
deploy.ps1 calls sops -d directly
```

SOPS age key: honor `SOPS_AGE_KEY` / `SOPS_AGE_KEY_FILE`, or keyring account `SOPS_AGE_KEY` under the same service before decrypt.

## YAML `secrets:`

```yaml
secrets:
  - API_TOKEN
  - env: ACCOUNT_ID
    required: false
```

**CONSTRAINT:** Hosts MUST resolve only secrets listed in config (or passed in `Options.Secrets`). MUST skip keyring and SOPS when the list is empty.

- Enforcement: Trace `Load` / host decode; empty `secrets:` does not call `Keyring.Get` or load SOPS
- Violation: STOP, bind secrets from YAML (or Options), re-verify

## Keyring

Service id = `Options.App`. Account = secret env name.

**CONSTRAINT:** Tests MUST pass `Options.Keyring` and `Options.SecretFile`. MUST NOT rely on `SetTestKeyring` as the primary DI seam.

- Enforcement: Unit tests inject `NewMemKeyring` and `MemSecretFile`
- Violation: STOP, add injection seams, re-verify

## Host integration

Hosts with custom YAML decoders (`KnownFields`, custom parsers) MUST:

1. `ResolveConfigPath`
2. Decode the file (including `secrets: []Secret`)
3. `ResolveSecrets` with those entries and `Options.Keyring`
4. Then expand `${VAR}` / validate backends

They MUST NOT Unmarshal the whole host config through Viper when they need KnownFields.

**CI and containers:** set process env on the host before `docker compose`; do not expect a keyring inside the image. Container-mounted config templates SHOULD omit `secrets:` (env injection only); host `serve` config MAY list `secrets:` for keyring-backed dev.
