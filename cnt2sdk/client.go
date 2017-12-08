package cnt2sdk

/*
Go 版本SDK

Finish:
	1. 启动注册获取NodeId，在etcd中创建临时节点
	2. 加载所有配置，并Watch
	5. 提供配置变更Listener的注册接口

	6. 配置变更下发通知处理逻辑
		6.1 全量发布处理逻辑
			6.1.1 判断当前配置version与下发的version关系，需要更新才更新
		6.2 灰度发布处理逻辑
			6.1.1 判断下发的nodeId是否包含自身，有才更新
		6.3 需要更新，调用查询接口获取配置信息
		6.4 调用Listener通知
		6.5 如果Listener没有返回errors，回调接口通知处理成功

// List:
	3. 加载配置初始化完毕，写本地文件
	4. 如果1 ~ 3 失败，读取本地配置。并延迟自动进行初始化（可选）
	7. 配置每次变更都异步写到本地缓存（可选）

*/
import (
	"context"
	_ "encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/yy-java/cnt2/service/register"

	"google.golang.org/grpc"

	"log"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/yy-java/cnt2/etcd"
	"github.com/yy-java/cnt2/grpc/pb"
	"github.com/yy-java/cnt2/service/models"
	"github.com/yy-java/cnt2/service/utils"

	grpcc "github.com/yy-java/cnt2/grpc/client"
)

type Config struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Version int64  `json:"version"`
}

type XClientConfig struct {
	etcd.EtcdClientConfig
	App     string
	Profile string
}

type XClient struct {
	EtcdClient   *etcd.EtcdClient
	nodeId       string
	listeners    map[string][]ConfigListener
	localCache   map[string]Config
	lock         sync.RWMutex
	listenerLock sync.RWMutex
	config       XClientConfig
	grpcHosts    grpcc.GrpcClient
	serverIp     string
}

type ConfigListener interface {
	Callback(key, newValue, oldValue string, event mvccpb.Event_EventType) error
}

