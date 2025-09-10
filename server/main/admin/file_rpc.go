package main

import (
	"context"

	"google.golang.org/grpc/peer"
	"mlib.com/confy"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

func (s *AdminService) GetFileUploadConfig(ctx context.Context, in *pbapi.GetFileUploadConfigRequest) (out *pbapi.GetFileUploadConfigReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetFileUploadConfig call:%#v", p.Addr.String(), in)

	// check_update_url := confy.Get[string]("check_update_url")
	out = &pbapi.GetFileUploadConfigReply{
		FileSizeLimit:           confy.GetWithDefault[int64]("upload.file_size_limit", 15),
		BatchCountLimit:         confy.GetWithDefault[int64]("upload.file_batch_limit", 5),
		ImageFileSizeLimit:      confy.GetWithDefault[int64]("upload.image_file_size_limit", 10),
		VideoFileSizeLimit:      confy.GetWithDefault[int64]("upload.video_file_size_limit", 100),
		AudioFileSizeLimit:      confy.GetWithDefault[int64]("upload.audio_file_size_limit", 50),
		WorkflowFileUploadLimit: confy.GetWithDefault[int64]("upload.workflow_file_limit", 10),
	}
	return
}
