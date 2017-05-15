package mdqi

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
	"github.com/pkg/errors"
)

var (
	ErrNotASlashCommand    = errors.New("there are no SlashCommand")
)

type App struct {
	// cmdPath is path to mdq command.
	cmdPath string

	// historyPath is path to command history file for liner.
	historyPath string
}

type Conf struct {
}

type Result struct {
	Database string
	Columns  []string
	Rows     []map[string]interface{}
}

type SlashCommand struct {
	Category string
	Name     string
	Args     []string
}

func NewApp(conf Conf) (*App, error) {
	// TODO: Check if mdq command exists by exec.LookPath.
	// TODO: Make historyPath configuarable.

	historyPath, err := createHistoryFile()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create history file")
	}

	return &App{
		cmdPath:     "mdq",
		historyPath: historyPath,
	}, nil
}

func createHistoryFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "failed to get current user")
	}

	path := filepath.Join(usr.HomeDir, ".mdqi_history")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if _, err := os.Create(path); err != nil {
			return "", errors.Wrap(err, "failed to create history file")
		}
	}

	return path, nil
}

func (app *App) Run() {
	app.runLiner()
}

func (app *App) runLiner() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(app.historyPath); err == nil {
		line.ReadHistory(f)
		f.Close()
	} else {
		fmt.Fprintln(os.Stderr, "failed to read command history: ", err)
	}

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

			line.AppendHistory(l)
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

		if f, err := os.Create(app.historyPath); err == nil {
			if _, err := line.WriteHistory(f); err != nil {
				fmt.Fprintln(os.Stderr, "failed to write history: ", err)
			}

			f.Close()
		} else {
			fmt.Fprintln(os.Stderr, "failed to create history file: ", err)
		}
	}
}
