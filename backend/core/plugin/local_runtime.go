package plugin

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	pluginentities "github.com/odysseythink/gofy/backend/entities/plugin"
	"github.com/odysseythink/mlog"
	"gopkg.in/yaml.v3"
)

// LocalRuntime manages a Python plugin subprocess.
type LocalRuntime struct {
	mu               sync.Mutex
	pluginID         string
	uniqueIdentifier string
	pythonPath       string
	pluginDir        string
	cmd              *exec.Cmd
	stdin            io.WriteCloser
	stdout           io.ReadCloser
	stderr           io.ReadCloser
	state            *PluginRuntimeState
	stopCh           chan struct{}
	messageCh        chan *SessionMessage
	ctx              context.Context
	cancel           context.CancelFunc

	// Per-session fan-out. readOutput routes incoming messages to the matching
	// sessions[sid] channel when registered, and also to messageCh for catch-all
	// consumers of Messages(). Listen() registers; closeSession() tears down.
	sessionsMu sync.Mutex
	sessions   map[string]chan SessionMessage

	// Parsed plugin manifest; nil until loaded.
	declarationOnce sync.Once
	declaration     *pluginentities.PluginDeclaration
}

// LocalRuntimeConfig configures the local runtime.
type LocalRuntimeConfig struct {
	PythonPath       string // Path to Python interpreter
	PluginDir        string // Plugin directory
	PluginID         string
	UniqueIdentifier string
	WorkingDir       string            // Working directory for the process
	EnvVars          map[string]string // Additional environment variables
}

// NewLocalRuntime creates a new local Python plugin runtime.
func NewLocalRuntime(cfg LocalRuntimeConfig) *LocalRuntime {
	ctx, cancel := context.WithCancel(context.Background())

	pythonPath := cfg.PythonPath
	if pythonPath == "" {
		pythonPath = "python3"
	}

	return &LocalRuntime{
		pluginID:         cfg.PluginID,
		uniqueIdentifier: cfg.UniqueIdentifier,
		pythonPath:       pythonPath,
		pluginDir:        cfg.PluginDir,
		state:            &PluginRuntimeState{Status: PLUGIN_RUNTIME_STATUS_PENDING},
		stopCh:           make(chan struct{}),
		messageCh:        make(chan *SessionMessage, 100),
		ctx:              ctx,
		cancel:           cancel,
	}
}

// Start launches the Python plugin process.
func (lr *LocalRuntime) Start() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	lr.state.Status = PLUGIN_RUNTIME_STATUS_LAUNCHING

	// Find the plugin entry point
	entryPoint := filepath.Join(lr.pluginDir, "main.py")
	if _, err := os.Stat(entryPoint); os.IsNotExist(err) {
		// Try alternative entry points
		entryPoint = filepath.Join(lr.pluginDir, "__init__.py")
		if _, err := os.Stat(entryPoint); os.IsNotExist(err) {
			return fmt.Errorf("no plugin entry point found in %s", lr.pluginDir)
		}
	}

	lr.cmd = exec.CommandContext(lr.ctx, lr.pythonPath, entryPoint)
	lr.cmd.Dir = lr.pluginDir
	lr.cmd.Env = append(os.Environ(),
		fmt.Sprintf("PLUGIN_ID=%s", lr.pluginID),
		fmt.Sprintf("PLUGIN_UNIQUE_IDENTIFIER=%s", lr.uniqueIdentifier),
	)

	var err error
	lr.stdin, err = lr.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	lr.stdout, err = lr.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	lr.stderr, err = lr.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := lr.cmd.Start(); err != nil {
		lr.state.Status = PLUGIN_RUNTIME_STATUS_STOPPED
		return fmt.Errorf("failed to start plugin process: %w", err)
	}

	now := time.Now()
	lr.state.Status = PLUGIN_RUNTIME_STATUS_ACTIVE
	lr.state.ActiveAt = &now

	// Read stdout messages in background
	go lr.readOutput()
	// Read stderr for logging
	go lr.readErrors()
	// Wait for process exit
	go lr.waitForExit()

	mlog.Infof("local runtime started for plugin %s (pid: %d)", lr.uniqueIdentifier, lr.cmd.Process.Pid)
	return nil
}

