package main

import (
	"context"

	"github.com/spf13/viper"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

func (s *AdminService) GetVersion(ctx context.Context, in *pbapi.GetVersionRequest) (out *pbapi.GetVersionReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetVersion call:%#v", p.Addr.String(), in)

	// check_update_url := viper.GetString("check_update_url")
	out = &pbapi.GetVersionReply{
		Version: viper.GetString("CURRENT_VERSION"),
		Features: &pbapi.VersionFeatures{
			CanReplaceLogo:            viper.GetBool("CAN_REPLACE_LOGO"),
			ModelLoadBalancingEnabled: viper.GetBool("MODEL_LB_ENABLED"),
		},
	}
	return
	// if check_update_url == "" {
	// 	return result
	// }
	//     try:
	//         response = requests.get(check_update_url, {"current_version": args.get("current_version")})
	//     except Exception as error:
	//         logging.warning("Check update version error: {}.".format(str(error)))
	//         result["version"] = args.get("current_version")
	//         return result

	//     content = json.loads(response.content)
	//     if _has_new_version(latest_version=content["version"], current_version=f"{args.get('current_version')}"):
	//         result["version"] = content["version"]
	//         result["release_date"] = content["releaseDate"]
	//         result["release_notes"] = content["releaseNotes"]
	//         result["can_auto_update"] = content["canAutoUpdate"]
	//     return result
	// return
}
