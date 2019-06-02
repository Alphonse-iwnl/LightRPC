package common

import "C"
import (
	"net/http"
	"rpc/LightBackend/utils"
	"strings"
)

type RpcHandlerFunc struct {
	f           func(http.ResponseWriter, *http.Request)
	ServiceName string
}

type ServiceResponseWriter struct {
	w      http.ResponseWriter
	status int
}

func (sr *ServiceResponseWriter) Write(b []byte) (int, error) {
	return sr.w.Write(b)
}

func (sr *ServiceResponseWriter) Header() http.Header {
	return sr.w.Header()
}

func (sr *ServiceResponseWriter) WriteHeader(statusCode int) {
	sr.status = statusCode
	sr.w.WriteHeader(statusCode)
}

func (hf RpcHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	serviceRW := &ServiceResponseWriter{w: w, status: 200}
	utils.LOG.Infof("Recv request from %s, request_url:%s,", r.RemoteAddr, r.URL)
	// 检查是否有header, 一般用于过滤直接用浏览器发起的链接
	form := r.Header.Get("Content-Type")

	if form == "" {
		hf.EmptyParamHandler(w, r)
		return
	}
	method := strings.Split(hf.ServiceName, ".")[1]
	// 熔断判断
	isFuse := CheckFuse(method)
	if isFuse {
		hf.FuseHandler(w, r)
		return
	}
	isDegrading := CheckDegrading(method)
	if isDegrading {
		hf.DegradingHandler(w, r)
		return
	}
	// 降级判断
	stat := utils.NewStateEntry(hf.ServiceName)
	hf.f(serviceRW, r)
	//time.Sleep(time.Millisecond*20)
	utils.EndStateEntry(stat, serviceRW.status)
}

func (hf RpcHandlerFunc) EmptyParamHandler(w http.ResponseWriter, r *http.Request) {
	utils.LOG.Warn("No param in request, bad request.")
	w.WriteHeader(400)
	_, err := w.Write([]byte("Method need param, please check your request."))
	if err != nil {
		utils.LOG.Errorf("Write response error:%v.", err)
	}
}

func (hf RpcHandlerFunc) FuseHandler(w http.ResponseWriter, r *http.Request) {
	// 这里也要监控一下service.method.fuse 下同
	fuseState := utils.NewStateEntry(hf.ServiceName + ".fuse")
	w.WriteHeader(503)
	utils.EndStateEntry(fuseState, 503)
}

func (hf RpcHandlerFunc) DegradingHandler(w http.ResponseWriter, r *http.Request) {
	degradingState := utils.NewStateEntry(hf.ServiceName + ".degrading")
	w.WriteHeader(503)
	utils.EndStateEntry(degradingState, 503)
}
