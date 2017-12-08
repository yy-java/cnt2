package globals

import (
	"sync"

	"github.com/astaxie/beego"
	"github.com/yy-java/cnt2/etcd"
)

var (
	initLock   sync.RWMutex
	etcdClient *etcd.EtcdClient
)

func GetEtcdClient() *etcd.EtcdClient {
	if etcdClient == nil {
		initLock.Lock()
		defer initLock.Unlock()

		if etcdClient == nil {
			endpoints := beego.AppConfig.Strings("etcd::endpoints")
			etcdCliConfig := etcd.EtcdClientConfig{Endpoints: endpoints}
			etcdClient, _ = etcd.NewEtcdClient(etcdCliConfig)
		}
	}

	return etcdClient
}
