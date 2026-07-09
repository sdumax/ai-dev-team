package main

import (
	"embed"
	"os"

	"devteam/cmd"
	"github.com/spf13/cobra"
)

//go:embed .ai
var aiFS embed.FS

func main() {
	root := &cobra.Command{
		Use:   "devteam",
		Short: "AI Dev Team — Multi-agent development workflow",
		Long: `A portable multi-agent development team that ships features on your command.

Supports: opencode, Claude Code, Cursor, Copilot, Windsurf, Codex CLI`,
	}

	root.AddCommand(cmd.NewInstallCmd(aiFS))
	root.AddCommand(cmd.NewActivateCmd())
	root.AddCommand(cmd.NewUpdateCmd())
	root.AddCommand(cmd.NewStatusCmd())
	root.AddCommand(cmd.NewListCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
