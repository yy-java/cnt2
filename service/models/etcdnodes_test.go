package models

import (
	"fmt"
	"testing"
)

func Test_GrpcServerRegisterInfo_Json(t *testing.T) {
	s := make(map[int]string)
	s[1] = "test"
	g := GrpcServerRegisterInfo{s, 123, 123}
	k, e := g.ToJson()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(k)
	var g2 GrpcServerRegisterInfo
	g2.FromJson(k)
	fmt.Println(g2)
}
