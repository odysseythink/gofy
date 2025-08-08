package python

import (
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"mlib.com/gofy/server/main/sandbox/global"
	"mlib.com/gofy/server/main/sandbox/runner"
	"mlib.com/gofy/server/main/sandbox/runner/types"
	"mlib.com/mlog"
)

type PythonRunner struct {
	runner.TempDirRunner
}

//go:embed prescript.py
var sandbox_fs []byte

func (p *PythonRunner) Run(
	code string,
	timeout time.Duration,
	stdin []byte,
	preload string,
	options *types.RunnerOptions,
) (chan []byte, chan []byte, chan bool, error) {
	// initialize the environment
	untrusted_code_path, key, err := p.InitializeEnvironment(code, preload, options)
	if err != nil {
		mlog.Errorf("InitializeEnvironment failed:%v", err)
		return nil, nil, nil, err
	}

	// capture the output
	output_handler := runner.NewOutputCaptureRunner()
	output_handler.SetTimeout(timeout)
	output_handler.SetAfterExitHook(func() {
		// remove untrusted code
		os.Remove(untrusted_code_path)
	})

	// create a new process
	cmd := exec.Command(
		viper.GetString("python_path"),
		untrusted_code_path,
		LIB_PATH,
		key,
	)
	cmd.Env = []string{}
	cmd.Dir = LIB_PATH

	if viper.GetString("proxy.socks5") != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("HTTPS_PROXY=%s", viper.GetString("proxy.socks5")))
		cmd.Env = append(cmd.Env, fmt.Sprintf("HTTP_PROXY=%s", viper.GetString("proxy.socks5")))
	} else if viper.GetString("proxy.https") != "" || viper.GetString("proxy.http") != "" {
		if viper.GetString("proxy.https") != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("HTTPS_PROXY=%s", viper.GetString("proxy.https")))
		}
		if viper.GetString("proxy.http") != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("HTTP_PROXY=%s", viper.GetString("proxy.http")))
		}
	}
	allowed_syscalls := viper.GetIntSlice("allowed_syscalls")
	if len(allowed_syscalls) > 0 {
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("ALLOWED_SYSCALLS=%s",
				strings.Trim(strings.Join(strings.Fields(fmt.Sprint(allowed_syscalls)), ","), "[]"),
			),
		)
	}

	err = output_handler.CaptureOutput(cmd)
	if err != nil {
		mlog.Errorf("CaptureOutput failed:%v", err)
		return nil, nil, nil, err
	}

	return output_handler.GetStdout(), output_handler.GetStderr(), output_handler.GetDone(), nil
}

func (p *PythonRunner) InitializeEnvironment(code string, preload string, options *types.RunnerOptions) (string, string, error) {
	if !checkLibAvaliable() {
		// ensure environment is reversed
		releaseLibBinary(false)
	}
	uid := viper.GetIntWithDefault("sandbox_user_uid", 65537)
	// create a tmp dir and copy the python script
	temp_code_name := strings.ReplaceAll(uuid.NewV4().String(), "-", "_")
	temp_code_name = strings.ReplaceAll(temp_code_name, "/", ".")

	script := strings.Replace(
		string(sandbox_fs),
		"{{uid}}", strconv.Itoa(uid), 1,
	)

	script = strings.Replace(
		script,
		"{{gid}}", strconv.Itoa(global.SANDBOX_GROUP_ID), 1,
	)

	if options.EnableNetwork {
		script = strings.Replace(
			script,
			"{{enable_network}}", "1", 1,
		)
	} else {
		script = strings.Replace(
			script,
			"{{enable_network}}", "0", 1,
		)
	}

	script = strings.Replace(
		script,
		"{{preload}}",
		fmt.Sprintf("%s\n", preload),
		1,
	)

	// generate a random 512 bit key
	key_len := 64
	key := make([]byte, key_len)
	_, err := rand.Read(key)
	if err != nil {
		return "", "", err
	}

	// encrypt the code
	encrypted_code := make([]byte, len(code))
	for i := 0; i < len(code); i++ {
		encrypted_code[i] = code[i] ^ key[i%key_len]
	}

	// encode code using base64
	code = base64.StdEncoding.EncodeToString(encrypted_code)
	// encode key using base64
	encoded_key := base64.StdEncoding.EncodeToString(key)

	code = strings.Replace(
		script,
		"{{code}}",
		code,
		1,
	)

	untrusted_code_path := fmt.Sprintf("%s/tmp/%s.py", LIB_PATH, temp_code_name)
	err = os.MkdirAll(path.Dir(untrusted_code_path), 0755)
	if err != nil {
		return "", "", err
	}
	err = os.WriteFile(untrusted_code_path, []byte(code), 0755)
	if err != nil {
		return "", "", err
	}

	return untrusted_code_path, encoded_key, nil
}
