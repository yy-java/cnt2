// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	//	"github.com/astaxie/beego/plugins/cors"
	"github.com/yy-java/cnt2/httpserver/controllers"
)

func init() {

	beego.AddNamespace(beego.NewNamespace("/config",
		beego.NSInclude(
			&controllers.ConfigController{},
		),
	))
	beego.AddNamespace(beego.NewNamespace("/app",
		beego.NSInclude(
			&controllers.AppController{},
		),
	))
	beego.AddNamespace(beego.NewNamespace("/node",
		beego.NSInclude(
			&controllers.NodeController{},
		),
	))
	beego.AddNamespace(beego.NewNamespace("/publish",
		beego.NSInclude(
			&controllers.PublishController{},
		),
	))
	beego.AddNamespace(beego.NewNamespace("/userauth",
		beego.NSInclude(
			&controllers.UserAuthController{},
		),
	))
	beego.AddNamespace(beego.NewNamespace("/profile",
		beego.NSInclude(
			&controllers.ProfileController{},
		),
	))
}
