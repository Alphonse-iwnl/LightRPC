package main

import (
	"rpc/LightRpc/client"
	"rpc/LightRpc/common"
	. "rpc/LightRpc/init"
)

func main() {
	InitConfig("", nil)
	common.InitBalancer()
	//var wg sync.WaitGroup
	//for i := 0; i < 10; i++ {
	//	wg.Add(1)
	//	go func() {
	//
	//		for j := 0; j < 10000; j++ {
	//			client.RPCRequestDemo("HelloService")
	//		}
	//		wg.Done()
	//	}()
	//}
	//wg.Wait()
	client.RPCRequestDemo("HelloService")
}
