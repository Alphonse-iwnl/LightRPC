package client

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"rpc/LightBackend/utils"
	"rpc/LightRpc/common"
	pb "rpc/LightRpc/proto"
)

// RPCRequestDemo need service name , usually in proto file
func RPCRequestDemo(serviceName string) {

	// build clientConn
	conn := common.NewClientConn(serviceName)
	// build input struct
	Input := &pb.Test01{
		Message: proto.String("Hello test"),
	}
	proxy := pb.NewHelloServiceClient(&conn)
	output, err := proxy.ExecService(context.Background(), Input)
	if err != nil {
		utils.LOG.Errorf("Send RPC request error:%v.", err)
		return
	}
	fmt.Printf("Recv RPC response, output message:%s,id:%s\n", *output.Message, *output.Id)
}
