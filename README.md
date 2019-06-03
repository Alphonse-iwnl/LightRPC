# LightRPC

<a href="https://github.com/996icu/996.ICU/edit/master/LICENSE">
    <img alt="996icu" src="https://camo.githubusercontent.com/a72e7743f15db219a6aba534f9de456e86268dd6/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f6c6963656e73652d416e74692532303939362d626c75652e7376673f7374796c653d666c61742d737175617265">
</a>


> 本项目是基于protobuf搭建的go语言轻量级RPC框架. 

* 作品是作者2019年6月的本科毕业设计,模仿gRPC的框架搭建流程并重写了protobuf插件;
* 框架流程简单,插件设计比较猎奇.
* 作者过去两个月实习都在用python,可能在代码中出现不优雅甚至错误的go语言编码示范= =.


## Get Start

> 所有的示例都在.git目录下演示,具体路径请自行调整  
> 目前框架有一个echo调用示例

* 框架基础库使用了一些开源库支撑,可以尝试先使用`rely.sh`安装部分依赖;其中protobuf安装的时间较长,请耐心等待.  


### 查看echo调用示例
```
# 启动RPC服务端
$ go run LightRpc/rpc_server.go
# 启动RPC客户端 调用RPC方法
$ go run LightRpc/rpc_client.go
```
* 用户使用框架时需要编写或关注
    * 配置文件
    * proto文件
    * rpc方法实现(server)
    * rpc方法调用(client)
* 用户可根据echo示例和后文的教程熟悉使用流程并自定义RPC服务.

### 安装protobuf插件
```
$ git clone https://github.com/Alphonse-iwnl/LightRPC.git
$ cd LightProtoPlug/
$ go install
```
### 编写.proto文件并使用插件
```
$ cd LightRpc/proto/
$ touch test.proto
# 编写proto内容 e.g.
# 可以这样定义RPC方法
# service HelloService{
#     rpc ExecService(<input struct>) returns (<output struct>);
# }
...
# opt可为server,cli,all;可以分别生成两端代码
$ LightProtoPlug --file test.proto --opt all
```

### 实现接口
* 编辑文件LightRpc/server/rpcServer.go
* 为**ExecService**结构体实现proto文件中定义的方法 e.g.
```
# ExecService结构体是固定的 使用该结构体可以避免用户自己去框架注册方法
func (es *ExecService) ExecService(ctx context.Context, in *<input struct>) (*<output struct>, error)
```
### 编写配置文件 启动RPC服务端
* 配置文件默认为LightRpc/conf/config.toml,框架自动搜索并读取;如果要修改文件名字需要在启动时修改输入参数.
* 启动框架
```
go run LightRpc/rpc_server.go
```

### 进行RPC调用
* LightRpc/client/rpcClient.go文件中调用示例,用户可根据proto文件中自定义的结构体和方法修改demo
```
go run LightRpc/rpc_client.go
```
---
文档待补充...
### 基础设施
* 监控
* 日志
* 配置文件

### 配置文件热更新
* 框架启动后同时监听运维接口(配置文件中为OpsPort),运维接口目前提供health check和配置文件热更新
* 框架启动后修改本地配置文件,调用<addr:port>/reconfig,框架会开始热更新;目前支持的热更新配置较少
* 观察者模式 

### 服务降级与熔断
