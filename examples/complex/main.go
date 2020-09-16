package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/brad-jones/goexec/v2"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cmd, err := goexec.Cmd("docker",
		goexec.Args(
			"run", "--rm",
			"-v", fmt.Sprintf("%s:%s", cwd, nixPath(cwd)),
			"-w", nixPath(cwd),
			"-e", "FOO",
			"alpine:latest",
			"sh", "-c", "ls -hal && env",
		),
		goexec.EnvCombined(map[string]string{"FOO": "BAR"}),
	)
	if err != nil {
		panic(err)
	}

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func nixPath(in string) string {
	out := strings.ReplaceAll(in, "\\", "/")
	out = strings.Replace(out, "C:", "/c", 1)
	return out
}
