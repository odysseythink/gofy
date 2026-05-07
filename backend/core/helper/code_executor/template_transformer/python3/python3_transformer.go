package python3

import (
	"fmt"

	"github.com/odysseythink/gofy/backend/core/helper/code_executor/template_transformer/base"
)

type Python3TemplateTransformer struct {
	*base.BaseTemplateTransformer
}

func (former *Python3TemplateTransformer) GetRunnerScript() string {
	runner_script := fmt.Sprintf(`
# declare main function
%s

import json
from base64 import b64decode

# decode and prepare input dict
inputs_obj = json.loads(b64decode('%s').decode('utf-8'))

# execute main function
output_obj = main(**inputs_obj)

# convert output to json and print
output_json = json.dumps(output_obj)
result = f'''<<RESULT>>{output_json}<<RESULT>>'''
print(result)
            `, base.CODE_PLACEHOLDER, base.INPUTS_PLACEHOLDER)
	return runner_script
}
