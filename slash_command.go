package mdqi

import (
	"regexp"
)

type SlashCommand struct {
	Category string
	Name     string
	Args     []string
}

type SlashCommandDefinition interface {
	Category() string
	Name() string
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
func (c SlashCommandExit) Handle(app *App, cmd *SlashCommand) error {
	app.Alive = false
	return nil
}

type SlashCommandTagAdd struct{}

func (c SlashCommandTagAdd) Category() string { return "tag" }
func (c SlashCommandTagAdd) Name() string     { return "add" }
func (c SlashCommandTagAdd) Handle(app *App, cmd *SlashCommand) error {
	if len(cmd.Args) != 1 {
		logger.Println("usage: /tag add {tagname}")
		return nil
	}

	app.AddTag(cmd.Args[0])
	debug.Printf("tag added: %+v\n", app.GetTags())

	return nil
}

type SlashCommandTagRemove struct{}

func (c SlashCommandTagRemove) Category() string { return "tag" }
func (c SlashCommandTagRemove) Name() string     { return "remove" }
func (c SlashCommandTagRemove) Handle(app *App, cmd *SlashCommand) error {
	if len(cmd.Args) != 1 {
		logger.Println("usage: /tag remove {tagname}")
		return nil
	}

	app.RemoveTag(cmd.Args[0])
	debug.Printf("tag removed: %+v\n", app.GetTags())

	return nil
}

type SlashCommandTagShow struct{}

func (c SlashCommandTagShow) Category() string { return "tag" }
func (c SlashCommandTagShow) Name() string     { return "show" }
func (c SlashCommandTagShow) Handle(app *App, cmd *SlashCommand) error {
	rows := []map[string]interface{}{}

	for _, tag := range app.GetTags() {
		rows = append(rows, map[string]interface{}{"tag": tag})
	}

	results := []Result{
		Result{
			Database: "(mdq)",
			Columns:  []string{"tag"},
			Rows:     rows,
		},
	}

	Print(results)

	return nil
}
