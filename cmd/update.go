package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var devteamRepo = "https://github.com/QODESQUARE/ai-dev-team.git"

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update AI Dev Team from git",
		Long:  `Checks for updates and rebuilds from source if a newer version is available.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate()
		},
	}
}

func runUpdate() error {
	dir, err := findTeamDir()
	if err != nil {
		return fmt.Errorf("AI Dev Team not found. Run 'devteam install' first")
	}

	origin := devteamRepo
	if data, err := os.ReadFile(filepath.Join(dir, ".git-origin")); err == nil {
		if s := strings.TrimSpace(string(data)); s != "" {
			origin = s
		}
	}

	remoteHash, err := getRemoteHEAD(origin)
	if err != nil {
		fmt.Printf("  Warning: could not check remote (%s)\n", err)
	} else if remoteHash == commitHash {
		fmt.Println("  Already up to date")
		return nil
	} else if commitHash == "dev" {
		fmt.Println("  Development build — updating...")
	} else {
		fmt.Printf("  New commits available (%.7s... -> %.7s...)\n", commitHash, remoteHash)
	}

	fmt.Printf("  Cloning from %s\n", origin)

	tmpDir, err := os.MkdirTemp("", "devteam-update")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	clone := exec.Command("git", "clone", "--depth", "1", origin, tmpDir)
	clone.Stdout = os.Stdout
	clone.Stderr = os.Stderr
	if err := clone.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	fmt.Println("  Building binary...")
	build := exec.Command("go", "build",
		"-ldflags", "-X devteam/cmd.commitHash="+remoteHash,
		"-o", filepath.Join(tmpDir, "devteam"), ".")
	build.Dir = tmpDir
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	exe, err := os.ReadFile(filepath.Join(tmpDir, "devteam"))
	if err != nil {
		return fmt.Errorf("read built binary: %w", err)
	}
	binDest := filepath.Join(dir, "devteam")
	if err := os.WriteFile(binDest, exe, 0755); err != nil {
		return fmt.Errorf("write binary: %w", err)
	}
	os.WriteFile(filepath.Join(dir, ".git-commit"), []byte(remoteHash+"\n"), 0644)
	fmt.Println("  ✓ devteam binary updated")

	aiDir := filepath.Join(tmpDir, ".ai")
	if _, err := os.Stat(aiDir); err == nil {
		installDir := filepath.Join(dir, ".ai")
		os.RemoveAll(installDir)
		if err := copyDir(aiDir, installDir); err != nil {
			return fmt.Errorf("copy .ai: %w", err)
		}
		fmt.Println("  ✓ .ai/ resources updated")
	}

	fmt.Println()
	fmt.Println("=== Update complete ===")
	fmt.Println("  Run 'devteam activate' in each project to refresh symlinks.")
	return nil
}

func getRemoteHEAD(origin string) (string, error) {
	cmd := exec.Command("git", "ls-remote", origin, "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git ls-remote: %w", err)
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return "", fmt.Errorf("empty response from remote")
	}
	return fields[0], nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}
