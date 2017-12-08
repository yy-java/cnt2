package server

import (
	"flag"
	"fmt"
	"testing"

	"github.com/astaxie/beego/orm"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"time"
	"github.com/yy-java/cnt2/grpc/pb"
)

var (
	serverPort = 20000
)

func init() {
	orm.RegisterDataBase("default", "mysql", "cnt2_db_user:q0NUVMca1@tcp(58.215.143.133:6307)/cnt2_db?charset=utf8&loc=Asia%2FShanghai", 5, 10)
}

func getConn() (*grpc.ClientConn, error) {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	serverAddr := fmt.Sprintf("127.0.0.1:%d", serverPort)
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	return conn, err
}

func startServer(abort chan int) {
	go Start(NewServer(serverPort))
	for {
		select {
		case <-abort:
			return
		}
	}
}

func stopServer(abort chan int) {
	abort <- 1
}

func Test_QueryAll(t *testing.T) {
	ch := make(chan int, 1)
	go startServer(ch)
	conn, err := getConn()
	defer conn.Close()
	defer stopServer(ch)
	client := pb.NewConfigCenterServiceClient(conn)

	req := &pb.QueryRequest{App: "app1", Profile: "dev"}

	resp, err := client.QueryAll(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func Test_QueryKey(t *testing.T) {
	ch := make(chan int, 1)
	go startServer(ch)
	conn, err := getConn()
	defer conn.Close()
	defer stopServer(ch)
	client := pb.NewConfigCenterServiceClient(conn)

	req := &pb.QueryRequest{App: "app1", Profile: "dev", Key: "key1", KeyVersion: 10}

	resp, err := client.QueryKey(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func Test_RegisterClient(t *testing.T) {
	ch := make(chan int, 1)
	go startServer(ch)
	conn, err := getConn()
	defer conn.Close()
	defer stopServer(ch)
	client := pb.NewConfigCenterServiceClient(conn)

	req := &pb.RegisterRequest{App: "1", Profile: "1", ServerIp: "127.0.0.1", Pid: "20000", RegisterTime: time.Now().Unix() * 1000}

	resp, err := client.RegisterClient(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

	time.Sleep(1000 * time.Second)
}

func Test_ValueChangeResultNotify(t *testing.T) {
	ch := make(chan int, 1)
	go startServer(ch)
	conn, err := getConn()
	defer conn.Close()
	defer stopServer(ch)
	client := pb.NewConfigCenterServiceClient(conn)

	req := &pb.ValueChangeResultRequest{NodeId: "1", App: "app1", Profile: "dev", Key: "key1", Result: pb.ValueChangeResultRequest_SUCCESS}

	resp, err := client.ValueChangeResultNotify(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
