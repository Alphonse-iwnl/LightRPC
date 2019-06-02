package init

import (
	"rpc/LightBackend/utils"
	"rpc/LightRpc/common"
	"strconv"
)

func InitService(srv *common.Server) {
	// init server info from config file
	// start listen and serve
	common.InitOps()
	port := utils.DefaultToml.Server.Port
	srv.S.Addr = ":" + strconv.Itoa(port)
	srv.S.Handler = srv.Mu
	//fmt.Println(utils.DefaultToml)

	utils.LOG.Infof("\nInit Server success.\nServer ListenAt:%d,\nOpsPort:%d,\nServiceName:%s,\nDegradingStatus:%s,\nLogLevel:%s",
		utils.DefaultToml.Server.Port,
		utils.DefaultToml.Server.OpsPort,
		utils.DefaultToml.Server.ServiceName,
		strconv.FormatBool(utils.DefaultToml.Server.Degrading),
		utils.DefaultToml.Log.Level)
	err := srv.S.ListenAndServe()
	if err != nil {
		utils.LOG.Fatalf("Start listen error:%v", err)
	}

}

func InitConfig(filePath string, config interface{}) {
	// read config file
	utils.LoadTomlConfig(config, filePath)
	utils.InitDefaultLogger(utils.DefaultToml.Log.Level, utils.DefaultToml.Server.ServiceName)
}

func InitFramework(configPath string, tomlConfig interface{}) {
	InitConfig(configPath, tomlConfig)
	srv := InitHandler()
	InitService(srv)
}
