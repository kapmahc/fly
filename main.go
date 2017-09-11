package main

import (
	"path"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/kapmahc/fly/routers"
)

func main() {
	logs.SetLogger(logs.AdapterConsole)
	logs.SetLogger(logs.AdapterFile, `{"filename":"`+path.Join("tmp", "www.log")+`"}`)
	beego.Run()
}
