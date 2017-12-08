package controllers

import (
	//	"encoding/json"
	//	"github.com/astaxie/beego"
	"log"
	. "github.com/yy-java/cnt2/db"
	configService "github.com/yy-java/cnt2/service/config"
	. "github.com/yy-java/cnt2/service/errors"

	userService "github.com/yy-java/cnt2/service/user"
)

// Operations about Config
type ConfigController struct {
	BaseController
}

// eg: http://localhost:8081/config/create?app=app2&profile=test&key=key&value=val&validator=validator&description=desc
// @router /create [post,get]
func (c *ConfigController) Create() {
	app := c.GetString("app")
	profile := c.GetString("profile")
	key := c.GetString("key")
	value := c.GetString("value")
	validator := c.GetString("validator")
	description := c.GetString("description")

	log.Printf("Create app:%s, profile:%s, key:%s, value:%s", app, profile, key, value)

	var er error

	if er == nil && (len(app) == 0 || len(profile) == 0 || len(key) == 0) {
		er = ErrInvalidParam
	}
	if er == nil {
		resultConfig, err := configService.CreateConfigByUser(c.Uid, c.Username, app, profile, key, value, validator, description)
		if err != nil {
			er = err
		} else {
			c.JsonResp(resultConfig.Id, nil)
		}
	}
	if er != nil {
		c.JsonResp(nil, er)
	}

	c.ServeJSON()
}

// eg: http://localhost:8081/config/update?app=app2&profile=test&key=key&value=val&validator=validator&description=desc
// @router /update [post,get]
func (c *ConfigController) Update() {
	defer c.ServeJSON()
	app := c.GetString("app")
	profile := c.GetString("profile")
	key := c.GetString("key")
	value := c.GetString("value")
	validator := c.GetString("validator")
	description := c.GetString("description")

	log.Printf("Update app:%s, profile:%s, key:%s, value:%s", app, profile, key, value)

	var er error

	if er == nil {
		uAuth, _ := userService.FindUserAuthByInput(&UserAuth{App: app, Uid: c.Uid})
		if uAuth == nil {
			er = ErrUnauthenticated
		}
	}

	if er == nil && (len(app) == 0 || len(profile) == 0 || len(key) == 0) {
		er = ErrInvalidParam
	}
	if er == nil {
		_, err := configService.UpdateConfigByUser(c.Uid, c.Username, app, profile, key, value, validator, description)
		if err != nil {
			er = err
		} else {
			c.JsonResp(nil, nil)
		}
	}
	if er != nil {
		c.JsonResp(nil, er)
	}
	c.ServeJSON()
}

// eg:http://localhost:8081/config/query?app=app1&profile=dev
// @router /query [post,get]
func (u *ConfigController) Query() {
	defer u.ServeJSON()
	app := u.GetString("app")
	profile := u.GetString("profile")

	log.Printf("Query app:%s , profile:%", app, profile)
	if len(app) > 0 {
		config, err := configService.FindAllConfig(app, profile)
		if err != nil {
			u.JsonResp(nil, err)
		} else {
			u.JsonResp(config, nil)
		}
	} else {
		u.JsonResp(nil, ErrInvalidParam)
	}

}

// eg:http://localhost:8081/config/queryWithKey?app=app1&profile=dev&key=key
// @router /queryWithKey [post,get]
func (u *ConfigController) QueryWithKey() {
	defer u.ServeJSON()
	app := u.GetString("app")
	profile := u.GetString("profile")
	key := u.GetString("key")

	log.Printf("Query app:%s , profile:%s , key:%s", app, profile, key)
	if len(app) > 0 {
		config, err := configService.FindConfig(app, profile, key)
		if err != nil {
			u.JsonResp(nil, err)
		} else {
			u.JsonResp(config, nil)
		}
	} else {
		u.JsonResp(nil, ErrInvalidParam)
	}

}

// eg:http://localhost:8081/config/del?app=app1&profile=dev1&key=key1
// @router /del [post,get]
func (u *ConfigController) Delete() {
	app := u.GetString("app")
	profile := u.GetString("profile")
	key := u.GetString("key")

	log.Printf("Delete app:%s profile:%s key:%s", app, profile, key)
	if app != "" {
		err := configService.DeleteConfig(u.Username, app, profile, key)
		u.JsonResp(nil, err)
	} else {
		u.JsonResp(nil, ErrInvalidParam)
	}
	u.ServeJSON()
}

// eg:http://localhost:8081/config/approve?app=app1&profile=dev1&key=key1&version=
// @router /approve [post,get]
func (u *ConfigController) Approve() {

	defer u.ServeJSON()

	app := u.GetString("app")
	profile := u.GetString("profile")
	key := u.GetString("key")

	version, _ := u.GetInt64("version", 0)

	log.Printf("approve app:%s,profile:%s,key:%s, version:%s", app, profile, key, version)

	if len(app) > 0 && len(profile) > 0 && len(key) > 0 && version > 0 {
		err := configService.ApproveConfigByUser(u.Uid, u.Username, app, profile, key, version)
		u.JsonResp(nil, err)
	} else {
		u.JsonResp(nil, ErrInvalidParam)
	}
}

// eg:http://localhost:8081/config/rollback?app=app1&profile=dev1&key=key1&version=
// @router /rollback [post,get]
func (u *ConfigController) Rollback() {

	defer u.ServeJSON()

	app := u.GetString("app")
	profile := u.GetString("profile")
	key := u.GetString("key")

	version, _ := u.GetInt64("version", 0)

	log.Printf("rollback app:%s,profile:%s,key:%s, version:%s", app, profile, key, version)

	if app != "" {
		_, err := configService.RollbackConfigByUser(u.Uid, u.Username, app, profile, key, version)
		u.JsonResp(nil, err)
	} else {
		u.JsonResp(nil, ErrInvalidParam)
	}

}

// eg:http://localhost:8081/config/profiles/app1
// @router /profiles/:app [post,get]
func (req *ConfigController) ListProfiles() {
	defer req.ServeJSON()
	app := req.GetString(":app")

	if len(app) <= 0 {
		req.JsonResp(nil, ErrInvalidParam)
		return
	}
	profiles, err := configService.FindAppProfiles(app)

	req.JsonResp(profiles, err)
}

// eg:http://localhost:8081/config/queryHistory?app=app1&profile=dev&key=key
// @router /queryHistory [post,get]
func (u *ConfigController) QueryConfigHistory() {
	defer u.ServeJSON()
	app := u.GetString("app")
	profile := u.GetString("profile")
	key := u.GetString("key")

	log.Printf("QueryConfigHistory app:%s , profile:%s , key:%s", app, profile, key)

	if len(app) > 0 {
		configHislist, err := configService.QueryConfigHistory(app, profile, key)
		u.JsonResp(configHislist, err)
	} else {
		u.JsonResp(nil, ErrInvalidParam)
	}
}
