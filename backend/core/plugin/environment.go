package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"mlib.com/mlog"
)

// EnvironmentManager handles Python virtual environment setup for plugins.
type EnvironmentManager struct {
	mu         sync.Mutex
	basePath   string
	pythonPath string
}

// NewEnvironmentManager creates an environment manager.
func NewEnvironmentManager(basePath, pythonPath string) *EnvironmentManager {
	if basePath == "" {
		basePath = filepath.Join(os.TempDir(), "gofy_plugins")
	}
	if pythonPath == "" {
		pythonPath = "python3"
	}
	os.MkdirAll(basePath, 0755)
	return &EnvironmentManager{basePath: basePath, pythonPath: pythonPath}
}

// SetupEnvironment creates a virtual environment for a plugin.
func (em *EnvironmentManager) SetupEnvironment(pluginID string) (string, error) {
	em.mu.Lock()
	defer em.mu.Unlock()

	envPath := filepath.Join(em.basePath, pluginID, "venv")
	if _, err := os.Stat(filepath.Join(envPath, "bin", "python")); err == nil {
		return envPath, nil // Already exists
	}

	mlog.Infof("creating virtual environment for plugin %s at %s", pluginID, envPath)

	// Try uv first (faster), fall back to venv
	if uvPath, err := exec.LookPath("uv"); err == nil {
		cmd := exec.Command(uvPath, "venv", envPath, "--python", em.pythonPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			mlog.Errorf("uv venv failed: %s", string(out))
		} else {
			return envPath, nil
		}
	}

	// Fall back to python -m venv
	cmd := exec.Command(em.pythonPath, "-m", "venv", envPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to create venv: %s: %w", string(out), err)
	}

	return envPath, nil
}

// InstallDependencies installs plugin Python dependencies.
func (em *EnvironmentManager) InstallDependencies(pluginID string, pluginDir string) error {
	envPath := filepath.Join(em.basePath, pluginID, "venv")
	pipPath := filepath.Join(envPath, "bin", "pip")

	// Check for requirements.txt
	reqFile := filepath.Join(pluginDir, "requirements.txt")
	if _, err := os.Stat(reqFile); err != nil {
		// Try pyproject.toml
		pyprojectFile := filepath.Join(pluginDir, "pyproject.toml")
		if _, err := os.Stat(pyprojectFile); err != nil {
			mlog.Infof("no requirements found for plugin %s", pluginID)
			return nil
		}
		// Install from pyproject.toml
		cmd := exec.Command(pipPath, "install", "-e", pluginDir)
		cmd.Env = em.buildEnv(envPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("pip install failed: %s: %w", string(out), err)
		}
		return nil
	}

	// Install from requirements.txt
	mlog.Infof("installing dependencies for plugin %s from %s", pluginID, reqFile)

	// Try uv pip first
	if uvPath, err := exec.LookPath("uv"); err == nil {
		cmd := exec.Command(uvPath, "pip", "install", "-r", reqFile, "--python", filepath.Join(envPath, "bin", "python"))
		cmd.Env = em.buildEnv(envPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			mlog.Errorf("uv pip install failed: %s, falling back to pip", string(out))
		} else {
			return nil
		}
	}

	// Fall back to pip
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, pipPath, "install", "-r", reqFile, "--no-cache-dir")
	cmd.Env = em.buildEnv(envPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pip install failed: %s: %w", string(out), err)
	}

	mlog.Infof("dependencies installed for plugin %s", pluginID)
	return nil
}

// GetPythonPath returns the Python path inside the plugin's venv.
func (em *EnvironmentManager) GetPythonPath(pluginID string) string {
	return filepath.Join(em.basePath, pluginID, "venv", "bin", "python")
}

// GetPluginDir returns the plugin installation directory.
func (em *EnvironmentManager) GetPluginDir(pluginID string) string {
	return filepath.Join(em.basePath, pluginID, "code")
}

// CleanupEnvironment removes a plugin's virtual environment.
func (em *EnvironmentManager) CleanupEnvironment(pluginID string) error {
	envDir := filepath.Join(em.basePath, pluginID)
	mlog.Infof("cleaning up environment for plugin %s", pluginID)
	return os.RemoveAll(envDir)
}

// ListInstalledPackages returns installed packages in a plugin's environment.
func (em *EnvironmentManager) ListInstalledPackages(pluginID string) ([]string, error) {
	pipPath := filepath.Join(em.basePath, pluginID, "venv", "bin", "pip")
	cmd := exec.Command(pipPath, "list", "--format=columns")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var packages []string
	for i, line := range lines {
		if i < 2 { // Skip header
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			packages = append(packages, fmt.Sprintf("%s==%s", parts[0], parts[1]))
		}
	}
	return packages, nil
}

func (em *EnvironmentManager) buildEnv(envPath string) []string {
	env := os.Environ()
	env = append(env,
		fmt.Sprintf("VIRTUAL_ENV=%s", envPath),
		fmt.Sprintf("PATH=%s:%s", filepath.Join(envPath, "bin"), os.Getenv("PATH")),
	)
	return env
}
