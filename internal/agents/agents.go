package agents

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type Agent struct {
	Name        string
	Description string
	SourceFile  string
}

func All() []Agent {
	return []Agent{
		{"pm", "Plan Agent \u2014 guides the user through requirements, presents summaries, delegates all writing to other agents", "agents/pm.md"},
		{"architect", "System Architect \u2014 designs technical solutions and ensures architectural integrity", "agents/architect.md"},
		{"team-lead", "Team Lead \u2014 translates requirements into tickets and manages the execution pipeline", "agents/team-lead.md"},
		{"developer", "Developer \u2014 implements tickets following project standards and creates pull requests", "agents/developer.md"},
		{"reviewer", "Code Reviewer \u2014 reviews PRs for correctness, architecture compliance, code quality, and completeness", "agents/reviewer.md"},
		{"qa", "QA Engineer \u2014 validates implementations meet acceptance criteria and don't break existing functionality", "agents/qa.md"},
		{"doc-writer", "Documentation Writer \u2014 keeps project documentation accurate and up to date", "agents/doc-writer.md"},
	}
}

func SetupAll(dir string, providerNames []string) error {
	for _, name := range providerNames {
		switch name {
		case "opencode":
			if err := setupOpenCode(dir); err != nil {
				return fmt.Errorf("opencode agents: %w", err)
			}
		case "claude-code":
			if err := setupClaudeCode(dir); err != nil {
				return fmt.Errorf("claude-code agents: %w", err)
			}
		case "cursor":
			if err := setupCursor(dir); err != nil {
				return fmt.Errorf("cursor agents: %w", err)
			}
		case "windsurf":
			if err := setupWindsurf(dir); err != nil {
				return fmt.Errorf("windsurf agents: %w", err)
			}
		case "copilot":
			if err := setupCopilot(dir); err != nil {
				return fmt.Errorf("copilot agents: %w", err)
			}
		case "codex-cli":
			if err := setupCodex(dir); err != nil {
				return fmt.Errorf("codex-cli agents: %w", err)
			}
		}
	}
	return nil
}

func setupOpenCode(dir string) error {
	oldDir := filepath.Join(dir, ".opencode", "agents")
	if fi, err := os.Lstat(oldDir); err == nil && fi.IsDir() {
		os.RemoveAll(oldDir)
		fmt.Println("  \u2713 opencode — cleaned up legacy .opencode/agents/")
	}

	agentsDir := filepath.Join(dir, ".opencode", "agent")
	os.MkdirAll(agentsDir, 0755)

	synced := 0
	skipped := 0
	for _, a := range All() {
		sourcePath := filepath.Join(dir, ".ai", a.SourceFile)
		targetPath := filepath.Join(agentsDir, a.Name+".md")

		srcData, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("read %s: %w", a.SourceFile, err)
		}

		existing, err := os.ReadFile(targetPath)
		if err == nil && bytes.Equal(srcData, existing) {
			skipped++
			continue
		}

		if err := os.WriteFile(targetPath, srcData, 0644); err != nil {
			return fmt.Errorf("write %s: %w", a.Name, err)
		}
		synced++
	}

	if synced > 0 {
		fmt.Printf("  \u2713 opencode \u2014 %d agents synced (.opencode/agent/)\n", len(All()))
	} else {
		fmt.Printf("  \u2713 opencode \u2014 %d agents up to date\n", len(All()))
	}
	return nil
}

func setupClaudeCode(dir string) error {
	skillsDir := filepath.Join(dir, ".claude", "skills")
	created := 0

	for _, a := range All() {
		skillDir := filepath.Join(skillsDir, a.Name)
		skillFile := filepath.Join(skillDir, "SKILL.md")

		if _, err := os.Stat(skillFile); err == nil {
			continue
		}

		os.MkdirAll(skillDir, 0755)

		content := fmt.Sprintf(`---
name: %s
description: %s
context: fork
---

Read .ai/%s in full, then follow its instructions.
`, a.Name, a.Description, a.SourceFile)

		if err := os.WriteFile(skillFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s SKILL.md: %w", a.Name, err)
		}
		created++
	}

	if created > 0 {
		fmt.Printf("  \u2713 claude-code \u2014 %d skills registered (.claude/skills/)\n", len(All()))
	} else {
		fmt.Println("  \u2713 claude-code \u2014 skills already registered")
	}
	return nil
}