// Stop terminates the plugin process.
func (lr *LocalRuntime) Stop() {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	lr.cancel()

	// Close stopCh only once
	select {
	case <-lr.stopCh:
	default:
		close(lr.stopCh)
	}

	if lr.cmd != nil && lr.cmd.Process != nil {
		_ = lr.cmd.Process.Signal(os.Interrupt)
		// Give it 5 seconds to gracefully shutdown
		done := make(chan error, 1)
		go func() { done <- lr.cmd.Wait() }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = lr.cmd.Process.Kill()
		}
	}

	now := time.Now()
	lr.state.Status = PLUGIN_RUNTIME_STATUS_STOPPED
	lr.state.StoppedAt = &now
	mlog.Infof("local runtime stopped for plugin %s", lr.uniqueIdentifier)
}

// Stopped returns whether the runtime has stopped.
func (lr *LocalRuntime) Stopped() bool {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.state.Status == PLUGIN_RUNTIME_STATUS_STOPPED
}

// RuntimeState returns the current runtime state.
func (lr *LocalRuntime) RuntimeState() PluginRuntimeState {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return *lr.state
}

// SendMessage sends a message to the plugin process via stdin.
func (lr *LocalRuntime) SendMessage(msg *SessionMessage) error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.stdin == nil {
		return fmt.Errorf("runtime not started")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Write message as newline-delimited JSON
	if _, err := lr.stdin.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write to stdin: %w", err)
	}

	return nil
}

// Messages returns a channel of messages from the plugin.
func (lr *LocalRuntime) Messages() <-chan *SessionMessage {
	return lr.messageCh
}

// readOutput reads newline-delimited JSON messages from stdout.
func (lr *LocalRuntime) readOutput() {
	scanner := bufio.NewScanner(lr.stdout)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer
	for scanner.Scan() {
		line := scanner.Bytes()
		var msg SessionMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			mlog.Errorf("plugin %s: invalid output: %s", lr.uniqueIdentifier, string(line))
			continue
		}
		// Route to per-session listener when registered; otherwise fall back
		// to the catch-all messageCh consumed via Messages().
		if lr.dispatchToSession(msg) {
			continue
		}
		select {
		case lr.messageCh <- &msg:
		case <-lr.stopCh:
			return
		default:
			mlog.Errorf("plugin %s: message channel full, dropping message", lr.uniqueIdentifier)
		}
	}
}

// readErrors logs stderr output from the plugin process.
func (lr *LocalRuntime) readErrors() {
	scanner := bufio.NewScanner(lr.stderr)
	for scanner.Scan() {
		mlog.Errorf("plugin %s stderr: %s", lr.uniqueIdentifier, scanner.Text())
	}
}

// waitForExit waits for the process to exit and updates state accordingly.
func (lr *LocalRuntime) waitForExit() {
	if lr.cmd == nil {
		return
	}
	err := lr.cmd.Wait()
	lr.mu.Lock()
	defer lr.mu.Unlock()
	if err != nil {
		mlog.Errorf("plugin %s process exited with error: %v", lr.uniqueIdentifier, err)
	}
	now := time.Now()
	lr.state.Status = PLUGIN_RUNTIME_STATUS_STOPPED
	lr.state.StoppedAt = &now
	lr.state.Restarts++
}

// --- PluginLifetime interface stubs ---
// These satisfy PluginLifetime for type-assertion callers; full semantics
// are implemented where the local runtime is actually wired end-to-end.

func (lr *LocalRuntime) Type() PluginRuntimeType { return PLUGIN_RUNTIME_TYPE_LOCAL }

