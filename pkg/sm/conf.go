package sm

import (
	yaml "gopkg.in/yaml.v2"
	"os"
)

var Conf *ModelConf

type ModelConf struct {
	Name string
	Http struct {
		Host       string //服务器IP，一般为内网IP，无则会自动获取IPv4地址
		Port       string //http端口
		Root       string //接口根路由，建议：/sysman
		Secret     string //如果配置了秘钥，所有接口都需要加密
		ClientRoot string //sysman调用各系统的根路由，建议：/sysman
	}
	JwtSecret string //登录加密秘钥

	Sdp struct {
		PingSecond int //ping间隔（秒），默认5秒，最小2秒，超过1个间隔周期+3秒未发送ping则会移除该实例
	}

	Mq struct{}

	Redis struct {
		Host string
		Pass string
	}

	MySql struct {
		Host string
		Pass string
		User string
		Db   string
	}

	Feishu struct {
		Webhookurl string
		Secret     string
	}
}

func initConf() {
	if Conf == nil {
		Conf = new(ModelConf)

		f_n := "./config.yaml"
		yamlData, err := os.ReadFile(f_n)
		if err != nil {
			panic("配置文件路径有误")
		}

		err = yaml.Unmarshal(yamlData, Conf)
		if err != nil {
			panic("配置文件信息有误")
		}
	}
}
