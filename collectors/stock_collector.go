package collectors

import (
	"api.frank.top/stockInfo/global"
	"api.frank.top/stockInfo/xueqiu"
	"api.frank.top/stockInfo/global"
	"api.frank.top/stockInfo/xueqiu"
	"github.com/ShawnRong/tushare-go"
	"strings"
)


type StockCollector struct {
	Api xueqiu.XueqiuApi
	LoginStatus bool
}

type TuShareReq struct {
	ApiName    string `json:"api_name"`
	Token string `json:"token"`
	Params string `json:"params"`
	Fields string `json:"fields"`
}

func (sc *StockCollector) InitStockCurrentInfo()  {
	var totalCount int
	global.GVA_DB.Raw("SELECT count(*) FROM stock_base_info").Scan(&totalCount)
	totalPage := totalCount/300+1
	for i := 1; i <= totalPage; i++ {
		start := 300*(i-1)
		var tsCodes []string
		global.GVA_DB.Raw("SELECT ts_code FROM stock_base_info limit ?,?",start,300).Scan(&tsCodes)
		for j := 0 ;j< len(tsCodes);j++{
			sc.GetStockRealInfo(tsCodes[j])
		}
	}
}

func (sc *StockCollector) GetStockRealInfo(tsCode string)  {
	baseInfo := sc.Api.GetDailyStockBase(tsCode)
	insertSql :="insert into stock_daily_base("
	var values []string
	for k, v := range baseInfo {
		insertSql=insertSql+k+","
		values = append(values, v)
	}
	//去掉最后一个字符
	insertSql = strings.TrimRight(insertSql,",")
	insertSql = insertSql+") values (?)"
	global.GVA_DB.Exec(insertSql,values)
}

func (sc * StockCollector) InitStockHistoryDailyData()  {
	var totalCount int
	global.GVA_DB.Raw("SELECT count(*) FROM stock_base_info").Scan(&totalCount)
	totalPage := totalCount/300+1
	for i := 1; i <= totalPage; i++ {
		start := 300*(i-1)
		var tsCodes []string
		global.GVA_DB.Raw("SELECT ts_code FROM stock_base_info limit ?,?",start,300).Scan(&tsCodes)
		for j := 0 ;j< len(tsCodes);j++{
			sc.DailyData(tsCodes[j],"20220101")
		}
	}
}

func (sc * StockCollector) DailyData(tsCode string,startDate string )  {
	c := tushare.New("98dc5435aad747016e74b9365c4d2d736fa22cb7e2b553fae480edab")
	// 参数
	params := make(map[string]string)
	params["ts_code"] = tsCode
	params["start_date"] = startDate
	// 字段
	fields :=[] string {"ts_code",
		"trade_date",
		"open",
		"high",
		"low",
		"close",
		"pre_close",
		"change",
		"pct_chg",
		"vol",
		"amount"}
	// 根据api 请求对应的接口
	data, _ := c.Daily(params, fields)
	d := data.Data
	f := d.Fields
	items := d.Items

	insertSql :="insert into stock_daily("
	for i := 0; i < len(f); i++ {
		if string(f[i]) == "change" {
			f[i] = "`change`"
		}
		if i==len(f)-1 {
			insertSql=insertSql+f[i]+") values (?)"
		}else{
			insertSql=insertSql+f[i]+","
		}
	}
	for i := 0; i < len(items); i++ {
		global.GVA_DB.Exec(insertSql,items[i])
	}
}
/**
获取股票基础信息数据
 */
func (* StockCollector) InitBaseStockInfo()  {

	c := tushare.New("98dc5435aad747016e74b9365c4d2d736fa22cb7e2b553fae480edab")
	// 参数
	params := make(map[string]string)
	// 字段
	fields :=[] string {"ts_code","symbol","name","area","industry","market","list_date","list_status","delist_date","is_hs","curr_type"}
	// 根据api 请求对应的接口
	data, _ := c.StockBasic(params, fields)
	d := data.Data
	f := d.Fields
	items := d.Items

	insertSql :="insert into stock_base_info("
	for i := 0; i < len(f); i++ {
		if i==len(f)-1 {
			insertSql=insertSql+f[i]+") values (?)"
		}else{
			insertSql=insertSql+f[i]+","
		}
	}
	for i := 0; i < len(items); i++ {
		global.GVA_DB.Exec(insertSql,items[i])
	}


}