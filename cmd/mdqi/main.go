package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/handlename/mdqi"
	"github.com/handlename/mdqi/slash"
)

var version string

func main() {
	var (
		confPath    string
		showVersion bool
	)

	flag.StringVar(&confPath, "conf", "", "path to configuration file.")
	flag.BoolVar(&showVersion, "version", false, "display version.")
	flag.Parse()

	if showVersion {
		fmt.Println("mdqi", version)
		os.Exit(0)
	}

	conf := mdqi.Conf{}
	if confPath != "" {
		var err error
		if conf, err = mdqi.ConfFromFile(confPath); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	app, err := mdqi.NewApp(conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	app.RegisterSlashCommandDefinition(slash.Display{})
	app.RegisterSlashCommandDefinition(slash.Exit{})
	app.RegisterSlashCommandDefinition(slash.Help{})
	app.RegisterSlashCommandDefinition(slash.TagClear{})
	app.RegisterSlashCommandDefinition(slash.TagSet{})
	app.RegisterSlashCommandDefinition(slash.TagShow{})
	app.RegisterSlashCommandDefinition(slash.ToggleDisplay{})

	app.Run()
}
