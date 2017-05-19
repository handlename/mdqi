package mdqi

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseSlashCommand(t *testing.T) {
	type Expected struct {
		Result *SlashCommand
		Error  error
	}

	type Pattern struct {
		Title    string
		Query    string
		Expected Expected
	}

	patterns := []Pattern{
		Pattern{
			Title: "command only category",
			Query: "/quit",
			Expected: Expected{
				Result: &SlashCommand{
					Category: "quit",
					Name:     "",
					Args:     []string{},
				},
				Error: nil,
			},
		},
		Pattern{
			Title: "command with name",
			Query: "/tag show",
			Expected: Expected{
				Result: &SlashCommand{
					Category: "tag",
					Name:     "show",
					Args:     []string{},
				},
				Error: nil,
			},
		},
		Pattern{
			Title: "command with options",
			Query: "/tag set db1",
			Expected: Expected{
				Result: &SlashCommand{
					Category: "tag",
					Name:     "set",
					Args:     []string{"db1"},
				},
				Error: nil,
			},
		},
		Pattern{
			Title: "not a command",
			Query: "select * from item",
			Expected: Expected{
				Result: nil,
				Error:  ErrNotASlashCommand,
			},
		},
	}

	for _, pattern := range patterns {
		t.Log("test:", pattern.Title)

		r, err := ParseSlashCommand(pattern.Query)

		if !reflect.DeepEqual(r, pattern.Expected.Result) {
			t.Errorf("unexpected result: %#v", r)
		}

		if err != pattern.Expected.Error {
			t.Error("unexpected error:", err)
		}
	}
}

func TestRegisterSlashCommandDefinition(t *testing.T) {
	ErrTest := fmt.Errorf("error for test")

	app, _ := NewApp(Conf{})

	handler := func(app *App, cmd *SlashCommand) error {
		return ErrTest
	}

	if err := app.RegisterSlashCommand("tag", "", handler); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if app.slashCommandDefinition["tag"] == nil {
		t.Fatal("there are no category map")
	}

	def := app.slashCommandDefinition["tag"][""]

	if c := def.Category; c != "tag" {
		t.Fatal("unexpected category:", c)
	}

	if n := def.Name; n != "" {
		t.Fatal("unexpected name:", n)
	}

	if err := def.Handler(app, nil); err != ErrTest {
		t.Fatal("unexpected handler")
	}
}

func TestFindSlashCommandDefinition(t *testing.T) {
	app, _ := NewApp(Conf{})

	app.slashCommandDefinition["tag"] = map[string]SlashCommandDefinition{}
	app.slashCommandDefinition["tag"]["set"] = SlashCommandDefinition{
		Category: "tag",
		Name:     "set",
	}

	// registered definition
	{
		def, err := app.FindSlashCommandDefinition("tag", "set")

		if err != nil {
			t.Error("unexpected error:", err)
		}

		if def.Category != "tag" || def.Name != "set" {
			t.Errorf("unexpected definition: %+v", def)
		}
	}

	// unknown definition
	{
		def, err := app.FindSlashCommandDefinition("tag", "remove")

		if err != ErrSlashCommandNotFound {
			t.Error("expected error")
		}

		if def.Category != "" {
			t.Errorf("unexpected definition: %#v", def)
		}
	}
}

func TestSlashCommandExit(t *testing.T) {
	app, _ := NewApp(Conf{})

	SlashCommandExit(app, nil)

	if app.Alive {
		t.Fatal("app.Alive must be false.")
	}
}

func TestSlashCommandTagAdd(t *testing.T) {
	app, _ := NewApp(Conf{})

	SlashCommandTagAdd(app, &SlashCommand{
		Args: []string{"db1"},
	})

	if tag := app.GetTags()[0]; tag != "db1" {
		t.Fatal("unexpected tag:", tag)
	}
}

func TestSlashCommandTagRemove(t *testing.T) {
	app, _ := NewApp(Conf{})

	app.AddTag("db1")
	app.AddTag("db2")

	SlashCommandTagRemove(app, &SlashCommand{
		Args: []string{"db1"},
	})

	if tags := app.GetTags(); !sortEqual(tags, []string{"db2"}) {
		t.Fatalf("failed to remove tag: %+v", tags)
	}
}
