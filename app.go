package mdqi

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"os/exec"

	"github.com/peterh/liner"
	"github.com/pkg/errors"
)

const (
	defaultMdqPath         = "mdq"
	defaultHistoryFilename = ".mdqi_history"
)

var Version string

var (
	ErrSlashCommandNotFound    = errors.New("unknown SlashCommand")
	ErrNotASlashCommand        = errors.New("there are no SlashCommand")
	ErrSlashCommandInvalidArgs = errors.New("invalid args")
	ErrUnknownPrinterName      = errors.New("unknown printer name")
)

type App struct {
	// Alive turns into false, mdqi will exit.
	Alive bool

	// mdqPath is path to mdq command.
	mdqPath string

	// mdqConfigPath is path to configuration file for mdq command.
	mdqConfigPath string

	// historyPath is path to command history file for liner.
	historyPath string

	// slashCommandDefinition holds SlashCommandDefinition.
	// app.slashCommandDefinition[category][name] = SlashCommandDefinition
	slashCommandDefinition map[string]map[string]SlashCommandDefinition

	// tag stores tag value for --tag option of mdq.
	tag string

	// Printer prints results.
	Printer Printer
}

type Result struct {
	Database string
	Columns  []string
	Rows     []map[string]interface{}
}

func init() {
	DefaultOutput = os.Stdout
}

func NewApp(conf Conf) (*App, error) {
	// validate mdq path
	mdqPath := defaultMdqPath
	if path := conf.Mdq.Bin; path != "" {
		Debug.Println("conf.Mdq.Bin =", path)

		abs, err := ExpandPath(path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to expand path to mdq")
		}

		if err := lookMdqPath(abs); err != nil {
			return nil, errors.Wrapf(err, "mdq command not found at %s", path)
		}

		mdqPath = abs
	}

	// mdq config path
	mdqConfPath := conf.Mdq.Config
	if path := conf.Mdq.Config; path != "" {
		Debug.Println("conf.Mdq.Config =", path)

		abs, err := ExpandPath(path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to expand path to mdq conf")
		}

		mdqConfPath = abs
	}

	// create history file
	historyPath := conf.Mdqi.History
	if historyPath == "" {
		Debug.Println("conf.Mdqi.History =", historyPath)

		var err error
		if historyPath, err = defaultHistoryPath(); err != nil {
			return nil, errors.Wrap(err, "failed to create history file at default path")
		}
	}
	{
		abs, err := ExpandPath(historyPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to expand path to history file")
		}
		historyPath = abs
	}
	if err := createHistoryFile(historyPath); err != nil {
		return nil, errors.Wrapf(err, "failed to create history file at %s", historyPath)
	}

	app := &App{
		Alive: true,

		mdqPath:                mdqPath,
		mdqConfigPath:          mdqConfPath,
		historyPath:            historyPath,
		slashCommandDefinition: map[string]map[string]SlashCommandDefinition{},
		Printer:                HorizontalPrinter{},
	}

	// set default tag
	if tag := conf.Mdqi.DefaultTag; tag != "" {
		app.SetTag(tag)
		Debug.Println("conf.Mdqi.DefaultTag =", tag)
	}

	// set default display
	if display := conf.Mdqi.DefaultDisplay; display != "" {
		if err := app.SetPrinterByName(display); err != nil {
			return nil, errors.Wrap(err, "failed to set default printer")
		}

		Debug.Println("conf.Mdqi.DefaultDisplay =", display)
	}

	return app, nil
}

func createHistoryFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if _, err := os.Create(path); err != nil {
			return errors.Wrapf(err, "failed to create history file at %s", path)
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

func lookMdqPath(path string) error {
	_, err := exec.LookPath(path)

	return err
}

func (app *App) SlashCommandCategories() []string {
	defs := app.slashCommandDefinition
	keys := make([]string, 0, len(defs))

	for key := range defs {
		keys = append(keys, key)
	}

	return keys
}

func (app *App) SlashCommandNames(category string) []string {
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

	// If input dosen't be match it, input is not finished.
	rgxLineFinisehd := regexp.MustCompile(";$")

	// If this value is false, input will be continued.
	lineFinished := true

	var input string // input from line
	var err error

LOOP:
	for {
		if !app.Alive {
			fmt.Println("bye")
			break LOOP
		}

		if lineFinished {
			input, err = line.Prompt(app.formatPrompt(false))
		} else {
			var l string
			l, err = line.Prompt(app.formatPrompt(true))
			input = strings.Join([]string{input, l}, " ")
		}

		switch err {
		case nil:
			input = strings.Trim(input, " \n")

			if input == "" {
				continue
			}

			// run slash command
			scmd, _ := ParseSlashCommand(input)
			if scmd != nil {
				app.runSlashCommand(scmd)
				continue
			}

			if lineFinished = rgxLineFinisehd.MatchString(input); lineFinished {
			} else {
				// If line is not finished, read next line as continue.
				continue
			}

			line.AppendHistory(input)

			// run mdq
			results, err := app.RunCmd(input, app.buildCmdArgs()...)
			if err != nil {
				Logger.Println(err.Error())
			}

			Print(app.Printer, results)
		case liner.ErrPromptAborted:
			Logger.Println("aborted")
			break LOOP
		case io.EOF:
			fmt.Println("bye")
			break LOOP
		default:
			Logger.Println("error on reading line: ", err)
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
		Logger.Printf("failed to read command history at %s: %s", app.historyPath, err)
	}
}

func (app *App) saveHistory(line *liner.State) {
	if f, err := os.Create(app.historyPath); err == nil {
		if _, err := line.WriteHistory(f); err != nil {
			Logger.Println("failed to write history: ", err)
		}

		f.Close()
	} else {
		Logger.Printf("failed to create history file at %s: %s", app.historyPath, err)
	}
}

func (app *App) buildCmdArgs() []string {
	args := []string{}

	// config
	if path := app.mdqConfigPath; path != "" {
		args = append(args, "--config="+path)
	}

	// tag
	if tag := app.tag; tag != "" {
		args = append(args, "--tag="+tag)
	}

	return args
}

func (app *App) runSlashCommand(scmd *SlashCommand) {
	sdef, err := app.FindSlashCommandDefinition(scmd.Category, scmd.Name)

	switch err {
	case nil:
		if err := sdef.Handle(app, scmd); err != nil {
			Logger.Println("failed to handle slash command:", err)
		}
	case ErrSlashCommandNotFound:
		Logger.Println("unknown slash command")
	}

	return
}

func (app *App) SetPrinterByName(name string) error {

	switch name {
	case "horizontal":
		app.Printer = HorizontalPrinter{}
	case "vertical":
		app.Printer = VerticalPrinter{}
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

func (app *App) formatPrompt(continues bool) string {
	if continues {
		return "   | "
	} else {
		if tag := app.GetTag(); tag == "" {
			return "mdq> "
		} else {
			return fmt.Sprintf("mdq[%s]> ", app.GetTag())
		}
	}
}

// ExpandPath expands file path like `~/path/to/foo`
func ExpandPath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "failed to get current user")
	}

	return strings.Replace(path, "~", usr.HomeDir, 1), nil
}
