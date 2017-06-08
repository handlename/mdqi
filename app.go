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

const (
	defaultHistoryFilename = ".mdqi_history"
)

var (
	ErrSlashCommandNotFound    = errors.New("unknown SlashCommand")
	ErrNotASlashCommand        = errors.New("there are no SlashCommand")
	ErrSlashCommandInvalidArgs = errors.New("invalid args")
	ErrUnknownPrinterName      = errors.New("unknown printer name")
)

type App struct {
	// Alive turns into false, mdqi will exit.
	Alive bool

	// cmdPath is path to mdq command.
	cmdPath string

	// historyPath is path to command history file for liner.
	historyPath string

	// slashCommandDefinition holds SlashCommandDefinition.
	// app.slashCommandDefinition[category][name] = SlashCommandDefinition
	slashCommandDefinition map[string]map[string]SlashCommandDefinition

	// tag stores tag value for --tag option of mdq.
	tag string

	printer Printer
}

type Result struct {
	Database string
	Columns  []string
	Rows     []map[string]interface{}
}

func init() {
	defaultOutput = os.Stdout
}

func NewApp(conf Conf) (*App, error) {
	// TODO: Check if mdq command exists by exec.LookPath.

	// create history file
	historyPath := conf.Mdqi.History
	if path := conf.Mdqi.History; path != "" {
		var err error
		if err = createHistoryFile(path); err != nil {
			return nil, errors.Wrapf(err, "failed to create history file at %s", path)
		}
		historyPath = path
		debug.Println("conf.Mdqi.History =", historyPath)
	}

	app := &App{
		Alive: true,

		cmdPath:                "mdq",
		historyPath:            historyPath,
		slashCommandDefinition: map[string]map[string]SlashCommandDefinition{},
		printer:                HorizontalPrinter{},
	}

	// set default tag
	if tag := conf.Mdqi.DefaultTag; tag != "" {
		app.SetTag(tag)
		debug.Println("conf.Mdqi.DefaultTag =", tag)
	}

	// set default display
	if display := conf.Mdqi.DefaultDisplay; display != "" {
		if err := app.SetPrinterByName(display); err != nil {
			return nil, errors.Wrap(err, "failed to set default printer")
		}

		debug.Println("conf.Mdqi.DefaultDisplay =", display)
	}

	return app, nil
}

func createHistoryFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if _, err := os.Create(path); err != nil {
			return errors.Wrap(err, "failed to create history file")
		}
	}

	return nil
}

func defaultHistoryPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "failed to get current user")
	}

	return filepath.Join(usr.HomeDir, defaultHistoryFilename), err
}

func (app *App) slashCommandCategories() []string {
	defs := app.slashCommandDefinition
	keys := make([]string, 0, len(defs))

	for key := range defs {
		keys = append(keys, key)
	}

	return keys
}

func (app *App) slashCommandNames(category string) []string {
	defs := app.slashCommandDefinition[category]
	keys := make([]string, 0, len(defs))

	for key := range defs {
		keys = append(keys, key)
	}

	return keys
}

func (app *App) Run() {
	app.runLiner()
}

func (app *App) runLiner() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	app.initHistory(line)

LOOP:
	for {
		if !app.Alive {
			fmt.Println("bye")
			break LOOP
		}

		l, err := line.Prompt("mdq> ")

		switch err {
		case nil:
			l = strings.Trim(l, " \n")

			if l == "" {
				continue
			}

			line.AppendHistory(l)

			scmd, _ := ParseSlashCommand(l)
			if scmd != nil {
				app.runSlashCommand(scmd)
				continue
			}

			results, err := app.RunCmd(l, app.buildCmdArgs()...)
			if err != nil {
				logger.Println(err.Error())
			}

			Print(app.printer, results)
		case liner.ErrPromptAborted:
			logger.Println("aborted")
			break LOOP
		case io.EOF:
			fmt.Println("bye")
			break LOOP
		default:
			logger.Println("error on reading line: ", err)
			break LOOP
		}

		app.saveHistory(line)
	}
}

func (app *App) initHistory(line *liner.State) {
	if f, err := os.Open(app.historyPath); err == nil {
		line.ReadHistory(f)
		f.Close()
	} else {
		logger.Println("failed to read command history: ", err)
	}
}

func (app *App) saveHistory(line *liner.State) {
	if f, err := os.Create(app.historyPath); err == nil {
		if _, err := line.WriteHistory(f); err != nil {
			logger.Println("failed to write history: ", err)
		}

		f.Close()
	} else {
		logger.Println("failed to create history file: ", err)
	}
}

func (app *App) buildCmdArgs() []string {
	args := []string{}

	// tag
	if app.tag != "" {
		args = append(args, "--tag="+app.tag)
	}

	return args
}

func (app *App) runSlashCommand(scmd *SlashCommand) {
	sdef, err := app.FindSlashCommandDefinition(scmd.Category, scmd.Name)

	switch err {
	case nil:
		if err := sdef.Handle(app, scmd); err != nil {
			logger.Println("failed to handle slash command:", err)
		}
	case ErrSlashCommandNotFound:
		logger.Println("unknown slash command")
	}

	return
}

func (app *App) SetPrinterByName(name string) error {

	switch name {
	case "horizontal":
		app.printer = HorizontalPrinter{}
	case "vertical":
		app.printer = VerticalPrinter{}
	default:
		return ErrUnknownPrinterName
	}

	return nil
}

func (app *App) GetTag() string {
	return app.tag
}

func (app *App) SetTag(tag string) {
	app.tag = tag
}

func (app *App) ClearTag() {
	app.tag = ""
}
