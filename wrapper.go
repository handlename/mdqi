package mdqi

import (
	"encoding/json"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

func (app *App) RunCmd(query string, args ...string) (results []Result, err error) {
	// Run mdq command
	out, err := runCmd(app.CmdPath, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run command")
	}

	// Parse result
	results, err = parseCmdOutput(out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse command output")
	}

	return results, nil
}

func runCmd(path string, query string, args ...string) (out []byte, err error) {
	cmd := exec.Command(path, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open stdin pipe for command %s", path)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, query)
	}()

	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run command %s", path)
	}

	return out, nil
}

func parseCmdOutput(out []byte) (results []Result, err error) {
	if err = json.Unmarshal(out, &results); err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", string(out))
	}

	return results, nil
}
