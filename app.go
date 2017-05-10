package mdqi

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peterh/liner"
)

type App struct {
	// cmdPath is path to mdq command.
	cmdPath string
}

type Conf struct {
}

type Result struct {
	Database string
	Columns  []string
	Rows     []map[string]interface{}
}

func NewApp(conf Conf) (*App, error) {
	// TODO: Check if mdq command exists by exec.LookPath.

	return &App{
		cmdPath: "mdq",
	}, nil
}

func (app *App) Run() {
	app.runLiner()
}

func (app *App) runLiner() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

LOOP:
	for {
		l, err := line.Prompt("mdq> ")

		switch err {
		case nil:
			results, err := app.RunCmd(strings.Trim(l, " \n"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}

			Print(results)
		case liner.ErrPromptAborted:
			fmt.Fprintln(os.Stderr, "aborted")
			break LOOP
		case io.EOF:
			fmt.Println("bye")
			break LOOP
		default:
			fmt.Fprintln(os.Stderr, "error on reading line: ", err)
			break LOOP
		}
	}
}
