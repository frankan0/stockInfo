package xueqiu

import (
	"api.frank.top/stockInfo/global"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type TencentApi struct{

}

func (*TencentApi) GetDailyStockBase(tsCodeTuShare string) map[string]string {
	//code 转换
	split := strings.Split(tsCodeTuShare, ".")
	tsCode := strings.ToLower(split[1]+split[0])

	url :="http://qt.gtimg.cn/q="+tsCode
	req, _ := http.NewRequest("GET", url, nil)
	// 比如说设置个token
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	resp, err := (&http.Client{}).Do(req)
	//resp, err := http.Get(serviceUrl + "/topic/query/false/lsj")
	if err != nil {
		global.GVA_LOG.Error("查询腾讯接口失败",zap.Error(err))
		
	}
	var info map[string]string
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)
	data := string(respByte)
	fmt.Print(data)
	stockArr := strings.Split(data, ";")
		stockInfoArr := stockArr[0]
		stockInfoData := strings.Split(stockInfoArr, "=")
		stockInfo := stockInfoData[1]
		stockInfo = strings.Trim(stockInfo, "\"")
		//parseData
		stockInfoDetail := strings.Split(stockInfo, "~")
		info = make(map[string]string)
		info["ts_code"] = tsCodeTuShare
		info["close"] =  stockInfoDetail[4]
		info["turnover_rate"] =  stockInfoDetail[38]
		info["pe"] =  stockInfoDetail[52]
		info["pe_ttm"] =  stockInfoDetail[39]
		info["ps"] =  stockInfoDetail[46]
		info["total_mv"] =  stockInfoDetail[45]
		info["circ_mv"] =  stockInfoDetail[44]
		info["volume_ratio"] =  stockInfoDetail[49]

		timestr := stockInfoDetail[30]
		tt, _ := time.Parse("20060102150405", timestr)
		info["trade_date"] =  tt.Format("20060102")
		info["total_share"] =  stockInfoDetail[73]
		info["float_share"] =  stockInfoDetail[72]

	return info
}

