package rt

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
		colly.AllowedDomains("www.rt.com","rt.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "article__heading") {
			content=element.Text
		}
	})
	collector.Visit(s)

	return content
}

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("www.rt.com","rt.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "article__text") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}

func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("\ndone: rt")


	n := model.News{}
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")

	collector := colly.NewCollector(
		colly.AllowedDomains("www.rt.com","rt.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "card__header") {
			n.URL = "https://www.rt.com"+element.ChildAttr("a", "href")
			n.Content = getcontent(n.URL)
			n.Content=strings.ReplaceAll(n.Content, "\"","'")
			n.Content=strings.ReplaceAll(n.Content, ".",". ")
			n.Content=strings.ReplaceAll(n.Content, "READ MORE: ","")
			n.Headline = getheadline (n.URL)
			n.Headline = strings.ReplaceAll(n.Headline, "\n\n        ", "")
			n.Headline = strings.ReplaceAll(n.Headline, "\n\n    ", "")

			newsstatus:=0
			if strings.Contains(n.URL, "show"){
				newsstatus = 1
			}
			if n.Content != "" && newsstatus == 0{
				News = append(News, n)
			}
		}
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println( "\nVisiting: ", request.URL.String())
	})

	collector.Visit("https://www.rt.com/search?q="+keyword)

	return News
}










