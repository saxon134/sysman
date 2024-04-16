package smclient

type Config struct {
	Redis struct {
		Host string `json:"host"`
		Pass string `json:"pass"`
		Db   int    `json:"db"`
	} `json:"redis"`
	MySql struct {
		Host string `json:"host"`
		User string `json:"user"`
		Pass string `json:"pass"`
		Db   string `json:"db"`
	} `json:"mySql"`
	Sysman struct {
		Root   string
		Secret string
	}
}
