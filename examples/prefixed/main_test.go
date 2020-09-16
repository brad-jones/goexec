package main_test

import (
	"os/exec"
	"regexp"
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
		assert.Regexp(t, regexp.QuoteMeta("\x1b[")+"(.*)mip1"+regexp.QuoteMeta("\x1b[0m |"), replies1.First().String())

		replies2 := actual.Filter(func(v string) bool { return strings.Contains(v, "from 127.0.0.2") })
		c2, err := replies2.Count()
		assert.Nil(t, err)
		assert.Equal(t, 4, c2)
		assert.Regexp(t, regexp.QuoteMeta("\x1b[")+"(.*)mip2"+regexp.QuoteMeta("\x1b[0m |"), replies2.First().String())
	}
}

func normaliseCmdOutput(in []byte) stream.Stream {
	out := string(in)
	return koazee.StreamOf(strings.Split(out, "\n"))
}
