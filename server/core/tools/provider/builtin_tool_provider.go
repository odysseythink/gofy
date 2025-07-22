package provider

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/tools/tool"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	"mlib.com/mlog"
)

type BuiltinToolProviderController struct {
	*BaseToolProviderController
}

func (controller *BuiltinToolProviderController) ProviderType() toolsenumtypes.ToolProviderType {
	return toolsenumtypes.ToolProvider_BUILT_IN
}

func (controller *BuiltinToolProviderController) _get_builtin_tools() []*tool.Tool {
	if controller.Tools != nil {
		return controller.Tools
	}
	if controller.Identity == nil {
		return nil
	}
	provider := controller.Identity.Name
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	tool_path := filepath.Join(dir, "builtin", provider, "tools")
	tool_files := []string{}
	err := filepath.Walk(tool_path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(path, ".yaml") && !strings.HasPrefix(path, "__") {
			tool_files = append(tool_files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	tools := []*tool.Tool{}
	for _, tool_file := range tool_files {
		// get tool name
		bindata, err := os.ReadFile(tool_file)
		if err != nil {
			mlog.Errorf("read file=%s failed:%v", tool_file, err)
			panic(exceptions.NewValueError("real tool file failed:" + err.Error()))
		}
		tool := new(tool.Tool)
		err = yaml.Unmarshal(bindata, tool)
		if err != nil {
			mlog.Errorf("yaml Unmarshal file=%s failed:%v", tool_file, err)
			panic(exceptions.NewValueError("yaml Unmarshal file failed:" + err.Error()))
		}
		// tool_name = tool_file.split(".")[0]
		// tool = load_yaml_file(path.join(tool_path, tool_file), ignore_error=False)

		// // get tool class, import the module
		// assistant_tool_class = load_single_subclass_from_source(
		//     module_name=f"core.tools.provider.builtin.{provider}.tools.{tool_name}",
		//     script_path=path.join(
		//         path.dirname(path.realpath(__file__)), "builtin", provider, "tools", f"{tool_name}.py"
		//     ),
		//     parent_type=BuiltinTool,
		// )
		// tool["identity"]["provider"] = provider
		// tools=append(tools, assistant_tool_class(**tool))
	}
	controller.Tools = tools
	return tools
}
