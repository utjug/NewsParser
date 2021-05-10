package fr

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
		colly.AllowedDomains("www.connexionfrance.com","connexionfrance.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("id"), "reactRenderer") {
			content=element.ChildText("h4")
		}
	})
	collector.Visit(s)

	return content
}

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("www.connexionfrance.com","connexionfrance.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("id"), "article") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}

func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("\ndone: france")


	n := model.News{}
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")

	collector := colly.NewCollector(
		colly.AllowedDomains("www.connexionfrance.com","connexionfrance.com"),
		)
		collector.OnHTML(`*`, func(element *colly.HTMLElement) {
			if strings.Contains(element.Attr("class"), "article-teaser article-teaser--separated") {
				n.URL = "https://www.connexionfrance.com"+element.ChildAttr("a", "href")
				n.Content = getcontent(n.URL)
				n.Content=strings.ReplaceAll(n.Content, "\"","'")
				n.Content=strings.ReplaceAll(n.Content, ".",". ")
				n.Content=strings.ReplaceAll(n.Content, "Read more: ", "")
				n.Headline = getheadline (n.URL)
				if len(n.Content)>0{
					News = append(News, n)
				}
			}
		})

		collector.OnRequest(func(request *colly.Request) {
			fmt.Println("\nVisiting: ", request.URL.String())
		})

		collector.Visit("https://www.connexionfrance.com/search?query="+keyword)

	return News
}
