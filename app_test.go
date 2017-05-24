package mdqi

import (
	"testing"
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
}