package mdqi

import (
	"fmt"
	"regexp"
	"sort"
)

type SlashCommand struct {
	Category string
	Name     string
	Args     []string
}

type SlashCommandDefinition interface {
	Category() string
	Name() string
	Example() string
	Help() string
	Handle(app *App, cmd *SlashCommand) error
}

// /category [ cmd [ name [ arg1 [ arg2 ]... ] ] ]
var slashCommandRegexp = regexp.MustCompile("^/([a-z]+)(?: ([a-z]+)(?: ([-_a-zA-Z0-9]+))*)? *$")

func ParseSlashCommand(query string) (*SlashCommand, error) {
	matches := slashCommandRegexp.FindAllStringSubmatch(query, -1)

	if len(matches) == 0 {
		return nil, ErrNotASlashCommand
	}

	category := matches[0][1]
	name := ""
	args := []string{}

	if 3 <= len(matches[0]) && matches[0][2] != "" {
		name = matches[0][2]
	}

	if 4 <= len(matches[0]) && matches[0][3] != "" {
		args = matches[0][3:]
	}

	return &SlashCommand{
		Category: category,
		Name:     name,
		Args:     args,
	}, nil
}

func (app *App) RegisterSlashCommandDefinition(d SlashCommandDefinition) error {
	category := d.Category()
	name := d.Name()

	if c := app.slashCommandDefinition[category]; c != nil {
		if c[name] != nil {
			logger.Printf("there are definition for same category(=%s) and name(=%s), so current one will be overwritten.", category, name)
		}
	}

	if app.slashCommandDefinition[category] == nil {
		app.slashCommandDefinition[category] = map[string]SlashCommandDefinition{}
	}

	app.slashCommandDefinition[category][name] = d

	return nil
}

func (app *App) FindSlashCommandDefinition(category, name string) (SlashCommandDefinition, error) {
	c := app.slashCommandDefinition[category]

	if c == nil {
		return nil, ErrSlashCommandNotFound
	}

	if c[name] == nil {
		return nil, ErrSlashCommandNotFound
	}

	return c[name], nil
}

type SlashCommandExit struct{}

func (c SlashCommandExit) Category() string { return "exit" }
func (c SlashCommandExit) Name() string     { return "" }
func (c SlashCommandExit) Example() string  { return "/exit" }

func (c SlashCommandExit) Help() string {
	return "Exit mdqi."
}

func (c SlashCommandExit) Handle(app *App, cmd *SlashCommand) error {
	app.Alive = false
	return nil
}

type SlashCommandTagSet struct{}

func (c SlashCommandTagSet) Category() string { return "tag" }
func (c SlashCommandTagSet) Name() string     { return "set" }
func (c SlashCommandTagSet) Example() string  { return "/tag set {tagname}" }

func (c SlashCommandTagSet) Help() string {
	return "Set tag to pass to mdq's --tag option."
}

func (c SlashCommandTagSet) Handle(app *App, cmd *SlashCommand) error {
	if len(cmd.Args) != 1 {
		fmt.Println(c.Help())
		return nil
	}

	app.SetTag(cmd.Args[0])
	debug.Println("tag stored:", app.GetTag())

	return nil
}

type SlashCommandTagClear struct{}

func (c SlashCommandTagClear) Category() string { return "tag" }
func (c SlashCommandTagClear) Name() string     { return "clear" }
func (c SlashCommandTagClear) Example() string  { return "/tag clear {tagname}" }

func (c SlashCommandTagClear) Help() string {
	return "Clear tag to pass to mdq's --tag option."
}

func (c SlashCommandTagClear) Handle(app *App, cmd *SlashCommand) error {
	if len(cmd.Args) != 1 {
		fmt.Println(c.Help())
		return nil
	}

	app.ClearTag()
	debug.Printf("tag cleard: %+v\n", app.GetTag())

	return nil
}

type SlashCommandTagShow struct{}

func (c SlashCommandTagShow) Category() string { return "tag" }
func (c SlashCommandTagShow) Name() string     { return "show" }
func (c SlashCommandTagShow) Example() string  { return "/tag show" }

func (c SlashCommandTagShow) Help() string {
	return "Show registered tags to pass to mdq's --tag option."
}

func (c SlashCommandTagShow) Handle(app *App, cmd *SlashCommand) error {
	results := []Result{
		Result{
			Database: "(mdq)",
			Columns:  []string{"tag"},
			Rows: []map[string]interface{}{
				map[string]interface{}{"tag": app.GetTag()},
			},
		},
	}

	Print(app.printer, results)

	return nil
}

type SlashCommandHelp struct{}

func (c SlashCommandHelp) Category() string { return "help" }
func (c SlashCommandHelp) Name() string     { return "" }
func (c SlashCommandHelp) Example() string  { return "/help" }

func (c SlashCommandHelp) Help() string {
	return "Show this help."
}

func (c SlashCommandHelp) Handle(app *App, cmd *SlashCommand) error {
	rows := []map[string]interface{}{}

	categories := app.slashCommandCategories()
	sort.Strings(categories)

	for _, category := range categories {
		inCategory := app.slashCommandDefinition[category]

		names := app.slashCommandNames(category)
		sort.Strings(names)

		for _, name := range names {
			def := inCategory[name]

			rows = append(rows, map[string]interface{}{
				"category": category,
				"name":     name,
				"example":  def.Example(),
				"help":     def.Help(),
			})
		}
	}

	results := []Result{
		Result{
			Database: "(mdq)",
			Columns:  []string{"category", "name", "example", "help"},
			Rows:     rows,
		},
	}

	Print(app.printer, results)

	return nil
}

type SlashCommandDisplay struct{}

func (c SlashCommandDisplay) Category() string { return "display" }
func (c SlashCommandDisplay) Name() string     { return "set" }
func (c SlashCommandDisplay) Example() string  { return "/display set (horizontal|vertical)" }

func (c SlashCommandDisplay) Help() string {
	return "Set display format."
}

func (c SlashCommandDisplay) Handle(app *App, cmd *SlashCommand) error {
	if len(cmd.Args) != 1 {
		return ErrSlashCommandInvalidArgs
	}

	switch cmd.Args[0] {
	case "horizontal":
		app.printer = HorizontalPrinter{}
	case "vertical":
		app.printer = VerticalPrinter{}
	default:
		return ErrSlashCommandInvalidArgs
	}

	debug.Printf("printer has been changed to %s", cmd.Args[0])

	return nil
}
