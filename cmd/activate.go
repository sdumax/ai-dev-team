package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"devteam/internal/agents"
	"devteam/internal/detect"
	"github.com/spf13/cobra"
)

func NewActivateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "activate [dir]",
		Short: "Activate AI Dev Team in a project",
		Long: `Detects AI tools in the project and configures them for /ship.
Creates .ai/ symlinks pointing to the global installation.
Safe to run multiple times — detects existing state.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			abs, err := filepath.Abs(dir)
			if err != nil {
				return err
			}
			return runActivate(abs)
		},
	}
}

func runActivate(dir string) error {
	fmt.Println("=== AI Dev Team — Activate ===")
	fmt.Printf("  Project: %s\n", dir)

	selected := resolve_providers(dir)

	if len(selected) == 0 {
		fmt.Println("  No providers selected.")
		return nil
	}

	if err := link_team(dir); err != nil {
		return fmt.Errorf("symlink .ai/: %w", err)
	}

	if err := detect.Configure(dir, selected); err != nil {
		return fmt.Errorf("configure providers: %w", err)
	}

	names := make([]string, len(selected))
	for i, p := range selected {
		names[i] = p.Name
	}
	if err := agents.SetupAll(dir, names); err != nil {
		return fmt.Errorf("setup agents: %w", err)
	}

	update_gitignore(dir)

	if err := agents.CreateOpenCodeGlobal(); err != nil {
		return fmt.Errorf("install agents: %w", err)
	}
	if err := agents.UpdateAGENTSMD(dir); err != nil {
		return fmt.Errorf("update AGENTS.md: %w", err)
	}

	fmt.Println()
	fmt.Println("=== Activation complete ===")
	fmt.Printf("  Project:   %s\n", dir)
	fmt.Printf("  Providers: %s\n", strings.Join(names, ", "))
	fmt.Printf("  .ai/     -> global installation\n")
	fmt.Println()
	fmt.Println("  Use /ship \"your feature\" to start.")
	return nil
}

func resolve_providers(dir string) []detect.Provider {
	all := detect.All()
	var found []detect.Provider
	for _, p := range all {
		if p.Detect(dir) {
			found = append(found, p)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	switch len(found) {
	case 0:
		fmt.Println("\n  No AI tool config files found.")
		fmt.Println("  Which provider(s) are you using?")
		for i, p := range all {
			fmt.Printf("    %d) %s\n", i+1, p.Name)
		}
		fmt.Printf("  Enter numbers (comma-separated, or 0 to skip): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" || text == "0" {
			return nil
		}
		var result []detect.Provider
		for _, s := range strings.Split(text, ",") {
			s = strings.TrimSpace(s)
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > len(all) {
				continue
			}
			result = append(result, all[n-1])
		}
		return result

	case 1:
		fmt.Printf("  Detected %s (%s)\n", found[0].Name, found[0].ConfigFile)
		fmt.Printf("  Configure /ship for %s? [Y/n]: ", found[0].Name)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" || strings.EqualFold(text, "y") || strings.EqualFold(text, "yes") {
			return found
		}
		return nil

	default:
		fmt.Println("  Multiple tools detected:")
		for i, p := range found {
			fmt.Printf("    %d) %s (%s)\n", i+1, p.Name, p.ConfigFile)
		}
		fmt.Printf("  Configure which? (numbers comma-separated, or 'all'): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if strings.EqualFold(text, "all") {
			return found
		}
		var result []detect.Provider
		for _, s := range strings.Split(text, ",") {
			s = strings.TrimSpace(s)
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > len(found) {
				continue
			}
			result = append(result, found[n-1])
		}
		return result
	}
}

func link_team(dir string) error {
	teamDir, err := findTeamDir()
	if err != nil {
		return err
	}
	aiSource := filepath.Join(teamDir, ".ai")
	aiTarget := filepath.Join(dir, ".ai")

	os.MkdirAll(filepath.Join(aiTarget, "tickets", "shipments"), 0755)

	type link struct {
		Name string
		Src  string
	}
	links := []link{
		{"agents", filepath.Join(aiSource, "agents")},
		{"templates", filepath.Join(aiSource, "templates")},
		{"ship.md", filepath.Join(aiSource, "ship.md")},
		{"init.sh", filepath.Join(aiSource, "init.sh")},
	}

	for _, l := range links {
		target := filepath.Join(aiTarget, l.Name)
		fi, err := os.Lstat(target)
		if err == nil && fi.Mode()&os.ModeSymlink != 0 {
			fmt.Printf("  ✓ %s already linked\n", l.Name)
			continue
		}
		os.RemoveAll(target)
		if err := os.Symlink(l.Src, target); err != nil {
			return fmt.Errorf("symlink %s: %w", l.Name, err)
		}
		fmt.Printf("  ✓ %s → symlink\n", l.Name)
	}

	return nil
}

func findTeamDir() (string, error) {
	exe, err := os.Executable()
	if err == nil {
		candidate := filepath.Dir(exe)
		if _, err := os.Stat(filepath.Join(candidate, ".ai")); err == nil {
			return candidate, nil
		}
	}
	candidate := "."
	if _, err := os.Stat(filepath.Join(candidate, ".ai")); err == nil {
		abs, _ := filepath.Abs(candidate)
		return abs, nil
	}
	return "", fmt.Errorf(".ai/ not found. Run 'devteam install' first")
}

func update_gitignore(dir string) {
	entries := []string{".ai/"}
	gf := filepath.Join(dir, ".gitignore")

	data, err := os.ReadFile(gf)
	if err != nil {
		var out strings.Builder
		for _, e := range entries {
			out.WriteString(e + "\n")
		}
		os.WriteFile(gf, []byte(out.String()), 0644)
		fmt.Println("  ✓ .gitignore — created")
		return
	}

	content := string(data)
	var missing []string
	for _, e := range entries {
		if strings.Contains(content, e) {
			fmt.Printf("  ✓ .gitignore — %s already ignored\n", e)
		} else {
			missing = append(missing, e)
		}
	}

	if len(missing) == 0 {
		return
	}

	f, _ := os.OpenFile(gf, os.O_APPEND|os.O_WRONLY, 0644)
	if f != nil {
		defer f.Close()
		for _, e := range missing {
			f.WriteString("\n" + e)
		}
		f.WriteString("\n")
		fmt.Printf("  ✓ .gitignore — added %s\n", strings.Join(missing, ", "))
	}
}
