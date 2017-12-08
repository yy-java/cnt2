package controllers

import (
	"log"
	"strings"

	"github.com/yy-java/cnt2/db"
	"github.com/yy-java/cnt2/etcd"
	configService "github.com/yy-java/cnt2/service/config"
	. "github.com/yy-java/cnt2/service/errors"
	userService "github.com/yy-java/cnt2/service/user"
)

type AppController struct {
	BaseController
}

// eg: http://localhost:8081/app/create?app=app10&appType=1&name=test-name2&charger=dyf&chargerUid=
// @router /create [post,get]
func (req *AppController) Create() {
	defer req.ServeJSON()

	app := req.GetString("app")
	appType, _ := req.GetInt8("appType")
	name := req.GetString("name")
	charger := req.GetString("charger")
	chargerUid, _ := req.GetInt64("chargerUid")

	if chargerUid <= 0 {
		chargerUid = req.Uid
	}

	log.Printf("Create app:%s, appType:%d, name:%s, charger:%s", app, appType, name, charger)

	appObj := db.App{App: app}
	appObj.Read()
	if len(appObj.Name) > 0 {
		req.JsonResp(nil, ErrExists)
		return
	}
	if appType < 0 || len(name) == 0 || len(charger) == 0 || chargerUid <= 0 {
		req.JsonResp(nil, ErrInvalidParam)
		return
	}
	if len(app) == 0 || strings.Index(app, etcd.KeySeprator) >= 0 || app == "grpcservers" {
		req.JsonResp(nil, ErrInvalidParam)
		return
	}

	appObj = db.App{App: app, AppType: appType, Name: name, Charger: charger, ChargerUid: chargerUid}
	er := appObj.Create()
	if er == nil {
		userAuth := db.UserAuth{Uid: chargerUid, Uname: charger, App: app, Permission: int8(db.PermissionType_Manage)}
		er = userService.SaveUserAuth(&userAuth)
	}
	req.JsonResp(nil, er)
}

// @router /test
func (req *AppController) Test() {
	defer req.ServeJSON()
	req.JsonResp("ok", nil)
}

// eg:http://localhost:8081/app/del/app1
// @router /del/:app [post,get]
func (req *AppController) Delete() {
	defer req.ServeJSON()
	app := req.GetString(":app")

	log.Printf("Delete app:%s", app)
	if len(app) > 0 {
		appObj := db.App{App: app}
		_, err := appObj.Delete()

		profile := db.Profile{App: app}
		profile.DeleteByApp()

		//		config := db.Config{App: app}
		//		config.DeleteByInput()
		//
		//		configHistory := db.ConfigHistory{App: app}
		//		configHistory.DeleteByInput()

		configService.DeleteConfigByApp(app)

		auth := db.UserAuth{App: app}
		auth.DeleteByInput()

		req.JsonResp(nil, err)
	} else {
		req.JsonResp(nil, ErrInvalidParam)
	}

}

// eg:http://localhost:8081/app/update?app=&appType=&name=&charger=
// @router /update [post,get]
func (req *AppController) Update() {
	defer req.ServeJSON()
	app := req.GetString("app")
	appType, _ := req.GetInt8("appType", -1)
	name := req.GetString("name")
	charger := req.GetString("charger")

	if len(app) <= 0 || len(name) <= 0 || len(charger) <= 0 || appType < 0 {
		req.JsonResp(nil, ErrInvalidParam)
		return
	}

	log.Printf("Update app:%s, appType:%d, name:%s, charger:%s", app, appType, name, charger)
	appObj := db.App{App: app, AppType: appType, Name: name, Charger: charger}
	_, err := appObj.Update()
	req.JsonResp(nil, err)
}

// eg:http://localhost:8081/app/list
// @router /list [post,get]
func (req *AppController) List() {
	defer req.ServeJSON()

	apps, err := (&db.App{}).FindAll(req.Uid)

	if err != nil || len(apps) <= 0 {
		req.JsonResp(apps, err)
		return
	}

	var list []*db.AppExt

	for _, v := range apps {
		vv := &db.AppExt{}
		vv.App = *v
		list = append(list, vv)
	}

	if req.GetSession("isAdmin") == true {
		for _, v := range list {
			v.Permission = int8(db.PermissionType_Admin)
		}
	} else {
		userAuths, err := userService.FindUserAuthByInput(&db.UserAuth{Uid: req.Uid})
		if err == nil && len(userAuths) > 0 {
			for _, v := range list {
				for _, ua := range userAuths {
					if v.App.App == ua.App {
						v.Permission = ua.Permission
					}
				}
			}
		}
	}

	req.JsonResp(list, err)
}
