package slash

import (
	"fmt"

	"github.com/handlename/mdqi"
)

type TagSet struct{}

func (c TagSet) Category() string { return "tag" }
func (c TagSet) Name() string     { return "set" }
func (c TagSet) Example() string  { return "/tag set {tagname}" }

func (c TagSet) Help() string {
	return "Set tag to pass to mdq's --tag option."
}

func (c TagSet) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	if len(cmd.Args) != 1 {
		fmt.Println(c.Help())
		return nil
	}

	app.SetTag(cmd.Args[0])
	mdqi.Debug.Println("tag stored:", app.GetTag())

	return nil
}

type TagClear struct{}

func (c TagClear) Category() string { return "tag" }
func (c TagClear) Name() string     { return "clear" }
func (c TagClear) Example() string  { return "/tag clear {tagname}" }

func (c TagClear) Help() string {
	return "Clear tag to pass to mdq's --tag option."
}

func (c TagClear) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	app.ClearTag()
	mdqi.Debug.Printf("tag cleard: %+v\n", app.GetTag())

	return nil
}

type TagShow struct{}

func (c TagShow) Category() string { return "tag" }
func (c TagShow) Name() string     { return "show" }
func (c TagShow) Example() string  { return "/tag show" }

func (c TagShow) Help() string {
	return "Show registered tags to pass to mdq's --tag option."
}

func (c TagShow) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	results := []mdqi.Result{
		mdqi.Result{
			Database: "(mdq)",
			Columns:  []string{"tag"},
			Rows: []map[string]interface{}{
				map[string]interface{}{"tag": app.GetTag()},
			},
		},
	}

	mdqi.Print(app.Printer, results)

	return nil
}
