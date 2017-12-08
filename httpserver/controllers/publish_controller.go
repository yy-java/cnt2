package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/yy-java/cnt2/httpserver/globals"
	"github.com/yy-java/cnt2/service/codes"
	configService "github.com/yy-java/cnt2/service/config"
	. "github.com/yy-java/cnt2/service/errors"
	publishService "github.com/yy-java/cnt2/service/publish"
	userService "github.com/yy-java/cnt2/service/user"
)

// Operations about App
type PublishController struct {
	BaseController
}

// eg:http://localhost:8081/publish?id=&nodes=1,2,3&pubType=1(1是灰度2是全量)
// @router / [post,get]
func (req *PublishController) Publish() {
	defer req.ServeJSON()
	//捕捉发布过程中异常
	defer func() {
		if err := recover(); err != nil {
			var errStr string
			if newErr, ok := err.(error); ok {
				errStr = newErr.Error()
			}
			req.JsonResp(nil, New(codes.ServerError, "发布失败！"+errStr))
		}
	}()

	id, _ := req.GetInt64("id", 0)
	nodes := strings.TrimSpace(req.GetString("nodes"))
	pubType, _ := req.GetInt("pubType")
	if id <= 0 || pubType <= 0 {
		req.JsonResp(nil, ErrInvalidParam)
	}
	log.Printf("Publish id:%d,nodes: %s---------------------", id, nodes)

	config, err := configService.FindConfigById(id)

	if err != nil || config == nil {
		req.JsonResp(nil, New(codes.NotExists, "not found config!"))
		return
	}
	//未审核不能发
	if config.ApproveType != 2 {
		req.JsonResp(nil, New(codes.PermissionDenied, "not approved!"))
		return
	}

	log.Printf("Publish id:%d, config:%v", id, config)

	//非项目授权不能发
	permission := userService.CheckPermission(req.Uid, config.App)
	if permission <= 0 {
		req.JsonResp(nil, ErrPermissionDenied)
		return
	}
	if pubType == 2 {
		publishService.AllPublishNode(globals.GetEtcdClient(), config.App, config.Profile, config.Key, 0, config.Version, strings.Split(nodes, ","))
	} else {
		publishService.PartitionPublishNode(globals.GetEtcdClient(), config.App, config.Profile, config.Key, 0, config.Version, strings.Split(nodes, ","))
	}
	req.JsonResp(nil, nil)
}

// eg:http://localhost:8081/publish/queryHistory?app=app1&profile=dev&key=key&version=
// @router /queryHistory [post,get]
func (u *PublishController) QueryPublishHistory() {
	defer u.ServeJSON()
	app := u.GetString("app")
	profile := u.GetString("profile")
	key := u.GetString("key")
	version, _ := u.GetInt64("version")
	log.Printf("QueryPublishHistory app:%s , profile:%s , key:%s, version:%s", app, profile, key, version)
	if len(app) == 0 || len(profile) == 0 || len(key) == 0 || version <= 0 {
		u.JsonResp(nil, ErrInvalidParam)
		return
	}
	publishList, err := publishService.QueryPublishHistory(app, profile, key, version)
	u.JsonResp(publishList, err)
}

// @router /queryIp [post,get]
func (u *PublishController) QueryIp() {
	defer u.ServeJSON()
	ip := u.GetString("ip")
	host := u.Ctx.Request.Host //u.Ctx.Request.Host
	var uri string

	if strings.LastIndex(host, "github.com/yy-java") != -1 {
		uri = "http://cmdb.sysop.duowan.com:8088/webservice/server/getServerInfos.do?ip=" + ip
		res, err := getIpRoom(uri, true)
		if err == nil {
			u.JsonResp(res, err)
			return
		}
	}
	uri = "http://int.dpool.sina.com.cn/iplookup/iplookup.php?format=js&ip=" + ip
	res, err := getIpRoom(uri, false)
	if err == nil {
		u.JsonResp(res, nil)
	} else {
		u.JsonResp(nil, err)
	}
}

type YYIpInfo struct {
	Obj []struct {
		RoomName string `json:"roomName"`
	} `json:"object"`
}
type SinaIpInfo struct {
	City string `json:"city"`
}

func getIpRoom(url string, isYY bool) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := string(body)
	log.Printf("url:%s , resp:%s ", url, result)
	if !isYY {
		result = strings.Replace(result, "var remote_ip_info =", "", 1)
		result = strings.Replace(result, ";", "", 1)
	}

	if isYY {
		var yyinfo YYIpInfo
		if err := json.Unmarshal([]byte(result), &yyinfo); err == nil {
			if len(yyinfo.Obj) > 0 {
				return yyinfo.Obj[0].RoomName, nil
			}
		}
	} else {
		var sinaIpInfo SinaIpInfo
		if err2 := json.Unmarshal([]byte(result), &sinaIpInfo); err2 == nil {
			return sinaIpInfo.City, nil
		}
	}
	return "", ErrInvalidParam
}
