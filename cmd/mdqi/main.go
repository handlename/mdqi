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

	results, err := app.RunCmd(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	app.Print(results)
}
