package request

import (
	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/core/exceptions"
)

type AppPageReq struct {
	BasePageReq
	Mode          string   `json:"mode" form:"mode"` // "chat", "workflow", "agent-chat", "channel", "all"
	Name          string   `json:"name" form:"name"`
	TagIDs        []string `json:"tag_ids" form:"tag_ids"`
	IsCreatedByMe bool     `json:"is_created_by_me" form:"is_created_by_me"`
}

func NewAppPageReq() *AppPageReq {
	return &AppPageReq{
		Mode: "all",
		BasePageReq: BasePageReq{
			Page:  1,
			Limit: 20,
		},
	}
}

func (req *AppPageReq) Valid() (bool, *exceptions.ValueError) {
	if req.Page < 1 || req.Page > 99999 {
		return false, exceptions.NewValueError("invalid page")
	}
	if req.Limit < 1 || req.Limit > 100 {
		return false, exceptions.NewValueError("invalid limit")
	}
	if req.Mode != "all" && req.Mode != "chat" && req.Mode != "workflow" && req.Mode != "agent-chat" && req.Mode != "channel" {
		return false, exceptions.NewValueError("invalid mode")
	}
	for _, v := range req.TagIDs {
		if _, err := uuid.FromString(v); err != nil {
			return false, exceptions.NewValueError("invalid tag_ids")
		}
	}
	return true, nil
}
