package pbapi

import (
	"github.com/odysseythink/gofy/backend/cluster"
	"google.golang.org/grpc"
)

func init() {
	for _, sd := range []*grpc.ServiceDesc{
		&Admin_ServiceDesc,
		&App_ServiceDesc,
		&Datasets_ServiceDesc,
		&Link_ServiceDesc,
		&Plugins_ServiceDesc,
		&Sandbox_ServiceDesc,
		&Tools_ServiceDesc,
	} {
		cluster.RegisterGrpcServiceDesc(sd)
	}
}
