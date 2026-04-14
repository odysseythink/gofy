package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odysseythink/mlog"
)

// MailPayload is the payload for email tasks.
type MailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Type    string `json:"type"` // registration, reset_password, invite, etc.
}

// HandleMailRegistration sends a registration verification email.
func HandleMailRegistration(ctx context.Context, task *Task) error {
	var payload MailPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	mlog.Infof("sending registration email to %s", payload.To)
	// TODO: Integrate with email service (config/email.go)
	return nil
}

// HandleMailResetPassword sends a password reset email.
func HandleMailResetPassword(ctx context.Context, task *Task) error {
	var payload MailPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	mlog.Infof("sending password reset email to %s", payload.To)
	// TODO: Integrate with email service
	return nil
}

// HandleMailInviteMember sends a member invitation email.
func HandleMailInviteMember(ctx context.Context, task *Task) error {
	var payload MailPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	mlog.Infof("sending invitation email to %s", payload.To)
	// TODO: Integrate with email service
	return nil
}
