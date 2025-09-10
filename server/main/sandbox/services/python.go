package services

import (
	"time"

	"mlib.com/confy"
	"mlib.com/gofy/server/main/sandbox/runner/python"
	runnertypes "mlib.com/gofy/server/main/sandbox/runner/types"
	"mlib.com/gofy/server/main/sandbox/types"
	"mlib.com/gofy/server/main/sandbox/utils"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type PythonService struct {
}
type RunCodeResponse struct {
	Stderr string `json:"error"`
	Stdout string `json:"stdout"`
}

func (s *PythonService) RunPython3Code(code string, preload string, options *runnertypes.RunnerOptions) *types.SandboxResponse {
	if err := utils.CheckOptions(options); err != nil {
		mlog.Errorf("check options failed:%v", err)
		return types.ErrorResponse(-400, err.Error())
	}

	if !confy.GetWithDefault[bool]("enable_preload", false) {
		preload = ""
	}

	timeout := time.Duration(
		confy.GetWithDefault[int]("worker_timeout", 0) * int(time.Second),
	)

	runner := python.PythonRunner{}
	stdout, stderr, done, err := runner.Run(
		code, timeout, nil, preload, options,
	)
	if err != nil {
		mlog.Errorf("run failed:%v", err)
		return types.ErrorResponse(-500, err.Error())
	}

	stdout_str := ""
	stderr_str := ""

	defer close(done)
	defer close(stdout)
	defer close(stderr)

	for {
		select {
		case <-done:
			return types.SuccessResponse(&RunCodeResponse{
				Stdout: stdout_str,
				Stderr: stderr_str,
			})
		case out := <-stdout:
			// fmt.Printf(string(out))
			stdout_str += string(out)
		case err := <-stderr:
			// fmt.Printf(string(err))
			stderr_str += string(err)
		}
	}
}

func (s *PythonService) ListPython3Dependencies() *pbapi.DependenciesReply {
	return &pbapi.DependenciesReply{
		Code:         0,
		Message:      "success",
		Dependencies: python.ListDependencies(),
	}
}
func (s *PythonService) UpdateDependencies() *pbapi.DependenciesReply {
	err := python.PreparePythonDependenciesEnv()
	if err != nil {
		mlog.Errorf("PreparePythonDependenciesEnv failed:%v", err)
		return &pbapi.DependenciesReply{
			Code:    -1,
			Message: err.Error(),
		}
	}
	return &pbapi.DependenciesReply{
		Code:    0,
		Message: "success",
	}
}

func (s *PythonService) RefreshPython3Dependencies() *pbapi.DependenciesReply {
	return &pbapi.DependenciesReply{
		Code:         0,
		Message:      "success",
		Dependencies: python.RefreshDependencies(),
	}
}
