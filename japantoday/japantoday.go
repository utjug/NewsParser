package japantoday

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
		colly.AllowedDomains("japantoday.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("itemprop"), "articleBody") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}

func getheadline(s string) string{
	var head string

	collector := colly.NewCollector(
		colly.AllowedDomains("japantoday.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("itemprop"), "headline") {
			head = element.Text
		}
	})
	collector.Visit(s)

	return head
}

func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("done: japan")


	n := model.News{}
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")

	collector := colly.NewCollector(
		colly.AllowedDomains("japantoday.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "media-body") {
			n.URL=element.ChildAttr("a","href")
			n.Headline = getheadline(n.URL)
			n.Content = strings.ReplaceAll(getcontent(n.URL), "\"", " ")
			n.Content = strings.ReplaceAll(n.Content, ".", ". ")
			if len(n.Content) > 200{
				News = append(News, n)
			}
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Keyword: ", keyword, "\nVisiting: ", request.URL.String())
	})
	collector.Visit("https://japantoday.com/search?keyword=" + keyword)

	return News

}

