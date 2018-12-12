package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
)

var domain = "https://www.qiushibaike.com"

type QSBKVideo struct {
	ArticleUrl string
	Title      string
	Laugh      int
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.qiushibaike.com"),
	)
	var rs = make([]*regexp.Regexp, 0)
	regx := regexp.MustCompile(`^.*?/article/[\d]+$`)
	regx2 := regexp.MustCompile(`^.*?/8hr/page/[\d]+/$`)
	regx3 := regexp.MustCompile(`^.*?/(hot|imgrank|text|history|pic|textnew)/$`)
	rs = append(rs, regx)
	rs = append(rs, regexp.MustCompile(`^https://www.qiushibaike.com/$`))
	rs = append(rs, regx2)
	rs = append(rs, regx3)
	c.URLFilters = rs
	c.Async = false
	c.OnHTML("div.index-head", func(element *colly.HTMLElement) {
		element.ForEach("li", func(i int, element *colly.HTMLElement) {
			href := element.ChildAttr("a", "href")
			if regx3.MatchString(href) {
				c.Visit(element.Request.AbsoluteURL(href))
			}
		})
	})
	c.OnHTML("li.item", func(element *colly.HTMLElement) {
		element.ForEach("div.recmd-right", func(i int, element *colly.HTMLElement) {
			video := new(QSBKVideo)
			href := element.ChildAttr("a.recmd-content", "href")
			video.ArticleUrl = element.Request.AbsoluteURL(href)
			element.ForEach("a.recmd-content", func(i int, element *colly.HTMLElement) {
				video.Title = element.Text
			})
			element.ForEach("div.recmd-num span:nth-child(1)", func(i int, element *colly.HTMLElement) {
				laugh, err := strconv.Atoi(element.Text)
				if err != nil {
					video.Laugh = 0
				} else {
					video.Laugh = laugh
				}
			})
			if video.Laugh < 1000 {
				return
			}
			fmt.Println("artilce: ", video)
		})
	})
	c.OnHTML("div#content-left", func(element *colly.HTMLElement) {
		element.ForEach("div.article", func(i int, element *colly.HTMLElement) {
			element.ForEach("div.author", func(i int, element *colly.HTMLElement) {
				articleAuthorImg := element.ChildAttr("img", "src")
				articleAuthor := element.ChildAttr("img", "alt")
				fmt.Println("article author image : ", articleAuthorImg)
				fmt.Println("article author : ", articleAuthor)
			})
			articleContent := element.ChildText("div.content")
			fmt.Println("article content: ", articleContent)
			laugh := element.ChildText("span.stats-vote i.number:nth-child(1)")
			fmt.Println("article laugh: ", laugh)
		})
	})
	c.OnHTML("ul.pagination", func(element *colly.HTMLElement) {
		element.ForEach("a[href]", func(i int, element *colly.HTMLElement) {
			pageHref := element.Attr("href")
			c.Visit(element.Request.AbsoluteURL(pageHref))
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting ", r.URL.String())
	})
	c.Visit("https://www.qiushibaike.com/")
}
