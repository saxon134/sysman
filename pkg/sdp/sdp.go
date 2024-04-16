package sdp

import (
	"fmt"
	"github.com/saxon134/sysman/pkg/sm"
	"time"
)

const RedisAppKey = "sdp:apps:%s"

type Instance struct {
	Host   string `json:"h"`
	Port   int    `json:"p"`
	Weight int    `json:"w,omitempty"` //权重
	Time   int64  `json:"t"`           //上次ping的时间
}

type Request struct {
	App  string  `json:"app" form:"app"`
	Host string  `json:"host" form:"host"`
	Port int     `json:"port" form:"port"`
	Cpu  float32 `json:"cpu" form:"cpu"`
	Memo float32 `json:"memo" form:"memo"`
}

var Chan chan Request

func Init() {
	Chan = make(chan Request, 10)
	go registerAndPing()
}

// 处理的时候必须保证要能拿到Redis数据
// 所以使用channel保证一致性，如果是部署多个实例，因为并发可能性极低，暂不考虑
func registerAndPing() {
	for {
		if in, ok := <-Chan; ok {
			var key = fmt.Sprintf(RedisAppKey, in.App)
			var sdpAry = make([]*Instance, 0, 10)
			_ = sm.Redis.GetObj(key, &sdpAry)

			var existed *Instance
			for _, c := range sdpAry {
				if c.Host == in.Host && c.Port == in.Port {
					existed = c
					break
				}
			}

			var now = time.Now().UnixMilli()

			//不存在是注册
			if existed == nil {
				var m = &Instance{
					Host:   in.Host,
					Port:   in.Port,
					Weight: 100,
					Time:   now,
				}
				if in.Cpu > 0.9 || in.Memo > 0.9 {
					m.Weight = 1
				} else if in.Cpu > 0.8 || in.Memo > 0.8 {
					m.Weight = 10
				}
				sdpAry = append(sdpAry, m)
			} else
			//存在是ping
			{
				//超过1秒超过200毫秒可能服务就慢了
				var delay = now - existed.Time - int64(sm.Conf.Sdp.PingSecond*1000)
				if delay > 1000 {
					existed.Weight = 1
				} else if delay > 200 {
					existed.Weight = 10
				} else {
					existed.Weight = 100
				}

				if in.Cpu > 0.9 || in.Memo > 0.9 {
					existed.Weight = 1
				} else if in.Cpu > 0.8 || in.Memo > 0.8 {
					existed.Weight = 10
				}
			}
			_ = sm.Redis.SetObj(key, sdpAry, time.Second*time.Duration(sm.Conf.Sdp.PingSecond+2))

			////新应用，入库
			//if existed == nil {
			//	var appObj = new(models.TblApp)
			//	err := db.MySql.Table(models.TBNApp).Where("app = ?", in.App).First(appObj).Error
			//	if db.MySql.IsError(err) {
			//		saLog.Err(saError.New(err))
			//	} else if appObj.Id <= 0 {
			//		appObj.App = in.App
			//		appObj.CreateAt = saTime.Now()
			//		err = db.MySql.Table(models.TBNApp).Save(appObj).Error
			//		if err != nil {
			//			saLog.Err(saError.New(err))
			//		}
			//	}
			//}
		}
	}
}
