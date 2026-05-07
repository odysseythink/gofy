package jinja2

import (
	"fmt"

	"github.com/odysseythink/gofy/backend/core/helper/code_executor/template_transformer/base"
)

type Jinja2TemplateTransformer struct {
	*base.BaseTemplateTransformer
}

func (transformer *Jinja2TemplateTransformer) TransformResponse(abs_transformer base.TemplateTransformer, response string) map[string]any {
	return map[string]any{"result": transformer.ExtractResultStrFromResponse(response)}

}
func (transformer *Jinja2TemplateTransformer) GetRunnerScript() string {
	runner_script := fmt.Sprintf(`
# declare main function
def main(**inputs):
    import jinja2
    template = jinja2.Template('''%s''')
    return template.render(**inputs)
    
import json
from base64 import b64decode

# decode and prepare input dict
inputs_obj = json.loads(b64decode('%s').decode('utf-8'))

# execute main function
output = main(**inputs_obj)

# convert output and print
result = f'''<<RESULT>>{{output}}<<RESULT>>'''
print(result)
            
`, base.CODE_PLACEHOLDER, base.INPUTS_PLACEHOLDER)
	return runner_script

}
func (transformer *Jinja2TemplateTransformer) GetPreloadScript() string {
	preload_script := `
import jinja2
from base64 import b64decode

def _jinja2_preload_():
    # prepare jinja2 environment, load template and render before to avoid sandbox issue
    template = jinja2.Template('{{s}}')
    template.render(s='a')

if __name__ == '__main__':
    _jinja2_preload_()

`

	return preload_script
}
