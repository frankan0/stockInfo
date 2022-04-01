package main

import (
	"api.frank.top/spider/collectors"
	"api.frank.top/spider/core"
	"api.frank.top/spider/global"
	"api.frank.top/spider/initialize"
	"go.uber.org/zap"
)


func main() {

	global.GVA_VP = core.Viper() // 初始化Viper
	global.GVA_LOG = core.Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	wallCollector := new(collectors.WallCollector)
	wallCollector.BingToday()
	//weiboCollector := new(collectors.WeiboCollector)
	//weiboCollector.Start()

}




