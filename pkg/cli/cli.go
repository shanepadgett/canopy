// Package cli provides a minimal, dependency-free CLI framework.
// It supports subcommands, flags, and auto-generated help.
// Designed to be lift-and-shift portable to any Go project.
package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

// App represents a CLI application with subcommands.
type App struct {
	Name        string
	Description string
	Version     string
	Commands    map[string]*Command
	Stdout      io.Writer
	Stderr      io.Writer
}

// New creates a new CLI application.
func New(name, description, version string) *App {
	return &App{
		Name:        name,
		Description: description,
		Version:     version,
		Commands:    make(map[string]*Command),
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
	}
}

// Command represents a subcommand with its own flags and action.
type Command struct {
	Name        string
	Usage       string
	Description string
	Flags       *FlagSet
	Subcommands map[string]*Command
	Action      func(ctx *Context) error
}

// NewCommand creates a new command.
func NewCommand(name, usage, description string) *Command {
	return &Command{
		Name:        name,
		Usage:       usage,
		Description: description,
		Flags:       NewFlagSet(name),
		Subcommands: make(map[string]*Command),
	}
}

// AddSubcommand adds a subcommand to this command.
func (c *Command) AddSubcommand(sub *Command) {
	c.Subcommands[sub.Name] = sub
}

// Context holds the parsed state available to command actions.
type Context struct {
	App     *App
	Command *Command
	Flags   *FlagSet
	Args    []string
}

// Add registers a command with the app.
func (a *App) Add(cmd *Command) {
	a.Commands[cmd.Name] = cmd
}

// Run parses arguments and executes the appropriate command.
func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.printHelp()
		return nil
	}

	cmdName := args[1]

	// Handle global flags
	if cmdName == "-h" || cmdName == "--help" || cmdName == "help" {
		if len(args) > 2 {
			return a.printCommandHelp(args[2])
		}
		a.printHelp()
		return nil
	}

	if cmdName == "-v" || cmdName == "--version" || cmdName == "version" {
		fmt.Fprintf(a.Stdout, "%s %s\n", a.Name, a.Version)
		return nil
	}

	cmd, ok := a.Commands[cmdName]
	if !ok {
		fmt.Fprintf(a.Stderr, "Unknown command: %s\n\n", cmdName)
		a.printHelp()
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	return a.runCommand(cmd, args[2:])
}

func (a *App) runCommand(cmd *Command, args []string) error {
	// Check for subcommands first
	if len(args) > 0 && len(cmd.Subcommands) > 0 {
		if sub, ok := cmd.Subcommands[args[0]]; ok {
			return a.runCommand(sub, args[1:])
		}
	}

	// Check for help flag
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			a.printCommandHelpDirect(cmd)
			return nil
		}
	}

	// Parse flags
	remaining, err := cmd.Flags.Parse(args)
	if err != nil {
		return fmt.Errorf("flag error: %w", err)
	}

	if cmd.Action == nil {
		a.printCommandHelpDirect(cmd)
		return nil
	}

	ctx := &Context{
		App:     a,
		Command: cmd,
		Flags:   cmd.Flags,
		Args:    remaining,
	}

	return cmd.Action(ctx)
}

func (a *App) printHelp() {
	w := tabwriter.NewWriter(a.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%s - %s\n\n", a.Name, a.Description)
	fmt.Fprintf(w, "Usage:\n")
	fmt.Fprintf(w, "  %s <command> [options]\n\n", a.Name)
	fmt.Fprintf(w, "Commands:\n")

	// Sort commands for consistent output
	names := make([]string, 0, len(a.Commands))
	for name := range a.Commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		cmd := a.Commands[name]
		fmt.Fprintf(w, "  %s\t%s\n", name, cmd.Description)
	}

	fmt.Fprintf(w, "\nRun '%s <command> --help' for more information on a command.\n", a.Name)
	w.Flush()
}

func (a *App) printCommandHelp(name string) error {
	cmd, ok := a.Commands[name]
	if !ok {
		return fmt.Errorf("unknown command: %s", name)
	}
	a.printCommandHelpDirect(cmd)
	return nil
}

func (a *App) printCommandHelpDirect(cmd *Command) {
	w := tabwriter.NewWriter(a.Stdout, 0, 0, 2, ' ', 0)

	if cmd.Usage != "" {
		fmt.Fprintf(w, "Usage:\n  %s %s\n\n", a.Name, cmd.Usage)
	} else {
		fmt.Fprintf(w, "Usage:\n  %s %s [options]\n\n", a.Name, cmd.Name)
	}

	fmt.Fprintf(w, "%s\n", cmd.Description)

	// Print subcommands if any
	if len(cmd.Subcommands) > 0 {
		fmt.Fprintf(w, "\nSubcommands:\n")
		names := make([]string, 0, len(cmd.Subcommands))
		for name := range cmd.Subcommands {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			sub := cmd.Subcommands[name]
			fmt.Fprintf(w, "  %s\t%s\n", name, sub.Description)
		}
	}

	// Print flags if any
	if cmd.Flags.Len() > 0 {
		fmt.Fprintf(w, "\nOptions:\n")
		cmd.Flags.PrintDefaults(w)
	}

	w.Flush()
}

