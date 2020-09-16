package goexec

import (
	"os"
	"os/exec"

	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
)

// Cmd provides a fluent, decorator based API for os/exec.
//
// 	Cmd("ping", Args("-c", "4", "1.1.1.1"))
func Cmd(cmd string, decorators ...func(*exec.Cmd) error) (c *exec.Cmd, err error) {
	defer goerr.Handle(func(e error) { c = nil; err = e })

	c = exec.Command(cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	for _, decorator := range decorators {
		goerr.Check(decorator(c))
	}

	if c.Env == nil {
		c.Env = os.Environ()
	}

	return
}

// MustCmd does the same as Cmd but panics if an error is encountered.
func MustCmd(cmd string, decorators ...func(*exec.Cmd) error) *exec.Cmd {
	c, e := Cmd(cmd, decorators...)
	goerr.Check(e)
	return c
}

// Run is a convenience function for simple cases.
//
// Instead of:
// 	Cmd("ping", Args("-c", "4", "8.8.8.8")).Run()
//
// You might write:
// 	Run("ping", "-c", "4", "8.8.8.8")
func Run(cmd string, args ...string) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	c, err := Cmd(cmd, Args(args...))
	goerr.Check(err)

	goerr.Check(c.Run())
	return
}

// MustRun does the same as Run but panics if an error is encountered.
func MustRun(cmd string, args ...string) {
	goerr.Check(Run(cmd, args...))
}

// RunAsync does the same as Run but does so asynchronously.
func RunAsync(cmd string, args ...string) *task.Task {
	return task.New(func(t *task.Internal) {
		MustRun(cmd, args...)
	})
}

// RunBuffered is a convenience function for simple cases.
//
// Instead of:
// 	RunBufferedCmd(Cmd("ping", Args("-c", "4", "8.8.8.8")))
//
// You might write:
// 	RunBuffered("ping", "-c", "4", "8.8.8.8")
//
// NOTE: `RunBuffered()` returns stdOut and stdErr as strings if you want the
// untouched byte arrays you could use `RunBufferedCmd()` instead.
func RunBuffered(cmd string, args ...string) (out *StdStrings, err error) {
	defer goerr.Handle(func(e error) { out = nil; err = e })

	c, err := Cmd(cmd, Args(args...))
	goerr.Check(err)

	o, err := RunBufferedCmd(c)
	goerr.Check(err)

	return &StdStrings{string(o.StdOut), string(o.StdErr)}, err
}

// MustRunBuffered does the same as RunBuffered but panics if an error is encountered.
func MustRunBuffered(cmd string, args ...string) *StdStrings {
	out, err := RunBuffered(cmd, args...)
	goerr.Check(err)
	return out
}

// RunBufferedAsync does the same as RunBuffered but does so asynchronously.
func RunBufferedAsync(cmd string, args ...string) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustRunBuffered(cmd, args...))
	})
}

// RunPrefixed is a convenience function for simple cases.
//
// Instead of:
// 	RunPrefixedCmd("foo", Cmd("ping", Args("-c", "4", "8.8.8.8")))
//
// You might write:
// 	RunPrefixed("foo", "ping", "-c", "4", "8.8.8.8")
func RunPrefixed(prefix, cmd string, args ...string) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	c, err := Cmd(cmd, Args(args...))
	goerr.Check(err)

	goerr.Check(RunPrefixedCmd(prefix, c))
	return
}

// MustRunPrefixed does the same as RunPrefixed but panics if an error is encountered.
func MustRunPrefixed(prefix, cmd string, args ...string) {
	goerr.Check(RunPrefixed(prefix, cmd, args...))
}

// RunPrefixedAsync does the same as RunPrefixed but does so asynchronously.
func RunPrefixedAsync(prefix, cmd string, args ...string) *task.Task {
	return task.New(func(t *task.Internal) {
		MustRunPrefixed(prefix, cmd, args...)
	})
}

// RunPrefixedCmd will prefix all StdOut and StdErr with given prefix.
// This is useful when running many commands concurrently,
// output will look similar to docker-compose.
func RunPrefixedCmd(prefix string, cmd *exec.Cmd) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	stdOutPipeR, stdOutPipeW, err := os.Pipe()
	goerr.Check(err)
	cmd.Stdout = stdOutPipeW

	stdErrPipeR, stdErrPipeW, err := os.Pipe()
	goerr.Check(err)
	cmd.Stderr = stdErrPipeW

	return prefixed(prefix,
		os.Stdout, os.Stderr,
		stdOutPipeR, stdErrPipeR,
		stdOutPipeW, stdErrPipeW, func() error {
			return cmd.Run()
		},
	)
}

// MustRunPrefixedCmd does the same as RunPrefixedCmd but panics if an error is encountered.
func MustRunPrefixedCmd(prefix string, cmd *exec.Cmd) {
	goerr.Check(RunPrefixedCmd(prefix, cmd))
}

// RunPrefixedCmdAsync does the same as RunPrefixedCmd but does so asynchronously.
func RunPrefixedCmdAsync(prefix string, cmd *exec.Cmd) *task.Task {
	return task.New(func(t *task.Internal) {
		RunPrefixedCmd(prefix, cmd)
	})
}

