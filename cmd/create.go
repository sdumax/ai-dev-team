package cmd

import (
	"fmt"
	"path/filepath"

	"devteam/internal/agents"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [dir]",
		Short: "Create opencode agents and AGENTS.md",
		Long: `Install devteam agents globally to ~/.config/opencode/agent/
and create or update AGENTS.md in the project with delegation rules.`,
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
			return runCreate(abs)
		},
	}
}

func runCreate(dir string) error {
	fmt.Println("=== AI Dev Team — Create ===")
	fmt.Printf("  Project: %s\n\n", dir)

	if err := agents.CreateOpenCodeGlobal(); err != nil {
		return fmt.Errorf("install agents: %w", err)
	}

	fmt.Println()
	if err := agents.UpdateAGENTSMD(dir); err != nil {
		return fmt.Errorf("update AGENTS.md: %w", err)
	}

	fmt.Println()
	fmt.Println("=== Create complete ===")
	fmt.Println("  Agents: ~/.config/opencode/agent/")
	fmt.Printf("  Config: %s/AGENTS.md\n", dir)
	return nil
}
