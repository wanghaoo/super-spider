package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
 for true {
	 crawlingZhihuRecommend()
	 time.Sleep(10 * time.Second)
 }
}

func crawlingZhihuRecommend() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.zhihu.com"),
	)
	c.OnHTML("div.Topstory-recommend", func(element *colly.HTMLElement) {
		element.ForEach("div.Card.TopstoryItem.TopstoryItem-isRecommend", func(i int, element *colly.HTMLElement) {
			title := element.ChildText("h2.ContentItem-title a")
			answerUrl := element.ChildAttr("h2.ContentItem-title a", "href")
			if len(title) <= 0 {
				return
			}
			up := element.ChildText("button.Button.VoteButton.VoteButton--up")
			up = strings.Replace(up, "赞同 ", "", -1)
			up = strings.Replace(up, "赞同", "", -1)
			up = strings.Replace(up, "\u200b", "", -1)
			var upNum int
			var err error
			if strings.Contains(up, "K") {
				up = strings.Replace(up, "K", "", -1)
				upNumFloat, err2 := strconv.ParseFloat(up, 10)
				if err2 != nil {
					return
				}
				upNum = int(upNumFloat * 1000)
			} else {
				upNum, err = strconv.Atoi(up)
			}
			if err != nil {
				return
			}
			if upNum < 1000 {
				return
			}
			fmt.Println(title)
			fmt.Println(element.Request.AbsoluteURL(answerUrl))
			fmt.Println(upNum)
		})
	})
	c.OnHTML("div.Topstory-content", func(element *colly.HTMLElement) {
		element.ForEach("div.Card.TopstoryItem", func(i int, element *colly.HTMLElement) {
			title := element.ChildText("h2.ContentItem-title a")
			answerUrl := element.ChildAttr("h2.ContentItem-title a", "href")
			fmt.Println(title)
			if len(title) <= 0 {
				return
			}
			up := element.ChildText("button.Button.VoteButton.VoteButton--up")
			up = strings.Replace(up, "赞同 ", "", -1)
			up = strings.Replace(up, "赞同", "", -1)
			up = strings.Replace(up, "\u200b", "", -1)
			var upNum int
			var err error
			if strings.Contains(up, "K") {
				up = strings.Replace(up, "K", "", -1)
				upNumFloat, err2 := strconv.ParseFloat(up, 10)
				if err2 != nil {
					return
				}
				upNum = int(upNumFloat * 1000)
			} else {
				upNum, err = strconv.Atoi(up)
			}
			if err != nil {
				return
			}
			if upNum < 10 {
				return
			}
			fmt.Println("content:", title)
			fmt.Println(element.Request.AbsoluteURL(answerUrl))
			fmt.Println(upNum)
		})
	})
	c.OnError(func(response *colly.Response, e error) {
		fmt.Println(string(response.Body))
		fmt.Println(response.StatusCode)
	})
	c.OnResponse(func(response *colly.Response) {
		fmt.Println(string(response.Body))
		fmt.Println(response.StatusCode)
	})
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	})
	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, &http.Cookie{Name:"_xsrf", Value:"b0902632-32c5-4557-8ee0-5f224c14fce1"})
	cookies = append(cookies, &http.Cookie{Name:"_zap", Value:"3d3c3c78-7525-43ef-a98b-b68ea5f51c2d"})
	cookies = append(cookies, &http.Cookie{Name:"capsion_ticket", Value:"2|1:0|10:1544582582|14:capsion_ticket|44:NDM0NDczMDEwMjZlNDUyNGIwOGJjNmMwMDVjNjBkNjA=|43d2df9291c143b6e02a90d1063d3e5965adcdedec3a3a313619450971054941"})
	cookies = append(cookies, &http.Cookie{Name:"d_c0", Value:"AECn7SD5RA6PToN9Pz7GTNRmo6cetXDStmw=|1537926878"})
	cookies = append(cookies, &http.Cookie{Name:"q_c1", Value:"3d8626f347dd439c8b89a0f5be57947f|1537932724000|1537932724000"})
	cookies = append(cookies, &http.Cookie{Name:"tgw_l7_route", Value:"23ddf1acd85bb5988efef95d7382daa0"})
	cookies = append(cookies, &http.Cookie{Name:"tst", Value:"r"})
	cookies = append(cookies, &http.Cookie{Name:"z_c0", Value:"2|1:0|10:1544582610|4:z_c0|92:Mi4xSUc4RkFBQUFBQUFBUUtmdElQbEVEaVlBQUFCZ0FsVk4wc1A5WEFCbzdjSEVsSFdWcFFyZERkU19IMFFlS0RndGlB|0e837215c42ba137f722b6cb0ea8f4b2ee97afac7d56d6f760a5d38c05041e97"})
	c.SetCookies("https://www.zhihu.com", cookies)
	c.Visit("https://www.zhihu.com/")
	c.Visit("https://www.zhihu.com/follow")
}
