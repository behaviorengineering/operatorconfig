# BOOTSTRAP — operatorconfig ai-copilots

**Audience:** Any AI agent in a host repo that depends on this Go module.

## Resolve module root

```bash
MODULE_ROOT="$(go list -m -f '{{.Dir}}' github.com/behaviorengineering/operatorconfig)"
test -f "$MODULE_ROOT/ai-copilots/BOOTSTRAP.md"
```

Honors `replace`, `go.work`, module cache, and vendor.

## Wire mode (symlink skills)

Ask IDE and OS if unknown. Prefer symlink on macOS/Linux; junction on Windows.

**Cursor (example):**

```bash
ln -snf "$MODULE_ROOT/ai-copilots/skills/operatorconfig" .cursor/skills/operatorconfig-lib
ln -snf "$MODULE_ROOT/ai-copilots/skills/export-env-operator" .cursor/skills/export-env-operator
```

**MUST NOT** copy skill bodies into the host unless links fail and the human approves.

After a module version bump, re-run wire (cache Dir changes).

## Load order

1. [ai-copilots/README.md](README.md)
2. [skills/operatorconfig/SKILL.md](skills/operatorconfig/SKILL.md)
3. [skills/export-env-operator/SKILL.md](skills/export-env-operator/SKILL.md) when resolving secrets for Compose or CI

Cross-product constraints MAY also exist in the host's shared operator-config pack skill; this module skill is canonical for resolve order and API.
