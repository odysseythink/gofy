package main

import (
	"context"
	"fmt"

	"mlib.com/gofy/server/cache"
	dbengine "mlib.com/gofy/server/db_engine"
	eventhandlers "mlib.com/gofy/server/events/event_handlers"

	_ "mlib.com/gofy/server/events"
	"mlib.com/gofy/server/proto/pbapi"

	_ "mlib.com/gofy/server/core/model_runtime/model_provides"
	"mlib.com/mlog"
)

type AdminService struct {
	pbapi.UnimplementedAdminServer
}

var (
	Admin AdminService
)

// go build -o admin.so -buildmode=plugin main.go
func (s *AdminService) Init(args ...any) error {
	mlog.Info("admin init......")
	err := cache.Instance().Init()
	if err != nil {
		mlog.Errorf("cahe init failed:%v", err)
		return fmt.Errorf("cahe init failed:%v", err)
	}
	err = dbengine.Instance().Init()
	if err != nil {
		mlog.Errorf("mysql init failed:%v", err)
		return fmt.Errorf("mysql init failed:%v", err)
	}
	eventhandlers.Init()

	return nil
}

func (s *AdminService) RunOnce(ctx context.Context) error {
	return nil
}
func (s *AdminService) Destroy() {

}

func (s *AdminService) UserData() any {
	return nil
}
