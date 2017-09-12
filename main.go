package main

import (
	"log"
	"path"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego/toolbox"
	"github.com/kapmahc/fly/plugins/nut"
	_ "github.com/kapmahc/fly/routers"
	_ "github.com/lib/pq"
)

func main() {
	logs.SetLogger(logs.AdapterConsole)
	logs.SetLogger(logs.AdapterFile, `{"filename":"`+path.Join("tmp", "www.log")+`"}`)

	orm.Debug = beego.BConfig.RunMode != beego.PROD
	orm.RegisterDataBase(
		"default",
		beego.AppConfig.String("databasedriver"),
		beego.AppConfig.String("databasesource"),
	)

	if err := nut.Open(); err != nil {
		log.Panic(err)
		return
	}

	toolbox.StartTask()
	defer toolbox.StopTask()

	beego.ErrorController(&nut.ErrorController{})
	beego.Run()
}
