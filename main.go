package main

import (
	"api.frank.top/spider/core"
	"api.frank.top/spider/global"
	"api.frank.top/spider/init"
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

type Wall struct {
	Title string
	Link string
	Date string
	PicName string
	//高清 https://cn.bing.com/th?id=OHR.ThousandSprings_ZH-CN7431323804_UHD.jpg
}

func main() {


	global.GVA_VP = core.Viper() // 初始化Viper
	global.GVA_LOG = core.Zap()  // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = init.Gorm() // gorm连接数据库

	walls := make([]Wall,100)
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

	c.OnHTML(".item div.card", func(element *colly.HTMLElement) {
		wall := Wall{}
		wall.Link = element.ChildAttr("img","src")
		wall.Title = element.ChildText("div.description h3")
		wall.Date = element.ChildText("div.description .calendar em.t")
		//del head
		httpArr := strings.Split(wall.Link,"/")
		//去掉后面的
		picNameArr :=strings.Split(httpArr[len(httpArr)-1],"_")
		wall.PicName = picNameArr[0]+"_"+picNameArr[1]
		walls = append(walls,wall)

	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// On every span tag with the class next-button
	c.OnHTML("div.page", func(h *colly.HTMLElement) {
		t := h.ChildAttr("a:last-child", "href")
		c.Visit("https://bing.ioliu.cn"+t)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://bing.ioliu.cn/?p=1")

	csvSave("walls1.csv", walls)

	fmt.Println(c)
}

func saveData() {

}

// 数据持久化
func csvSave(fName string, data []Wall) error {
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Date", "Title", "Link", "PicName"})
	for _, v := range data {
		writer.Write([]string{v.Date, v.Title, v.Link, v.PicName})
	}
	return nil
}



