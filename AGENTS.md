# Agents

This module is a Go library for operator config discovery and secret resolution. Humans read [README.md](README.md).

**Load before wiring config paths or secrets:**

1. [ai-copilots/README.md](ai-copilots/README.md)
2. [ai-copilots/skills/operatorconfig/SKILL.md](ai-copilots/skills/operatorconfig/SKILL.md)
3. [ai-copilots/skills/export-env-operator/SKILL.md](ai-copilots/skills/export-env-operator/SKILL.md) (Compose / CI dotenv export)

## Wire host discovery

Execute [ai-copilots/BOOTSTRAP.md](ai-copilots/BOOTSTRAP.md) in wire mode.

```bash
go list -m -f '{{.Dir}}' github.com/behaviorengineering/operatorconfig
```
