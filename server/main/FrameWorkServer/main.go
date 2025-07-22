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

	// 工作模块名称
	modulename := viper.GetString("ServerCom.ModuleName")
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