// RunBufferedCmd will buffer all StdOut and StdErr, returning the buffers.
// This is useful when you wish to parse the results of a command.
func RunBufferedCmd(cmd *exec.Cmd) (out *StdBytes, err error) {
	defer goerr.Handle(func(e error) { out = nil; err = e })

	stdOutPipeR, stdOutPipeW, err := os.Pipe()
	goerr.Check(err)
	cmd.Stdout = stdOutPipeW

	stdErrPipeR, stdErrPipeW, err := os.Pipe()
	goerr.Check(err)
	cmd.Stderr = stdErrPipeW

	return buffered(
		os.Stdout, os.Stderr,
		stdOutPipeR, stdErrPipeR,
		stdOutPipeW, stdErrPipeW, func() error {
			return cmd.Run()
		},
	)
}

// MustRunBufferedCmd does the same as RunBufferedCmd but panics if an error is encountered.
func MustRunBufferedCmd(cmd *exec.Cmd) *StdBytes {
	out, err := RunBufferedCmd(cmd)
	goerr.Check(err)
	return out
}

// RunBufferedCmdAsync does the same as RunBufferedCmd but does so asynchronously.
func RunBufferedCmdAsync(cmd *exec.Cmd) *task.Task {
	return task.New(func(t *task.Internal) {
		t.Resolve(MustRunBufferedCmd(cmd))
	})
}

// Pipe will send the output of the first command
// to the input of the second and so on.
func Pipe(cmds ...*exec.Cmd) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	for key, cmd := range cmds {
		if key > 0 {
			cmds[key-1].Stdout = nil
			pipe, err := cmds[key-1].StdoutPipe()
			goerr.Check(err)
			cmd.Stdin = pipe
		}
	}

	for _, cmd := range cmds {
		goerr.Check(cmd.Start())
	}

	for _, cmd := range cmds {
		goerr.Check(cmd.Wait())
	}

	return
}

// MustPipe does the same as Pipe but panics if an error is encountered.
func MustPipe(cmds ...*exec.Cmd) {
	goerr.Check(Pipe(cmds...))
}

// PipeAsync does the same as Pipe but does so asynchronously.
func PipeAsync(cmds ...*exec.Cmd) *task.Task {
	return task.New(func(t *task.Internal) {
		MustPipe(cmds...)
	})
}

// PipePrefixed will prefix all StdOut and StdErr with given prefix.
// This is useful when running many pipes concurrently,
// output will look similar to docker-compose.
func PipePrefixed(prefix string, cmds ...*exec.Cmd) (err error) {
	defer goerr.Handle(func(e error) { err = e })

	stdOutPipeR, stdOutPipeW, err := os.Pipe()
	goerr.Check(err)

	stdErrPipeR, stdErrPipeW, err := os.Pipe()
	goerr.Check(err)

	for _, cmd := range cmds {
		cmd.Stderr = stdErrPipeW
	}

	cmds[len(cmds)-1].Stdout = stdOutPipeW

	return prefixed(prefix,
		os.Stdout, os.Stderr,
		stdOutPipeR, stdErrPipeR,
		stdOutPipeW, stdErrPipeW, func() error {
			return Pipe(cmds...)
		},
	)
}

// MustPipePrefixed does the same as PipePrefixed but panics if an error is encountered.
func MustPipePrefixed(prefix string, cmds ...*exec.Cmd) {
	goerr.Check(PipePrefixed(prefix, cmds...))
}

// PipePrefixedAsync does the same as PipePrefixed but does so asynchronously.
func PipePrefixedAsync(prefix string, cmds ...*exec.Cmd) *task.Task {
	return task.New(func(t *task.Internal) {
		MustPipePrefixed(prefix, cmds...)
	})
}

// PipeBuffered will buffer all StdOut and StdErr, returning the buffers.
// This is useful when you wish to parse the results of a pipe.
func PipeBuffered(cmds ...*exec.Cmd) (out *StdBytes, err error) {
	defer goerr.Handle(func(e error) { out = nil; err = e })

	stdOutPipeR, stdOutPipeW, err := os.Pipe()
	goerr.Check(err)

	stdErrPipeR, stdErrPipeW, err := os.Pipe()
	goerr.Check(err)

	for _, cmd := range cmds {
		cmd.Stderr = stdErrPipeW
	}

	cmds[len(cmds)-1].Stdout = stdOutPipeW

	return buffered(
		os.Stdout, os.Stderr,
		stdOutPipeR, stdErrPipeR,
		stdOutPipeW, stdErrPipeW, func() error {
			return Pipe(cmds...)
		},
	)
}

// MustPipeBuffered does the same as PipeBuffered but panics if an error is encountered.
func MustPipeBuffered(cmds ...*exec.Cmd) *StdBytes {
	out, err := PipeBuffered(cmds...)
	goerr.Check(err)
	return out
}

// PipeBufferedAsync does the same as PipeBuffered but does so asynchronously.
func PipeBufferedAsync(cmds ...*exec.Cmd) *task.Task {
	return task.New(func(t *task.Internal) {
		MustPipeBuffered(cmds...)
	})
}
