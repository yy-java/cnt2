package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/yy-java/cnt2/etcd"
)

type defaultWatchCallback struct {
}

func (d *defaultWatchCallback) Callback(kv *mvccpb.KeyValue, event mvccpb.Event_EventType) {
	fmt.Printf("watcherLogger %v,%v\n", kv, event)
}

var (
	testEndpoints = []string{"172.27.141.11:2379", "172.27.141.12:2379", "172.27.141.13:2379"}
	defaultConfig = etcd.EtcdClientConfig{Endpoints: testEndpoints, DialTimeout: 3 * time.Second, RequestTimeout: 3 * time.Second, WatchCallbackFunc: &defaultWatchCallback{}}
	etcdClient, _ = etcd.NewEtcdClient(defaultConfig)
	grpcClient    = GrpcClient{EtcdClient: etcdClient}
)

func Test_GetServer(t *testing.T) {
	err := grpcClient.Init()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Return Server ", grpcClient.GetServer())
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Return Server ", grpcClient.GetServer())
	time.Sleep(10 * time.Second)
	fmt.Println("Return Server ", grpcClient.GetServer())
	time.Sleep(30 * time.Second)
	fmt.Println("Return Server ", grpcClient.GetServer())
	time.Sleep(10 * time.Millisecond)
	fmt.Println("All Servers ", grpcClient.getServersWithLock())
}
