package goexec

import (
	"bytes"
	"fmt"
	"io"

	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
	"github.com/brad-jones/goprefix/v2/pkg/prefixer"
)

func prefixed(prefix string,
	stdOut, stdErr io.Writer,
	stdOutPipeR, stdErrPipeR io.Reader,
	stdOutPipeW, stdErrPipeW io.WriteCloser,
	fn func() error) error {

	errorCh := make(chan error)

	// Make the prefix colorful
	p := prefixer.New(colorchooser.Sprint(prefix) + " | ")

	// Run the function, pipeing all StdOut and StdErr to our scanners
	go func() {
		defer stdOutPipeW.Close()
		defer stdErrPipeW.Close()
		errorCh <- fn()
	}()

	// Prefix all StdOut
	go func() {
		if err := p.ReadFrom(stdOutPipeR).WriteTo(stdOut); err != nil {
			errorCh <- fmt.Errorf("prefixing standard out: %s", err)
		}
	}()

	// Prefix all StdErr
	go func() {
		if err := p.ReadFrom(stdErrPipeR).WriteTo(stdErr); err != nil {
			errorCh <- fmt.Errorf("prefixing standard out: %s", err)
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
