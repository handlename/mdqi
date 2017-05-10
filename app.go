package mdqi

type App struct {
	// cmdPath is path to mdq command.
	cmdPath string
}

type Conf struct {
}

type Result struct {
	Database string
	Columns  []string
	Rows     []map[string]interface{}
}

func NewApp(conf Conf) (*App, error) {
	// TODO: Check if mdq command exists by exec.LookPath.

	return &App{
		cmdPath: "mdq",
	}, nil
}
