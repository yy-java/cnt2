package controllers

import (
	"log"
	"github.com/yy-java/cnt2/db"
	configService "github.com/yy-java/cnt2/service/config"
	. "github.com/yy-java/cnt2/service/errors"
)

type ProfileController struct {
	BaseController
}

// eg: http://localhost:8081/profile/create?app=app10&profile=prof&name=test-name2
// @router /create [post,get]
func (c *ProfileController) Create() {
	defer c.ServeJSON()

	app := c.GetString("app")
	name := c.GetString("name")
	profile := c.GetString("profile")
	srcProfile := c.GetString("srcProfile")

	log.Printf("Create app:%s, profile:%s, name:%s ,srcProfile:%s", app, profile, name, srcProfile)

	if len(app) == 0 || len(profile) == 0 || len(name) == 0 {
		c.JsonResp(nil, ErrInvalidParam)
		return
	}

	appObj := db.Profile{App: app, Profile: profile}
	appObj.Read()

	if len(appObj.Name) > 0 {
		c.JsonResp(nil, ErrExists)
		return
	}
	appObj.Name = name
	err := appObj.Create()
	if err == nil && len(srcProfile) > 0 {
		err = configService.CopyConfig(app, srcProfile, profile, c.Username, c.Uid)
	}
	c.JsonResp(nil, err)
}

// eg:http://localhost:8081/profile/del?app=app1&profile=
// @router /del [post,get]
func (u *ProfileController) Delete() {
	var err error
	defer u.ServeJSON()

	app := u.GetString("app")
	profile := u.GetString("profile")

	log.Printf("Delete app:%s,profile:%s", app, profile)

	if len(app) > 0 && len(profile) > 0 {
		appObj := db.Profile{App: app, Profile: profile}
		_, err = appObj.DeleteByAppAndProfile()

		//		config := db.Config{App: app, Profile: profile}
		//		config.DeleteByInput()
		//
		//		configHistory := db.ConfigHistory{App: app, Profile: profile}
		//		configHistory.DeleteByInput()
		configService.DeleteConfigByAppProfile(app, profile)

	} else {
		err = ErrInvalidParam
	}
	u.JsonResp(nil, err)
}

// eg:http://localhost:8081/profile/update?app=&profile=&name=
// @router /update [post,get]
func (u *ProfileController) Update() {
	var err error
	defer u.ServeJSON()

	app := u.GetString("app")
	profile := u.GetString("profile")
	name := u.GetString("name")

	log.Printf("Update app:%s, profile:%s, name:%s", app, profile, name)
	if len(app) <= 0 || len(profile) <= 0 || len(name) <= 0 {
		err = ErrInvalidParam
	}
	appObj := db.Profile{App: app, Profile: profile, Name: name}
	_, err = appObj.Update()
	u.JsonResp(nil, err)
}

// eg:http://localhost:8081/profile/list?app=app1
// @router /list [post,get]
func (u *ProfileController) List() {
	defer u.ServeJSON()
	app := u.GetString("app")
	profiles, err := (&db.Profile{}).FindByApp(app)
	u.JsonResp(profiles, err)
}
