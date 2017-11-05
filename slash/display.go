package slash

import "github.com/handlename/mdqi"

type Display struct{}

func (c Display) Category() string { return "display" }
func (c Display) Name() string     { return "set" }
func (c Display) Example() string  { return "/display set (horizontal|vertical)" }

func (c Display) Help() string {
	return "Set display format."
}

func (c Display) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	if len(cmd.Args) != 1 {
		return mdqi.ErrSlashCommandInvalidArgs
	}

	if err := app.SetPrinterByName(cmd.Args[0]); err != nil {
		mdqi.Debug.Printf("error on SetPrinterByName: %s", err)
		return mdqi.ErrSlashCommandInvalidArgs
	}

	mdqi.Debug.Printf("printer has been changed to %s", cmd.Args[0])

	return nil
}

type ToggleDisplay struct{}

func (c ToggleDisplay) Category() string { return "v" }
func (c ToggleDisplay) Name() string     { return "" }
func (c ToggleDisplay) Example() string  { return "/v" }

func (c ToggleDisplay) Help() string {
	return "Toggle display type between horizontal and vertical."
}

func (c ToggleDisplay) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	name := "horizontal"

	switch app.Printer.(type) {
	case mdqi.HorizontalPrinter:
		name = "vertical"
	case mdqi.VerticalPrinter:
		name = "horizontal"
	}

	app.SetPrinterByName(name)
	mdqi.Debug.Println("printer has been changed to", name)

	return nil
}
