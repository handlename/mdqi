package slash

import (
	"testing"

	"github.com/handlename/mdqi"
)

func TestDisplay(t *testing.T) {
	app, _ := mdqi.NewApp(mdqi.Conf{})

	def := Display{}

	{
		// invalid format

		if err := def.Handle(app, &mdqi.SlashCommand{Args: []string{"foo"}}); err != mdqi.ErrSlashCommandInvalidArgs {
			t.Error("must be error")
		}
	}

	{
		// format: vertical

		if err := def.Handle(app, &mdqi.SlashCommand{Args: []string{"vertical"}}); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		switch ty := app.Printer.(type) {
		case mdqi.VerticalPrinter:
		default:
			t.Errorf("unexpected printer type: %s", ty)
		}
	}

	{
		// format: horizontal

		if err := def.Handle(app, &mdqi.SlashCommand{Args: []string{"horizontal"}}); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		switch ty := app.Printer.(type) {
		case mdqi.HorizontalPrinter:
		default:
			t.Errorf("unexpected printer type: %s", ty)
		}
	}
}

func TestToggleDisplay(t *testing.T) {
	app, _ := mdqi.NewApp(mdqi.Conf{})
	app.SetPrinterByName("horizontal")

	def := ToggleDisplay{}

	if err := def.Handle(app, &mdqi.SlashCommand{}); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if n := app.Printer.Name(); n != "vertical" {
		t.Error("unexpected printer:", n)
	}

	if err := def.Handle(app, &mdqi.SlashCommand{}); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if n := app.Printer.Name(); n != "horizontal" {
		t.Error("unexpected printer:", n)
	}
}
