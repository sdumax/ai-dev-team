package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update AI Dev Team from git",
		Long:  `Pulls the latest version from git in the global installation directory.`,
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

	if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
		return fmt.Errorf("not a git repository at %s", dir)
	}

	cmd := exec.Command("git", "-C", dir, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	fmt.Println()
	fmt.Println("=== Update complete ===")
	fmt.Println("  Run 'devteam activate' in each project to refresh symlinks.")
	return nil
}
