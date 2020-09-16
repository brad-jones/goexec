package main_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wesovilabs/koazee"
	"github.com/wesovilabs/koazee/stream"
)

func TestBuffered(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		actual := normaliseCmdOutput(out)

		replies := actual.Filter(func(v string) bool { return strings.Contains(v, "from 1.2.3.4") })
		c, err := replies.Count()
		assert.Nil(t, err)
		assert.Equal(t, 4, c)
	}
}

func normaliseCmdOutput(in []byte) stream.Stream {
	out := string(in)
	return koazee.StreamOf(strings.Split(out, "\n"))
}
