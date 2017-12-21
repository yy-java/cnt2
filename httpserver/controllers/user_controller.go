package controllers

import (
	"log"
	"github.com/yy-java/cnt2/db"
	. "github.com/yy-java/cnt2/service/errors"
	"errors"
)

type UserController struct {
	BaseController
}

// @router /login [post,get]
func (req *UserController) Login() {
	defer req.ServeJSON()

	username := req.GetString("username")
	pwd := req.GetString("pwd")

	log.Printf("Login username:%s, pwd:%d", username,pwd)

	user := db.User{Username: username,Pwd:pwd}
	user.Login()
	if user.Uid <=  0 {
		req.JsonResp(nil, errors.New("Not found"))
		return
	}
	req.SetSession("uid",user.Uid)
	req.SetSession("username",user.Username)

	req.JsonResp(nil, nil)
}


// @router /create [post,get]
func (req *UserController) Create() {
	defer req.ServeJSON()


	username := req.GetString("username")
	pwd := req.GetString("pwd")

	log.Printf("Create username:%s", username)

	user := db.User{Username: username,Pwd:pwd}
	user.Create()
	if user.Uid <=  0 {
		req.JsonResp(nil, errors.New("Create error"))
		return
	}
	req.JsonResp(nil, nil)
}

// @router /del/:uid [post,get]
func (req *UserController) Delete() {
	defer req.ServeJSON()
	uid,_ := req.GetInt64(":uid")

	log.Printf("Delete uid:%s", uid)
	if uid > 0 {
		user := db.User{Uid: uid}
		_, err := user.Delete()

		req.JsonResp(nil, err)
	} else {
		req.JsonResp(nil, ErrServerErr)
	}

}

// @router /list [post,get]
func (req *UserController) List() {
	defer req.ServeJSON()

	users, err := (&db.User{}).FindAll()

	if err != nil || len(users) <= 0 {
		req.JsonResp(nil, err)
		return
	}

	req.JsonResp(users, err)
}
