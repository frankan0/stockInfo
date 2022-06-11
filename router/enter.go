package router

type RouterGroup struct {
	Stock StockRouter
}

var RouterGroupApp = new(RouterGroup)

