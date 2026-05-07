package response

import (
	"encoding/json"
	"time"

	"github.com/odysseythink/gofy/backend/core/file"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type SimpleAccountResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountResponse struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Avatar            string `json:"avatar"`
	AvatarUrl         string `json:"avatar_url"`
	Email             string `json:"email"`
	IsPasswordSet     bool   `json:"is_password_set"`
	InterfaceLanguage string `json:"interface_language"`
	InterfaceTheme    string `json:"interface_theme"`
	Timezone          string `json:"timezone"`
	LastLoginAt       int64  `json:"last_login_at"`
	LastLoginIp       string `json:"last_login_ip"`
	CreatedAt         int64  `json:"created_at"`
}

func NewAccountResponse(args any) *AccountResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(AccountResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AccountResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(AccountResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AccountResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if acc, ok := args.(*models.Account); ok && acc != nil {
		ret := &AccountResponse{
			Id:                acc.ID,
			Name:              acc.Name,
			Avatar:            acc.Avatar,
			Email:             acc.Email,
			IsPasswordSet:     acc.IsPasswordSet(),
			InterfaceLanguage: acc.InterfaceLanguage,
			InterfaceTheme:    acc.InterfaceTheme,
			Timezone:          acc.Timezone,
			LastLoginIp:       acc.LastLoginIP,
		}
		now := time.Now()
		ret.AvatarUrl = file.GetSignedFileURL(acc.Avatar)
		if acc.LastLoginAt != nil {
			ret.LastLoginAt = acc.LastLoginAt.Unix()
		} else {
			ret.LastLoginAt = now.Unix()
		}
		if acc.CreatedAt != nil {
			ret.CreatedAt = acc.CreatedAt.Unix()
		} else {
			ret.CreatedAt = now.Unix()
		}
		return ret
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(AccountResponse)
}

type AccountWithRoleResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	AvatarUrl    string `json:"avatar_url"`
	Email        string `json:"email"`
	LastLoginAt  int64  `json:"last_login_at"`
	LastActiveAt int64  `json:"last_active_at"`
	CreatedAt    int64  `json:"created_at"`
	Role         string `json:"role"`
	Status       string `json:"status"`
}

func NewAccountWithRoleResponse(args any) *AccountWithRoleResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(AccountWithRoleResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AccountWithRoleResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(AccountWithRoleResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AccountWithRoleResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if acc, ok := args.(*models.Account); ok && acc != nil {
		ret := &AccountWithRoleResponse{
			Id:     acc.ID,
			Name:   acc.Name,
			Avatar: acc.Avatar,
			Email:  acc.Email,
			Role:   acc.Role(),
			Status: string(acc.Status),
		}
		now := time.Now()
		ret.AvatarUrl = file.GetSignedFileURL(acc.Avatar)
		if acc.LastLoginAt != nil {
			ret.LastLoginAt = acc.LastLoginAt.Unix()
		} else {
			ret.LastLoginAt = now.Unix()
		}
		if acc.LastActiveAt != nil {
			ret.LastActiveAt = acc.LastActiveAt.Unix()
		} else {
			ret.LastActiveAt = now.Unix()
		}
		if acc.CreatedAt != nil {
			ret.CreatedAt = acc.CreatedAt.Unix()
		} else {
			ret.CreatedAt = now.Unix()
		}
		return ret
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(AccountWithRoleResponse)
}

type AccountWithRoleListResponse struct {
	Accounts []*AccountWithRoleResponse `json:"accounts"`
}

func NewAccountWithRoleListResponse(args any) *AccountWithRoleListResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(AccountWithRoleListResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AccountWithRoleListResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(AccountWithRoleListResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to AccountWithRoleListResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else if accs, ok := args.([]*models.Account); ok && len(accs) > 0 {
		rsp := &AccountWithRoleListResponse{}
		for _, acc := range accs {
			if acc == nil {
				continue
			}
			if rsp.Accounts == nil {
				rsp.Accounts = make([]*AccountWithRoleResponse, 0)
			}
			rsp.Accounts = append(rsp.Accounts, NewAccountWithRoleResponse(acc))
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(AccountWithRoleListResponse)
}
