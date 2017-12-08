package etcd

import (
	"strconv"
	"strings"
)

const (
	// KeyType_Unknown not standard key patter
	KeyType_Unknown = iota
	// KeyType_Profile ProfileKey , /appName/profile/profileName/key
	KeyType_Profile
	// KeyType_Deploy Deploy key, /appName/deploy/nodeId/key
	KeyType_Deploy

	// KeySeprator for path seperate
	KeySeprator = "/"

	// KeySepratorEnd use for range search in etcd
	KeySepratorEnd = "0"

	KeyProfile = "profiles"

	KeyDeploy = "nodes"

	KeyGrpcServer = "/grpcservers"
)

// EKeyType keytype enum
type EKeyType int

type Key struct {
	AppName     string
	ProfileName string
	FullPath    string
}

// ConfigKey etcd key property
type ConfigKey struct {
	Key
	KeyName string
}

// NodeKey node lease key
type NodeKey struct {
	Key
	NodeId string
}

func split(key string) []string {
	return strings.Split(key, KeySeprator)[1:] //去掉第一个空字符，默认格式是/xxx/xx/xx
}

func ParseConfigKey(key string) ConfigKey {
	var c ConfigKey
	tmp := split(key)
	c.AppName = tmp[0]
	c.FullPath = key
	c.ProfileName = tmp[2]
	c.KeyName = tmp[len(tmp)-1]
	return c
}

func ParseNodeKey(key string) NodeKey {
	var c NodeKey
	tmp := split(key)
	c.AppName = tmp[0]
	c.FullPath = key
	c.ProfileName = tmp[2]
	c.NodeId = tmp[len(tmp)-1]
	return c
}

func NewConfigKey(app, profile, key string) ConfigKey {
	var c ConfigKey
	c.AppName = app
	c.KeyName = profile
	c.ProfileName = profile
	c.FullPath = KeySeprator + app + KeySeprator + KeyProfile + KeySeprator + profile + KeySeprator + key
	return c
}

func NewNodeKey(app, profile, nodeId string) NodeKey {
	var c NodeKey
	c.AppName = app
	c.NodeId = nodeId
	c.ProfileName = profile
	c.FullPath = KeySeprator + app + KeySeprator + KeyDeploy + KeySeprator + profile + KeySeprator + nodeId
	return c
}

// GetKeyPrefix return range of /xx/xx/xx/ , use to list range keys in etcd
func GetKeyPrefix(app, profile string, keyType EKeyType) string {
	var prefix string
	mid := KeyProfile
	if EKeyType(KeyType_Deploy) == keyType {
		mid = KeyDeploy
	}
	prefix = KeySeprator + app + KeySeprator + mid + KeySeprator + profile + KeySeprator
	return prefix
}

func GetCommonNodesKeyPrefix(app string) string {
	return KeySeprator + app + KeySeprator + KeyDeploy + KeySeprator
}

func GetGrpcServerKey(serverip string, port int) string {
	return KeyGrpcServer + KeySeprator + serverip + ":" + strconv.Itoa(port)
}
