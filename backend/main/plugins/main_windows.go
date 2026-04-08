//go:build windows
// +build windows

package main

import (
	"flag"
	"log"
	"os"

	"mlib.com/confy"
	"mlib.com/gofy/server/cluster"
	"mlib.com/mrun"

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
	confy.SetConfigFile(cfgfile)
	confy.SetConfigType("yaml")
	err := confy.ReadInConfig()
	if err != nil {
		log.Printf("read config file(%s) failed: %v\n", cfgfile, err)
		return
	}
	{
		logpath := confy.GetWithDefault[string]("log.path", "logs")
		loglevel := confy.GetWithDefault[uint32]("log.log_level", 1)
		log.Println("******loglevel=", loglevel)
		if loglevel >= 4 {
			loglevel = 1
		}
		mlog.SetLogLevel(loglevel)
		mlog.SetLogDir(logpath)
	}
	defer mlog.Flush()
	confy.WatchConfig()

	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&Plugins})

	err = mrun.Run(&Plugins)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
