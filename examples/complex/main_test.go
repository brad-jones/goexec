package main_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wesovilabs/koazee"
	"github.com/wesovilabs/koazee/stream"
)

func TestComplex(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		actual := normaliseCmdOutput(out)

		c1, err := actual.Filter(func(v string) bool { return strings.HasSuffix(v, "README.md") }).Count()
		assert.Nil(t, err)
		assert.Equal(t, 1, c1)

		c2, err := actual.Filter(func(v string) bool { return strings.HasSuffix(v, "main.go") }).Count()
		assert.Nil(t, err)
		assert.Equal(t, 1, c2)

		c3, err := actual.Filter(func(v string) bool { return strings.HasSuffix(v, "main_test.go") }).Count()
		assert.Nil(t, err)
		assert.Equal(t, 1, c3)

		envOutput := actual.Filter(func(v string) bool { return strings.Contains(v, "FOO=BAR") })
		c, err := envOutput.Count()
		assert.Nil(t, err)
		assert.Equal(t, 1, c)
	}
}

func normaliseCmdOutput(in []byte) stream.Stream {
	out := string(in)
	return koazee.StreamOf(strings.Split(out, "\n"))
}
