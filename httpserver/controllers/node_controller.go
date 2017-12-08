package controllers

import (
	"log"
	"github.com/yy-java/cnt2/httpserver/globals"
	"github.com/yy-java/cnt2/service/errors"
	nodeService "github.com/yy-java/cnt2/service/node"
	"github.com/yy-java/cnt2/service/register"
)

type NodeController struct {
	BaseController
}

// eg:http://localhost:8081/node/list?app=app1&profile=dev
// @router /list [post,get]
func (u *NodeController) List() {
	defer u.ServeJSON()
	app := u.GetString("app")
	profile := u.GetString("profile")
	key := u.GetString("key")
	version, _ := u.GetInt64("version")

	log.Printf("listNodes app:%s , profile:%s, key:%s, version：%s ", app, profile, key, version)

	if len(app) <= 0 || len(profile) <= 0 {
		u.JsonResp(nil, errors.ErrInvalidParam)
		return
	}
	result := make(map[string]interface{})
	//	nodes, err := nodeService.FindByAppAndProfile(u.Uid, app, profile)
	nodes, err := register.ListAllOnlineNodes(globals.GetEtcdClient(), app, profile)

	if err == nil {
		if len(key) > 0 {
			publishedNodes, err2 := nodeService.FindPublishedNode(app, profile, key, version)
			//			publishedNodes, err2 := register.ListPublishedNode(globals.GetEtcdClient(), app, profile, key)
			if err2 == nil {
				result["nodes"] = nodes
				result["publishedNodes"] = publishedNodes
			} else {
				err = err2
			}
		} else {
			result["nodes"] = nodes
		}

	}
	u.JsonResp(result, err)
}
