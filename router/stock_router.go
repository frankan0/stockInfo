package router

import (
	stock "api.frank.top/stockInfo/api/v1/stock"
	"github.com/gin-gonic/gin"
)

type StockRouter struct {

}

func (s *StockRouter) InitStockRouter(Router *gin.RouterGroup)  {

	stockRouter := Router.Group("stock")
	stockApi := new(stock.StockApi)
	{
		stockRouter.GET("/queryAmplifyVol", stockApi.QueryAmplifyVol)
		stockRouter.GET("/listDaily", stockApi.ListDaily)
	}
}
