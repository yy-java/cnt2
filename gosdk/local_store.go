package gosdk

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/coreos/etcd/pkg/fileutil"
	"log"
	"sync"
	"time"
)

const (
	BUCKET_NAME = "cnt2_bucket"
)

var (
	inited          bool
	lstore          *LocalStore
	localStore_Lock sync.Mutex
)

func InitLocalStore(app, path string) (*LocalStore, error) {
	if inited {
		return lstore, nil
	}
	if len(path) <= 0 {
		path = "/data/file/cnt2"
	}
	if !fileutil.Exist(path) {
		fileutil.CreateDirAll(path)
	}
	path = path + "/cnt2_" + app + ".db"
	localStore_Lock.Lock()
	defer localStore_Lock.Unlock()
	if inited {
		return lstore, nil
	}
	//方便测试使用
	//db, err := bolt.Open("cnt2_"+strconv.Itoa(os.Getpid())+".db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var bucket *bolt.Bucket
	//init bucket
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err = tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		return err
	})
	if err != nil {
		log.Println("createBucket error:%s", err)
		return nil, err
	}
	lstore = &LocalStore{Db: db}
	inited = true
	return lstore, nil
}
func (ls *LocalStore) Get(app, profile, key string) string {
	var result string
	ls.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		resultBytes := bucket.Get([]byte(genKey(app, profile, key)))
		if resultBytes != nil {
			result = string(resultBytes)
		}
		return nil
	})
	return result
}
func (ls *LocalStore) Has(app, profile, key string) bool {
	var result bool
	ls.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		resultBytes := bucket.Get([]byte(genKey(app, profile, key)))
		result = resultBytes != nil
		return nil
	})
	return result
}
func (ls *LocalStore) GetConfig(app, profile, key string) (*Config, error) {
	var config Config
	err := ls.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		val := bucket.Get([]byte(genKey(app, profile, key)))
		if val != nil {
			err := json.Unmarshal(val, &config)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &config, nil
}
func (ls *LocalStore) Put(app, profile, key, value string) error {
	return ls.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		return bucket.Put([]byte(genKey(app, profile, key)), []byte(value))
	})
}
func (ls *LocalStore) Del(app, profile, key string) error {
	return ls.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		return bucket.Delete([]byte(genKey(app, profile, key)))
	})
}
func genKey(app, profile, key string) string {
	return app + "/" + profile + "/" + key
}
