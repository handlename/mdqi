package mdqi

import "github.com/morikuni/mdq"

type App struct {
	// CmdPath is path to mdq command.
	CmdPath string
}

type Conf struct {
}

type Result struct {
	mdq.Result
}

func NewApp(conf Conf) (*App, error) {
	// TODO: Check if mdq command exists by exec.LookPath.

	return &App{
		CmdPath: "mdq",
	}, nil
}
