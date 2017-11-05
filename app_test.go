package mdqi

import (
	"testing"

	"github.com/handlename/mdqi/test"
)

func TestManageTag(t *testing.T) {
	app, _ := NewApp(Conf{})

	if tag := app.GetTag(); tag != "" {
		t.Error("unexpected tag:", tag)
	}

	app.SetTag("foo")

	if tag := app.GetTag(); tag != "foo" {
		t.Error("unexpected tag:", tag)
	}

	app.SetTag("bar")

	if tag := app.GetTag(); tag != "bar" {
		t.Error("unexpected tag:", tag)
	}

	app.ClearTag()

	if tag := app.GetTag(); tag != "" {
		t.Error("unexpected tag:", tag)
	}
}

func TestBulidCmdArgs(t *testing.T) {
	app, _ := NewApp(Conf{})

	if args := app.buildCmdArgs(); !test.SortEqual(args, []string{}) {
		t.Errorf("unexpected args: %#v", args)
	}

	app.SetTag("foo")

	if args := app.buildCmdArgs(); !test.SortEqual(args, []string{"--tag=foo"}) {
		t.Errorf("unexpected args: %+v", args)
	}
}

func TestSetPrinterByName(t *testing.T) {
	app, _ := NewApp(Conf{})

	// horizontal
	{
		if err := app.SetPrinterByName("horizontal"); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		switch ty := app.Printer.(type) {
		case HorizontalPrinter:
		default:
			t.Errorf("unexpected printer type: %s", ty)
		}
	}

	// horizontal
	{
		if err := app.SetPrinterByName("vertical"); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		switch ty := app.Printer.(type) {
		case VerticalPrinter:
		default:
			t.Errorf("unexpected printer type: %s", ty)
		}
	}

	// unknown
	{
		if err := app.SetPrinterByName("unknown"); err != ErrUnknownPrinterName {
			t.Errorf("unexpected error: %s", err)
		}
	}
}

func TestNewAppDefaultTag(t *testing.T) {
	conf := Conf{
		Mdqi: ConfMdqi{
			DefaultTag: "db1",
		},
	}

	app, err := NewApp(conf)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if tag := app.GetTag(); tag != "db1" {
		t.Error("unexpected tag:", tag)
	}
}

func TestNewAppDefaultDisplay(t *testing.T) {
	conf := Conf{
		Mdqi: ConfMdqi{
			DefaultDisplay: "vertical",
		},
	}

	app, err := NewApp(conf)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	switch ty := app.Printer.(type) {
	case VerticalPrinter:
	default:
		t.Errorf("unexpected printer type: %s", ty)
	}
}
