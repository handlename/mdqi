package slash

import (
	"testing"

	"github.com/handlename/mdqi"
)

func TestExit(t *testing.T) {
	app, _ := mdqi.NewApp(mdqi.Conf{})

	def := Exit{}
	def.Handle(app, nil)

	if app.Alive {
		t.Fatal("app.Alive must be false.")
	}
}
