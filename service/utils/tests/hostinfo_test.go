package tests

import (
	"fmt"
	"testing"

	"github.com/yy-java/cnt2/service/utils"
)

func Test_GetMyIPInfo(t *testing.T) {
	for i, r := range utils.GetMyIPInfo() {
		fmt.Println(i, r)
	}
}

func Test_PropertiesReader(t *testing.T) {
	fmt.Println(utils.GetHostInfo("E:\\home\\dspeak\\yyms\\hostinfo.ini"))
}
