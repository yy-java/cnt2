package etcd

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"errors"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
)

const (
	DefaultLeaseTTL = 30
	DefaultTimeout  = 5 * time.Second
)

type WatchCallback interface {
	Callback(*mvccpb.KeyValue, mvccpb.Event_EventType)
}

type EtcdClientShared struct {
	Status  atomic.Value
	LeaseId atomic.Value
}

type EtcdClient struct {
	client          *clientv3.Client
	status          EtcdClientShared
	stopCtx         context.Context
	stopFunc        context.CancelFunc
	startLeaseMutex sync.Mutex
	Config          EtcdClientConfig
}

type EtcdClientConfig struct {
	Endpoints         []string
	DialTimeout       time.Duration
	RequestTimeout    time.Duration
	LeaseTtl          int64
	WatchCallbackFunc WatchCallback
}

func makeClient(dialTimeout, requestTimeout time.Duration, endpoints []string) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Printf("Create etcd client error: %v", err)
	}

	return cli, nil
}

func closeClient(cli *clientv3.Client) {
	if cli == nil {
		return
	}

	err := cli.Close()
	if err != nil {
		log.Printf("Close etcd client error: %v", err)
	}
}

func handleOpError(err error, op string, key string) {
	if err == nil {
		return
	}

	switch err {
	case context.Canceled:
		log.Printf("ctx is canceled by another routine: %v. op: %s, key: %s", err, op, key)
	case context.DeadlineExceeded:
		log.Printf("ctx is attached with a deadline is exceeded: %v. op: %s, key: %s", err, op, key)
	case rpctypes.ErrEmptyKey:
		log.Printf("client-side error: %v\n", err)
	default:
		log.Printf("bad cluster endpoints, which are not etcd servers: %v. op: %s, key: %s", err, op, key)
	}

}

func NewEtcdClient(config EtcdClientConfig) (*EtcdClient, error) {
	var client *EtcdClient
	if len(config.Endpoints) == 0 {
		return nil, errors.New("EndPoints must be set")
	}
	cli, err := makeClient(config.DialTimeout, config.RequestTimeout, config.Endpoints)
	if err != nil {
		return nil, err
	}
	if config.DialTimeout.Seconds() <= 0 {
		config.DialTimeout = DefaultTimeout
	}
	if config.RequestTimeout.Seconds() <= 0 {
		config.RequestTimeout = DefaultTimeout
	}
	if config.LeaseTtl <= 0 {
		config.LeaseTtl = DefaultLeaseTTL
	}
	ctx, cancel := context.WithCancel(context.Background())
	client = &EtcdClient{client: cli, Config: config, stopCtx: ctx, stopFunc: cancel}
	client.status.Status.Store(true)
	client.status.LeaseId.Store(clientv3.LeaseID(0))
	return client, nil
}

func (e *EtcdClient) Close() {
	status := e.GetStatus()
	if status {
		e.status.Status.Store(false)
		e.stopFunc()
		closeClient(e.client)
	}
}

func (e *EtcdClient) GetStatus() bool {
	status := false
	es := e.status.Status.Load()
	status, _ = es.(bool)
	return status
}

func (e *EtcdClient) GetLeaseId() clientv3.LeaseID {
	var leaseId clientv3.LeaseID
	es := e.status.LeaseId.Load()
	leaseId, ok := es.(clientv3.LeaseID)
	if ok {
		return leaseId
	} else {
		return clientv3.LeaseID(0)
	}
}

func (e *EtcdClient) Put(key, val string) (*clientv3.PutResponse, error) {
	cli := e.client
	ctx, cancel := context.WithTimeout(context.Background(), e.Config.RequestTimeout)
	resp, err := cli.Put(ctx, key, val)
	cancel()

	handleOpError(err, "Put", key)

	return resp, err
}

func (e *EtcdClient) Get(key string) (*clientv3.GetResponse, error) {
	cli := e.client
	ctx, cancel := context.WithTimeout(context.Background(), e.Config.RequestTimeout)
	resp, err := cli.Get(ctx, key)
	cancel()

	handleOpError(err, "Get", key)

	return resp, err
}

func (e *EtcdClient) Delete(key string) (*clientv3.DeleteResponse, error) {
	cli := e.client

	ctx, cancel := context.WithTimeout(context.Background(), e.Config.RequestTimeout)
	resp, err := cli.Delete(ctx, key)
	cancel()

	handleOpError(err, "Delete", key)

	return resp, err
}

