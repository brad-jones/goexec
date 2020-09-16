package main_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wesovilabs/koazee"
	"github.com/wesovilabs/koazee/stream"
)

func TestPrefixed(t *testing.T) {
	out, err := exec.Command("go", "run", ".").CombinedOutput()
	if assert.NoError(t, err) {
		actual := normaliseCmdOutput(out)

		replies1 := actual.Filter(func(v string) bool { return strings.Contains(v, "from 127.0.0.1") })
		c1, err := replies1.Count()
		assert.Nil(t, err)
		assert.Equal(t, 4, c1)
		assert.True(t, strings.HasPrefix(replies1.First().String(), "ip1 |"))

		replies2 := actual.Filter(func(v string) bool { return strings.Contains(v, "from 127.0.0.2") })
		c2, err := replies2.Count()
		assert.Nil(t, err)
		assert.Equal(t, 4, c2)
		assert.True(t, strings.HasPrefix(replies2.First().String(), "ip2 |"))
	}
}

func normaliseCmdOutput(in []byte) stream.Stream {
	out := string(in)
	return koazee.StreamOf(strings.Split(out, "\n"))
}
