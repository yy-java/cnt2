package controllers

import (
	"github.com/astaxie/beego"
	"log"
	"net/url"
	"strings"
	"time"
	"github.com/yy-java/cnt2/db"
	"github.com/yy-java/cnt2/httpserver/utils"
	"github.com/yy-java/cnt2/service/errors"
	userService "github.com/yy-java/cnt2/service/user"
)

// Operations about App
type BaseController struct {
	beego.Controller
	Uid      int64
	Username string
}

var whitelist = map[string]int{
	"/app/test":        0,
	"/app/create":      0,
	"/publish/queryIp": 0,
}

func getReferHost(uri string) string {
	u, err := url.Parse(uri)
	if err == nil {
		return u.Scheme + "://" + u.Host
	}
	return uri
}
func (c *BaseController) Prepare() {
	//添加response
	header := c.Ctx.ResponseWriter.Header()
	header.Set("Access-Control-Allow-Origin", getReferHost(c.Ctx.Request.Header.Get("Referer")))
	header.Set("Access-Control-Allow-Credentials", "true")

	uid := c.GetSession("uid")
	if uid != nil{
		c.Uid = uid.(int64)
		c.Username = c.GetSession("username").(string)
	}
	username := c.GetSession("username")
	if username != nil{
		c.Username = username.(string)
	}

	c.Data["startTime"] = time.Now().UnixNano()

	reqURI := c.Ctx.Request.RequestURI
	log.Printf("path:%s, uid:%s, username:%s", reqURI, c.Uid, c.Username)
	if c.Uid <= 0 && strings.IndexAny(reqURI,"/user/login") == -1 {
		c.JsonResp(nil, errors.ErrUnauthenticated)
		c.ServeJSON()
	}
	//权限控制
	isAdmin := c.GetSession("isAdmin")
	log.Printf("isAdmin：%s", isAdmin)
	if nil == isAdmin {
		auth, err := userService.FindUserAuthByInput(&db.UserAuth{Uid: c.Uid, Permission: int8(db.PermissionType_Admin)})
		log.Printf("auth____________________________________：%s", auth)
		if err == nil {
			if len(auth) > 0 {
				c.SetSession("isAdmin", true)
			} else {
				c.SetSession("isAdmin", false)
			}
		}
	}
	_, ok := whitelist[strings.Split(reqURI, "?")[0]]
	if ok {
		return
	}
	log.Printf("it's not white list")
	app := c.GetString("app") //query
	if len(app) < 0 {
		app = c.GetString(":app") //path
	}
	if len(app) > 0 {
		if c.HasPermission(app) == false {
			c.JsonResp(nil, errors.ErrPermissionDenied)
			c.ServeJSON()
		}
	}

}

func (c *BaseController) Finish() {
	endTime := time.Now().UnixNano()
	startTime, ok := c.Data["startTime"].(int64)
	if !ok {
		return
	}
	delete(c.Data, "startTime")

	log.Printf("%s cost:%dms", c.Ctx.Request.RequestURI, (endTime-startTime)/1000000)
}
func (c *BaseController) JsonpResp(data interface{}, err error) {
	c.Data["jsonp"] = utils.Resp(data, err)
}
func (c *BaseController) JsonResp(data interface{}, err error) {
	c.Data["json"] = utils.Resp(data, err)
}

func (c *BaseController) HasPermission(app string) bool {
	if c.GetSession("isAdmin") == true {
		return true
	}
	auth, err := userService.FindUserAuthByInput(&db.UserAuth{Uid: c.Uid, App: app})
	if err == nil {
		if len(auth) > 0 {
			return true
		}
	}
	return false
}
