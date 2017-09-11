package routers

import (
	"github.com/kapmahc/fly/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
