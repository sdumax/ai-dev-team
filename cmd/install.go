package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"devteam/internal/agents"
	"github.com/spf13/cobra"
)

func NewInstallCmd(aiFS fs.FS) *cobra.Command {
	return &cobra.Command{
		Use:   "install [path]",
		Short: "Install AI Dev Team globally",
		Long: `Install the AI Dev Team to ~/.ai-dev-team (or a custom path).
Extracts embedded team resources and adds 'devteam' to your PATH.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir := filepath.Join(os.Getenv("HOME"), ".ai-dev-team")
			if len(args) > 0 {
				targetDir = args[0]
			}
			return runInstall(aiFS, targetDir)
		},
	}
}

func runInstall(srcFS fs.FS, target string) error {
	fmt.Println("=== AI Dev Team Install ===")

	os.MkdirAll(target, 0755)

	if origin, err := getGitOrigin(); err == nil {
		os.WriteFile(filepath.Join(target, ".git-origin"), []byte(origin+"\n"), 0644)
	}
	if commitHash != "dev" {
		os.WriteFile(filepath.Join(target, ".git-commit"), []byte(commitHash+"\n"), 0644)
	}

	if err := extractFS(srcFS, ".ai", filepath.Join(target, ".ai")); err != nil {
		return fmt.Errorf("extract .ai: %w", err)
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable: %w", err)
	}
	src, err := os.ReadFile(exe)
	if err != nil {
		return fmt.Errorf("read binary: %w", err)
	}
	binDest := filepath.Join(target, "devteam")
	if err := os.WriteFile(binDest, src, 0755); err != nil {
		return fmt.Errorf("write binary: %w", err)
	}

	os.Chmod(filepath.Join(target, ".ai", "init.sh"), 0755)

	if err := setupOpenCodeGlobal(target); err != nil {
		return fmt.Errorf("setup opencode agents: %w", err)
	}

	addToPath(target)

	fmt.Printf("  Installed: %s\n", target)
	fmt.Println("  Restart your terminal or run: source ~/.zshrc (or ~/.bashrc)")
	fmt.Println()
	fmt.Println("  Then in any project: devteam activate")
	return nil
}

func extractFS(srcFS fs.FS, prefix, dst string) error {
	return fs.WalkDir(srcFS, prefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, prefix)
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := fs.ReadFile(srcFS, path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}

func getGitOrigin() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func setupOpenCodeGlobal(target string) error {
	opencodeDir := filepath.Join(target, ".opencode", "agent")
	os.MkdirAll(opencodeDir, 0755)

	written := 0
	for _, a := range agents.All() {
		src := filepath.Join(target, ".ai", a.SourceFile)
		dst := filepath.Join(opencodeDir, a.Name+".md")

		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read %s: %w", a.SourceFile, err)
		}

		if err := os.WriteFile(dst, data, 0644); err != nil {
			return fmt.Errorf("write %s: %w", a.Name, err)
		}
		written++
	}

	fmt.Printf("  \u2713 opencode \u2014 %d agents registered globally\n", written)
	return nil
}

func addToPath(dir string) {
	rc := filepath.Join(os.Getenv("HOME"), ".bashrc")
	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".zshrc")); err == nil {
		rc = filepath.Join(os.Getenv("HOME"), ".zshrc")
	}
	opencodeDir := filepath.Join(dir, ".opencode")
	line := fmt.Sprintf("\n# AI Dev Team\nexport PATH=\"$PATH:%s\"\nexport OPENCODE_CONFIG_DIR=\"%s\"\n", dir, opencodeDir)

	data, err := os.ReadFile(rc)
	if err != nil {
		fmt.Printf("  Warning: could not read %s\n", rc)
		return
	}
	if strings.Contains(string(data), dir) {
		fmt.Println("  Already in PATH")
		return
	}

	f, err := os.OpenFile(rc, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("  Warning: could not write %s\n", rc)
		return
	}
	defer f.Close()
	f.WriteString(line)
	fmt.Printf("  Added to PATH in %s\n", rc)
}
