package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [ship-NNN]",
		Short: "Show shipment status",
		Long:  `Shows status of AI Dev Team shipments. With no argument, lists all shipments.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			dir := findProjectRoot()
			if dir == "" {
				return fmt.Errorf("no .ai/ found in current or parent directories.\nRun 'devteam activate' first.")
			}
			if len(args) > 0 {
				return showShipment(dir, args[0])
			}
			return listShipments(dir)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all shipments",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			dir := findProjectRoot()
			if dir == "" {
				return fmt.Errorf("no .ai/ found.\nRun 'devteam activate' first.")
			}
			return listShipments(dir)
		},
	}
	cmd.AddCommand(listCmd)

	return cmd
}

func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	for dir := wd; dir != "/"; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, ".ai")); err == nil {
			return dir
		}
	}
	return ""
}

func shipmentsDir(projectDir string) string {
	return filepath.Join(projectDir, ".ai", "tickets", "shipments")
}

func listShipments(projectDir string) error {
	sd := shipmentsDir(projectDir)
	entries, err := os.ReadDir(sd)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("  No shipments found.")
			return nil
		}
		return fmt.Errorf("read shipments: %w", err)
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "ship-") {
			dirs = append(dirs, e.Name())
		}
	}

	if len(dirs) == 0 {
		fmt.Println("  No shipments found.")
		return nil
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i] < dirs[j]
	})

	fmt.Println("  Shipments:")
	for _, d := range dirs {
		rawStatus, title := readManifestSummary(filepath.Join(sd, d, "manifest.md"))
		status := formatStatus(rawStatus)
		fmt.Printf("    %-12s %s %s\n", d, status, title)
	}
	return nil
}

func showShipment(projectDir string, shipment string) error {
	sd := shipmentsDir(projectDir)

	if !strings.HasPrefix(shipment, "ship-") {
		shipment = "ship-" + shipment
	}

	shipDir := filepath.Join(sd, shipment)
	if _, err := os.Stat(shipDir); os.IsNotExist(err) {
		return fmt.Errorf("shipment %s not found", shipment)
	}

	manifestPath := filepath.Join(shipDir, "manifest.md")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}

	manifest := string(manifestData)
	status, title := parseManifestFields(manifest)
	tickets := parseTicketTable(manifest)

	fmt.Printf("\n")
	fmt.Printf("  Shipment: %s\n", shipment)
	fmt.Printf("  Title:    %s\n", title)
	fmt.Printf("  Status:   %s\n", formatStatus(status))

	total := len(tickets)
	done, _, _, _, _, _ := countTicketStatus(tickets)

	fmt.Printf("  Progress: %d/%d tickets\n", done, total)

	if done > 0 || total > 0 {
		bar := progressBar(done, total, 20)
		fmt.Printf("            %s\n", bar)
	}

	fmt.Printf("\n")
	if len(tickets) > 0 {
		fmt.Printf("  Tickets:\n")
		fmt.Printf("    %-8s %-20s %-8s %-12s %s\n", "Ticket", "Title", "Area", "Status", "PR")
		fmt.Printf("    %s\n", strings.Repeat("─", 72))
		for _, t := range tickets {
			pr := t.pr
			if pr == "" {
				pr = "—"
			} else {
				pr = shortURL(pr)
			}
			fmt.Printf("    %-8s %-20s %-8s %-12s %s\n", t.id, truncate(t.title, 20), t.area, formatStatus(t.status), pr)
		}
	}

	activityPath := filepath.Join(shipDir, "activity.log")
	if data, err := os.ReadFile(activityPath); err == nil {
		lines := strings.TrimSpace(string(data))
		if lines != "" {
			activity := strings.Split(lines, "\n")
			show := activity
			if len(show) > 5 {
				show = activity[len(activity)-5:]
			}
			fmt.Printf("\n  Recent Activity:\n")
			for _, line := range show {
				fmt.Printf("    %s\n", line)
			}
		}
	}

	fmt.Printf("\n")
	return nil
}

type ticketInfo struct {
	id     string
	title  string
	area   string
	status string
	pr     string
}

func readManifestSummary(path string) (status, title string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "unknown", ""
	}
	s, t := parseManifestFields(string(data))
	if s == "" {
		s = "unknown"
	}
	return s, t
}

var manifestTitleRe = regexp.MustCompile(`\*\*Title:\*\*\s*(.*)`)
var manifestStatusRe = regexp.MustCompile(`\*\*Status:\*\*\s*(.*)`)

func parseManifestFields(manifest string) (status, title string) {
	if m := manifestTitleRe.FindStringSubmatch(manifest); len(m) > 1 {
		title = strings.TrimSpace(m[1])
	}
	if m := manifestStatusRe.FindStringSubmatch(manifest); len(m) > 1 {
		status = strings.TrimSpace(m[1])
		status = strings.SplitN(status, "\n", 2)[0]
	}
	return
}

func parseTicketTable(manifest string) []ticketInfo {
	// Columns: ID | Title | Area | Status | Assignee | PR | Depends On
	re := regexp.MustCompile(`\| (T-\d+) \| (.*?) \| (\w+) \| ([\w ]+?) \| \w+ \| (\S+|—) \|`)
	matches := re.FindAllStringSubmatch(manifest, -1)
	var tickets []ticketInfo
	for _, m := range matches {
		pr := strings.TrimSpace(m[5])
		if pr == "—" {
			pr = ""
		}
		tickets = append(tickets, ticketInfo{
			id:     strings.TrimSpace(m[1]),
			title:  strings.TrimSpace(m[2]),
			area:   strings.TrimSpace(m[3]),
			status: strings.TrimSpace(m[4]),
			pr:     pr,
		})
	}
	return tickets
}

func countTicketStatus(tickets []ticketInfo) (todo, inprog, review, qa, done, blocked int) {
	for _, t := range tickets {
		switch strings.ToUpper(t.status) {
		case "TODO":
			todo++
		case "IN PROGRESS":
			inprog++
		case "REVIEW":
			review++
		case "QA":
			qa++
		case "DONE":
			done++
		case "BLOCKED":
			blocked++
		}
	}
	return
}

func formatStatus(s string) string {
	switch strings.ToUpper(s) {
	case "PLANNING":
		return "\033[33mPLANNING\033[0m"
	case "IN PROGRESS":
		return "\033[34mIN PROGRESS\033[0m"
	case "REVIEW":
		return "\033[35mREVIEW\033[0m"
	case "QA":
		return "\033[36mQA\033[0m"
	case "READY TO MERGE":
		return "\033[32mREADY TO MERGE\033[0m"
	case "DONE":
		return "\033[32mDONE\033[0m"
	case "TODO":
		return "\033[33mTODO\033[0m"
	case "BLOCKED":
		return "\033[31mBLOCKED\033[0m"
	default:
		return s
	}
}

func progressBar(current, total, width int) string {
	if total == 0 {
		return strings.Repeat(" ", width) + " 0%"
	}
	filled := current * width / total
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	pct := current * 100 / total
	return fmt.Sprintf("%s %d%%", bar, pct)
}

func shortURL(url string) string {
	// Show "#NN" for PR URLs
	re := regexp.MustCompile(`/(\d+)/?$`)
	if m := re.FindStringSubmatch(url); len(m) > 1 {
		return "#" + m[1]
	}
	return truncate(url, 12)
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n-1]) + "…"
}

func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all shipments",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := findProjectRoot()
			if dir == "" {
				return fmt.Errorf("no .ai/ found.\nRun 'devteam activate' first.")
			}
			return listShipments(dir)
		},
	}
}

func init() {
	formatStatus("") // ensure import isn't flagged; colors are ANSI escape codes for terminal output
}
