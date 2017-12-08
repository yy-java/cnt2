package client

import (
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"math/rand"
	"strings"
	"sync"
	"github.com/yy-java/cnt2/etcd"
	"github.com/yy-java/cnt2/service/models"
	"github.com/yy-java/cnt2/service/utils"
)

type GrpcClient struct {
	EtcdClient *etcd.EtcdClient
	Hostinfo   utils.HostInfo
	servers    sync.Map
}

func (gc *GrpcClient) Init() error {
	kv, err := gc.EtcdClient.ListKeys(etcd.KeyGrpcServer)
	if err != nil {
		return err
	}
	gc.EtcdClient.Watch2(etcd.KeyGrpcServer, true, gc)
	for _, k := range kv {
		s := parseKV(k)
		if s != nil {
			gc.servers.LoadOrStore(string(k.Key), *s)
		}
	}
	return nil
}

func parseKV(kv *mvccpb.KeyValue) *models.GrpcServerRegisterInfo {
	var s models.GrpcServerRegisterInfo
	err := s.FromJson(string(kv.Value))
	if err != nil {
		return nil
	}
	return &s
}

func (gc *GrpcClient) getServerGroupWithLock(groupId int) []models.GrpcServerRegisterInfo {
	grps := make([]models.GrpcServerRegisterInfo, 0)
	gc.servers.Range(func(key, value interface{}) bool {
		var v = value.(models.GrpcServerRegisterInfo)
		if v.GroupId == groupId {
			grps = append(grps, v)
		}
		return true
	})

	return grps
}

func (gc *GrpcClient) getServersWithLock() []models.GrpcServerRegisterInfo {
	grps := make([]models.GrpcServerRegisterInfo, 0)
	gc.servers.Range(func(key, value interface{}) bool {
		var v = value.(models.GrpcServerRegisterInfo)

		grps = append(grps, v)

		return true
	})
	return grps
}

func (gc *GrpcClient) addServerWithLock(key string, info models.GrpcServerRegisterInfo) {
	gc.servers.Store(key, info)
}

func (gc *GrpcClient) deleteServerWithLock(key string) {
	gc.servers.Delete(key)
}

func (gc *GrpcClient) Callback(kv *mvccpb.KeyValue, etype mvccpb.Event_EventType) {
	key := string(kv.Key)

	if strings.Index(key, etcd.KeyGrpcServer) != 0 {
		return
	}
	switch etype {
	case mvccpb.PUT:
		s := parseKV(kv)
		if s != nil {
			gc.addServerWithLock(key, *s)
		}
	case mvccpb.DELETE:
		gc.deleteServerWithLock(key)
	}
}

func (gc *GrpcClient) GetServer() string {
	ret := ""
	grpId := gc.Hostinfo.GroupId
	grps := gc.getServerGroupWithLock(grpId)

	var target []models.GrpcServerRegisterInfo
	if len(grps) > 0 {
		target = grps
	} else {
		target = gc.getServersWithLock()
	}

	//random select 1
	size := len(target)
	if size == 0 {
		return ret
	}
	var got *models.GrpcServerRegisterInfo
	if size == 1 {
		got = &target[0]
	} else {
		got = &target[rand.Intn(size)]
	}

	var sip string
	if len(gc.Hostinfo.Ips) > 0 {
		// 找有相同类型的ip
		for _, s := range gc.Hostinfo.Ips {
			ip, ok := got.ServerIP[s.Type.Val]
			if ok {
				sip = ip
				break
			}
		}
		// 没找到 ,再遍历一次
		if len(sip) == 0 {
			for _, s := range got.ServerIP {
				sip = s
				break
			}
		}
	} else {
		for _, ip := range got.ServerIP {
			sip = ip
			break
		}
	}

	ret = fmt.Sprintf("%s:%d", sip, got.Port)

	return ret
}
