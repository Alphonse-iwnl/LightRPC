package init

import (
	"rpc/LightRpc/common"
	. "rpc/LightRpc/proto"
	"rpc/LightRpc/server"
)

func InitHandler() *common.Server{
	// register rpc method to server
	handler := server.ExecService{}
	srv := common.NewServer()
	RegisterServiceServer(ServiceMiddle{S: srv}, &handler)
	return srv
}
