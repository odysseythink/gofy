package main

import (
	"context"

	"github.com/spf13/viper"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

func (s *AdminService) GetFileUploadConfig(ctx context.Context, in *pbapi.GetFileUploadConfigRequest) (out *pbapi.GetFileUploadConfigReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetFileUploadConfig call:%#v", p.Addr.String(), in)

	// check_update_url := viper.GetString("check_update_url")
	out = &pbapi.GetFileUploadConfigReply{
		FileSizeLimit:           viper.GetInt64WithDefault("upload.file_size_limit", 15),
		BatchCountLimit:         viper.GetInt64WithDefault("upload.file_batch_limit", 5),
		ImageFileSizeLimit:      viper.GetInt64WithDefault("upload.image_file_size_limit", 10),
		VideoFileSizeLimit:      viper.GetInt64WithDefault("upload.video_file_size_limit", 100),
		AudioFileSizeLimit:      viper.GetInt64WithDefault("upload.audio_file_size_limit", 50),
		WorkflowFileUploadLimit: viper.GetInt64WithDefault("upload.workflow_file_limit", 10),
	}
	return
}