func (e *EtcdClient) DeleteKeys(keyPrefix string) (*clientv3.DeleteResponse, error) {
	cli := e.client

	ctx, cancel := context.WithTimeout(context.Background(), e.Config.RequestTimeout)

	resp, err := cli.Delete(ctx, keyPrefix, clientv3.WithPrefix())
	log.Printf("DeleteKeys key:%s, resp:%v\n", keyPrefix, resp)

	cancel()

	handleOpError(err, "DeleteKeys", keyPrefix)

	return resp, err
}

func (e *EtcdClient) ListKeys(keyPrefix string) ([]*mvccpb.KeyValue, error) {
	cli := e.client

	ctx, cancel := context.WithTimeout(context.Background(), e.Config.RequestTimeout)
	resp, err := cli.Get(ctx, keyPrefix, clientv3.WithPrefix())

	cancel()

	log.Printf("ListKeys %s, %v\n", resp, err)

	handleOpError(err, "ListKeyPrefix ", keyPrefix)
	if err != nil {
		return make([]*mvccpb.KeyValue, 0), err
	}

	return resp.Kvs, nil
}

/*
Watch
	key	      watch key
	recursive   watch key as Prefix
*/
func (e *EtcdClient) Watch(key string, recursive bool) int {
	return e.Watch2(key, recursive, e.Config.WatchCallbackFunc)
}

/*
Watch2
	key	        watch key
	recursive   watch key as Prefix
	callback    interface WatchCallback
*/
func (e *EtcdClient) Watch2(key string, recursive bool, callback WatchCallback) int {
	cli := e.client

	var rch *clientv3.WatchChan
	if recursive {
		t := cli.Watch(context.Background(), key, clientv3.WithPrefix())
		rch = &t
	} else {
		t := cli.Watch(context.Background(), key)
		rch = &t
	}

	if rch != nil {
		go e.watchThread(key, rch, callback)
	} else {
		return 0
	}
	return 1
}

func (e *EtcdClient) watchThread(key string, rch *clientv3.WatchChan, callback WatchCallback) {
	for {
		select {
		case wresp := <-*rch:
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					log.Printf("Received PUT event. event: %s, key: %q, value: %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
				case mvccpb.DELETE:
					log.Printf("Received DELETE event. event: %s, key: %q, value: %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
				default:
					log.Printf("Unresolved event. event: %s, key: %q, value: %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
				}
				if callback != nil {
					callback.Callback(ev.Kv, ev.Type)
				}
			}
		case <-e.stopCtx.Done():
			log.Printf("Stop Watch on %s \n", key)
			return
		}
	}
}

func (e *EtcdClient) startLease() {
	firstStart := false
	e.startLeaseMutex.Lock()
	if e.GetLeaseId() > 0 {
		return
	}
	cli := e.client
	ttl := e.Config.LeaseTtl
	resp, err := cli.Grant(context.TODO(), ttl)
	if err != nil {
		log.Fatalf("start lease error %v", err)
		return
	}
	leaseId := resp.ID
	e.status.LeaseId.Store(leaseId)
	firstStart = true
	e.startLeaseMutex.Unlock()
	if firstStart {
		go e.startKeepAlive(e.GetLeaseId())
	}
}

func (e *EtcdClient) startKeepAlive(leaseId clientv3.LeaseID) {
	cli := e.client
	for {
		kch, err := cli.KeepAlive(context.TODO(), leaseId)
		if err != nil {
			s := e.Config.LeaseTtl / 2
			if s == 0 {
				s = 1
			}
			time.Sleep(time.Second * time.Duration(s))
		}
	Out:
		for {
			select {
			case ka := <-kch:
				log.Println("keepaLive Response ", ka)
				if ka == nil {
					break Out // stop this select loop ,create new Lease chan
				}
			case <-e.stopCtx.Done():
				e.status.LeaseId.Store(clientv3.LeaseID(0))
				cli.Revoke(context.TODO(), leaseId)
				log.Println("exit Lease [", leaseId, ") loop")
				return
			}
		}
	}
}
func (e *EtcdClient) Lease(k, v string) error {
	if e.GetLeaseId() <= 0 {
		e.startLease()
	}
	cli, leaseId := e.client, e.GetLeaseId()

	rsp, err := cli.Put(context.TODO(), k, v, clientv3.WithLease(leaseId))

	if err != nil {
		return err
	}
	log.Println(k, v, " with leaseId ", leaseId, rsp)
	return nil

}
