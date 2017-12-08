package publish

import (
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/yy-java/cnt2/db"
	"github.com/yy-java/cnt2/etcd"
	"github.com/yy-java/cnt2/service/config"
	"github.com/yy-java/cnt2/service/models"
)

func deleteStr(s []string, item string) []string {
	var r []string
	for _, str := range s {
		if str != item {
			r = append(r, str)
		}
	}
	return r
}
func containStr(s []string, e string) bool {
	for _, str := range s {
		if str == e {
			return true
		}
	}
	return false
}

/**
找到profile下的节点
*/
func ListAllNode(etcdClient *etcd.EtcdClient, app, profile string) []models.RegisterNode {
	nodes, err := etcdClient.ListKeys(etcd.GetKeyPrefix(app, profile, etcd.KeyType_Deploy)) // key:app名称/nodes/profile名称/
	if err != nil {
		log.Printf("ListAllNode error,app:%v,profile:%v ,error=%v", app, profile, err)
	}
	var retValue []models.RegisterNode
	for index, element := range nodes {
		var registerNode models.RegisterNode
		err := registerNode.FromJson(string(element.Value))
		if err != nil {
			continue
		}
		retValue = append(retValue, registerNode)
		log.Printf("node,%d,%v", index, registerNode)
	}
	return retValue
}
func UpdateConfigAfterAllPublished(app, profile, key string, version int64) {
	configHistory := config.FindConfigHistoryByKeyAndVersion(app, profile, key, version)
	log.Printf("update config,app:%s,profile:%s,key:%s,version:%d,val:%s", app, profile, key, version, configHistory.CurValue)
	newConfig := db.Config{App: app, Profile: profile, Key: key}
	newConfig.UpdateConfigPublishedValue(configHistory.CurVersion, configHistory.CurValue)
}
func getConfigValue(etcdClient *etcd.EtcdClient, app string, profile string, key string, configNode *models.ConfigNode) {
	configPath := etcd.GetKeyPrefix(app, profile, etcd.KeyType_Profile)
	keyPath := configPath + key
	val, _ := etcdClient.Get(keyPath)
	err := configNode.FromJson(string(val.Kvs[0].Value))
	if err != nil {
		log.Printf("Parse configNode error %v", err)
	}
}
func NodeUpdateCallback(etcdClient *etcd.EtcdClient, app, profile, key string, publishId int64, nodeId string) {
	var configNode models.ConfigNode
	getConfigValue(etcdClient, app, profile, key, &configNode)
	nodeIdList := configNode.PublishInfo.PublishNodes
	log.Printf("app:%v,profile:%v,key:%v,nodeIdList:%v", app, profile, key, configNode.PublishInfo.PublishNodes)
	isContained := containStr(nodeIdList, nodeId)
	if isContained {
		update := deleteStr(nodeIdList, nodeId)
		configNode.PublishInfo.PublishNodes = update
		PartitionPublishNode(etcdClient, app, profile, key, publishId, configNode.Version, configNode.PublishInfo.PublishNodes)
		if len(configNode.PublishInfo.PublishNodes) == 0 { //全量后更新
			UpdateConfigAfterAllPublished(app, profile, key, configNode.Version)
		}
	}
}

/**
灰度发布
*/
func PartitionPublishNode(etcdClient *etcd.EtcdClient, app, profile, key string, publishId int64, version int64, nodeIdList []string) {
	configKey := etcd.NewConfigKey(app, profile, key)
	configPublishInfo := models.ConfigPublishInfo{PublishId: publishId, Key: key, Version: version, PublishNodes: nodeIdList}
	configValue := models.ConfigNode{App: app, Profile: profile, Key: key, Version: version, PublishInfo: configPublishInfo}
	configJson, _ := configValue.ToJson()

	log.Printf("PartitionPublishNode update key:%s,value:%s", configKey.FullPath, configJson)

	insertNodePublish(nodeIdList, app, profile, key, version, publishId, int8(db.PublishType_Gray))
	etcdClient.Put(configKey.FullPath, configJson)
}

/**
全量发布
*/
func AllPublishNode(etcdClient *etcd.EtcdClient, app, profile, key string, publishId int64, version int64, nodeIdList []string) {
	configKey := etcd.NewConfigKey(app, profile, key)
	var emptyNodeIdList []string
	configPublishInfo := models.ConfigPublishInfo{PublishId: publishId, Key: key, Version: version, PublishNodes: emptyNodeIdList}
	configValue := models.ConfigNode{App: app, Profile: profile, Key: key, Version: version, PublishInfo: configPublishInfo}
	configJson, _ := configValue.ToJson()
	log.Printf("AllPublishNode update key:%s,value:%s", configKey.FullPath, configJson)
	UpdateConfigAfterAllPublished(app, profile, key, version)

	/*var node db.Node
	configs, err := node.FindAllNode(app, profile)
	if err != nil {
		panic(err)
	}
	var nodeIdList []string
	for _, config := range configs {
		nodeIdList = append(nodeIdList, strconv.FormatInt(config.Id, 10))
	}*/
	insertNodePublish(nodeIdList, app, profile, key, version, publishId, int8(db.PublishType_Full))
	etcdClient.Put(configKey.FullPath, configJson)
}
func insertNodePublish(nodeIdList []string, app string, profile string, key string, version int64, publishId int64, publishType int8) {
	for _, nodeId := range nodeIdList {
		if strings.TrimSpace(nodeId) == "" {
			continue
		}
		i, err := strconv.ParseInt(nodeId, 10, 64)
		if err != nil {
			panic(err)
		}
		nodePublish := db.NodePublish{NodeId: i, App: app, Profile: profile, Key: key, Version: version,
			PublishTime: time.Now(), PublishResult: int8(db.PublishResult_Create), PublishType: publishType}
		nodePublish.Create()
		id := nodePublish.Id
		log.Println(id)
		log.Printf("insert to db,app:%v,profile:%v,key:%v,publishId:%v,version:%v,nodeId:%v,insertId:%v", app, profile, key, publishId, version, i, id)
	}
}
func QueryPublishHistory(app string, profile string, key string, version int64) ([]*db.NodePublishExt, error) {
	nodePublish := db.NodePublish{}
	result, err := nodePublish.QueryPublishHistory(app, profile, key, version)
	return result, err
}
