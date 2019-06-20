package goexec

import (
	"os"
	"os/exec"

	"github.com/brad-jones/goerr"
	"github.com/go-errors/errors"
)

// Cmd provides a fluent, decorator based API for os/exec.
//
// 	Cmd("ping", Args("-c", "4", "1.1.1.1"))
func Cmd(cmd string, decorators ...func(*exec.Cmd) error) (c *exec.Cmd, err error) {
	defer goerr.Handle(func(e error) {
		c = nil
		err = errors.Wrap(e, 0)
	})

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

// Run is a convenience function for simple cases.
//
// Instead of:
// 	Cmd("ping", Args("-c", "4", "8.8.8.8")).Run()
//
// You might write:
// 	Run("ping", "-c", "4", "8.8.8.8")
func Run(cmd string, args ...string) (err error) {
	defer goerr.Handle(func(e error) {
		err = errors.Wrap(e, 0)
	})

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
func RunAsync(cmd string, args ...string) (done <-chan struct{}, err <-chan error) {
	doneCh := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go func() {
		err := Run(cmd, args...)
		if err == nil {
			close(doneCh)
		} else {
			errCh <- errors.Wrap(err, 0)
		}
	}()
	return doneCh, errCh
}

// RunBuffered is a convenience function for simple cases.
//
// Instead of:
// 	RunBufferedCmd(Cmd("ping", Args("-c", "4", "8.8.8.8")))
//
// You might write:
// 	stdOut, stdErr, err := RunBuffered("ping", "-c", "4", "8.8.8.8")
func RunBuffered(cmd string, args ...string) (stdOutBuf, stdErrBuf string, err error) {
	defer goerr.Handle(func(e error) {
		stdOutBuf = ""
		stdErrBuf = ""
		err = errors.Wrap(e, 0)
	})

	c, err := Cmd(cmd, Args(args...))
	goerr.Check(err)

	o, e, err := RunBufferedCmd(c)
	if err != nil {
		err = errors.Wrap(e, 0)
	}

	return string(o), string(e), err
}

// MustRunBuffered does the same as RunBuffered but panics if an error is encountered.
func MustRunBuffered(cmd string, args ...string) (stdOutBuf, stdErrBuf string) {
	stdOutBuf, stdOutErr, err := RunBuffered(cmd, args...)
	goerr.Check(err)
	return stdOutBuf, stdOutErr
}

// RunBufferedAsync does the same as RunBuffered but does so asynchronously.
func RunBufferedAsync(cmd string, args ...string) (stdOutBuf, stdErrBuf <-chan string, err <-chan error) {
	stdOutCh := make(chan string, 1)
	stdErrCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		stdOut, stdErr, err := RunBuffered(cmd, args...)
		stdOutCh <- stdOut
		stdErrCh <- stdErr
		if err != nil {
			errCh <- errors.Wrap(err, 0)
		}
	}()
	return stdOutCh, stdErrCh, errCh
}

// RunPrefixed is a convenience function for simple cases.
//
// Instead of:
// 	RunPrefixedCmd("foo", Cmd("ping", Args("-c", "4", "8.8.8.8")))
//
// You might write:
// 	RunPrefixed("foo", "ping", "-c", "4", "8.8.8.8")
func RunPrefixed(prefix, cmd string, args ...string) (err error) {
	defer goerr.Handle(func(e error) {
		err = errors.Wrap(e, 0)
	})

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
func RunPrefixedAsync(prefix, cmd string, args ...string) (done <-chan struct{}, err <-chan error) {
	doneCh := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go func() {
		err := RunPrefixed(prefix, cmd, args...)
		if err == nil {
			close(doneCh)
		} else {
			errCh <- errors.Wrap(err, 0)
		}
	}()
	return doneCh, errCh
}

// RunPrefixedCmd will prefix all StdOut and StdErr with given prefix.
// This is useful when running many commands concurrently,
// output will look similar to docker-compose.
func RunPrefixedCmd(prefix string, cmd *exec.Cmd) (err error) {
	defer goerr.Handle(func(e error) {
		err = errors.Wrap(e, 0)
	})

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
func RunPrefixedCmdAsync(prefix string, cmd *exec.Cmd) (done <-chan struct{}, err <-chan error) {
	doneCh := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go func() {
		err := RunPrefixedCmd(prefix, cmd)
		if err == nil {
			close(doneCh)
		} else {
			errCh <- errors.Wrap(err, 0)
		}
	}()
	return doneCh, errCh
}

// RunBufferedCmd will buffer all StdOut and StdErr, returning the buffers.
// This is useful when you wish to parse the results of a command.
func RunBufferedCmd(cmd *exec.Cmd) (stdOut, stdErr []byte, err error) {
	defer goerr.Handle(func(e error) {
		stdOut = nil
		stdErr = nil
		err = errors.Wrap(e, 0)
	})

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
func MustRunBufferedCmd(cmd *exec.Cmd) (stdOut, stdErr []byte) {
	stdOut, stdErr, err := RunBufferedCmd(cmd)
	goerr.Check(err)
	return stdOut, stdErr
}

// RunBufferedCmdAsync does the same as RunBufferedCmd but does so asynchronously.
func RunBufferedCmdAsync(cmd *exec.Cmd) (stdOut, stdErr <-chan []byte, err <-chan error) {
	stdOutCh := make(chan []byte, 1)
	stdErrCh := make(chan []byte, 1)
	errCh := make(chan error, 1)
	go func() {
		stdOut, stdErr, err := RunBufferedCmd(cmd)
		stdOutCh <- stdOut
		stdErrCh <- stdErr
		if err != nil {
			errCh <- errors.Wrap(err, 0)
		}
	}()
	return stdOutCh, stdErrCh, errCh
}

// Pipe will send the output of the first command
// to the input of the second and so on.
func Pipe(cmds ...*exec.Cmd) (err error) {
	defer goerr.Handle(func(e error) {
		err = errors.Wrap(e, 0)
	})

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
func PipeAsync(cmds ...*exec.Cmd) (done <-chan struct{}, err <-chan error) {
	doneCh := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go func() {
		err := Pipe(cmds...)
		if err == nil {
			close(doneCh)
		} else {
			errCh <- errors.Wrap(err, 0)
		}
	}()
	return doneCh, errCh
}

// PipePrefixed will prefix all StdOut and StdErr with given prefix.
// This is useful when running many pipes concurrently,
// output will look similar to docker-compose.
func PipePrefixed(prefix string, cmds ...*exec.Cmd) (err error) {
	defer goerr.Handle(func(e error) {
		err = errors.Wrap(e, 0)
	})

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
func PipePrefixedAsync(prefix string, cmds ...*exec.Cmd) (done <-chan struct{}, err <-chan error) {
	doneCh := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go func() {
		err := PipePrefixed(prefix, cmds...)
		if err == nil {
			close(doneCh)
		} else {
			errCh <- errors.Wrap(err, 0)
		}
	}()
	return doneCh, errCh
}

// PipeBuffered will buffer all StdOut and StdErr, returning the buffers.
// This is useful when you wish to parse the results of a pipe.
func PipeBuffered(cmds ...*exec.Cmd) (stdOut, stdErr []byte, err error) {
	defer goerr.Handle(func(e error) {
		stdOut = nil
		stdErr = nil
		err = errors.Wrap(e, 0)
	})

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
func MustPipeBuffered(cmds ...*exec.Cmd) (stdOut, stdErr []byte) {
	stdOut, stdErr, err := PipeBuffered(cmds...)
	goerr.Check(err)
	return stdOut, stdErr
}

// PipeBufferedAsync does the same as PipeBuffered but does so asynchronously.
func PipeBufferedAsync(cmds ...*exec.Cmd) (stdOut, stdErr <-chan []byte, err <-chan error) {
	stdOutCh := make(chan []byte, 1)
	stdErrCh := make(chan []byte, 1)
	errCh := make(chan error, 1)
	go func() {
		stdOut, stdErr, err := PipeBuffered(cmds...)
		stdOutCh <- stdOut
		stdErrCh <- stdErr
		if err != nil {
			errCh <- errors.Wrap(err, 0)
		}
	}()
	return stdOutCh, stdErrCh, errCh
}
