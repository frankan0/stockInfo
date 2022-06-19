package stock

import (
	"api.frank.top/stockInfo/model/request"
	"api.frank.top/stockInfo/model/response"
	"api.frank.top/stockInfo/model/stock"
	"api.frank.top/stockInfo/global"
	"api.frank.top/stockInfo/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type StockApi struct {

}

func (s *StockApi) ListDaily(c *gin.Context)  {
	var pageInfo request.PageInfo
	if c.ShouldBind(&pageInfo) != nil{
		response.FailWithMessage("获取参数失败",c)
		return
	}
	if err, list, total := service.ServiceGroupApp.StockService.QueryLatestDailyData(pageInfo);err!=nil{
		global.GVA_LOG.Error("QueryLatestDailyData!", zap.Error(err))
		response.FailWithMessage("获取日线行情失败", c)
	}else{
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}


func (s *StockApi) QueryAmplifyVol(c *gin.Context)  {
	dayType := c.Query("dayType")
	muti := c.Query("multiple")
	multiple, _ := strconv.ParseFloat(muti,  64)
	data := service.ServiceGroupApp.StockService.QueryAmplifyVol(dayType, multiple)
	var stockLists []StockDailyInfo
	for i := range data {
		var dayInfo StockDailyInfo
		dayInfo.Avg = data[i]
		daily := service.ServiceGroupApp.StockService.QueryLatestDaily(data[i].TsCode)
		dayInfo.Daily = daily
		dayInfo.BaseInfo = service.ServiceGroupApp.StockService.QueryBaseInfo(data[i].TsCode)
		dayInfo.DailyBase = service.ServiceGroupApp.StockService.QueryLatestDailyBase(data[i].TsCode)
		stockLists = append(stockLists, dayInfo)
	}
	response.OkWithData(stockLists,c)
}

type StockDailyInfo struct {
	Daily stock.Daily `json:"daily"`
	BaseInfo stock.BaseInfo `json:"base_info"`
	Avg stock.AvgVol `json:"avg"`
	DailyBase stock.DailyBase `json:"daily_base"`
}