package canada

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
	)

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("globalnews.ca"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("aria-label"), "Article text") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}


func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("\ndone: canada")
	var (
		//i int
	)

	n := model.News{}
	keyword:=strings.ReplaceAll(key.Keyword, " ", "%20")

	//for i=1; i < pages; i++ {
	//for i = 1; i < 2; i++ {


		collector := colly.NewCollector(
			colly.AllowedDomains("globalnews.ca"),
		)
		collector.OnHTML(`*`, func(element *colly.HTMLElement) {
			if strings.Contains(element.Attr("class"), "story-h") {
				n.Headline = element.ChildText("a")
			}
			if strings.Contains(element.Attr("class"), "story-h") {
				n.URL = element.ChildAttr("a", "href")
				n.Content=strings.ReplaceAll(getcontent(n.URL), "\n", "")
				n.Content=strings.ReplaceAll(n.Content, "\t", " ")
				n.Content=strings.ReplaceAll(n.Content, "   Read more:     ", " ")
				n.Content=strings.ReplaceAll(n.Content, "     ", ".")
				n.Content=strings.ReplaceAll(n.Content, "  READ MORE:.", "")
				n.Content=strings.ReplaceAll(n.Content, ".", ". ")
				n.Content=strings.ReplaceAll(n.Content, "U. S.", "US")
				if len(n.Content)>0{
					News = append(News, n)
				}
			}
		})

		collector.OnRequest(func(request *colly.Request) {
			fmt.Println( "\nVisiting: ", request.URL.String())
		})

		collector.Visit("https://globalnews.ca/?s="+keyword)
	//}
	return News
}
