package journal

import (
	"fmt"
)

// Priority of a journal message
type Priority int

const (
	PriEmerg Priority = iota
	PriAlert
	PriCrit
	PriErr
	PriWarning
	PriNotice
	PriInfo
	PriDebug
)

// Print prints a message to the local systemd journal using Send().
func Print(priority Priority, format string, a ...interface{}) error {
	return Send(fmt.Sprintf(format, a...), priority, nil)
}
