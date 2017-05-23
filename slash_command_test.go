package mdqi

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

var ErrTest = errors.New("error for test")

type TestDefFoo struct{}

func (d TestDefFoo) Category() string { return "test" }
func (d TestDefFoo) Name() string     { return "foo" }
func (d TestDefFoo) Example() string  { return "/test foo" }
func (d TestDefFoo) Help() string     { return "I'm foo." }
func (d TestDefFoo) Handle(app *App, cmd *SlashCommand) error {
	return ErrTest
}

type TestDefBar struct{}

func (d TestDefBar) Category() string { return "test" }
func (d TestDefBar) Name() string     { return "bar" }
func (d TestDefBar) Example() string  { return "/test bar" }
func (d TestDefBar) Help() string     { return "I'm bar." }
func (d TestDefBar) Handle(app *App, cmd *SlashCommand) error {
	return ErrTest
}

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
	app, _ := NewApp(Conf{})

	if err := app.RegisterSlashCommandDefinition(TestDefFoo{}); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if app.slashCommandDefinition["test"] == nil {
		t.Fatal("there are no category map")
	}

	def := app.slashCommandDefinition["test"]["foo"]

	if def == nil {
		t.Fatal("there are no definition")
	}

	if c := def.Category(); c != "test" {
		t.Fatal("unexpected category:", c)
	}

	if n := def.Name(); n != "foo" {
		t.Fatal("unexpected name:", n)
	}

	if err := def.Handle(app, nil); err != ErrTest {
		t.Fatalf("unexpected handler returns error: %s", err)
	}
}

func TestFindSlashCommandDefinition(t *testing.T) {
	app, _ := NewApp(Conf{})
	app.slashCommandDefinition["test"] = map[string]SlashCommandDefinition{}
	app.slashCommandDefinition["test"]["foo"] = TestDefFoo{}

	// registered definition
	{
		def, err := app.FindSlashCommandDefinition("test", "foo")

		if err != nil {
			t.Error("unexpected error:", err)
		}

		if def == nil {
			t.Error("definition must be returned.")
		}

		if def.Category() != "test" || def.Name() != "foo" {
			t.Errorf("unexpected definition: %#v", def)
		}
	}

	// unknown definition
	{
		def, err := app.FindSlashCommandDefinition("test", "bar")

		if err != ErrSlashCommandNotFound {
			t.Error("expected error")
		}

		if def != nil {
			t.Errorf("unexpected definition: %#v", def)
		}
	}
}

func TestSlashCommandExit(t *testing.T) {
	app, _ := NewApp(Conf{})

	def := SlashCommandExit{}
	def.Handle(app, nil)

	if app.Alive {
		t.Fatal("app.Alive must be false.")
	}
}

func TestSlashCommandTagAdd(t *testing.T) {
	app, _ := NewApp(Conf{})

	def := SlashCommandTagAdd{}
	def.Handle(app, &SlashCommand{
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

	def := SlashCommandTagRemove{}
	def.Handle(app, &SlashCommand{
		Args: []string{"db1"},
	})

	if tags := app.GetTags(); !sortEqual(tags, []string{"db2"}) {
		t.Fatalf("failed to remove tag: %+v", tags)
	}
}

func TestSlashCommandTagShow(t *testing.T) {
	orgOutput := defaultOutput
	var out bytes.Buffer
	defaultOutput = &out
	defer func() {
		defaultOutput = orgOutput
	}()

	app, _ := NewApp(Conf{})

	app.AddTag("db1")
	app.AddTag("db2")

	def := SlashCommandTagShow{}
	def.Handle(app, &SlashCommand{
		Args: []string{"db1"},
	})

	expect := `
+-------+-----+
|  db   | tag |
+-------+-----+
| (mdq) | db1 |
| (mdq) | db2 |
+-------+-----+
`

	if s := out.String(); !compareAfterTrim(s, expect, " \n") {
		t.Fatalf("unexpected output:\n%s", s)
	}
}

func TestSlashCommandHelp(t *testing.T) {
	orgOutput := defaultOutput
	var out bytes.Buffer
	defaultOutput = &out
	defer func() {
		defaultOutput = orgOutput
	}()

	app, _ := NewApp(Conf{})
	app.slashCommandDefinition = map[string]map[string]SlashCommandDefinition{}

	app.RegisterSlashCommandDefinition(TestDefFoo{})
	app.RegisterSlashCommandDefinition(TestDefBar{})

	def := SlashCommandHelp{}
	def.Handle(app, &SlashCommand{})

	expect := `
+-------+----------+------+-----------+----------+
|  db   | category | name |  example  |   help   |
+-------+----------+------+-----------+----------+
| (mdq) | test     | foo  | /test foo | I'm foo. |
| (mdq) | test     | bar  | /test bar | I'm bar. |
+-------+----------+------+-----------+----------+
`

	if s := out.String(); !compareAfterTrim(s, expect, " \n") {
		t.Fatalf("unexpected output:\n%s", s)
	}
}
