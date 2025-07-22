//go:build windows
// +build windows

package main

import (
	"flag"
	"log"
	"os"

	"mlib.com/gofy/server/cluster"
	"mlib.com/mrun"

	"github.com/spf13/viper"

	"mlib.com/mlog"
)

func main() {
	var cfgfile string
	flag.StringVar(&cfgfile, "c", "", "choose config file.")
	flag.Parse()
	if cfgfile == "" {
		log.Println("usage: ./server -c config.yml")
		return
	}
	viper.SetConfigFile(cfgfile)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("read config file(%s) failed: %v\n", cfgfile, err)
		return
	}
	{
		logpath := viper.GetStringWithDefault("log.path", "logs")
		loglevel := viper.GetUint32WithDefault("log.log_level", 1)
		log.Println("******loglevel=", loglevel)
		if loglevel >= 4 {
			loglevel = 1
		}
		mlog.SetLogLevel(loglevel)
		mlog.SetLogDir(logpath)
	}
	defer mlog.Flush()
	viper.WatchConfig()
	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&Link})

	err = mrun.Run(&Link)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
