package slash

import (
	"sort"

	"github.com/handlename/mdqi"
)

type Help struct{}

func (c Help) Category() string { return "help" }
func (c Help) Name() string     { return "" }
func (c Help) Example() string  { return "/help" }

func (c Help) Help() string {
	return "Show this help."
}

func (c Help) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	rows := []map[string]interface{}{}

	categories := app.SlashCommandCategories()
	sort.Strings(categories)

	for _, category := range categories {
		// inCategory := app.SlashCommandDefinition[category]

		names := app.SlashCommandNames(category)
		sort.Strings(names)

		for _, name := range names {
			// def := inCategory[name]
			def, _ := app.FindSlashCommandDefinition(category, name)

			rows = append(rows, map[string]interface{}{
				"category": category,
				"name":     name,
				"example":  def.Example(),
				"help":     def.Help(),
			})
		}
	}

	results := []mdqi.Result{
		mdqi.Result{
			Database: "(mdq)",
			Columns:  []string{"category", "name", "example", "help"},
			Rows:     rows,
		},
	}

	mdqi.Print(app.Printer, results)

	return nil
}
