package openingstatement

type OpeningStatementConfigManager struct{}

func (mgr *OpeningStatementConfigManager) Convert(config map[string]any) (string, []string) {
	/*
	   Convert model config to model config

	   :param config: model config args
	*/
	// opening statement
	var opening_statement string
	if _, ok := config["opening_statement"]; ok {
		if _, ok := config["opening_statement"].(string); ok {
			opening_statement = config["opening_statement"].(string)
		}
	}

	// suggested questions
	var suggested_questions []string
	if _, ok := config["suggested_questions"]; ok {
		if _, ok := config["suggested_questions"].([]string); ok {
			suggested_questions = config["suggested_questions"].([]string)
		}
	}

	return opening_statement, suggested_questions
}
func (mgr *OpeningStatementConfigManager) ValidateAndSetDefaults(config map[string]any) (map[string]any, []string) {
	/*
	   Validate and set defaults for opening statement feature

	   :param config: app model config args
	*/
	if _, ok := config["opening_statement"]; ok {
		if _, ok := config["opening_statement"].(string); !ok {
			config["opening_statement"] = ""
		}
	} else {
		config["opening_statement"] = ""
	}

	// suggested_questions
	if _, ok := config["suggested_questions"]; ok {
		if _, ok := config["suggested_questions"].([]string); !ok {
			config["suggested_questions"] = []string{}
		}
	} else {
		config["suggested_questions"] = []string{}
	}

	return config, []string{"opening_statement", "suggested_questions"}
}
