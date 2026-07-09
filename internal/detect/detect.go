package detect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Provider struct {
	Name       string
	ConfigFile string
	Detect     func(dir string) bool
	Block      string
}

func All() []Provider {
	return []Provider{
		{
			Name:       "opencode",
			ConfigFile: "AGENTS.md",
			Detect:     hasFile("AGENTS.md"),
			Block:      opencodeBlock(),
		},
		{
			Name:       "claude-code",
			ConfigFile: "CLAUDE.md",
			Detect:     hasFile("CLAUDE.md"),
			Block:      simpleBlock("CLAUDE.md"),
		},
		{
			Name:       "cursor",
			ConfigFile: ".cursorrules",
			Detect:     hasFile(".cursorrules"),
			Block:      simpleBlock(".cursorrules"),
		},
		{
			Name:       "windsurf",
			ConfigFile: ".windsurfrules",
			Detect:     hasFile(".windsurfrules"),
			Block:      simpleBlock(".windsurfrules"),
		},
		{
			Name:       "copilot",
			ConfigFile: ".github/copilot-instructions.md",
			Detect:     hasFile(".github/copilot-instructions.md"),
			Block:      simpleBlock(".github/copilot-instructions.md"),
		},
		{
			Name:       "codex-cli",
			ConfigFile: "CONVENTIONS.md",
			Detect:     hasFile("CONVENTIONS.md"),
			Block:      simpleBlock("CONVENTIONS.md"),
		},
	}
}

func hasFile(path string) func(dir string) bool {
	return func(dir string) bool {
		_, err := os.Stat(filepath.Join(dir, path))
		return err == nil
	}
}

func Configure(dir string, providers []Provider) error {
	fmt.Println()
	fmt.Println("Configuring providers")

	for _, p := range providers {
		fpath := filepath.Join(dir, p.ConfigFile)

		if p.ConfigFile == ".github/copilot-instructions.md" {
			os.MkdirAll(filepath.Join(dir, ".github"), 0755)
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			os.WriteFile(fpath, []byte(p.Block), 0644)
			fmt.Printf("  ✓ %s — created\n", p.ConfigFile)
			continue
		}

		if strings.Contains(string(data), "/ship") || strings.Contains(string(data), "ai-dev-team") {
			fmt.Printf("  ✓ %s — already configured\n", p.ConfigFile)
			continue
		}

		f, err := os.OpenFile(fpath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("open %s: %w", p.ConfigFile, err)
		}
		f.WriteString(p.Block)
		f.Close()
		fmt.Printf("  ✓ %s — added /ship command\n", p.ConfigFile)
	}

	return nil
}

func opencodeBlock() string {
	b := "`"
	return "\n\n## Slash Commands\n\n### /ship\n" +
		"When the user types " + b + "/ship <requirement>" + b + ":\n" +
		"1. Read " + b + ".ai/ship.md" + b + " in full.\n" +
		"2. Execute the Ship Workflow defined there.\n" +
		"3. All phases (PM \u2192 Architect \u2192 Team Lead \u2192 Execution \u2192 Review \u2192 QA \u2192 Merge) run in sequence.\n\n" +
		"This is the multi-agent development workflow. It creates tickets, assigns work to sub-agents, reviews, QAs, and prepares PRs for merge.\n"
}

func simpleBlock(file string) string {
	return fmt.Sprintf("\n# AI Dev Team\nUse /ship to start a multi-agent development workflow.\nSee .ai/ship.md for details.\n")
}
