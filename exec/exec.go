package exec

import (
	"os/exec"
	"strings"
)

type Output struct {
	Out string
	Err string
}

func RunCommand(args ...string) (Output, error) {
	baseCmd := args[0]
	cmdArgs := args[1:]

	sout, serr := new(strings.Builder), new(strings.Builder)

	cmd := exec.Command(baseCmd, cmdArgs...)
	cmd.Stdout = sout
	cmd.Stderr = serr

	err := cmd.Run()
	if err != nil {
		return Output{}, err
	}

	return Output{
		Out: sout.String(),
		Err: serr.String(),
	}, nil
}
