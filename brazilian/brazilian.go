package brazilian

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

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("riotimesonline.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "td-ss-main-content") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}


func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("\ndone: braz")
	var (
		i int
		//p string
	)

	n := model.News{}

	for i = 1; i < 2; i++ {

		//p = strconv.Itoa(i)
		keyword:=strings.ReplaceAll(key.Keyword, " ", "+")

		collector := colly.NewCollector(
			colly.AllowedDomains("riotimesonline.com"),
		)
		collector.OnHTML(`*`, func(element *colly.HTMLElement) {
			if strings.Contains(element.Attr("class"), "entry-title") {
					n.Headline = element.ChildAttr("a", "title")
			}
			if strings.Contains(element.Attr("class"), "entry-title") {
				n.URL = element.ChildAttr("a", "href")
				n.Content = strings.ReplaceAll(getcontent(n.URL), ". . .To read the full NEWS and much more, Subscribe to our Premium Membership Plan. Already Subscribed? Login Here", "")
				n.Content = strings.ReplaceAll(n.Content, "\"", "'")
				n.Content = strings.ReplaceAll(n.Content, ".", ". ")
				n.Content+="."
				if len(n.Content)>0{
					News = append(News, n)
				}

			}
		})

		collector.OnRequest(func(request *colly.Request) {
			fmt.Println("\nVisiting: ", request.URL.String())
		})

		collector.Visit("https://riotimesonline.com/?s="+keyword)
	}
	return News
}
