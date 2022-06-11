package service

import "api.frank.top/stockInfo/service/stock"

type ServiceGroup struct {
	StockService stock.StockService
}

var ServiceGroupApp = new(ServiceGroup)
