package mdqi

import (
	"regexp"
)

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
