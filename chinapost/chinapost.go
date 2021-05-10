package chinapost

import (
	"example.com/parser/key"
	"example.com/parser/model"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"sync"
)

var (
	News []model.News
	pages int
	ch int

)

func getheadline(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("china.timesofnews.com"),
	)
	collector.OnHTML("h1", func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "heading") {
			content=element.Text
		}
	})
	collector.Visit(s)

	return content
}

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("china.timesofnews.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "blogtext") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}

func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("\ndone: china")


	n := model.News{}
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")

	collector := colly.NewCollector(
		colly.AllowedDomains("china.timesofnews.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "text") {
			if len(element.ChildAttr("a", "href"))>70 {
				n.URL = element.ChildAttr("a", "href")
				n.Content = strings.ReplaceAll(getcontent(n.URL), "    ", " ")
				n.Content = strings.ReplaceAll(n.Content, ".", ". ")
				n.Content = strings.ReplaceAll(n.Content, "Download our app or subscribe to our Telegram channel for the latest updates on the coronavirus outbreak: https://cna.asia/telegramArticle source:", "")
				n.Headline = getheadline (n.URL)
				if n.Content != "" && string([]rune(n.Content)[14]) != ":"{
					News = append(News, n)
				}
			}
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println( "\nVisiting: ", request.URL.String())
	})

	collector.Visit("https://china.timesofnews.com/?s="+keyword)

	return News
}


