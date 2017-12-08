package register

import (
	"errors"
	"strconv"
	"time"

	"encoding/json"

	"github.com/yy-java/cnt2/db"
	"github.com/yy-java/cnt2/etcd"
	"github.com/yy-java/cnt2/service/models"

	"log"
)

var (
	AppNotExistErr = errors.New("app not exits")
	ParamErr       = errors.New("param error")
)

// RegisterClient for grpcServer to save and client register info
func RegisterClient(app, profile, serverIp string, pid int) (int64, error) {
	var node db.Node
	if len(app) == 0 || len(profile) == 0 {
		return 0, ParamErr
	}
	node.App = app
	node.Pid = pid
	node.Profile = profile
	node.RegisterTime = time.Now()
	node.Sip = serverIp
	err := node.Create()
	if err != nil {
		return 0, err
	}
	return node.Id, err
}

// ClientOnline register info into etcd ,hold by lease
func ClientOnline(etcdClient *etcd.EtcdClient, app, profile, serverIp string, nodeId int64, pid int) error {
	snodeId := strconv.FormatInt(nodeId, 10)
	nodeKey := etcd.NewNodeKey(app, profile, snodeId)
	node := models.RegisterNode{NodeId: snodeId, App: app, Profile: profile, RegisterTime: time.Now().Unix() * 1000, ServerIP: serverIp, Sip: serverIp, Pid: pid}
	jsonNode, _ := json.Marshal(node)
	err := etcdClient.Lease(nodeKey.FullPath, string(jsonNode))
	return err
}

// ListAllOnlineNodes for console server to list all online nodes to do Publishing
func ListAllOnlineNodes(etcdClient *etcd.EtcdClient, app, profile string) ([]models.RegisterNode, error) {
	var ky string
	if profile == "common" {
		ky = etcd.GetCommonNodesKeyPrefix(app)
	} else {
		ky = etcd.GetKeyPrefix(app, profile, etcd.KeyType_Deploy)
	}
	kvs, err := etcdClient.ListKeys(ky)
	log.Printf("ListAllOnlineNodes key:%s, values:%s", ky, kvs)

	var result []models.RegisterNode
	if err != nil {
		return result, err
	}
	//	dict := make(map[string]models.RegisterNode)

	for _, kv := range kvs {
		k, v := string(kv.Key), kv.Value
		log.Printf("foreach key:%s, values:%s", k, v)
		var vNode models.RegisterNode
		if err := json.Unmarshal(v, &vNode); err == nil {
			//dict[k] = vNode
			result = append(result, vNode)
		}
	}
	return result, nil
}

func ListPublishedNode(etcdClient *etcd.EtcdClient, app, profile, key string) ([]string, error) {
	configKey := etcd.NewConfigKey(app, profile, key)
	getResp, err := etcdClient.Get(configKey.FullPath)

	if err != nil {
		return nil, err
	}

	for _, kv := range getResp.Kvs {
		k, v := string(kv.Key), kv.Value
		log.Printf("foreach key:%s, values:%s", k, v)
		var vNode models.ConfigNode
		if err := json.Unmarshal(v, &vNode); err == nil {
			//dict[k] = vNode
			return vNode.PublishInfo.PublishNodes, nil
		}
	}
	return nil, nil
}
