package goexec

import (
	"io"
	"os"
	"os/exec"
)

// Args allows you to define the arguments sent to the command to be run
func Args(args ...string) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Args = append([]string{c.Path}, args...)
		return nil
	}
}

// Cwd allows you to configure the working directory of the command to be run
func Cwd(dir string) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Dir = dir
		return nil
	}
}

// Env allows you to set the exact environment in which the command will be run
func Env(env map[string]string) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		e := []string{}
		for k, v := range env {
			e = append(e, k+"="+v)
		}
		c.Env = e
		return nil
	}
}

// EnvCombined will add the variables you provide to the existing environment
func EnvCombined(env map[string]string) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		e := os.Environ()
		for k, v := range env {
			e = append(e, k+"="+v)
		}
		c.Env = e
		return nil
	}
}

// SetIn allows you to set a custom StdIn stream
func SetIn(in io.Reader) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Stdin = in
		return nil
	}
}

// SetOut allows you to set a custom StdOut stream
func SetOut(out io.Writer) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Stdout = out
		return nil
	}
}

// SetErr allows you to set a custom StdErr stream
func SetErr(err io.Writer) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Stderr = err
		return nil
	}
}
