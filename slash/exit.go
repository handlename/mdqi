package slash

import "github.com/handlename/mdqi"

type Exit struct{}

func (c Exit) Category() string { return "exit" }
func (c Exit) Name() string     { return "" }
func (c Exit) Example() string  { return "/exit" }

func (c Exit) Help() string {
	return "Exit mdqi."
}

func (c Exit) Handle(app *mdqi.App, cmd *mdqi.SlashCommand) error {
	app.Alive = false
	return nil
}
