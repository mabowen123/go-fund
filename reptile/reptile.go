package reptile

import (
	"encoding/json"
	"fmt"
	"fund/mysql"
	"fund/xRedis"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var wg sync.WaitGroup

func requestEstimate(code string) {
	response, _ := http.Get(fmt.Sprintf("https://fundgz.1234567.com.cn/js/%v.js?rt=%v", code, time.Now().Unix()))
	robots, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	r, _ := regexp.Compile("\\{+.*\\}")
	jsonData := []byte(r.FindString(string(robots)))
	var data map[string]interface{}
	json.Unmarshal(jsonData, &data)
	xRedis.HmSet(xRedis.FundName(code), data, xRedis.Ttl)
	defer wg.Done()
}

func requestActual(code string) {
	defer wg.Done()

	if time.Now().Hour() < 19 {
		return
	}

	data := xRedis.HGetAll(xRedis.FundName(code))
	rActualDate, _ := time.Parse("2006-01-02", data["actual_date"])
	if (time.Now().Day() == rActualDate.Day()) &&
		(time.Now().Month() == rActualDate.Month() &&
			(time.Now().Year() == rActualDate.Year())) {
		return
	}

	res, _ := http.Get(fmt.Sprintf("http://fund.eastmoney.com/%v.html", code))
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	date := doc.Find(".dataItem02").Find("p").Text()
	date = date[14:24]
	actualDate, _ := time.Parse("2006-01-02", date)
	if actualDate.Day() == time.Now().Day() {
		node := doc.Find(".dataItem02").Find(".dataNums").Children()
		actualGsz := node.Eq(0).Text()
		actualGszl := node.Eq(1).Text()
		data := map[string]interface{}{
			"actual_gsz":  actualGsz,
			"actual_gszl": actualGszl,
			"actual_date": date,
		}

		xRedis.HmSet(xRedis.FundName(code), data, xRedis.Ttl)
	}
}

func Run(needLock bool) {
	isLock := xRedis.Get(xRedis.FundLock())
	if needLock && (isLock == "true") {
		return
	}

	fundIds := new([]string)
	mysql.Db.Model(&mysql.UserFund{}).Pluck("distinct fund_id", fundIds)
	for _, code := range *fundIds {
		wg.Add(2)
		go requestActual(code)
		go requestEstimate(code)
	}
	wg.Wait()
	xRedis.Set(xRedis.FundLock(), "true", time.Minute)
}
