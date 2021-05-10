package kremlin

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
		colly.AllowedDomains("www.en.kremlin.ru","en.kremlin.ru"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "width_limiter") {
			content=element.ChildText("h1")
		}
	})
	collector.Visit(s)

	return content
}

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("www.en.kremlin.ru","en.kremlin.ru"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "entry-content") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}

func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("\ndone: kremlin")


	n := model.News{}
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")

	collector := colly.NewCollector(
		colly.AllowedDomains("www.en.kremlin.ru","en.kremlin.ru"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "hentry__title") {
			n.URL = "http://www.en.kremlin.ru"+element.ChildAttr("a", "href")
			n.Content = getcontent(n.URL)
			n.Content = strings.ReplaceAll(n.Content, "\n", " ")
			n.Headline = getheadline (n.URL)
			if n.Content != "" && len(n.Content)<10000{
				News = append(News, n)
			}
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println( "\nVisiting: ", request.URL.String())
	})

	collector.Visit("http://www.en.kremlin.ru/search?query="+keyword)

	return News
}