func setupCursor(dir string) error {
	rulesDir := filepath.Join(dir, ".cursor", "rules")
	os.MkdirAll(rulesDir, 0755)

	created := 0
	for _, a := range All() {
		ruleFile := filepath.Join(rulesDir, a.Name+".mdc")

		if _, err := os.Stat(ruleFile); err == nil {
			continue
		}

		content := fmt.Sprintf(`---
description: AI Dev Team -- %s
alwaysApply: false
---

You are the **%s** agent of the AI Dev Team.
%s

Read .ai/%s for your full instructions.
`, a.Description, a.Name, a.Description, a.SourceFile)

		if err := os.WriteFile(ruleFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s.mdc: %w", a.Name, err)
		}
		created++
	}

	if created > 0 {
		fmt.Printf("  \u2713 cursor \u2014 %d rules registered (.cursor/rules/)\n", len(All()))
	} else {
		fmt.Println("  \u2713 cursor \u2014 rules already registered")
	}
	return nil
}

func setupWindsurf(dir string) error {
	rulesDir := filepath.Join(dir, ".windsurf", "rules")
	os.MkdirAll(rulesDir, 0755)

	created := 0
	for _, a := range All() {
		ruleFile := filepath.Join(rulesDir, a.Name+".md")

		if _, err := os.Stat(ruleFile); err == nil {
			continue
		}

		content := fmt.Sprintf(`---
trigger: manual
description: AI Dev Team -- %s
---

You are the **%s** agent of the AI Dev Team.
%s

Read .ai/%s for your full instructions.
`, a.Description, a.Name, a.Description, a.SourceFile)

		if err := os.WriteFile(ruleFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s.md: %w", a.Name, err)
		}
		created++
	}

	if created > 0 {
		fmt.Printf("  \u2713 windsurf \u2014 %d rules registered (.windsurf/rules/)\n", len(All()))
	} else {
		fmt.Println("  \u2713 windsurf \u2014 rules already registered")
	}
	return nil
}

func setupCopilot(dir string) error {
	instructionsDir := filepath.Join(dir, ".github", "instructions")
	os.MkdirAll(instructionsDir, 0755)

	created := 0
	for _, a := range All() {
		instFile := filepath.Join(instructionsDir, a.Name+".instructions.md")

		if _, err := os.Stat(instFile); err == nil {
			continue
		}

		content := fmt.Sprintf(`---
applyTo: "**"
---

You are the **%s** agent of the AI Dev Team.
%s

Read .ai/%s for your full instructions.
`, a.Name, a.Description, a.SourceFile)

		if err := os.WriteFile(instFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s.instructions.md: %w", a.Name, err)
		}
		created++
	}

	if created > 0 {
		fmt.Printf("  \u2713 copilot \u2014 %d instructions registered (.github/instructions/)\n", len(All()))
	} else {
		fmt.Println("  \u2713 copilot \u2014 instructions already registered")
	}
	return nil
}

func setupCodex(dir string) error {
	skillsDir := filepath.Join(dir, ".agents", "skills")
	created := 0

	for _, a := range All() {
		skillDir := filepath.Join(skillsDir, a.Name)
		skillFile := filepath.Join(skillDir, "SKILL.md")

		if _, err := os.Stat(skillFile); err == nil {
			continue
		}

		os.MkdirAll(skillDir, 0755)

		content := fmt.Sprintf(`---
name: %s
description: %s
---

Read .ai/%s in full, then follow its instructions.
`, a.Name, a.Description, a.SourceFile)

		if err := os.WriteFile(skillFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s SKILL.md: %w", a.Name, err)
		}
		created++
	}

	if created > 0 {
		fmt.Printf("  \u2713 codex-cli \u2014 %d skills registered (.agents/skills/)\n", len(All()))
	} else {
		fmt.Println("  \u2713 codex-cli \u2014 skills already registered")
	}
	return nil
}
