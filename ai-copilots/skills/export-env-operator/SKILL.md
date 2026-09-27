---
name: export-env-operator
description: >-
  Operate operatorconfig export-env: resolve secrets (env, keyring, optional SOPS)
  and write a dotenv for Docker Compose. Use for homelab deploy scripts, not for
  in-container secret storage.
---

# export-env operator

**Moral:** Inspect with `--dry-run`, then write `--out` or rely on process env. MUST NOT shell out to `sops`, `security`, or Credential Manager from deploy repos.

## Binary

Built from `cmd/export-env` in the operatorconfig module.

Bare invoke (no args) prints the agent guide and exits 0. Human flags: `export-env --help`.

## Order (locked)

1. Process environment
2. OS keyring (`Options.App` service, account = env name). On Windows, `OSKeyring.Get` decodes Credential Manager blobs stored as UTF-16 LE (Control Panel / cmdkey).
3. Optional SOPS file (`SecretsEncPath` or `~/.config/<app>/secrets.enc.yaml`)

Declared `secrets:` only. Required secrets fail closed.

## Inspect

```bash
export-env --app <app> --secrets API_TOKEN,ACCOUNT_ID --dry-run
export-env --app <app> --config ~/.config/<app>/config.yaml --dry-run
```

**CONSTRAINT:** MUST use `--dry-run` before first `--out` on a new machine.

- Enforcement: Read deploy script; dry-run precedes compose
- Violation: STOP, add dry-run step

## Execute

```bash
export-env --app polypus --secrets CF_AI_API_KEY,CF_ACCOUNT_ID --out .env.deploy
docker compose --env-file .env.deploy -f docker-compose.deploy.yml up -d
```

**CONSTRAINT:** MUST NOT print secret values in logs.

- Enforcement: Grep scripts for echo of `CF_` or `TOKEN`
- Violation: STOP, remove prints

## CI and containers

Inject env only inside images. Run `export-env` on the host runner before `docker compose`.
