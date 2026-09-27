package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/behaviorengineering/operatorconfig/pkg/operatorconfig"
)

const version = "0.1.0"

func main() {
	if len(os.Args) == 1 {
		printAgentGuide()
		os.Exit(0)
	}

	var (
		appName    string
		outPath    string
		configPath string
		dryRun     bool
		showHelp   bool
		showVer    bool
		secretCSV  string
	)
	fs := flag.NewFlagSet("export-env", flag.ExitOnError)
	fs.StringVar(&appName, "app", "", "XDG app name and keyring service (required)")
	fs.StringVar(&outPath, "out", "", "Write resolved secrets as dotenv (mode 0600)")
	fs.StringVar(&configPath, "config", "", "Load secrets: list from this YAML (optional)")
	fs.StringVar(&secretCSV, "secrets", "", "Comma-separated env names when --config is omitted")
	fs.BoolVar(&dryRun, "dry-run", false, "Resolve secrets only; do not write --out or exec")
	fs.BoolVar(&showHelp, "help", false, "Show flag help")
	fs.BoolVar(&showVer, "version", false, "Print version")
	_ = fs.Parse(os.Args[1:])

	if showVer {
		fmt.Println(version)
		return
	}
	if showHelp || fs.NArg() > 0 && fs.Args()[0] == "help" {
		fs.Usage()
		return
	}

	if strings.TrimSpace(appName) == "" {
		fmt.Fprintln(os.Stderr, "export-env: --app is required (or run with no args for agent guide)")
		os.Exit(2)
	}

	secrets, err := resolveSecretList(appName, configPath, secretCSV)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(secrets) == 0 {
		fmt.Fprintln(os.Stderr, "export-env: no secrets to resolve (use --config or --secrets)")
		os.Exit(2)
	}

	opts := operatorconfig.Options{
		App:     appName,
		Secrets: secrets,
	}
	if err := operatorconfig.ResolveSecrets(opts, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if dryRun {
		fmt.Println("export-env: dry-run ok (secrets resolved in process env)")
		return
	}

	if out := strings.TrimSpace(outPath); out != "" {
		if err := operatorconfig.WriteDotenv(out, secrets); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("export-env: wrote", out)
	}
}

func resolveSecretList(app, configPath, secretCSV string) ([]operatorconfig.Secret, error) {
	if p := strings.TrimSpace(configPath); p != "" {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("export-env: read config: %w", err)
		}
		return operatorconfig.ParseSecretsYAML(data)
	}
	if secretCSV == "" {
		return nil, nil
	}
	var out []operatorconfig.Secret
	for _, part := range strings.Split(secretCSV, ",") {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		out = append(out, operatorconfig.Secret{Env: name, Required: true})
	}
	return out, nil
}

func printAgentGuide() {
	fmt.Printf(`export-env %s — operatorconfig secret resolver for Compose and CI

Purpose: Resolve declared secrets (env → keyring → optional SOPS file) and write a dotenv or leave env set for a child process.

Agent harness:
  AGENTS.md and ai-copilots/ live in the operatorconfig module root.
  Load ai-copilots/skills/operatorconfig/SKILL.md before inventing Keychain or sops shell steps.

Inspect (safe):
  export-env --app <app> --secrets NAME --dry-run
  export-env --config /path/to/config.yaml --app <app> --dry-run

Execute:
  export-env --app <app> --secrets CF_AI_API_KEY,CF_ACCOUNT_ID --out .env.deploy
  docker compose --env-file .env.deploy -f docker-compose.deploy.yml up -d

Human flags: export-env --help

`, version)
}
