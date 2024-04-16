package sm

type Status int

const (
	StatusNull             Status = 0  //未知，初始
	StatusProcessing       Status = 1  //处理中
	StatusSuccess          Status = 2  //成功，正常
	StatusPartiallySuccess Status = 3  //部分成功
	StatusFail             Status = -1 //失败，有异常

	StatusEnable  Status = 2  //激活
	StatusDisable Status = -1 //禁止
	StatusError   Status = -1 //失败
)

func (m Status) Int() int {
	return int(m)
}
