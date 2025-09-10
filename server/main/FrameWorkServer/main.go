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

	// 工作模块名称
	modulename := confy.Get[string]("ServerCom.ModuleName")
	if modulename == "" {
		mlog.Errorf("ModuleName is not exist")
		return
	}

	mlog.Infof("ModuleName: %s", modulename)
	err = mrun.RegisterLibsoWithModule(modulename+".so", cluster.GetExportModulename(modulename), nil, nil)
	if err != nil {
		mlog.Errorf("mrun.RegisterLibso(%s, nil, nil) failed: %v", cluster.GetExportModulename(modulename), err)
		return
	}
	// fmt.Println("+++++++++++++++++++=", reflect.TypeOf(cluster.Instance()).Elem().PkgPath())
	err = mrun.Run(cluster.Instance())
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
