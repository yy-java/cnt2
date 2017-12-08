package controllers

import (
	"log"
	"github.com/yy-java/cnt2/db"
	. "github.com/yy-java/cnt2/service/errors"
	"github.com/yy-java/cnt2/service/user"
)

type UserAuthController struct {
	BaseController
}

// eg: http://localhost:8081/userauth/create?app=app10&permission=1&uid=&username=
// @router /create [post,get]
func (req *UserAuthController) Create() {

	defer req.ServeJSON()
	app := req.GetString("app")
	uid, _ := req.GetInt64("uid", 0)
	username := req.GetString("username")
	permission, _ := req.GetInt8("permission", 0)

	log.Printf("Create uid:%d, username: %s, app:%s, permission:%d", uid, username, app, permission)

	var er error

	if !user.CheckManagePermission(req.Uid, app) {
		er = ErrPermissionDenied
	}
	log.Printf("er = %v", er)
	if er == nil && (len(app) == 0 || permission <= 0 || len(username) <= 0 || uid <= 0) {
		er = ErrInvalidParam
	}
	if er == nil {
		userAuth := db.UserAuth{Uid: uid, Uname: username, App: app, Permission: permission}
		er = user.SaveUserAuth(&userAuth)
	}
	req.JsonResp(nil, er)
}

// eg: http://localhost:8081/userauth/del?id=
// @router /del [post,get]
func (u *UserAuthController) Delete() {
	id, _ := u.GetInt64("id", 0)

	log.Printf("Delete id:%d", id)

	err := user.RemoveUserAuthById(id)
	u.JsonResp(nil, err)
	u.ServeJSON()
}

// eg:http://localhost:8081/userauth/update_permission?id=&permission=&app=
// @router /update_permission [post,get]
func (u *UserAuthController) UpdatePermission() {
	defer u.ServeJSON()
	id, _ := u.GetInt64("id", 0)
	permission, _ := u.GetInt8("permission", 0)
	app := u.GetString("app")

	log.Printf("UpdatePermission id:%d, permission:%s", id, permission)

	if !user.CheckManagePermission(u.Uid, app) {
		u.JsonResp(nil, ErrPermissionDenied)
		return
	}
	if id > 0 {
		userAuth, err := user.FindUserAuthById(id)
		if err == nil {
			userAuth.Permission = permission
			err = user.SaveUserAuth(userAuth)
		}
		u.JsonResp(nil, err)
	} else {
		u.JsonResp(nil, ErrInvalidParam)
	}
}

// eg:http://localhost:8081/userauth/list/app1
// @router /list/:app [post,get]
func (u *UserAuthController) List() {
	defer u.ServeJSON()
	app := u.GetString(":app")

	log.Printf("List app:%s", app)
	userAuths, err := user.FindUserAuthByInput(&db.UserAuth{App: app})
	u.JsonResp(userAuths, err)
}

// eg:http://localhost:8081/userauth/querycuruser?app=app1  99是最高权限 1是开发 9是管理
// @router /querycuruser [post,get]
func (u *UserAuthController) QueryCurUserWithApp() {
	defer u.ServeJSON()

	if u.GetSession("isAdmin") == true {
		u.JsonResp(99, nil)
		return
	}

	app := u.GetString("app")

	log.Printf("QueryCurrentUserWithApp app:%s", app)
	userAuths, err := user.FindUserAuthByInput(&db.UserAuth{App: app, Uid: u.Uid})
	if err != nil {
		u.JsonResp(-1, nil)
		return
	}
	if len(userAuths) <= 0 {
		u.JsonResp(-1, nil)
	} else {
		u.JsonResp(userAuths[0].Permission, nil)
	}
}

// eg:http://localhost:8081/userauth/logout
// @router /logout [post,get]
func (u *UserAuthController) Logout() {
	defer u.ServeJSON()
	u.DestroySession()
	u.JsonResp(nil, nil)
}
