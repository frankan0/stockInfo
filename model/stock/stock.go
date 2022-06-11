package stock

import "time"

type AvgVol struct {
	ID int64
	TsCode string `json:"ts_code"`
	AvgThree float64 `json:"avg_three"`
	AvgFive float64 `json:"avg_five"`
	CurrentVol float64 `json:"current_vol"`
	CurrentPrice float64 `json:"current_price"`
	PctChg float64 `json:"pct_chg"`
	DataTime time.Time `json:"data_time"`
	TradeDate int `json:"trade_date"`
}
func (AvgVol) TableName() string {
	return "stock_vol_avg"
}

type BaseInfo struct {
	ID int64 `json:"id"`
	TsCode string `json:"ts_code"`
	Symbol string `json:"symbol"`
	Name string `json:"name"`
	Area string `json:"area"`
	ListStatus string `json:"list_status"`
}

func (BaseInfo) TableName() string {
	return "stock_base_info"
}

type Daily struct {
	ID int64 `json:"id"`
	TsCode string `json:"ts_code"`
	TradeDate int `json:"trade_date"`
	Open float64 `json:"open"`
	High float64 `json:"high"`
	Low float64 `json:"low"`
	Close float64 `json:"close"`
	PreClose float64 `json:"pre_close"`
	Vol float64 `json:"vol"`
	Amount float64 `json:"amount"`
	Change float64 `json:"change"`
	PctChg float64 `json:"pct_chg"`
	TurnoverRate float64 `json:"turnover_rate"`
	PeTtm float64 `json:"pe_ttm"`
	TotalMv float64 `json:"total_mv"`
}
func (Daily) TableName() string {
	return "stock_daily"
}
