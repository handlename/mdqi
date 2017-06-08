package mdqi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func (app *App) RunCmd(query string, args ...string) (results []Result, err error) {
	// Run mdq command
	out, err := runCmd(app.mdqPath, query, args...)
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
	debug.Printf("runs command: echo '%s' | %s %s", query, path, strings.Join(args, " "))

	cmd := exec.Command(path, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open stdin pipe for command %s", path)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get stdout pipe for command %s", path)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get stderr pipe for command %s", path)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, query)
	}()

	if err = cmd.Start(); err != nil {
		return nil, errors.Wrapf(err, "failed to run command %s", path)
	}

	if o, _ := ioutil.ReadAll(stderr); 0 < len(o) {
		logger.Println(string(o))
	}

	if out, err = ioutil.ReadAll(stdout); err != nil {
		return nil, errors.Wrap(err, "failed to read stdout")
	}

	return out, nil
}

func parseCmdOutput(out []byte) (results []Result, err error) {
	if err = json.Unmarshal(out, &results); err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", string(out))
	}

	return results, nil
}
