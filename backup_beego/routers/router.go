package routers

import (
	"github.com/astaxie/beego"
	"go-samba4/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
}
