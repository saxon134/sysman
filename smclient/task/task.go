package task

import (
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saUrl"
	"github.com/saxon134/go-utils/saHttp"
	"github.com/saxon134/go-utils/saLog"
	"time"
)

type Client struct {
	root   string
	secret string
}

type Handler func(key string, params string) (resp *RunResult, err error)

var _client *Client
var _tasks map[string]Handler

// Init
// remoteUrl sysman http接口地址
// secret sysman http接口秘钥
// app 应用名称
// host 应用IP
// port 应用端口
func Init(root string, secret string) *Client {
	_client = &Client{root: saUrl.ConnectUri(root, "") + "/", secret: secret}
	return _client
}

// Register 注册任务
func (m *Client) Register(tasks map[string]Handler) {
	if _client == nil {
		panic("未初始化")
	}
	_tasks = tasks
}

func (m *Client) Run(in *RunRequest) (resp map[string]interface{}) {
	resp = map[string]interface{}{
		"key": in.Key, "logId": in.LogId,
	}

	//参数校验
	if in == nil || in.Key == "" || in.LogId <= 0 {
		resp["success"] = false
		resp["err"] = saError.ErrParams
		return
	}

	//签名校验
	if checkSign(in.Sign, in.Timestamp) == false {
		resp["success"] = false
		resp["err"] = saError.ErrUnauthorized
		return
	}

	//查询任务
	if _tasks == nil || _tasks[in.Key] == nil {
		resp["success"] = false
		resp["err"] = saError.ErrNotExisted
		return
	}

	//异步执行任务
	if in.Async == true {
		resp["success"] = true
		go func() {
			//异步执行结果
			var res = map[string]interface{}{"key": in.Key, "logId": in.LogId}

			//执行
			var result, err = _tasks[in.Key](in.Key, in.Params)
			if err != nil {
				res["success"] = false
				res["err"] = err.Error()
			} else {
				resp["success"] = result.Success
				resp["result"] = result.Result
				resp["err"] = result.Err
			}

			//最多回调3次
			for i := 0; i < 3; i++ {
				var sign, timestamp = genSign()
				err = saHttp.Do(saHttp.Params{
					Method: "POST",
					Url:    _client.root + "task/callback",
					Header: map[string]interface{}{"sign": sign, "timestamp": timestamp},
					Body:   map[string]interface{}{"code": 0, "result": res},
				}, nil)
				if err != nil {
					saLog.Err(saError.Stack(err))
					time.Sleep(time.Second * 30 * time.Duration(i+1))
				} else {
					return
				}
			}
		}()
		return
	} else
	//同步执行任务
	{
		var result, err = _tasks[in.Key](in.Key, in.Params)
		if err != nil {
			resp["success"] = false
			resp["err"] = err.Error()
			return
		}

		resp["success"] = result.Success
		resp["result"] = result.Result
		resp["err"] = result.Err
		return
	}
}
