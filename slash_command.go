package mdqi

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
)

type SlashCommandHandler func(app *App, cmd *SlashCommand) error

// /category [ cmd [ arg1 [ arg2 ]... ] ]
var slashCommandRegexp = regexp.MustCompile("^/([a-z]+)(?: ([a-z]+)(?: ([-_a-zA-Z0-9]+))+)? *$")

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

func (app *App) RegisterSlashCommand(category, name string, fn SlashCommandHandler) error {
	if c := app.slashCommandDefinition[category]; c != nil {
		if !reflect.DeepEqual(c[name], SlashCommandDefinition{}) {
			fmt.Fprintf(os.Stderr, "there are definition for same category(=%s) and name(=%s), so current one will be overwritten.", category, name)
		}
	}

	if app.slashCommandDefinition[category] == nil {
		app.slashCommandDefinition[category] = map[string]SlashCommandDefinition{}
	}

	app.slashCommandDefinition[category][name] = SlashCommandDefinition{
		Category: category,
		Name:     name,
		Handler:  fn,
	}

	return nil
}

func (app *App) FindSlashCommandDefinition(category, name string) (SlashCommandDefinition, error) {
	c := app.slashCommandDefinition[category]

	if c == nil {
		return SlashCommandDefinition{}, ErrSlashCommandNotFound
	}

	if reflect.DeepEqual(c[name], SlashCommandDefinition{}) {
		return SlashCommandDefinition{}, ErrSlashCommandNotFound
	}

	return c[name], nil
}
