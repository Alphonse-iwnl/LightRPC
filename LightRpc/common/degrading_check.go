package common

import . "rpc/LightBackend/utils"

var FuseMethods []string
var DegradingMethods []string

// CheckFuse 网络请求达到时判断是否熔断
// service实际是method
func CheckFuse(method string) bool {
	return isInFusing(method)
}

func InitDegradingConfig(){
	// 加载配置文件
	// 降级分为两种
	// degrading=true      1.整个服务提供GET接口
	// degradingApi=method 2.降级某个method
}

func CheckDegrading(method string) bool {

	return DefaultToml.Server.Degrading
}

// FuseCheck 协程检查是否有成功率低的服务
// append并非线程安全 所以要保证只有一个协程能修改FuseMethods
func FuseCheck() {
	for method, rate := range Manager.MethodStatus {
		if rate < 0.9 && !isInFusing(method) {
			FuseMethods = append(FuseMethods, method)
		}
		if rate >= 0.9 && isInFusing(method) {
			for index, fuseService := range FuseMethods {
				if fuseService == method {
					FuseMethods = append(FuseMethods[:index], FuseMethods[index+1:]...)
				}
			}
		}
	}
}

func isInFusing(method string) bool {
	methods := FuseMethods
	for _, item := range methods {
		if method == item {
			return true
		}
	}
	return false
}
