package task

// RunResponse
// @Description: 发起执行步骤请求的返回
type RunResponse struct {
	Code   int          `json:"code"`
	Msg    string       `json:"msg"`
	Result CallbackData `json:"result"`
}

// CallbackData
// @Description: 异步任务回调结果
type CallbackData struct {
	Success bool   `json:"success"`
	Err     string `json:"err"`
	Result  string `json:"result"`
	Key     string `json:"key"`
	LogId   int64  `json:"logId"` //本次task的日志ID
}