// Configuration parses pluginDir/manifest.{yaml,yml,json} on first call and
// caches the result. Returns nil if pluginDir is empty or parsing fails.
func (lr *LocalRuntime) Configuration() *pluginentities.PluginDeclaration {
	lr.declarationOnce.Do(func() {
		if lr.pluginDir == "" {
			return
		}
		candidates := []string{"manifest.yaml", "manifest.yml", "manifest.json"}
		for _, name := range candidates {
			path := filepath.Join(lr.pluginDir, name)
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			var d pluginentities.PluginDeclaration
			if strings.HasSuffix(name, ".json") {
				if err := json.Unmarshal(data, &d); err != nil {
					mlog.Warningf("plugin %s: parse %s failed: %v", lr.uniqueIdentifier, name, err)
					continue
				}
			} else {
				if err := yaml.Unmarshal(data, &d); err != nil {
					mlog.Warningf("plugin %s: parse %s failed: %v", lr.uniqueIdentifier, name, err)
					continue
				}
			}
			lr.declaration = &d
			return
		}
		mlog.Warningf("plugin %s: no manifest found in %s", lr.uniqueIdentifier, lr.pluginDir)
	})
	return lr.declaration
}

func (lr *LocalRuntime) Identity() (string, error) {
	return lr.uniqueIdentifier, nil
}

func (lr *LocalRuntime) HashedIdentity() (string, error) {
	return HashedIdentity(lr.uniqueIdentifier), nil
}

// Checksum is a stable hash of the plugin's on-disk bytes. It walks pluginDir,
// hashes file contents in lexicographic path order, and returns the composite
// SHA-256 hex. Falls back to hashing the unique identifier if the directory
// cannot be read (e.g. runtime not yet started or path missing).
func (lr *LocalRuntime) Checksum() (string, error) {
	if lr.pluginDir == "" {
		return HashedIdentity(lr.uniqueIdentifier), nil
	}
	h := sha256.New()
	err := filepath.Walk(lr.pluginDir, func(path string, info os.FileInfo, werr error) error {
		if werr != nil {
			return werr
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(lr.pluginDir, path)
		h.Write([]byte(rel))
		h.Write([]byte{0})
		f, ferr := os.Open(path)
		if ferr != nil {
			return ferr
		}
		_, cerr := io.Copy(h, f)
		f.Close()
		return cerr
	})
	if err != nil {
		return HashedIdentity(lr.uniqueIdentifier), nil
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Listen returns a channel that delivers only messages tagged with the given
// session ID. Subsequent Listen calls for the same sessionID return the same
// channel. CloseSession tears the registration down.
func (lr *LocalRuntime) Listen(sessionID string) (<-chan SessionMessage, error) {
	lr.sessionsMu.Lock()
	defer lr.sessionsMu.Unlock()
	if lr.sessions == nil {
		lr.sessions = make(map[string]chan SessionMessage)
	}
	if existing, ok := lr.sessions[sessionID]; ok {
		return existing, nil
	}
	ch := make(chan SessionMessage, 64)
	lr.sessions[sessionID] = ch
	return ch, nil
}

// CloseSession drops the per-session channel for sessionID and closes it.
func (lr *LocalRuntime) CloseSession(sessionID string) {
	lr.sessionsMu.Lock()
	defer lr.sessionsMu.Unlock()
	if ch, ok := lr.sessions[sessionID]; ok {
		close(ch)
		delete(lr.sessions, sessionID)
	}
}

// Write wraps (action, data) as an envelope and forwards to the plugin with the
// given sessionID. The envelope format is JSON: {"action": "...", "data": ...}
// so plugins can route based on action without parsing raw Data twice.
func (lr *LocalRuntime) Write(sessionID string, action string, data []byte) error {
	envelope, err := json.Marshal(map[string]any{
		"action": action,
		"data":   json.RawMessage(data),
	})
	if err != nil {
		return fmt.Errorf("envelope marshal: %w", err)
	}
	return lr.SendMessage(&SessionMessage{SessionID: sessionID, Data: envelope})
}

// dispatchToSession routes an inbound SessionMessage to a per-session channel
// if one has been registered via Listen. Returns true when delivered to a
// session channel. Non-blocking: drops if the session channel is full.
func (lr *LocalRuntime) dispatchToSession(msg SessionMessage) bool {
	if msg.SessionID == "" {
		return false
	}
	lr.sessionsMu.Lock()
	ch, ok := lr.sessions[msg.SessionID]
	lr.sessionsMu.Unlock()
	if !ok {
		return false
	}
	select {
	case ch <- msg:
		return true
	default:
		mlog.Warningf("plugin %s: session %s channel full, dropping", lr.uniqueIdentifier, msg.SessionID)
		return true
	}
}
