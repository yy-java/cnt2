package gosdk

import (
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
 * 流程：
 * 1.到ETCD拿GrpcServer的信息
 * 2.到ConfigCenter注册，拿到nodeId
 * 3.到etcd注册临时节点(/${appName}/nodes/${nodeId})
 * 4.到ConfigCenter查询当前profile下所有配置,写入本地副本
 * 5.订阅节点变化(/${appName}/profiles/${profile})，在服务端配置发生变化时更新本地副本
 * 6.配置发生变化是通知业务
 * 7.如果注册失败则去读本地文件
 *
 */
const (
	defaultDialTimeout = 5 * time.Second
	Node_TTL           = int64(30)
)

var (
	hostInfo           HostInfo
	configCenterClient ConfigCenterServiceClient
	localStore         *LocalStore
	nodeId             string //当前进程注册得到的节点
	clientConfig       *ClientConfig
	etcdClient         *clientv3.Client
	listeners          map[string][]ConfigListener
	listenerLock       sync.RWMutex
	grpcClientConn     *grpc.ClientConn
)

func init() {
	hostInfo = InitHostInfo()
	log.Printf("hostInfo %v", hostInfo.Ips)
}
func Start(ccfg *ClientConfig) (*Cnt2Service, error) {
	clientConfig = ccfg
	listeners = make(map[string][]ConfigListener)
	if clientConfig.DialTimeout <= 0 {
		clientConfig.DialTimeout = defaultDialTimeout
	}
	cfg := clientv3.Config{
		Endpoints:   clientConfig.Endpoints,
		DialTimeout: clientConfig.DialTimeout,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	//初始化本地存储
	localStore, err = InitLocalStore(clientConfig.App, clientConfig.LocalFilePath)
	if err != nil {
		return nil, err
	}
	configCenterClient, err = InitGrpcServer(etcdClient)
	if err != nil {
		return nil, err
	}
	//监听grpcServer
	go WatchGrpcServer(etcdClient)
	//register
	if err = registerClient(etcdClient); err != nil {
		return nil, err
	}

	//启动后全量更新本地数据库
	if queryResp, err := configCenterClient.QueryAll(context.TODO(), &QueryRequest{App: clientConfig.App, Profile: clientConfig.Profile}); err == nil {
		for _, res := range queryResp.GetResult() {
			if resJson, err := json.Marshal(res); err == nil {
				log.Printf("queryResp %s", string(resJson))
				localStore.Put(clientConfig.App, clientConfig.Profile, res.GetKey(), string(resJson))
			}
		}
	} else {
		log.Printf("queryAllResp:%s, err:%v", queryResp, err)
	}

	// 监听配置信息
	go WatchCnt2(etcdClient, configCenterClient, clientConfig.App, clientConfig.Profile)

	if clientConfig.EnableCommon {
		go WatchCnt2Common(etcdClient, configCenterClient, clientConfig.App)
	}

	return &Cnt2Service{clientConfig}, nil
}
func registerClient(etcdClient *clientv3.Client) error {
	//去grpc服务获得节点ID
	pid := os.Getpid()
	registerReq := &RegisterRequest{App: clientConfig.App, Profile: clientConfig.Profile, ServerIp: getIp(), Pid: strconv.Itoa(pid)}
	registerResp, err := configCenterClient.RegisterClient(context.TODO(), registerReq)
	log.Printf("regiesterReq:%s, registerResp: %s ", registerReq, registerResp)
	if err != nil {
		return err
	}
	nodeId = registerResp.GetNodeId()
	//在etcd创建临时节点
	if leaseResp, err := etcdClient.Grant(context.TODO(), Node_TTL); err == nil {
		log.Printf("leaseResp:%s", leaseResp.ID)
		nodeInfo := &Node{NodeId: nodeId, App: clientConfig.App, Profile: clientConfig.Profile, Pid: pid, Sip: registerReq.ServerIp, RegisterTime: time.Now().Unix() * 1000}
		nodeInfoJson, _ := json.Marshal(nodeInfo)
		etcdClient.Put(context.TODO(), GenRegisterNodeKey(clientConfig.App, clientConfig.Profile, nodeId), string(nodeInfoJson), clientv3.WithLease(leaseResp.ID))
		_, err := etcdClient.KeepAlive(context.Background(), leaseResp.ID)
		if err != nil {
			return err
		}

	} else {
		log.Printf("leaseResp:%v", err)
	}
	return nil
}

func InitGrpcServer(client *clientv3.Client) (ConfigCenterServiceClient, error) {

	grpcServerInfos, err := ParseGrpcServerInfoFromEtcd(client)
	if err != nil {
		return nil, err
	}
	bestAddress, _ := ChooseBestAddress(grpcServerInfos)

	log.Printf("init...the best address: %s", bestAddress)

	r := &GrpcServerResolver{grpcServerInfos}

	rb := grpc.RoundRobin(r)

	grpcClientConn, err := grpc.Dial(strings.Join(bestAddress, ";"), grpc.WithTimeout(time.Second*5), grpc.WithInsecure(), grpc.WithBalancer(rb))

	if err != nil {
		return nil, err
	}
	log.Printf("connected grpc server: %v", grpcClientConn)

	configServiceClient := NewConfigCenterServiceClient(grpcClientConn)

	return configServiceClient, nil
}

func getIp() string {
	if len(hostInfo.Ips) > 0 {
		for _, ip := range hostInfo.Ips {
			return ip.Ip
		}
	} else {
		s := GetMyIPInfo()
		if len(s) > 0 {
			return s[0]
		}
	}
	return ""
}

func (Cnt2Service) RegisterListener(listener ConfigListener, keys ...string) {
	listenerLock.Lock()
	defer listenerLock.Unlock()
	if listener == nil {
		return
	}
	for _, key := range keys {
		list := listeners[key]
		var hasExist bool
		for i := range list {
			if list[i] == listener {
				hasExist = true
				break
			}
		}
		if !hasExist {
			list = append(list, listener)
			listeners[key] = list
		}
	}
}

func (Cnt2Service) UnregisterListener(listener ConfigListener, keys ...string) {
	listenerLock.Lock()
	defer listenerLock.Unlock()
	for _, key := range keys {
		if listener == nil {
			delete(listeners, key)
		} else {
			list := listeners[key]
			for i := range list {
				if list[i] == listener {
					list = append(list[:i], list[i+1:]...)
					break
				}
			}
			listeners[key] = list
		}
	}
}

func getListenersWithLock(key string) ([]ConfigListener, bool) {
	listenerLock.RLock()
	defer listenerLock.RUnlock()

	s, ok := listeners[key]

	t, ok1 := listeners["*"]
	if ok1 {
		for _, e := range t {
			s = append(s, e)
			ok = true
		}
	}
	return s, ok
}
func (Cnt2Service) Close() {
	if etcdClient != nil {
		etcdClient.Close()
		etcdClient = nil
	}
	if localStore != nil {
		localStore.Db.Close()
		localStore = nil
	}
	if grpcClientConn != nil {
		grpcClientConn.Close()
		grpcClientConn = nil
	}
	close(grpcServerChan.ch)
}
func (c *Cnt2Service) GetConfig(key string) (string, bool) {
	if localStore == nil {
		localStore2, err := InitLocalStore(c.App, c.LocalFilePath)
		if err != nil {
			return "", false
		}
		localStore = localStore2
	}
	return localStore.Get(c.App, c.Profile, key), true
}
