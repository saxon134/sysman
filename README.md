# sysman

### 介绍


system management

系统运维包含：
>1. 服务管理(sdp Service Discovery Protocol)：服务注册、广播消息、定点消息、应用状态监测；
>2. 任务管理(task)：开启关闭任务、执行日志、支持任务流配置
>3. 角色权限控制 RBAC（Role-Based Access Control）：权限、菜单、角色、用户、查看角色哪些用户有


### 使用须知

>1. MySQL、Redis必须配置
>2. 如果配置了秘钥，则所有接口都需要加密
> > 加密方式是：sign = lowercase( MD5(secret + timestamp) )，请求header带上sign和t
>3. sysman/smclient 目录下是接入时引用代码
>4. web是前端代码


### 配置文件说明（config.yaml）
具体参考conf.go中模型的定义
```
name: sysman //项目名称
http:
  port: 6231
  root: /sysman
  secret: 
  clientroot: /sysman

sdp:
  pingsecond: 5
  secret:

redis:
  host: 127.0.0.1:6379
  pass:

mysql:
  host:
  pass:
  user:
  db: sysman

```



### 安装教程

```
go get github.com/saxon134/sysman/sysman;
go mod tidy;
go run main.go;
```


### SDP服务管理说明

#### 1. 初始化SDP Client

```
var client = NewClient(remoteHost string, secret string) 
```
  
返回client实例


#### 2. 注册服务
```
// Register 注册服务
// app、host、port: 服务信息
client.Register(app string, host string, port int)  
```

#### 3. 发现服务
```
host, port := client.Discovery(app string)

```

#### 4. 发送消息

```
// app: 必选，指定应用
// host/port: 可选，指定服务器；空则向所有实例的应用发送消息
func (m *Client) SendMsg(app string, host string, port int, msg interface{}) (err error)  

```
如果要对指定应用所有实例进行广播，则host设置为空


#### 5. 各应用应该提供HTTP接口，以便接收消息

```
POST {{conf.http.clientRoot}}/sdp/msg
```


### Task服务说明

> 支持查看任务状态、执行记录、最近执行时间
> 
> 支持关闭、开启任务
> 
> 支持触发一次执行一次任务
>
> 服务启动的时候，会去读取配置，确定是否开启任务
> 
> 不支持指定服务器配置任务，一般执行任务coding阶段就已经确定了执行服务器的，没必要再动态指定

#### 1. 各应用应该提供HTTP接口，以便接收任务消息处理

```
// 任务变更接口
// 在该接口中调用task.Event
POST {{conf.http.clientRoot}}/task/event

// 获取任务状态接口
// 在该接口中调用task.Status
GET {{conf.http.clientRoot}}/task/status
```



### 名词解释

#### resource
> 资源，如ECS，Redis，以host + port区分
>
> 实例运行注册，记录信息入库，会调用阿里云接口获取资源状态

#### app
> 系统，如电商系统、ORM等，以名词区分
>
> 实例运行注册，记录信息入库

#### case
> 实例，即部署的系统运行实例，以host + port区分
>
> 实例运行注册，记录信息入库

#### task
> 任务以key区分
>
> 实例运行注册，记录信息入库

