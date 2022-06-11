package stock

import (
	"api.frank.top/stockInfo/global"
	"api.frank.top/stockInfo/model/request"
	"api.frank.top/stockInfo/model/stock"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type StockService struct {

}


func (ss *StockService) QueryLatestDailyData(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&stock.Daily{})
	var stockDailyData []stock.Daily
	var tradeDate int
	global.GVA_DB.Raw("select trade_date from stock_daily where ts_code='000001.SZ' order by trade_date desc limit 1").Scan(&tradeDate)
	err = db.Where("trade_date=?",tradeDate).Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Where("trade_date=?",tradeDate).Order("pct_chg desc ").Find(&stockDailyData).Error
	return err, stockDailyData, total
}

func (*StockService) QueryBaseInfo(tsCode string) stock.BaseInfo{
	baseKey := fmt.Sprintf("stockBaseKey:%s",tsCode)
	result, err := global.GVA_REDIS.Get(context.Background(), baseKey).Result()
	if err != nil {
		//read from db
		baseInfo := queryBaseInfoFromDB(tsCode)
		marshal, _ := json.Marshal(baseInfo)
		err := global.GVA_REDIS.Set(context.Background(), baseKey, marshal, 0).Err()
		if err!=nil {
			global.GVA_LOG.Error("set baseinfo value to redis error",zap.Error(err))
		}
		return baseInfo
	}
	var data stock.BaseInfo
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal redis data error",zap.Error(err))
	}
	return data
}

func queryBaseInfoFromDB(code string) stock.BaseInfo{
	var baseInfo stock.BaseInfo
	err := global.GVA_DB.Where("ts_code = ?", code).First(&baseInfo).Error
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		global.GVA_LOG.Error("queryBaseInfoFromDB error ",zap.Error(err))
	}
	return baseInfo
}

func queryDailyInfoFromDB(code string,day int) stock.Daily{
	var baseInfo stock.Daily
	err := global.GVA_DB.Where("ts_code = ? and trade_date=?", code,day).First(&baseInfo).Error
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		global.GVA_LOG.Error("queryDailyInfoFromDB error ",zap.Error(err))
	}
	return baseInfo
}

func ( ss *StockService) QueryLatestDaily(tsCode string) stock.Daily{
	var tradeDate int
	global.GVA_DB.Raw("select trade_date from stock_daily where ts_code='000001.SZ' order by trade_date desc limit 1").Scan(&tradeDate)
	return ss.QueryDaily(tsCode, tradeDate)
}


func (*StockService) QueryDaily(tsCode string,tradeDate int) stock.Daily{
	dailyKey := fmt.Sprintf("stockDailyKey:%s:%s",tsCode,string(tradeDate))
	result, err := global.GVA_REDIS.Get(context.Background(), dailyKey).Result()
	if err != nil {
		//read from db
		baseInfo := queryDailyInfoFromDB(tsCode,tradeDate)
		marshal, _ := json.Marshal(baseInfo)
		err := global.GVA_REDIS.Set(context.Background(), dailyKey, marshal, 0).Err()
		if err!=nil {
			global.GVA_LOG.Error("set QueryDaily value to redis error",zap.Error(err))
		}
		return baseInfo
	}
	var data stock.Daily
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal redis data error",zap.Error(err))
	}
	return data
}



func (*StockService) QueryAmplifyVol(dayType int , multiple float64) []stock.AvgVol {
	var colName string
	if dayType == 1 {
		//3天
		colName="avg_three"
	}else{
		//5天
		colName="avg_five"
	}
	var stocks []stock.AvgVol
	err := global.GVA_DB.Raw("select * from stock_vol_avg where (current_vol/"+colName+") > ?", multiple).Scan(&stocks).Error
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		global.GVA_LOG.Error("QueryAmplifyVol error ",zap.Error(err))
	}
	return stocks
}


func (*StockService) ComputeAvgVolData() {
	db := global.GVA_DB.Model(&stock.BaseInfo{})
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return
	}
	var stockList []stock.BaseInfo
	totalPage := int((total + 300 - 1)/300)
	for i := 1; i <= totalPage; i++ {
		offset := 300 * (i - 1)
		err = db.Limit(300).Offset(offset).Find(&stockList).Error
		//计算该股票的平均成交量
		for si := range stockList {
			info := stockList[si]
			computeStockAvgVol(info.TsCode,5)
		}

	}
}

/**
  同时计算5天平均值和3天平均值
 */
func computeStockAvgVol(tsCode string,days int) {
	var dailyList []stock.Daily
	global.GVA_DB.Where("ts_code = ? order by trade_date desc limit ?", tsCode, days).Find(&dailyList)
	var sum3 float64
	var sum5 float64
	var currentVol float64
	var currentPrice float64
	var pctChg float64
	var tradeDate int
	for i := range dailyList {
		daily := dailyList[i]
		if i == 0 {
			currentVol = daily.Vol
			currentPrice = daily.Close
			pctChg = daily.PctChg
			tradeDate = daily.TradeDate
		}
		if i <= 2 {
			sum3 = sum3+daily.Vol
		}
		sum5 = sum5 + daily.Vol
	}
	var avg3 = sum3/3
	var avg5 = sum5/5

	var currentAvgVol stock.AvgVol
	err := global.GVA_DB.Where("ts_code = ? ", tsCode).Find(&currentAvgVol).Error
	if errors.Is(err, gorm.ErrRecordNotFound){
		//update
		currentAvgVol.AvgFive = avg5
		currentAvgVol.AvgThree = avg3
		currentAvgVol.DataTime = time.Now()
		currentAvgVol.CurrentPrice = currentPrice
		currentAvgVol.PctChg = pctChg
		currentAvgVol.TradeDate = tradeDate
		currentAvgVol.CurrentVol = currentVol
		global.GVA_DB.Save(&currentAvgVol)
	}else {
		//insert
		avgVol := stock.AvgVol{TsCode: tsCode, AvgThree: avg3, AvgFive: avg5,CurrentVol:currentVol,DataTime:time.Now(),CurrentPrice:currentPrice,PctChg: pctChg,TradeDate: tradeDate}
		global.GVA_DB.Create(&avgVol)
	}
}