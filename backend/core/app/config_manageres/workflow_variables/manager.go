package workflowvariables

import (
	"encoding/json"

	"github.com/odysseythink/mlog"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	"mlib.com/gofy/server/models"
)

type WorkflowVariablesConfigManager struct{}

func (mgr *WorkflowVariablesConfigManager) Convert(wf *models.Workflow) []*appconfigentities.VariableEntity {
	/*
	   Convert workflow start variables to variables

	   :param workflow: workflow instance
	*/
	vs := []*appconfigentities.VariableEntity{}

	// find start node
	user_input_form := wf.UserInputForm(false)

	// variables
	for _, v := range user_input_form {
		tmpdata, _ := json.Marshal(v)
		ve := new(appconfigentities.VariableEntity)
		err := json.Unmarshal(tmpdata, ve)
		if err != nil {
			mlog.Errorf("json unmarshal failed:%v", err)
		} else {
			vs = append(vs, ve)
		}

	}
	return vs
}
