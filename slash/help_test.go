package slash

import (
	"bytes"
	"errors"
	"testing"

	"github.com/handlename/mdqi"
	"github.com/handlename/mdqi/test"
)

var ErrTest = errors.New("error for test")

type TestDefFoo struct{}

func (d TestDefFoo) Category() string { return "test" }
func (d TestDefFoo) Name() string     { return "foo" }
func (d TestDefFoo) Example() string  { return "/test foo" }
func (d TestDefFoo) Help() string     { return "I'm foo." }
func (d TestDefFoo) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	return ErrTest
}

type TestDefBar struct{}

func (d TestDefBar) Category() string { return "test" }
func (d TestDefBar) Name() string     { return "bar" }
func (d TestDefBar) Example() string  { return "/test bar" }
func (d TestDefBar) Help() string     { return "I'm bar." }
func (d TestDefBar) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	return ErrTest
}

func TestSlashCommandHelp(t *testing.T) {
	orgOutput := mdqi.DefaultOutput
	var out bytes.Buffer
	mdqi.DefaultOutput = &out
	defer func() {
		mdqi.DefaultOutput = orgOutput
	}()

	app, _ := mdqi.NewApp(mdqi.Conf{})
	app.RegisterSlashCommandDefinition(TestDefFoo{})
	app.RegisterSlashCommandDefinition(TestDefBar{})

	def := Help{}
	def.Handle(app, &mdqi.SlashCommand{})

	expect := `
+-------+----------+------+-----------+----------+
|  db   | category | name |  example  |   help   |
+-------+----------+------+-----------+----------+
| (mdq) | test     | bar  | /test bar | I'm bar. |
| (mdq) | test     | foo  | /test foo | I'm foo. |
+-------+----------+------+-----------+----------+
`

	if s := out.String(); !test.CompareAfterTrim(s, expect) {
		t.Fatalf("unexpected output:\n%s", s)
	}
}
