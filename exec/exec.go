package exec

import (
	"os"
	"os/exec"
	"strings"
)

type Executor interface {
	Run() (*Result, error)
}

type CommandOpt func(o *Command)

type Command struct {
	command string
	args    []string
	cwd     string
	env     []string
	stdin   *strings.Reader
	stdout  strings.Builder
	stderr  strings.Builder
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// NewCommandOpt crea una nueva instancia de Command
func NewCommandOpt(opts ...CommandOpt) *Command {
	c := &Command{
		stdin:  strings.NewReader(""),
		stdout: strings.Builder{},
		stderr: strings.Builder{},
		env:    os.Environ(),
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Command) Run() (*Result, error) {
	cmd := exec.Command(c.command, c.args...)
	cmd.Dir = c.cwd
	cmd.Env = c.env
	cmd.Stdin = c.stdin
	cmd.Stdout = &c.stdout
	cmd.Stderr = &c.stderr

	err := cmd.Run()
	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		exitCode = 1
	}

	return &Result{
		Stdout:   c.stdout.String(),
		Stderr:   c.stderr.String(),
		ExitCode: exitCode,
	}, err
}

// WithCommand asigna command a Command
func WithCommand(command string) CommandOpt {
	return func(c *Command) {
		c.command = command
	}
}

// WithArgs asigna args a Command
func WithArgs(args []string) CommandOpt {
	return func(c *Command) {
		c.args = args
	}
}

// WithCwd asigna cwd a Command
func WithCwd(cwd string) CommandOpt {
	return func(c *Command) {
		c.cwd = cwd
	}
}

// WithEnv agrega env a Command
func WithEnv(env []string) CommandOpt {
	return func(c *Command) {
		c.env = append(c.env, env...)
	}
}

// WithStdin asigna stdin a Command
func WithStdin(stdin *strings.Reader) CommandOpt {
	return func(c *Command) {
		c.stdin = stdin
	}
}

// WithStdout asigna stdout a Command
func WithStdout(stdout strings.Builder) CommandOpt {
	return func(c *Command) {
		c.stdout = stdout
	}
}
