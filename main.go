package main

import (
	"api.frank.top/stockInfo/core"
	"api.frank.top/stockInfo/global"
	"api.frank.top/stockInfo/initialize"
	"api.frank.top/stockInfo/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


func main() {

	global.GVA_VP = core.Viper() // 初始化Viper
	global.GVA_LOG = core.Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库

	//var collectorManagerApp = new(collectors.CollectorManager)
	//collectorManagerApp.Start()
	//init redis
	if global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}
	//拉取股票基础数据
	//stockCollector := new(collectors.StockCollector)
	//stockCollector.InitStockHistoryDailyData("20220617")
	//stockCollector.InitStockCurrentInfo()
	r := gin.Default()
	publicGroup := r.Group("")
	router.RouterGroupApp.Stock.InitStockRouter(publicGroup)
	address := fmt.Sprintf(":%d",global.GVA_CONFIG.System.Addr)
	r.Run(address)
	//stockService := new(stock.StockService)
	//stockService.ComputeAvgVolData()
}




