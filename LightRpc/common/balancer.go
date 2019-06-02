package common

import (
	"net/http"
	. "rpc/LightBackend/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ServiceInstance struct {
	EndPoint string
	Port     int
	Ops      bool
	OpsPort  int
}

type ServiceManager struct {
	Offset              int
	AvailableInstance   map[string][]ServiceInstance
	UnAvailableInstance map[string][]ServiceInstance
	lock                *sync.Mutex
}

var LBManager *ServiceManager

func InitBalancer() {
	LBManager = &ServiceManager{}
	initServiceManager(DefaultToml)
	//go healthCheck()
	DefaultConfigManager.HotUpdateTarget = append(DefaultConfigManager.HotUpdateTarget, LBManager)
}

func initServiceManager(conf *DefaultTomlConfig) {
	LBManager.AvailableInstance = make(map[string][]ServiceInstance)
	LBManager.UnAvailableInstance = make(map[string][]ServiceInstance)
	LBManager.lock = &sync.Mutex{}
	LBManager.Offset = 0
	for _, item := range conf.ServerClient {
		endPoints := strings.Split(item.EndPoint, ",")
		var NewInstance ServiceInstance
		NewInstance.Port = item.Port
		if item.OpsPort != 0 {
			NewInstance.Ops = true
			NewInstance.OpsPort = item.OpsPort
		} else {
			NewInstance.Ops = false
		}
		for _, ep := range endPoints {
			if ep == ""{
				continue
			}
			NewInstance.EndPoint = ep
			LBManager.AvailableInstance[item.ServiceName] = append(LBManager.AvailableInstance[item.ServiceName], NewInstance)
		}
	}
}

func (sm *ServiceManager) UpdateConfig() {
	// reload lb
	initServiceManager(DefaultToml)
}

func GetServiceEndPoint(serviceName string) *ServiceInstance {
	return LBManager.RoundRobin(serviceName)
}

func (sm *ServiceManager) RoundRobin(serviceName string) *ServiceInstance {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if _, ok := sm.AvailableInstance[serviceName]; !ok {
		LOG.Error("No Service endpoint in config file.")
		return nil
	}
	resultInstance := sm.AvailableInstance[serviceName][sm.Offset]
	sm.Offset = sm.Offset + 1

	if sm.Offset == len(sm.AvailableInstance[serviceName]) {
		sm.Offset = 0
	}
	return &resultInstance
}

func healthCheck() {
	for {
		// check unavailable endpoint in availableList
		for serviceName, instances := range LBManager.AvailableInstance {
			for index, item := range instances {
				if !_sendHeartBeat(item) {
					LBManager.lock.Lock()
					LBManager.UnAvailableInstance[serviceName] = append(LBManager.UnAvailableInstance[serviceName], item)
					LBManager.AvailableInstance[serviceName] = append(LBManager.AvailableInstance[serviceName][:index], LBManager.AvailableInstance[serviceName][index+1:]...)
					LBManager.lock.Unlock()
				}
			}
		}
		// check available endpoint in unavailableList
		for serviceName, instances := range LBManager.UnAvailableInstance {
			for index, item := range instances {
				if _sendHeartBeat(item) {
					LBManager.lock.Lock()
					LBManager.AvailableInstance[serviceName] = append(LBManager.AvailableInstance[serviceName], item)
					LBManager.UnAvailableInstance[serviceName] = append(LBManager.UnAvailableInstance[serviceName][:index], LBManager.UnAvailableInstance[serviceName][index+1:]...)
					LBManager.lock.Unlock()
				}
			}
		}
		time.Sleep(time.Second * 10)
	}

}

func _sendHeartBeat(instance ServiceInstance) bool {
	resp, err := http.Get(instance.EndPoint + ":" + strconv.Itoa(instance.OpsPort) + "/health")
	if err != nil {
		LOG.Errorf("send http health check error:%v", err)
		return true
	}
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}
