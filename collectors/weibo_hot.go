package collectors

import (
	"api.frank.top/spider/global"
	"api.frank.top/spider/model"
	"api.frank.top/spider/service"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
	"time"
)

type WeiboCollector struct {

}

func (w *WeiboCollector) Start()  {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	timeStr := time.Now().Format("20060102")
	hots := []model.HotRank{}
	c.OnHTML("table tbody tr", func(element *colly.HTMLElement) {
		hot := model.HotRank{}
		hot.Link = element.ChildAttr("a","href")
		hot.Title = element.ChildText("a")
		hot.HotValue = element.ChildText("span")
		hot.Day = timeStr
		hot.Source = "weibo_hot_search"
		hot.CreateTime = time.Now()
		hots = append(hots,hot)
	})

	header := map[string]string{
		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Connection": "keep-alive",
		"Host":       "weibo.com",
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
		"Cookie": "SINAGLOBAL=1367227292248.0588.1635751267191; _ga=GA1.2.365741.1635930085; ALF=1678496955; SUB=_2AkMVYfTOf8NxqwJRmfscyG7qa4h3yw_EieKjPQUVJRMxHRl-yT9jqnQztRB6PuHaITFKD9SV888LudV9Tm0JQA9FMsvb; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9WFXAjC_4Ovxf7wKobT37.C6; UOR=,,tophub.today; _s_tentry=weibo.com; Apache=5978685592805.342.1648787245134; ULV=1648787245139:4:1:1:5978685592805.342.1648787245134:1646959471759; WBtopGlobal_register_version=2022040112",
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
		global.GVA_LOG.Info("Request URL: failed with response:",zap.String("url",r.Request.URL.String()))
	})
	c.Visit("https://s.weibo.com/top/summary/summary")
	if len(hots) >0 {
		err :=service.ServiceGroupApp.WallServiceGroup.BatchAddWeibo(hots)
		if err!=nil {
			global.GVA_LOG.Error("insert batch hots error.",zap.Error(err))
		}
	}

}
