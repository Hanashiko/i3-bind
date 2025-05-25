package main

import (
	// "bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "os/exec"
	"path/filepath"
	"regexp"
	"sort"
	// "strconv"
	"strings"
	// "syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	VERSION = "1.0.0"
)

var (
	configPath string
	noColor bool

	keyColor = color.New(color.FgCyan, color.Bold)
	actionColor = color.New(color.FgGreen)
	commentColor = color.New(color.FgYellow)
	errorColor = color.New(color.FgRed, color.Bold)
	successColor = color.New(color.FgGreen, color.Bold)
)

type Binding struct {
	Key string
	Action string
	Comment string
	Line int
	Raw string
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "i3-bind",
		Short: "A CLI/TUI utility to manage i3 window manager keybindings",
		Version: VERSION,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if noColor {
				color.NoColor = true
			}
			if configPath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					log.Fatal("Cannot datermine home directory")
				}
				configPath = filepath.Join(home, ".config", "i3", "config")
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to i3 config file (default: ~/.config/i3/config)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	var addCmd = &cobra.Command{
		Use: "add [key] [action...]",
		Short: "Add a new keybinding",
		Long: "Add a new keybinding to the i3 config file",
		Example: `  i3-bind add mod4+Enter exec alacritty
  i3-bind add mod4+d exec dmenu_run
  i3-bind add mod4+shift+q kill
  i3-bind add '$mod+shift+k' keepassxc`,
		Args: cobra.MinimumNArgs(2),
		Run: addBinding,
	}

	var removeCmd = &cobra.Command{
		Use: "remove [key]",
		Short: "Remove a keybinding",
		Long: "Remove a keybinding from the i3 config file",
		Example: `  i3-bind remove mod4+q
  i3-bind remove mod4+Enter`,
		Args: cobra.ExactArgs(1),
		Run: removeBinding,
	}

	var listCmd = &cobra.Command{
		Use: "list",
		Short: "List all keybindings",
		Long: "List all keybinding in the i3 config file with syntax highlighting",
		Run: listBindings,
	}

	var findCmd = &cobra.Command{
		Use: "find [search_term]",
		Short: "Find keybinding by action or key",
		Long: "Search for keybinding by action or key pattern",
		Example: `  i3-bind find firefox
  i3-bind finx exec
  i3-bind find mod4+shift
  i3-bind find '$mod+return'`,
		Args: cobra.ExactArgs(1),
		Run: findBindings,
	}

	var commentCmd = &cobra.Command{
		Use: "comment [key] [comment]",
		Short: "Add or update comment for a keybinding",
		Long: "Add or update a comment for an existing keybinding",
		Example: `  i3-bind comment mod4+r "restart i3"
  i3-bind comment mod4+shift+3 "exit i3"
  i3-bind comment "$mod+return" "run terminal"`,
		Args: cobra.ExactArgs(2),
		Run: commentBinding,
	}


	rootCmd.AddCommand(addCmd, removeCmd, listCmd, findCmd, commentCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readConfig() ([]string, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("i3 config file not found at %s", configPath)
	}
	
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v",err)
	}
	return strings.Split(string(content), "\n"),nil
}

func writeConfig(lines []string) error {
	content := strings.Join(lines, "\n")

	backupPath := configPath + ".backup"
	if err := ioutil.WriteFile(backupPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	if err := ioutil.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func parseBindings(lines []string) []Binding {
	var bindings []Binding
	bindRegex := regexp.MustCompile(`^\s*bindsym\s+([^\s]+)\s+(.+?)(?:\s*#\s*(.*))?$`)

	for i, line := range lines {
		matches := bindRegex.FindStringSubmatch(line)
		if matches != nil {
			comment := strings.TrimSpace(matches[3])

			if comment == "" && i > 0 {
				previousLine := strings.TrimSpace(lines[i-1])
				if strings.HasPrefix(previousLine, "#") {
					trimmed := strings.TrimSpace(strings.TrimPrefix(previousLine, "#"))
					
					if !strings.HasSuffix(trimmed, ":") {
						comment = trimmed
					}
				}
			}
			binding := Binding{
				Key: matches[1],
				Action: strings.TrimSpace(matches[2]),
				Comment: comment,
				Line: i+1,
				Raw: line,
			}
			bindings = append(bindings, binding)
		}
	}
	return bindings
}

func addBinding(cmd *cobra.Command, args []string) {
	key := args[0]
	action := strings.Join(args[1:], " ")

	lines, err := readConfig()
	if err != nil {
		errorColor.Printf("Error: %v\n",err)
		os.Exit(1)
	}

	bindings := parseBindings(lines)
	for _, binding := range bindings {
		if strings.EqualFold(binding.Key, key){
			errorColor.Printf("Error: Keybinding %s already exists\n", key)
			fmt.Printf("Current bindig: %s -> %s\n", keyColor.Sprint(binding.Key), actionColor.Sprint(binding.Action))
			fmt.Println("Use 'i3-bind remove' first or modify the config manually")
			os.Exit(1)
		}
	}

	newBinding := fmt.Sprintf("bindsym %s %s", key, action)

	insertIndex := len(lines)
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(strings.TrimSpace(lines[i]), "bindsym") {
			insertIndex = i + 1
			break
		}
	}

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:insertIndex]...)
	newLines = append(newLines, newBinding)
	newLines = append(newLines, lines[insertIndex:]...)

	if err := writeConfig(newLines); err != nil {
		errorColor.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	successColor.Printf("✓ Added keybinding: %s -> %s\n", keyColor.Sprint(key), actionColor.Sprint(action))
}

