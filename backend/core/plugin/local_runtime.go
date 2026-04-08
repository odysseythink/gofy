package plugin

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"mlib.com/mlog"
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
}

// LocalRuntimeConfig configures the local runtime.
type LocalRuntimeConfig struct {
	PythonPath       string            // Path to Python interpreter
	PluginDir        string            // Plugin directory
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
