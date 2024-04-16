package task

type RunRequest struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Key       string `json:"key"`
	LogId     int64  `json:"logId"`
	Async     bool   `json:"async"`
	Params    string `json:"params"`
}

type RunResult struct {
	Success bool   `json:"success"`
	Err     string `json:"err"`
	Result  string `json:"result"`
}
