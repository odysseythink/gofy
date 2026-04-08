package plugin

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"hash/fnv"
	"time"

	pluginentities "mlib.com/gofy/server/entities/plugin"
)

// PluginRuntimeType represents the type of plugin runtime.
type PluginRuntimeType string

const (
	PLUGIN_RUNTIME_TYPE_LOCAL      PluginRuntimeType = "local"
	PLUGIN_RUNTIME_TYPE_REMOTE     PluginRuntimeType = "remote"
	PLUGIN_RUNTIME_TYPE_SERVERLESS PluginRuntimeType = "serverless"
)

// PluginRuntimeStatus constants
const (
	PLUGIN_RUNTIME_STATUS_ACTIVE     = "active"
	PLUGIN_RUNTIME_STATUS_LAUNCHING  = "launching"
	PLUGIN_RUNTIME_STATUS_STOPPED    = "stopped"
	PLUGIN_RUNTIME_STATUS_RESTARTING = "restarting"
	PLUGIN_RUNTIME_STATUS_PENDING    = "pending"
)

// PluginRuntimeState holds the runtime state of a plugin instance.
type PluginRuntimeState struct {
	Restarts    int        `json:"restarts"`
	Status      string     `json:"status"`
	WorkingPath string     `json:"working_path"`
	ActiveAt    *time.Time `json:"active_at"`
	StoppedAt   *time.Time `json:"stopped_at"`
	Verified    bool       `json:"verified"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	Logs        []string   `json:"logs"`
}

// Hash computes a FNV-1a hash of the runtime state.
func (s *PluginRuntimeState) Hash() (uint64, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		return 0, err
	}
	j := fnv.New64a()
	_, err = j.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}
	return j.Sum64(), nil
}

// PluginRuntime is the base runtime struct that holds config and state.
type PluginRuntime struct {
	State  PluginRuntimeState               `json:"state"`
	Config pluginentities.PluginDeclaration `json:"config"`
}

func (r *PluginRuntime) Stopped() bool {
	return r.State.Status == PLUGIN_RUNTIME_STATUS_STOPPED
}

func (r *PluginRuntime) Stop() {
	r.State.Status = PLUGIN_RUNTIME_STATUS_STOPPED
}

func (r *PluginRuntime) Configuration() *pluginentities.PluginDeclaration {
	return &r.Config
}

func (r *PluginRuntime) RuntimeState() PluginRuntimeState {
	return r.State
}

func (r *PluginRuntime) InitState() {
	r.State = PluginRuntimeState{
		Restarts:    0,
		Status:      PLUGIN_RUNTIME_STATUS_PENDING,
		ActiveAt:    nil,
		StoppedAt:   nil,
		Verified:    false,
		ScheduledAt: nil,
		Logs:        []string{},
	}
}

func (r *PluginRuntime) SetActive() {
	r.State.Status = PLUGIN_RUNTIME_STATUS_ACTIVE
}

func (r *PluginRuntime) SetLaunching() {
	r.State.Status = PLUGIN_RUNTIME_STATUS_LAUNCHING
}

func (r *PluginRuntime) SetRestarting() {
	r.State.Status = PLUGIN_RUNTIME_STATUS_RESTARTING
}

func (r *PluginRuntime) SetPending() {
	r.State.Status = PLUGIN_RUNTIME_STATUS_PENDING
}

func (r *PluginRuntime) SetActiveAt(t time.Time) {
	r.State.ActiveAt = &t
}

func (r *PluginRuntime) SetScheduledAt(t time.Time) {
	r.State.ScheduledAt = &t
}

// HashedIdentity computes a SHA-256 hash of a string identity.
func HashedIdentity(identity string) string {
	hash := sha256.New()
	hash.Write([]byte(identity))
	return hex.EncodeToString(hash.Sum(nil))
}

// PluginBasicInfoInterface provides basic read-only information about a plugin.
type PluginBasicInfoInterface interface {
	Type() PluginRuntimeType
	Configuration() *pluginentities.PluginDeclaration
	Identity() (string, error)
	HashedIdentity() (string, error)
	Checksum() (string, error)
}

// SessionMessage represents a message from a plugin session.
type SessionMessage struct {
	SessionID string `json:"session_id"`
	Data      []byte `json:"data"`
}

// PluginRuntimeSessionIOInterface handles session-level I/O with a plugin.
type PluginRuntimeSessionIOInterface interface {
	PluginBasicInfoInterface
	// Listen listens for messages from the plugin for a given session.
	Listen(sessionID string) (<-chan SessionMessage, error)
	// Write writes a message to the plugin for a given session.
	Write(sessionID string, action string, data []byte) error
}

// PluginClusterLifetime provides cluster-level runtime state and lifecycle control.
type PluginClusterLifetime interface {
	RuntimeState() PluginRuntimeState
	Stopped() bool
	Stop()
}

// PluginLifetime combines basic info, session I/O, and cluster lifetime.
type PluginLifetime interface {
	PluginBasicInfoInterface
	PluginRuntimeSessionIOInterface
	PluginClusterLifetime
}

// PluginFullDuplexLifetime extends PluginLifetime with environment and state management.
type PluginFullDuplexLifetime interface {
	PluginLifetime

	InitEnvironment() error
	Cleanup()
	SetActive()
	SetLaunching()
	SetRestarting()
	SetPending()
	SetActiveAt(t time.Time)
	SetScheduledAt(t time.Time)
}

// PluginServerlessLifetime extends PluginLifetime for serverless runtimes.
type PluginServerlessLifetime interface {
	PluginLifetime

	InitEnvironment() error
	UploadPlugin() error
}