func removeBinding(cmd *cobra.Command, args []string) {
	key := args[0]

	lines, err := readConfig()
	if err != nil {
		errorColor.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	bindings := parseBindings(lines)
	found := false
	var removedBinding Binding

	for _, binding := range bindings {
		if strings.EqualFold(binding.Key, key) {
			found = true
			removedBinding = binding
			break
		}
	}

	if !found {
		errorColor.Printf("Error: Keybinding %s not found\n", key)
		os.Exit(1)
	}

	newLines := make([]string, 0, len(lines)-1)
	bindRegex := regexp.MustCompile(`(?i)^\s*bindsym\s+` + regexp.QuoteMeta(key) + `\s+`)

	for _, line := range lines {
		if !bindRegex.MatchString(line) {
			newLines = append(newLines, line)
		}
	}

	if err := writeConfig(newLines); err != nil {
		errorColor.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	successColor.Printf("✓ Removed keybinding: %s -> %s\n", keyColor.Sprint(removedBinding.Key), actionColor.Sprint(removedBinding.Action))
}

func listBindings(cmd *cobra.Command, args []string) {
	lines, err := readConfig()
	if err != nil {
		errorColor.Printf("Error: %v\n",err)
		os.Exit(1)
	}

	bindings := parseBindings(lines)
	if len(bindings) == 0 {
		fmt.Println("No keybindings found in config file")
		return
	}

	sort.Slice(bindings, func(i, j int) bool {
		return bindings[i].Key < bindings[j].Key
	})

	fmt.Printf("Found %d keybindings in %s:\n\n", len(bindings), configPath)

	for _, binding := range bindings {
		fmt.Printf("  %s -> %s", keyColor.Sprint(binding.Key),actionColor.Sprint(binding.Action))
		if binding.Comment != "" {
			fmt.Printf(" %s", commentColor.Sprintf("# %s",binding.Comment))
		}
		fmt.Println()
	}
}

func findBindings(cmd *cobra.Command, args []string) {
	searchTerm := args[0]

	lines, err := readConfig()
	if err != nil {
		errorColor.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	bindings := parseBindings(lines)
	var matches []Binding

	searchLower := strings.ToLower(searchTerm)
	for _, binding := range bindings {
		if strings.Contains(strings.ToLower(binding.Key), searchLower) ||
		   strings.Contains(strings.ToLower(binding.Action), searchLower) ||
		   strings.Contains(strings.ToLower(binding.Comment), searchLower) {
			matches = append(matches, binding)
		}
	}

	if len(matches) == 0 {
		fmt.Printf("No keybindings found matching '%s'\n", searchTerm)
		return
	}

	fmt.Printf("Found %d keybinding(s) matching '%s':\n\n",len(matches),searchTerm)

	for _, binding := range matches {
		fmt.Printf("  %s -> %s", keyColor.Sprint(binding.Key), actionColor.Sprint(binding.Action))
		if binding.Comment != "" {
			fmt.Printf(" %s", commentColor.Sprintf("# %s", binding.Comment))
		}
		fmt.Printf(" %s\n", color.New(color.FgBlack, color.Bold).Sprintf("(line %d)", binding.Line))
	}
}

func commentBinding(cmd *cobra.Command, args []string) {
	key := args[0]
	comment := args[1]

	lines, err := readConfig()
	if err != nil {
		errorColor.Printf("Error: %v\n",err)
		os.Exit(1)
	}

	bindings := parseBindings(lines)
	found := false

	for _, binding := range bindings {
		if strings.EqualFold(binding.Key, key) {
			found = true
			break
		}
	}

	if !found {
		errorColor.Printf("Error: Keybinding %s not found\n", key)
		os.Exit(1)
	}

	bindRegex := regexp.MustCompile(`^\s*bindsym\s+` + regexp.QuoteMeta(key) + `\s+(.+?)(?:\s*#.*)?$`)

	for i, line := range lines {
		if bindRegex.MatchString(line) {
			matches := bindRegex.FindStringSubmatch(line)
			if matches != nil {
				lines[i] = fmt.Sprintf("bindsym %s %s # %s", key, matches[1], comment)
				break
			}
		}
	}

	if err := writeConfig(lines); err != nil {
		errorColor.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	successColor.Printf("✓ Added comment to keybinding: %s # %s\n",keyColor.Sprint(key), commentColor.Sprint(comment))
}