// FlagSet is a minimal flag parser.
type FlagSet struct {
	name    string
	flags   map[string]*Flag
	ordered []string
}

// Flag represents a single flag.
type Flag struct {
	Name     string
	Short    string // single char alias
	Usage    string
	DefValue string
	Value    Value
}

// Value is the interface for flag values.
type Value interface {
	String() string
	Set(string) error
}

// NewFlagSet creates a new flag set.
func NewFlagSet(name string) *FlagSet {
	return &FlagSet{
		name:  name,
		flags: make(map[string]*Flag),
	}
}

// Len returns the number of flags.
func (f *FlagSet) Len() int {
	return len(f.flags)
}

// String defines a string flag.
func (f *FlagSet) String(name, short, defValue, usage string) *string {
	p := new(string)
	*p = defValue
	f.Var(&stringValue{p}, name, short, defValue, usage)
	return p
}

// Bool defines a bool flag.
func (f *FlagSet) Bool(name, short string, defValue bool, usage string) *bool {
	p := new(bool)
	*p = defValue
	def := "false"
	if defValue {
		def = "true"
	}
	f.Var(&boolValue{p}, name, short, def, usage)
	return p
}

// Int defines an int flag.
func (f *FlagSet) Int(name, short string, defValue int, usage string) *int {
	p := new(int)
	*p = defValue
	f.Var(&intValue{p}, name, short, fmt.Sprintf("%d", defValue), usage)
	return p
}

// Var registers a custom flag value.
func (f *FlagSet) Var(value Value, name, short, defValue, usage string) {
	flag := &Flag{
		Name:     name,
		Short:    short,
		Usage:    usage,
		DefValue: defValue,
		Value:    value,
	}
	f.flags[name] = flag
	if short != "" {
		f.flags[short] = flag
	}
	f.ordered = append(f.ordered, name)
}

// Parse parses arguments and returns remaining positional args.
func (f *FlagSet) Parse(args []string) ([]string, error) {
	var remaining []string
	i := 0

	for i < len(args) {
		arg := args[i]

		if !strings.HasPrefix(arg, "-") {
			remaining = append(remaining, arg)
			i++
			continue
		}

		// Strip leading dashes
		name := strings.TrimLeft(arg, "-")

		// Handle --flag=value
		value := ""
		if idx := strings.Index(name, "="); idx >= 0 {
			value = name[idx+1:]
			name = name[:idx]
		}

		flag, ok := f.flags[name]
		if !ok {
			return nil, fmt.Errorf("unknown flag: %s", arg)
		}

		// Bool flags don't require a value
		if _, isBool := flag.Value.(*boolValue); isBool {
			if value == "" {
				value = "true"
			}
		} else if value == "" {
			// Need next arg as value
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("flag %s requires a value", arg)
			}
			value = args[i]
		}

		if err := flag.Value.Set(value); err != nil {
			return nil, fmt.Errorf("invalid value for %s: %w", arg, err)
		}

		i++
	}

	return remaining, nil
}

// Get returns the value of a flag by name.
func (f *FlagSet) Get(name string) string {
	if flag, ok := f.flags[name]; ok {
		return flag.Value.String()
	}
	return ""
}

// GetBool returns a bool flag value.
func (f *FlagSet) GetBool(name string) bool {
	return f.Get(name) == "true"
}

// PrintDefaults writes flag defaults to w.
func (f *FlagSet) PrintDefaults(w io.Writer) {
	seen := make(map[string]bool)
	for _, name := range f.ordered {
		if seen[name] {
			continue
		}
		flag := f.flags[name]
		seen[name] = true

		leader := "--" + name
		if flag.Short != "" {
			leader = "-" + flag.Short + ", --" + name
		}

		defNote := ""
		if flag.DefValue != "" && flag.DefValue != "false" {
			defNote = fmt.Sprintf(" (default: %s)", flag.DefValue)
		}

		fmt.Fprintf(w, "  %s\t%s%s\n", leader, flag.Usage, defNote)
	}
}

// --- Value implementations ---

type stringValue struct{ p *string }

func (s *stringValue) String() string     { return *s.p }
func (s *stringValue) Set(v string) error { *s.p = v; return nil }

type boolValue struct{ p *bool }

func (b *boolValue) String() string {
	if *b.p {
		return "true"
	}
	return "false"
}
func (b *boolValue) Set(v string) error {
	switch strings.ToLower(v) {
	case "true", "1", "yes", "on":
		*b.p = true
	case "false", "0", "no", "off":
		*b.p = false
	default:
		return fmt.Errorf("invalid bool: %s", v)
	}
	return nil
}

type intValue struct{ p *int }

func (i *intValue) String() string { return fmt.Sprintf("%d", *i.p) }
func (i *intValue) Set(v string) error {
	var n int
	_, err := fmt.Sscanf(v, "%d", &n)
	if err != nil {
		return err
	}
	*i.p = n
	return nil
}
