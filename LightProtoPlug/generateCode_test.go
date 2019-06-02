package main

import (
	"fmt"
	"rpc/LightProtoPlug/common"
	"rpc/LightProtoPlug/protocol"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	opt := "all"
	info := common.MethodInfo{
		Protocol:   "rpc",
		MethodName: "ExecService",
		InputName:  "Test01",
		OutputName: "Test",
	}
	var infos []common.MethodInfo
	infos = append(infos, info)

	var total int64 = 0
	for i := 0; i < 1000; i++ {
		startTime:=time.Now().Unix()
		protocol.TypeServiceCode(opt, infos, "/Users/evan/golang/src/rpc/LightRpc/proto/test_service.pb.go")
		endTime:=time.Now().Unix()
		current:=endTime-startTime
		total+=current
	}
	fmt.Println(float64(total)/float64(1000))
}
