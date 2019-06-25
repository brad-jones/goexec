package goexec

// StdBytes represents the buffered output from a process.
type StdBytes struct {
	StdOut []byte
	StdErr []byte
}

// StdStrings represents the buffered output from a process after having been casted to strings.
type StdStrings struct {
	StdOut string
	StdErr string
}
