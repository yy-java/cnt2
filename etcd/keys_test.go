package etcd

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	} else {
		message = fmt.Sprintf("%s %v != %v", message, a, b)
	}
	fmt.Println(message)
	t.Fatal(message)
}

func Test_GetKeyPrefix(t *testing.T) {
	ret := GetKeyPrefix("test", "sh", EKeyType(KeyType_Deploy))
	fmt.Println(ret)
	assertEqual(t, ret, "/test/nodes/sh/", "")
}

func Test_NewConfigKey(t *testing.T) {
	ret := NewConfigKey("test", "sh", "key1")
	fmt.Println(ret)
	assertEqual(t, ret.FullPath, "/test/profiles/sh/key1", "")
}

func Test_NewNodeKey(t *testing.T) {
	ret := NewNodeKey("test", "sh", "12312")
	fmt.Println(ret)
	assertEqual(t, ret.FullPath, "/test/nodes/sh/12312", "")
}

func Test_ParseConfigKey(t *testing.T) {
	s := ParseConfigKey("/test/profiles/sh/key1")
	assertEqual(t, s.AppName, "test", "appName")
	assertEqual(t, s.ProfileName, "sh", "ProfileName")
	assertEqual(t, s.KeyName, "key1", "KeyName")
}

func Test_ParseNodeKey(t *testing.T) {
	s := ParseNodeKey("/test/nodes/sh/key1")
	assertEqual(t, s.AppName, "test", "appName")
	assertEqual(t, s.ProfileName, "sh", "ProfileName")
	assertEqual(t, s.NodeId, "key1", "NodeId")
}
