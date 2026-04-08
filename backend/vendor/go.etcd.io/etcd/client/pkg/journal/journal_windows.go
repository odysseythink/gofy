// Package journal provides write bindings to the local systemd journal.
// It is implemented in pure Go and connects to the journal directly over its
// unix socket.
//
// To read from the journal, see the "sdjournal" package, which wraps the
// sd-journal a C API.
//
// http://www.freedesktop.org/software/systemd/man/systemd-journald.service.html
package journal

import (
	"errors"
)

func Enabled() bool {
	return false
}

func Send(message string, priority Priority, vars map[string]string) error {
	return errors.New("could not initialize socket to journald")
}

func StderrIsJournalStream() (bool, error) {
	return false, nil
}

func StdoutIsJournalStream() (bool, error) {
	return false, nil
}
