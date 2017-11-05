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
			Logger.Printf("there are definition for same category(=%s) and name(=%s), so current one will be overwritten.", category, name)
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
