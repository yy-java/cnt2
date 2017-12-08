package etcd

import (
	"testing"
	"time"

	"fmt"

	"strconv"

	"github.com/coreos/etcd/mvcc/mvccpb"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yy-java/cnt2/service/models"
)

type defaultWatchCallback struct {
}

func (d *defaultWatchCallback) Callback(kv *mvccpb.KeyValue, event mvccpb.Event_EventType) {
	fmt.Printf("watcherLogger %v,%v\n", kv, event)
}

var (
	testEndpoints = []string{"61.147.187.152:2379", "61.147.187.142:2379", "61.147.187.150:2379"}
	defaultConfig = EtcdClientConfig{Endpoints: testEndpoints, DialTimeout: 3 * time.Second, RequestTimeout: 3 * time.Second, WatchCallbackFunc: &defaultWatchCallback{}}
	client, _     = NewEtcdClient(defaultConfig)
)

func TestKV(t *testing.T) {
	//defer client.Close()
	Convey("Subject: Test etcd operation", t, func() {
		key := "pairTest"
		client.Put(key, "value1")

		Convey("Test get operation", func() {
			resp, _ := client.Get(key)
			So(string(resp.Kvs[0].Value), ShouldEqual, "value1")

			Convey("Test delete operation", func() {
				resp, _ := client.Delete(key)
				So(resp.Deleted, ShouldEqual, 1)

				Convey("Test get for confirm", func() {
					resp, _ := client.Get(key)
					So(resp.Count, ShouldEqual, 0)
				})
			})
		})

	})
}

func Test_Watch(t *testing.T) {
	client.Watch("myname", false)
	client.Watch("myname2", false)
	client.Put("myname2", "test2")
	client.Put("myname", "test")
	client.Delete("myname")
	client.Delete("myname2")
	time.Sleep(1 * time.Second)
}

func Test_ListKeys(t *testing.T) {

	keyPrefix := "/testGoProxy/nodes/dev"
	client.Watch(keyPrefix, true)
	s, _ := client.ListKeys(keyPrefix)
	for _, v := range s {
		node := models.RegisterNode{}
		node.FromJson(string(v.Value))
		fmt.Println(node.ToJson())
	}
}

func Test_Lease(t *testing.T) {

	key := "leaseTest"
	val := "testLease"

	client.Watch(key, false)
	client.Lease(key, val)

	rsp, _ := client.Get(key)
	fmt.Println("Get 1", rsp.Kvs)
	time.Sleep(5 * time.Second)
	rsp, _ = client.Get(key)
	fmt.Println("Get 2", rsp.Kvs)
	client.stopFunc()
	time.Sleep(5 * time.Second)
	fmt.Println("current leaseId after stop ", client.GetLeaseId())
	fmt.Println("stop lease")
	time.Sleep(10 * time.Second)
	rsp, _ = client.Get(key)
	fmt.Println("Get 3", rsp.Kvs)
	rsp, _ = client.Get(key)

}

func concurrentGet(abort chan int, errs chan error) {
	for {
		select {
		case <-abort:
			return
		default:
			resp, err := client.Get("test")
			if err != nil {
				fmt.Println("Get Error ", err)
				errs <- err
			}
			fmt.Println("Get on ", resp.Kvs)
			time.Sleep(time.Millisecond * 5)
		}
	}
}

func concurrentWatch(abort chan int, errs chan error) {
	fmt.Println("Wathcing")
	client.Watch("test", false)
	i := 1
	for {
		select {
		case <-abort:
			return
		default:
			i += 1
			ai := strconv.Itoa(i)
			resp, err := client.Put("test", ai)
			if err != nil {
				fmt.Println("Put Error ", err)
				errs <- err
			}
			fmt.Println("Put on ", resp)
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func Test_ConcurrenClient(t *testing.T) {
	abort := make(chan int, 2)
	errs := make(chan error)

	go concurrentGet(abort, errs)
	go concurrentWatch(abort, errs)

	time.Sleep(2 * time.Second)
	abort <- 1
	abort <- 1
	for {
		select {
		case err := <-errs:
			fmt.Println("Error ", err)
		default:
			fmt.Println("Finish")
			return
		}
	}
}
