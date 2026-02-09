package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newSkillCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "skill",
		Short: "Generate a skill file (markdown) listing all commands",
		Example: `  clicker skill
  # Outputs markdown skill file to stdout

  clicker skill > vibium-skill.md
  # Save to file`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			root := cmd.Root()
			generateSkillFile(root)
		},
	}
}

func generateSkillFile(root *cobra.Command) {
	// Commands to skip (diagnostic/internal)
	skip := map[string]bool{
		"version":          true,
		"paths":            true,
		"install":          true,
		"launch-test":      true,
		"ws-test":          true,
		"bidi-test":        true,
		"daemon":           true,
		"serve":            true,
		"mcp":              true,
		"check-actionable": true,
		"skill":            true,
		"help":             true,
		"completion":       true,
	}

	// Categories and their commands
	categories := []struct {
		Name     string
		Commands []string
	}{
		{"Navigation", []string{"navigate", "url", "title"}},
		{"Reading", []string{"text", "html", "find", "find-all", "eval", "screenshot"}},
		{"Interaction", []string{"click", "type", "hover", "scroll", "keys", "select"}},
		{"Waiting", []string{"wait"}},
		{"Tabs", []string{"tabs", "tab-new", "tab-switch", "tab-close"}},
	}

	// Build command lookup
	cmdMap := map[string]*cobra.Command{}
	for _, c := range root.Commands() {
		cmdMap[c.Name()] = c
	}

	fmt.Println("# Vibium Clicker - Browser Automation Commands")
	fmt.Println()
	fmt.Println("Browser automation for AI agents and humans.")
	fmt.Println("Commands default to daemon mode (persistent browser). Use `--oneshot` for single-shot execution.")
	fmt.Println()

	for _, cat := range categories {
		fmt.Printf("## %s\n\n", cat.Name)
		for _, name := range cat.Commands {
			c, ok := cmdMap[name]
			if !ok {
				continue
			}
			fmt.Printf("- `clicker %s` — %s\n", c.Use, c.Short)
			delete(skip, name) // Mark as categorized
		}
		fmt.Println()
	}

	// Session management (not in categories above)
	fmt.Println("## Session")
	fmt.Println()
	if c, ok := cmdMap["quit"]; ok {
		fmt.Printf("- `clicker %s` — %s\n", c.Use, c.Short)
	}
	fmt.Println("- `clicker daemon status` — Show daemon status")
	fmt.Println()

	// Flags section
	fmt.Println("## Flags")
	fmt.Println()
	fmt.Println("- `--json` — Output as JSON")
	fmt.Println("- `--oneshot` — One-shot mode (no daemon)")
	fmt.Println("- `--headless` — Hide browser window")
	fmt.Println("- `--verbose` — Enable debug logging")
	fmt.Println()

	// List any uncategorized automation commands
	var uncategorized []string
	categorized := map[string]bool{}
	for _, cat := range categories {
		for _, name := range cat.Commands {
			categorized[name] = true
		}
	}
	categorized["quit"] = true

	for _, c := range root.Commands() {
		name := c.Name()
		if !skip[name] && !categorized[name] && !strings.HasPrefix(name, "_") {
			uncategorized = append(uncategorized, name)
		}
	}

	if len(uncategorized) > 0 {
		fmt.Println("## Other")
		fmt.Println()
		for _, name := range uncategorized {
			c := cmdMap[name]
			fmt.Printf("- `clicker %s` — %s\n", c.Use, c.Short)
		}
		fmt.Println()
	}
}
