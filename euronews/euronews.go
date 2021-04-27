package euronews

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
		colly.AllowedDomains("www.euronews.com","euronews.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "article-title") {
			content=element.Text
			content = strings.ReplaceAll(content, "\n", "")
		}
	})
	collector.Visit(s)

	return content
}

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("www.euronews.com","euronews.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "js-article-content") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}

func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("done: euronews")


	n := model.News{}
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")

	collector := colly.NewCollector(
		colly.AllowedDomains("www.euronews.com","euronews.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "qa-article-title") {
			n.URL = "https://www.euronews.com"+element.ChildAttr("a", "href")
			n.Content=strings.ReplaceAll(getcontent(n.URL), "\"","'")
			n.Content=strings.ReplaceAll(n.Content, ".",". ")
			n.Headline=strings.ReplaceAll(getheadline (n.URL), "    ","")
			News = append(News, n)
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Keyword: ", keyword, "\nVisiting: ", request.URL.String())
	})

	collector.Visit("https://www.euronews.com/search?query="+keyword)

	return News
}

