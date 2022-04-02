package collectors

import (
	"api.frank.top/spider/global"
	"api.frank.top/spider/model"
	"api.frank.top/spider/service"
	"encoding/json"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

type WallCollector struct {

}

func (w *WallCollector) BingToday() {
	resp, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN")
	if err != nil {
		global.GVA_LOG.Error("获取bing今日壁纸接口失败")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		global.GVA_LOG.Error("获取bing今日壁纸接口失败")
		return
	}
	var bingWalls model.BingResponse
	err = json.Unmarshal(body, &bingWalls)
	if err!=nil {
		global.GVA_LOG.Error("获取bing今日壁纸接口失败")
		return
	}
	bingWall := bingWalls.Images[0]
	//del head
	nameArr := strings.Split(bingWall.Urlbase,".")
	//入库
	wall := model.Wall{}
	wall.PicName= nameArr[1]
	wall.Title = bingWall.Title
	wall.Date = bingWall.Enddate
	wall.Link = bingWall.Url
	wall.Channel = "bing"
	service.ServiceGroupApp.WallServiceGroup.AddWall(wall)
}

func (w *WallCollector) Start()  {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: reddit.com
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.AllowedDomains("bing.ioliu.cn"),
	)

	//c.Limit(&colly.LimitRule{
	//	DomainGlob:  "*ioliu.*",
	//	Parallelism: 2,
	//	RandomDelay: 5 * time.Second,
	//})
	walls := []model.Wall{}
	c.OnHTML(".item div.card", func(element *colly.HTMLElement) {
		wall := model.Wall{}
		wall.Link = element.ChildAttr("img","src")
		wall.Title = element.ChildText("div.description h3")
		wall.Date = element.ChildText("div.description .calendar em.t")
		//del head
		httpArr := strings.Split(wall.Link,"/")
		//去掉后面的
		picNameArr :=strings.Split(httpArr[len(httpArr)-1],"_")
		wall.PicName = picNameArr[0]+"_"+picNameArr[1]
		if wall.PicName == "" {
			return
		}
		walls = append(walls,wall)
		if len(walls) >=200 {
			saveData2db(walls)
			walls = walls[:0]
		}

	})


	// On every span tag with the class next-button
	c.OnHTML("div.page", func(h *colly.HTMLElement) {
		t := h.ChildAttr("a:last-child", "href")
		c.Visit("https://bing.ioliu.cn"+t)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		global.GVA_LOG.Info("Request URL: failed with response:",zap.String("url",r.Request.URL.String()))
	})

	c.Visit("https://bing.ioliu.cn/?p=1")

	//csvSave("walls1.csv", walls)
	if len(walls)>0 {
		saveData2db(walls)
	}
}


func saveData2db(walls []model.Wall) {
	err :=service.ServiceGroupApp.WallServiceGroup.BatchAdd(walls)
	if err!=nil {
		global.GVA_LOG.Error("insert batch wall error.")
	}
}
