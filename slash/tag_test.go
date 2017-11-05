package slash

import (
	"bytes"
	"testing"

	"github.com/handlename/mdqi"
	"github.com/handlename/mdqi/test"
)

func TestTagAdd(t *testing.T) {
	app, _ := mdqi.NewApp(mdqi.Conf{})

	def := TagSet{}
	def.Handle(app, &mdqi.SlashCommand{
		Args: []string{"db1"},
	})

	if tag := app.GetTag(); tag != "db1" {
		t.Fatal("unexpected tag:", tag)
	}
}

func TestTagClear(t *testing.T) {
	app, _ := mdqi.NewApp(mdqi.Conf{})

	app.SetTag("db1")

	def := TagClear{}
	def.Handle(app, &mdqi.SlashCommand{})

	if tag := app.GetTag(); tag != "" {
		t.Fatal("failed to remove tag:", tag)
	}
}

func TestTagShow(t *testing.T) {
	orgOutput := mdqi.DefaultOutput
	var out bytes.Buffer
	mdqi.DefaultOutput = &out
	defer func() {
		mdqi.DefaultOutput = orgOutput
	}()

	app, _ := mdqi.NewApp(mdqi.Conf{})

	app.SetTag("db1")

	def := TagShow{}
	def.Handle(app, &mdqi.SlashCommand{})

	expect := `
+-------+-----+
|  db   | tag |
+-------+-----+
| (mdq) | db1 |
+-------+-----+
`

	if s := out.String(); !test.CompareAfterTrim(s, expect) {
		t.Fatalf("unexpected output:\n%s", s)
	}
}
