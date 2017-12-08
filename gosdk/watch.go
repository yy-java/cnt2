package gosdk

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"log"
)

//监听Cnt2配置
func WatchCnt2(client *clientv3.Client, configCenterClient ConfigCenterServiceClient, app, profile string) {
	watchChan := client.Watch(context.Background(), GenWatchCnt2Key(app, profile), clientv3.WithPrefix())
	for notifyResp := range watchChan {
		deal(notifyResp, configCenterClient, app, profile)
	}
}

//监听Cnt2公共配置
func WatchCnt2Common(client *clientv3.Client, configCenterClient ConfigCenterServiceClient, app string) {
	watchChan := client.Watch(context.Background(), GenWatchCnt2CommonKey(app), clientv3.WithPrefix())
	for notifyResp := range watchChan {
		deal(notifyResp, configCenterClient, app, COMMON)
	}
}

func notify(config *Config, evType mvccpb.Event_EventType) bool {
	listenerList, ok := getListenersWithLock(config.Key)
	if !ok {
		return true
	}
	if clientConfig.EnableCommon && config.Profile == COMMON && localStore.Has(config.App, clientConfig.Profile, config.Key) {
		return true
	}
	ret := true
	for _, listener := range listenerList {
		var err error
		switch evType {
		case mvccpb.PUT:
			err = listener.HandlePutEvent(config)
		case mvccpb.DELETE:
			err = listener.HandleDeleteEvent(config)
		default:
		}
		if err != nil {
			ret = false
		}
	}
	return ret
}
func deal(watchResp clientv3.WatchResponse, configCenterClient ConfigCenterServiceClient, app, profile string) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("deal ignore; panic recover! p: %v", p)
		}
	}()
	for _, ev := range watchResp.Events {
		log.Printf("watchCnt2 %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

		var needNotify, needFeedBack, feedBackResult bool
		var recvConfig *Config
		cnt2App, cnt2Profile, cnt2Key, _ := ExtractCntInfo(string(ev.Kv.Key))
		switch ev.Type {
		case mvccpb.PUT:
			recvConfig = (&Config{}).FromJson(string(ev.Kv.Value))
			if recvConfig == nil {
				continue
			}
			//判断是否灰度发布
			if len(recvConfig.PublishInfo.PublishNodes) > 0 && !recvConfig.PublishInfo.PublishNodes.Contains(nodeId) {
				continue
			}
			queryResp, err := configCenterClient.QueryKey(context.TODO(), &QueryRequest{App: app, Profile: profile, Key: recvConfig.Key, KeyVersion: recvConfig.Version})
			if err != nil {
				needFeedBack = true
				continue
			}
			needNotify = true

			recvConfig.Value = queryResp.GetValue()

			if ev.IsCreate() {
				localStore.Put(recvConfig.App, recvConfig.Profile, recvConfig.Key, recvConfig.ToJson())
			} else if ev.IsModify() {
				localConfig, err := localStore.GetConfig(recvConfig.App, recvConfig.Profile, recvConfig.Key)
				if err == nil && recvConfig.Version > localConfig.Version {
					localStore.Put(recvConfig.App, recvConfig.Profile, recvConfig.Key, recvConfig.ToJson())
				}
			}
		case mvccpb.DELETE:
			needNotify = true
			localStore.Del(cnt2App, cnt2Profile, cnt2Key)
		default:
		}
		if needNotify {
			//是PUT事件才需要回调通知
			if ev.Type == mvccpb.PUT {
				needFeedBack = true
			} else {
				recvConfig = &Config{App: cnt2App, Profile: cnt2Profile, Key: cnt2Key}
			}
			feedBackResult = notify(recvConfig, ev.Type)

		}
		if needFeedBack {
			vr := ValueChangeResultRequest_FAILED
			if feedBackResult {
				vr = ValueChangeResultRequest_SUCCESS
			}
			configCenterClient.ValueChangeResultNotify(context.TODO(),
				&ValueChangeResultRequest{
					NodeId:  nodeId,
					App:     recvConfig.App,
					Profile: recvConfig.Profile,
					Key:     recvConfig.Key,
					Version: recvConfig.Version,
					Result:  vr})
		}
	}
}

//监听grpcServer
func WatchGrpcServer(client *clientv3.Client) {
	watchChan := client.Watch(context.Background(), GenWatchGrpcServerKey(), clientv3.WithPrefix())
	for notifyResp := range watchChan {
		for _, ev := range notifyResp.Events {
			log.Printf("watchGrpcServer %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
		grpcServerInfos, err := ParseGrpcServerInfoFromEtcd(client)
		if err == nil {
			grpcServerChan.ch <- grpcServerInfos
		}
	}
}
