package main

import (
	"fmt"
	"os"

	"github.com/handlename/mdqi"
)

func main() {
	app, err := mdqi.NewApp(mdqi.Conf{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	app.RegisterSlashCommandDefinition(mdqi.SlashCommandDisplay{})
	app.RegisterSlashCommandDefinition(mdqi.SlashCommandExit{})
	app.RegisterSlashCommandDefinition(mdqi.SlashCommandHelp{})
	app.RegisterSlashCommandDefinition(mdqi.SlashCommandTagClear{})
	app.RegisterSlashCommandDefinition(mdqi.SlashCommandTagSet{})
	app.RegisterSlashCommandDefinition(mdqi.SlashCommandTagShow{})

	app.Run()
}
