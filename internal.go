package goexec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// We need a basic random source, so lets set it up here
var randGen *rand.Rand

func init() {
	randGen = rand.New(rand.NewSource(time.Now().Unix()))
}

var prefixToColorMap = &sync.Map{}
var chosenColors = &sync.Map{}
var chosenCount = 0
var availableColors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgHiRed,
	color.FgHiGreen,
	color.FgHiYellow,
	color.FgHiBlue,
	color.FgHiMagenta,
	color.FgHiCyan,
	color.BgRed,
	color.BgGreen,
	color.BgYellow,
	color.BgBlue,
	color.BgMagenta,
	color.BgCyan,
	color.BgHiRed,
	color.BgHiGreen,
	color.BgHiYellow,
	color.BgHiBlue,
	color.BgHiMagenta,
	color.BgHiCyan,
}

func colorChooser(prefix string) color.Attribute {

	// Return the cached value if it exists
	c, exists := prefixToColorMap.Load(prefix)
	if exists {
		return c.(color.Attribute)
	}

	// Randomly choose a new color
	c = availableColors[randGen.Intn(len(availableColors))]

	// Check if we reached the maximum number of available colors.
	// If so we will just have to reuse a color.
	if chosenCount == len(availableColors) {
		prefixToColorMap.Store(prefix, c)
		return c.(color.Attribute)
	}

	// Check if this color has already been chosen,
	// running ourselves again to select hopefully a different color.
	if _, chosen := chosenColors.Load(c); chosen == true {
		return colorChooser(prefix)
	}

	// Cache the result for next time
	chosenCount = chosenCount + 1
	chosenColors.Store(c, true)
	prefixToColorMap.Store(prefix, c)

	return c.(color.Attribute)
}

func prefixed(prefix string,
	stdOut, stdErr io.Writer,
	stdOutPipeR, stdErrPipeR io.Reader,
	stdOutPipeW, stdErrPipeW io.WriteCloser,
	fn func() error) error {

	errorCh := make(chan error)
	stdOutScanner := bufio.NewScanner(stdOutPipeR)
	stdErrScanner := bufio.NewScanner(stdErrPipeR)

	// Make the prefix colorful
	prefix = color.New(colorChooser(prefix)).Sprint(prefix + " | ")

	// Run the function, pipeing all StdOut and StdErr to our scanners
	go func() {
		defer stdOutPipeW.Close()
		defer stdErrPipeW.Close()
		errorCh <- fn()
	}()

	// Prefix all StdOut
	go func() {
		for stdOutScanner.Scan() {
			fmt.Fprintln(stdOut, prefix+strings.TrimSpace(stdOutScanner.Text())+"\r")
		}
		if err := stdOutScanner.Err(); err != nil {
			errorCh <- fmt.Errorf("prefxing standard out: %s", err)
		}
	}()

	// Prefix all StdErr
	go func() {
		for stdErrScanner.Scan() {
			fmt.Fprintln(stdErr, prefix+color.RedString(strings.TrimSpace(stdErrScanner.Text()))+"\r")
		}
		if err := stdErrScanner.Err(); err != nil {
			errorCh <- fmt.Errorf("prefxing standard err: %s", err)
		}
	}()

	// Catch any errors
	if err := <-errorCh; err != nil {
		return err
	}

	return nil
}

func buffered(fstdOut, stdErr io.Writer,
	stdOutPipeR, stdErrPipeR io.Reader,
	stdOutPipeW, stdErrPipeW io.WriteCloser,
	fn func() error) (out *StdBytes, err error) {

	errorCh := make(chan error)

	// Run the function, pipeing all StdOut and StdErr to our buffers
	go func() {
		defer stdOutPipeW.Close()
		defer stdErrPipeW.Close()
		errorCh <- fn()
	}()

	// Read all StdOut into our buffer
	stdOutC := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, stdOutPipeR); err != nil {
			errorCh <- err
		} else {
			stdOutC <- buf.Bytes()
		}
	}()

	// Read all StdErr into our buffer
	stdErrC := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, stdErrPipeR); err != nil {
			errorCh <- err
		} else {
			stdErrC <- buf.Bytes()
		}
	}()

	// Catch any errors
	if err := <-errorCh; err != nil {
		return &StdBytes{<-stdOutC, <-stdErrC}, err
	}

	return &StdBytes{<-stdOutC, <-stdErrC}, nil
}
