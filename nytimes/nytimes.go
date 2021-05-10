package nytimes

import (
	"example.com/parser/key"
	"example.com/parser/model"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"sync"
)
var News []model.News

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("nytimes.com", "www.nytimes.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "meteredContent") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}

func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("\ndone: ny")

	n := model.News{}
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")
	collector := colly.NewCollector(
		colly.AllowedDomains("nytimes.com", "www.nytimes.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("data-testid"), "search-bodega-result") {
			n.Headline = element.ChildText("h4")
			n.URL="https://www.nytimes.com"+element.ChildAttr("a","href")
			n.Content=getcontent(n.URL)
			n.Content=strings.ReplaceAll(n.Content, ".", ". ")
			if len(n.Content)>400{
				News = append(News, n)
			}
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("\nVisiting: ", request.URL.String())
	})
	collector.Visit("https://www.nytimes.com/search?dropmab=true&query=" + keyword + "&sort=best")


	return News

}