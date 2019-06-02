package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	. "rpc/LightBackend/utils"
	"strconv"
	"strings"
)

type HttpHandler func(http.ResponseWriter, *http.Request)
type ServiceInterface func(context.Context, interface{}) interface{}

const URL_PREFIX string = "/rpc/"

type ServiceDesc struct {
	ServiceName string
	HandlerType interface{}
	Methods     map[string]HttpHandler
	MetaData    string
}

type MethodDesc struct {
	MethodName string
	Handler    interface{}
}

type Server struct {
	S  http.Server
	Mu *http.ServeMux
	SD ServiceDesc
	M  map[string]interface{}
}

type ClientConn struct {
	C               http.Client
	ServiceName     string
	ServiceEndpoint string
	ServicePort     int
}

func NewServer() *Server {
	NewS := &Server{}
	NewS.S = http.Server{}
	NewS.M = make(map[string]interface{})
	NewS.SD = ServiceDesc{}
	NewS.Mu = http.NewServeMux()
	DefaultConfigManager.HotUpdateTarget = append(DefaultConfigManager.HotUpdateTarget, NewS)
	return NewS
}

func (s *Server) RegisterService(sd ServiceDesc, handler interface{}) {
	ServiceName := sd.ServiceName
	s.SD = sd
	s.M[ServiceName] = handler
	for methodName, handler := range sd.Methods {
		RegisterMux(s.Mu, URL_PREFIX+ServiceName+"/"+methodName, handler)
	}
}

func (s *Server) Errors(w http.ResponseWriter, r *http.Request, err error, status int) {
	w.WriteHeader(status)
	// fmt.Println("func write request exception to w and log")
	if status == 500 {
		LOG.Errorf("Coding rpc response error:%v , request_url:%s, data format:%s", err, r.URL, r.Header.Get("Content-Type"))
		_, _ = fmt.Fprint(w, "coding rpc response error.")
	} else if status == 400 {
		LOG.Errorf("Decode rpc request error:%v , request_url:%s, data format:%s", err, r.URL, r.Header.Get("Content-Type"))
		_, _ = fmt.Fprint(w, "decode rpc request error.")
	}
}

func (s *Server) UpdateConfig() {
	// update degrading info

}

func RegisterMux(mu *http.ServeMux, url string, f func(http.ResponseWriter, *http.Request)) {
	mu.Handle(url, RpcHandlerFunc{f: f, ServiceName: strings.Replace(url, "/", ".", -1)})
}

func NewClientConn(serviceName string) ClientConn {
	var cc ClientConn
	cc.ServiceName = serviceName
	cc.C = http.Client{}
	return cc
}

func (cc *ClientConn) Invoke(ctx context.Context, url string, body []byte) ([]byte, error) {
	// build http request to server with url
	infos := strings.Split(url, "/")
	serviceName := infos[0]
	method := infos[1]
	state := NewStateEntry(serviceName + "." + method)
	output, err := cc.invoke(url, body)
	if err != nil {
		LOG.Errorf("send http request to server error:%v", err)
		EndStateEntry(state, 400)
		return nil, err
	} else {
		EndStateEntry(state, 200)
		return output, nil
	}
}

func (cc *ClientConn) invoke(url string, body []byte) ([]byte, error) {
	// load balance
	endpoint := GetServiceEndPoint(cc.ServiceName)
	cc.ServiceEndpoint = endpoint.EndPoint
	cc.ServicePort = endpoint.Port
	//cc.ServiceEndpoint = "127.0.0.1"
	//cc.ServicePort = 19090

	// http post request
	requestUrl := "http://" + cc.ServiceEndpoint + ":" + strconv.Itoa(cc.ServicePort) + URL_PREFIX + url
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type","application/proto")
	if err != nil {
		LOG.Errorf("build http request error:%v", err)
		return nil, err
	}

	resp, err := cc.C.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		output, err := ioutil.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()
		if err != nil {
			return nil, err
		}
		return output, nil
	} else if resp.StatusCode == 503 {
		return nil, errors.New("service degrading or fusing")
	} else if resp.StatusCode == 400 {
		return nil, errors.New("bad request, check request args")
	} else {
		return nil, errors.New("service handler error or unavailable now")
	}
}
