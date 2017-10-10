package cli

import (
	"bytes"
	"os/exec"
	"syscall"
)

type ErrExec struct {
	Output Output
}

func (e *ErrExec) Error() string {
	return e.Output.StdErr
}

type Output struct {
	ExitCode int         `json:"exitCode"`
	StdOut   string      `json:"stdOut"`
	StdErr   string      `json:"stdErr"`
	Err      error       `json:"error"`
	Data     interface{} `json:"data"`
}

func Execute(cmd string, args ...string) Output {
	var outStream bytes.Buffer
	var errStream bytes.Buffer

	c := exec.Command(cmd, args...)
	c.Stdin = nil
	c.Stdout = &outStream
	c.Stderr = &errStream

	output := Output{}
	output.Err = c.Run()
	output.StdOut = string(outStream.Bytes())
	output.StdErr = string(errStream.Bytes())

	if exitErr, ok := output.Err.(*exec.ExitError); ok {
		if waitStatus, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			output.ExitCode = waitStatus.ExitStatus()
		}
	}

	return output
}
