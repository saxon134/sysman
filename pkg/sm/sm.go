package sm

import (
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saLog"
)

func Init() {
	//set error pkg
	saError.SetPkg("/go.yf/sysman/", "pkg/api/api.go", "pkg/api/middleware/")

	//init log
	saLog.Init(saLog.WarnLevel, saLog.ZapType)
	saLog.SetPkg("/go.yf/sysman/", "pkg/api/api.go", "pkg/api/middleware/")

	//init config
	initConf()

	//init database
	initDB()
}