// NewXClient 创建并初始化
func NewXClient(config XClientConfig) (*XClient, error) {
	client := XClient{config: config, localCache: make(map[string]Config), listeners: make(map[string][]ConfigListener)}
	err := client.initEtcdClient()
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// 初始化EtcdClient
func (cli *XClient) initEtcdClient() error {
	etcdConfig := etcd.EtcdClientConfig{Endpoints: cli.config.Endpoints,
		DialTimeout:       cli.config.DialTimeout,
		RequestTimeout:    cli.config.RequestTimeout,
		LeaseTtl:          cli.config.LeaseTtl,
		WatchCallbackFunc: cli,
	}
	etcdClient, err := etcd.NewEtcdClient(etcdConfig)
	if err != nil {
		return err
	}
	cli.EtcdClient = etcdClient

	hostinfo := utils.GetHostInfo("")
	if len(hostinfo.Ips) > 0 {
		for _, ip := range hostinfo.Ips {
			cli.serverIp = ip.Ip
			break
		}
	} else {
		s := utils.GetMyIPInfo()
		if len(s) > 0 {
			cli.serverIp = s[0]
		}
	}
	cli.grpcHosts = grpcc.GrpcClient{EtcdClient: etcdClient, Hostinfo: hostinfo}
	cli.grpcHosts.Init()

	err = cli.register()
	if err != nil {
		cli.Close()
		return err
	}

	err = cli.initConfigs()
	if err != nil {
		log.Fatalln("init Configs error ", err)
	}
	return nil
}

func (cli *XClient) getConn() (pb.ConfigCenterServiceClient, *grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	server := cli.grpcHosts.GetServer()
	log.Printf("Use grpcServer %s", server)
	if server != "" {
		conn, err := grpc.Dial(server, opts...)
		if err == nil {
			return pb.NewConfigCenterServiceClient(conn), conn, nil
		}
	}

	return nil, nil, errors.New("No Grpc Server Found")
}

// 向服务端注册获取nodeId
func (cli *XClient) register() error {
	pbCli, conn, err := cli.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	pid := os.Getpid()
	resp, err := pbCli.RegisterClient(context.Background(), &pb.RegisterRequest{App: cli.config.App, Profile: cli.config.Profile, ServerIp: cli.serverIp, Pid: strconv.Itoa(pid)})
	if err != nil {
		return err
	}
	if resp.Result != 1 {
		return errors.New("error registerClient result= " + strconv.FormatInt(int64(resp.Result), 10))
	}
	if len(resp.NodeId) <= 0 {
		return errors.New("error nodeId returns " + resp.NodeId)
	}
	cli.nodeId = resp.NodeId

	inode, err := strconv.ParseInt(cli.nodeId, 10, 64)
	if err != nil {
		return err
	}
	err = register.ClientOnline(cli.EtcdClient, cli.config.App, cli.config.Profile, cli.serverIp, inode, pid)
	if err != nil {
		return err
	}

	cli.EtcdClient.Watch(etcd.GetKeyPrefix(cli.config.App, cli.config.Profile, etcd.KeyType_Profile), true)
	return nil
}

// 通过接口查询配置并初始化本地缓存
func (cli *XClient) initConfigs() error {
	pbCli, conn, err := cli.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := pbCli.QueryAll(context.Background(), &pb.QueryRequest{App: cli.config.App, Profile: cli.config.Profile})
	if err != nil {
		return err
	}

	for _, msg := range resp.Result {
		c := Config{Key: msg.Key, Value: msg.Value, Version: msg.Version}
		cli.localCache[msg.Key] = c
	}
	return nil
}

// 注册Listeners
func (cli *XClient) AddListener(key string, callback ConfigListener) {
	cli.listenerLock.Lock()
	defer cli.listenerLock.Unlock()

	if callback != nil {
		cli.listeners[key] = append(cli.listeners[key], callback)
	}
}

// 反注册Listener
func (cli *XClient) RemoveListener(key string) {
	cli.listenerLock.Lock()
	defer cli.listenerLock.Unlock()

	delete(cli.listeners, key)
}

func (cli *XClient) getListenersWithLock(key string) ([]ConfigListener, bool) {
	cli.listenerLock.RLock()
	defer cli.listenerLock.RUnlock()

	s, ok := cli.listeners[key]

	t, ok1 := cli.listeners["*"]
	if ok1 {
		for _, e := range t {
			s = append(s, e)
			ok = true
		}
	}

	return s, ok
}

func (cli *XClient) getFromCacheWithLock(key string) (Config, bool) {
	cli.lock.RLock()
	defer cli.lock.RUnlock()

	local, isNew := cli.localCache[key]
	return local, isNew
}

func (cli *XClient) putToCacheWithLock(key string, config Config) {
	cli.lock.Lock()
	defer cli.lock.Unlock()

	cli.localCache[key] = config
}

func (cli *XClient) deleteCacheWithLock(key string) {
	cli.lock.Lock()
	defer cli.lock.Unlock()

	delete(cli.localCache, key)
}

// 实现Etcd的Callback接口
func (cli *XClient) Callback(kv *mvccpb.KeyValue, etype mvccpb.Event_EventType) {
	var val models.ConfigNode
	err := val.FromJson(string(kv.Value))
	if err != nil {
		log.Fatalln("parse ", kv, " to Json models.ConfigNode error")
		return
	}
	key := etcd.ParseConfigKey(string(kv.Key))
	if key.AppName != cli.config.App && key.ProfileName != cli.config.Profile {
		log.Fatalln("not match ", kv, " to Json models.ConfigNode error")
		return
	}
	local, isNew := cli.getFromCacheWithLock(key.KeyName)
	updated, succ := false, true
	switch etype {
	case mvccpb.PUT:
		goUpdate := false
		// 有配置更新
		if len(val.PublishInfo.PublishNodes) == 0 { //全量发布
			goUpdate = (isNew && local.Version < val.Version)
		} else { // 灰度发布
			for _, pnode := range val.PublishInfo.PublishNodes {
				if pnode == cli.nodeId { // 灰度处理
					goUpdate = (isNew && local.Version < val.Version)
					break
				}
			}
		}

		if goUpdate {
			updated = true
			pbCli, conn, err := cli.getConn()
			if err != nil {
				updated = false
				return
			}
			defer conn.Close()
			resp, err := pbCli.QueryKey(context.Background(), &pb.QueryRequest{App: cli.config.App, Profile: cli.config.Profile, Key: key.KeyName, KeyVersion: val.PublishInfo.Version})
			if err != nil {
				updated = false
				succ = false
			}
			c := Config{Value: resp.Value, Version: resp.Version, Key: resp.Key}
			cli.putToCacheWithLock(key.KeyName, c)
			var oldValue string
			if !isNew {
				oldValue = local.Value
			}
			succ = cli.notifyAllListeners(key.KeyName, c.Value, oldValue, etype)
		}
	case mvccpb.DELETE:
		// 有配置删除,没有灰度逻辑
		if !isNew {
			updated = true
			cli.deleteCacheWithLock(key.KeyName)
			succ = cli.notifyAllListeners(key.KeyName, "", local.Value, etype)
		}
	}
	if updated { //回调通知
		pbCli, conn, err := cli.getConn()
		if err != nil {
			updated = false
			return
		}
		defer conn.Close()
		t := pb.ValueChangeResultRequest_FAILED
		if succ {
			t = pb.ValueChangeResultRequest_SUCCESS
		}
		pbCli.ValueChangeResultNotify(context.Background(),
			&pb.ValueChangeResultRequest{NodeId: cli.nodeId,
				App:      cli.config.App,
				Profile:  cli.config.Profile,
				Key:      key.KeyName,
				DeployId: strconv.FormatInt(val.PublishInfo.PublishId, 10),
				Version:  cli.config.Version,
				Result:   t,
				Msg:      "",
			})
	}
}

func (cli *XClient) notifyAllListeners(key, newValue, oldValue string, event mvccpb.Event_EventType) bool {
	callbacks, ok := cli.getListenersWithLock(key)

	if !ok {
		return true
	}
	ret := true
	for _, cb := range callbacks {
		err := cb.Callback(key, newValue, oldValue, event)
		if err != nil {
			ret = false
		}
	}
	return ret
}

func (cli *XClient) GetConfig(key string) (string, bool) {
	if cli.localCache == nil {
		return "", false
	}
	v, ok := cli.localCache[key]
	if ok {
		return v.Value, true
	}
	return "", false
}

func (cli *XClient) GetAllConfig() map[string]string {
	if cli.localCache == nil {
		return nil
	}
	var copy map[string]string = make(map[string]string)
	for k, v := range cli.localCache {
		copy[k] = v.Value
	}
	return copy
}

func (cli *XClient) Close() {
	if cli.EtcdClient != nil {
		cli.EtcdClient.Close()
		cli.EtcdClient = nil
	}
}
