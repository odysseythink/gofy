package python

import (
	_ "embed"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/viper"
	"mlib.com/gofy/server/main/sandbox/global"
	"mlib.com/gofy/server/main/sandbox/runner"
	"mlib.com/mlog"
)

//go:embed env.sh
var env_script string

func PreparePythonDependenciesEnv() error {
	python_lib_paths := viper.GetStringSlice("python_lib_path")
	if len(python_lib_paths) == 0 {
		python_lib_paths = global.DEFAULT_PYTHON_LIB_REQUIREMENTS
	}

	runner := runner.TempDirRunner{}
	err := runner.WithTempDir("/", []string{}, func(root_path string) error {
		err := os.WriteFile(path.Join(root_path, "env.sh"), []byte(env_script), 0755)
		if err != nil {
			return err
		}

		for _, lib_path := range python_lib_paths {
			// check if the lib path is available
			if _, err := os.Stat(lib_path); err != nil {
				mlog.Warningf("python lib path %s is not available", lib_path)
				continue
			}
			exec_cmd := exec.Command(
				"bash",
				path.Join(root_path, "env.sh"),
				lib_path,
				LIB_PATH,
			)
			exec_cmd.Stderr = os.Stderr

			if err := exec_cmd.Run(); err != nil {
				return err
			}
		}

		os.RemoveAll(root_path)
		os.Remove(root_path)
		return nil
	})

	return err
}
