package common

import (
	"fmt"
	"net/http"
	"rpc/LightBackend/utils"
	"strconv"
)

type Handler func(http.ResponseWriter, *http.Request)

func InitOps() {
	OpsServer := Server{}
	OpsServer.Mu = &http.ServeMux{}
	OpsServer.S = http.Server{}
	OpsServer.S.Addr = ":" + strconv.Itoa(utils.DefaultToml.Server.OpsPort)
	OpsServer.Mu.Handle("/health", Handler(healthCheckHandler))
	OpsServer.Mu.Handle("/reconfig", Handler(reConfigHandler))
	//RegisterMux(OpsServer.Mu, "/health", healthCheckHandler)
	//RegisterMux(OpsServer.Mu, "/reconfig", reConfigHandler)
	go func() {
		//_ = OpsServer.S.ListenAndServe()
		_ = http.ListenAndServe(":"+strconv.Itoa(utils.DefaultToml.Server.OpsPort), OpsServer.Mu)
	}()
	fmt.Println("Ops Handler Listen at ", OpsServer.S.Addr)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// health check理论上只需要处理两种情况:正常请求和服务熔断
	// 因为网络不通或是服务挂掉都不会进到这里来
	w.WriteHeader(200)
	return
}

func reConfigHandler(w http.ResponseWriter, r *http.Request) {
	// 配置文件热更新
	// 首先重读一次配置文件 然后触发Manager的更新 让所有已经注册的模块执行更新
	utils.DefaultConfigManager.HotUpdateConfig()
	fmt.Println("reconfig...")
	_, _ = w.Write([]byte("Config Hot Config Success."))
}
