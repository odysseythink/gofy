package v1

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
)

type TagApi struct {
}

func (api *TagApi) List(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}
	tag_type := c.Query("type")
	keyword := c.Query("keyword")

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetTagList(c, &pbapi.GetTagListRequest{TenantId: tenant_id, Type: tag_type, Keyword: keyword})
		if err != nil {
			mlog.Errorf("remote call GetTagList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetTagList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetTagList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.Tags == nil {
					pbrsp.Tags = make([]*pbapi.TagField, 0)
				}
				c.JSON(http.StatusOK, pbrsp.Tags)
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *TagApi) Add(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !(acc.IsEditor() || acc.IsDatasetEditor()) {
		mlog.Errorf("user=%#v must be editor or dataset editor")
		response.Forbidden(c)
		return
	}
	in := pbapi.AddTagRequest{AccountId: acc.ID}
	if err := c.ShouldBindJSON(&in); err != nil {
		mlog.Errorf("bind param failed:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if len(in.Name) < 1 || len(in.Name) > 50 {
		mlog.Errorf("Name must be between 1 to 50 characters.")
		response.InvalidArgErrorWithDetail(c, "Name must be between 1 to 50 characters.")
		return
	}
	if !slices.Contains(models.TAG_TYPE_LIST, in.Type) {
		mlog.Errorf("Invalid tag type=%s.", in.Type)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("Invalid tag type=%s.", in.Type))
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).AddTag(c, &in)
		if err != nil {
			mlog.Errorf("remote call AddTag failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call AddTag failed",
			})
			return
		} else {
			mlog.Infof("remote call AddTag return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, pbrsp)
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *TagApi) Update(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !(acc.IsEditor() || acc.IsDatasetEditor()) {
		mlog.Errorf("user=%#v must be editor or dataset editor")
		response.Forbidden(c)
		return
	}
	tag_id := c.Param("tag_id")
	if tag_id == "" {
		mlog.Errorf("missing tag_id")
		response.InvalidArgErrorWithDetail(c, "missing tag_id")
		return
	}
	in := pbapi.UpdateTagRequest{AccountId: acc.ID, TagId: tag_id}
	if err := c.ShouldBindJSON(&in); err != nil {
		mlog.Errorf("bind param failed:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if len(in.Name) < 1 || len(in.Name) > 50 {
		mlog.Errorf("Name must be between 1 to 50 characters.")
		response.InvalidArgErrorWithDetail(c, "Name must be between 1 to 50 characters.")
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).UpdateTag(c, &in)
		if err != nil {
			mlog.Errorf("remote call UpdateTag failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call UpdateTag failed",
			})
			return
		} else {
			mlog.Infof("remote call UpdateTag return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, pbrsp)
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *TagApi) Del(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !(acc.IsEditor() || acc.IsDatasetEditor()) {
		mlog.Errorf("user=%#v must be editor or dataset editor")
		response.Forbidden(c)
		return
	}
	tag_id := c.Param("tag_id")
	if tag_id == "" {
		mlog.Errorf("missing tag_id")
		response.InvalidArgErrorWithDetail(c, "missing tag_id")
		return
	}
	in := pbapi.UpdateTagRequest{AccountId: acc.ID, TagId: tag_id}
	if err := c.ShouldBindJSON(&in); err != nil {
		mlog.Errorf("bind param failed:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if len(in.Name) < 1 || len(in.Name) > 50 {
		mlog.Errorf("Name must be between 1 to 50 characters.")
		response.InvalidArgErrorWithDetail(c, "Name must be between 1 to 50 characters.")
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).DelTag(c, &pbapi.DelTagRequest{AccountId: acc.ID, TagId: tag_id})
		if err != nil {
			mlog.Errorf("remote call DelTag failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call DelTag failed",
			})
			return
		} else {
			mlog.Infof("remote call DelTag return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, nil)
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *TagApi) AddTagBinding(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !(acc.IsEditor() || acc.IsDatasetEditor()) {
		mlog.Errorf("user=%#v must be editor or dataset editor")
		response.Forbidden(c)
		return
	}

	in := pbapi.AddTagBindingRequest{AccountId: acc.ID}
	if err := c.ShouldBindJSON(&in); err != nil {
		mlog.Errorf("bind param failed:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if !slices.Contains(models.TAG_TYPE_LIST, in.Type) {
		mlog.Errorf("Invalid tag type=%s.", in.Type)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("Invalid tag type=%s.", in.Type))
		return
	}
	if in.TargetId == "" {
		mlog.Errorf("missing target_id")
		response.InvalidArgErrorWithDetail(c, "missing target_id")
		return
	}
	if len(in.TagIds) < 1 {
		mlog.Errorf("Tag IDs is required.")
		response.InvalidArgErrorWithDetail(c, "Tag IDs is required.")
		return
	}
	for _, v := range in.TagIds {
		if v == "" {
			mlog.Errorf("item of Tag IDs=%#v can't be empty", in.TagIds)
			response.InvalidArgErrorWithDetail(c, fmt.Sprintf("item of Tag IDs=%#v can't be empty", in.TagIds))
			return
		}
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).AddTagBinding(c, &in)
		if err != nil {
			mlog.Errorf("remote call AddTagBinding failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call AddTagBinding failed",
			})
			return
		} else {
			mlog.Infof("remote call AddTagBinding return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, nil)
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
func (api *TagApi) DelTagBinding(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !(acc.IsEditor() || acc.IsDatasetEditor()) {
		mlog.Errorf("user=%#v must be editor or dataset editor")
		response.Forbidden(c)
		return
	}

	in := pbapi.DelTagBindingRequest{AccountId: acc.ID}
	if err := c.ShouldBindJSON(&in); err != nil {
		mlog.Errorf("bind param failed:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if !slices.Contains(models.TAG_TYPE_LIST, in.Type) {
		mlog.Errorf("Invalid tag type=%s.", in.Type)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("Invalid tag type=%s.", in.Type))
		return
	}
	if in.TargetId == "" {
		mlog.Errorf("missing target_id")
		response.InvalidArgErrorWithDetail(c, "missing target_id")
		return
	}
	if in.TagId == "" {
		mlog.Errorf("missing tag_id")
		response.InvalidArgErrorWithDetail(c, "missing tag_id")
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).DelTagBinding(c, &in)
		if err != nil {
			mlog.Errorf("remote call DelTagBinding failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call DelTagBinding failed",
			})
			return
		} else {
			mlog.Infof("remote call DelTagBinding return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, nil)
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
