package sensitivewordavoidance

import (
	"mlib.com/gofy/server/core/exceptions"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
)

type SensitiveWordAvoidanceConfigManager struct{}

func (mgr *SensitiveWordAvoidanceConfigManager) Convert(config map[string]any) *appconfigentities.SensitiveWordAvoidanceEntity {
	var sensitive_word_avoidance_dict map[string]any
	if _, ok := config["sensitive_word_avoidance"]; ok {
		if _, ok := config["sensitive_word_avoidance"].(map[string]any); ok {
			sensitive_word_avoidance_dict = config["sensitive_word_avoidance"].(map[string]any)
		}
	}
	if len(sensitive_word_avoidance_dict) == 0 {
		return nil
	}
	var enabled bool
	if _, ok := sensitive_word_avoidance_dict["enabled"]; ok {
		if _, ok := sensitive_word_avoidance_dict["enabled"].(bool); ok {
			enabled = sensitive_word_avoidance_dict["enabled"].(bool)
		}
	}
	var strtype string
	if _, ok := sensitive_word_avoidance_dict["type"]; ok {
		if _, ok := sensitive_word_avoidance_dict["type"].(string); ok {
			strtype = sensitive_word_avoidance_dict["type"].(string)
		}
	}
	var sconfig map[string]any
	if _, ok := sensitive_word_avoidance_dict["config"]; ok {
		if _, ok := sensitive_word_avoidance_dict["config"].(map[string]any); ok {
			sconfig = sensitive_word_avoidance_dict["config"].(map[string]any)
		}
	}
	if enabled {
		return &appconfigentities.SensitiveWordAvoidanceEntity{
			Type:   strtype,
			Config: sconfig,
		}
	} else {
		return nil
	}
}
func (mgr *SensitiveWordAvoidanceConfigManager) ValidateAndSetDefaults(
	tenant_id string, config map[string]any, only_structure_validate bool, /* = False*/
) (map[string]any, []string, error) {
	var sensitive_word_avoidance_dict map[string]any
	if _, ok := config["sensitive_word_avoidance"]; ok {
		if _, ok := config["sensitive_word_avoidance"].(map[string]any); ok {
			sensitive_word_avoidance_dict = config["sensitive_word_avoidance"].(map[string]any)
		} else {
			return nil, nil, exceptions.NewValueError("sensitive_word_avoidance must be of dict type")
		}
	} else {
		sensitive_word_avoidance_dict = map[string]any{"enabled": false}
		config["sensitive_word_avoidance"] = sensitive_word_avoidance_dict
	}
	if v, ok := config["sensitive_word_avoidance"].(map[string]any)["enabled"]; !ok || v == nil {
		config["sensitive_word_avoidance"].(map[string]any)["enabled"] = false
	} else {
		if _, ok := v.(bool); !ok {
			config["sensitive_word_avoidance"].(map[string]any)["enabled"] = false
		}
	}

	if config["sensitive_word_avoidance"].(map[string]any)["enabled"].(bool) {
		if v, ok := config["sensitive_word_avoidance"].(map[string]any)["type"]; !ok || v == nil {
			return nil, nil, exceptions.NewValueError("sensitive_word_avoidance.type is required")
		} else {
			if _, ok := v.(string); !ok {
				return nil, nil, exceptions.NewValueError("sensitive_word_avoidance.type must be string")
			}
		}

		if !only_structure_validate {
			// typ = config["sensitive_word_avoidance"]["type"].(string)
			// sensitive_word_avoidance_config = config["sensitive_word_avoidance"]["config"]

			// ModerationFactory.validate_config(name=typ, tenant_id=tenant_id, config=sensitive_word_avoidance_config)
		}
	}
	return config, []string{"sensitive_word_avoidance"}, nil
}
