package xueqiu

import (
	"api.frank.top/stockInfo/global"
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type XueqiuApi struct{

}

type IpProxyItem struct {
	Proxy string `json:"proxy"`
	lastStatus bool `json:"last_status"`
}

func getProxyList() []string {
	var ipPool []string
	response, err := http.Get("http://127.0.0.1:5010/all/")
	if err!=nil {
		global.GVA_LOG.Error("获取代理IP池失败",zap.Error(err))
		return ipPool
	}
	defer response.Body.Close()
	body,err := ioutil.ReadAll(response.Body)
	if err !=nil {
		global.GVA_LOG.Error("获取代理IP池失败",zap.Error(err))
		return ipPool
	}
	if response.StatusCode == http.StatusOK {
		var ips []IpProxyItem
		json.Unmarshal(body,&ips)
		for i := range ips {
			ipPool = append(ipPool,"http://"+ips[i].Proxy)
		}
	}
	return ipPool
}

func (*XueqiuApi) GetDailyStockBase(tsCodeTuShare string) map[string]string {
	//code 转换
	split := strings.Split(tsCodeTuShare, ".")
	tsCode := split[1]+split[0]
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	ipProxyList := getProxyList()
	if len(ipProxyList) >0 {
		switcher, err := proxy.RoundRobinProxySwitcher(ipProxyList...)
		if err != nil {
			global.GVA_LOG.Error("查询IP代理列表, 异常", zap.Error(err))
		}
		c.SetProxyFunc(switcher)
		global.GVA_LOG.Info("设置代理池成功")
	}


	var info map[string]string

	header := map[string]string{
		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Connection": "keep-alive",
		"Host":       "weibo.com",
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
		"Cookie": "Hm_lvt_1db88642e346389874251b5a1eded6e3=1654146253; device_id=d2d394607153586456bea67ce764be22; s=dn17syzban; __utmc=1; __utmz=1.1654866151.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); bid=81cc8689f2e1100ca956bb57a4a0ff75_l497qoj7; __utma=1.751048509.1654866151.1654866151.1654911773.2; acw_tc=2760826216549192486076382edcec1bf7879e64a9eddba74047c365707eca; is_overseas=0; remember=1; xq_a_token=a5792a63aa26be86c8f5f81e811e32da880f29bc; xqat=a5792a63aa26be86c8f5f81e811e32da880f29bc; xq_id_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1aWQiOjc3MzcxOTgzODksImlzcyI6InVjIiwiZXhwIjoxNjU3NTAzODQzLCJjdG0iOjE2NTQ5MTk2NTI2NDgsImNpZCI6ImQ5ZDBuNEFadXAifQ.cBcZXWypTXQdEQnXWKlKXtWgxpsgVgfVQiwdcGmOwqjsGcwqHwqzFvp-IL8MSirfeGYZUKd3eIB_B6Q1nHj0Jv9LkoMMIJNvZu5yIpqYbWQFMno4j4APASxcDdQVE7l2AiPysPza-ygRyMr9hlOvoXXwvXlNkkz_RQjxYC8pDSd-vyhIR7Ds6F6z9zq1bkk4FndgYZ4vZZmdeD5RxGjPQ5SXht_NkCkG1MEdXnNoJ_y7Bk3bKRZJPTdz7XmEvc8WfXVNXGYEAC3jbOjyci_Qa2unZ9k0sTcqSZeQGpf5RxmRzDq0gMla0VcQsi9AjxCClyRfYDQvLhHapYy-bEw_cg; xq_r_token=5dcd67636bce76fd8df57a6c5284ebe71021fec7; xq_is_login=1; u=7737198389; Hm_lpvt_1db88642e346389874251b5a1eded6e3=1654919657",
	}
	// 在提出请求之前打印 "访问…"
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting: ", r.URL.String())
		for key, value := range header {
			r.Headers.Add(key, value)
		}

	})
	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		global.GVA_LOG.Error("Request URL: failed with response:",zap.String("url",r.Request.URL.String()),zap.Error(err))
	})
	c.OnResponse(func(response *colly.Response) {
		var result map[string]interface{}
		json.Unmarshal(response.Body, &result)
		data :=result[tsCode].(map[string]interface{})
		info = make(map[string]string)
		timestr := data["time"].(string)
		tt, _ := time.Parse(time.RubyDate, timestr)
		info["ts_code"] = tsCode
		info["trade_date"] = tt.Format("20060102")
		info["close"] = decimalFormString(data["close"].(string))
		info["turnover_rate"] = decimalFormString(strings.TrimRight(data["turnover_rate"].(string),"%"))
		pe:=data["pe_lyr"].(string)
		if pe == "" {
			info["pe"] = "0"
		}else{
			info["pe"] = pe
		}
		s := data["pe_ttm"].(string)
		if s == "" {
			info["pe_ttm"] = "0"
		}else{
			info["pe_ttm"] = s
		}
		pb := data["pb"].(string)
		if pb == "" {
			info["pb"] ="0"
		}else{
			info["pb"] = pb
		}
		info["total_mv"] = decimalFormString(data["marketCapital"].(string))
		info["circ_mv"] = decimalFormString(data["float_market_capital"].(string))
		info["total_share"] = decimalFormString(data["totalShares"].(string))
		info["float_share"] = decimalFormString(data["float_shares"].(string))
	})
	c.Visit("https://xueqiu.com/v4/stock/quote.json?code="+tsCode)
	return info
}

func decimalFormString(data string) string {
	fromString, err := decimal.NewFromString(data)

	if err != nil {
		global.GVA_LOG.Error("transfer Data error ",zap.String("data",data),zap.Error(err))
		return "0"
	}
	return fromString.String()
}

