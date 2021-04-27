package india

import (
	"example.com/parser/key"
	"example.com/parser/model"
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
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
		colly.AllowedDomains("timesofindia.indiatimes.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "ga-headlines") {
			content=element.Text
		}
	})
	collector.Visit(s)

	return content
}

func howmany() int{

	collector:=colly.NewCollector(
		colly.AllowedDomains("timesofindia.indiatimes.com"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "title mb-4"){
			ch, err := strconv.Atoi(element.ChildText("a"))
			if err != nil {
				panic(err)
			}
			fmt.Printf("Scanned: %d\n", ch)
			pages=ch
		}
	})

	collector.Visit("https://timesofindia.indiatimes.com/topic/"+key.Key+"/")

	fmt.Println(pages)

	return pages
}

func GetNews(wg *sync.WaitGroup) []model.News {
	defer wg.Done()
	defer fmt.Println("done: india")
	var (
		i int
		p string
	)

	n := model.News{}
	//howmany()

	//for i=1; i < pages; i++ {
	for i = 1; i < 2; i++ {

		keyword:=strings.ReplaceAll(key.Keyword, " ", "-")
		p = strconv.Itoa(i)

		collector := colly.NewCollector(
			colly.AllowedDomains("timesofindia.indiatimes.com"),
		)
		collector.OnHTML(`*`, func(element *colly.HTMLElement) {
			if strings.Contains(element.Attr("class"), "content") {
				n.Headline = strings.ReplaceAll(element.ChildText("span"), "\n", " ")
			}
			if strings.Contains(element.Attr("class"), "content") {
				if len(n.Headline) < 400 && len(n.Headline) > 10{
					n.URL = "https://timesofindia.indiatimes.com/"+element.ChildAttr("a", "href")
					n.Content = strings.ReplaceAll(getcontent(n.URL), "\"", "'")
					if strings.Contains(n.Content, "function"){
						n.Content=""
					}
					if len(n.Content)>0 {
						News = append(News, n)
					}
				}
			}
		})

		collector.OnRequest(func(request *colly.Request) {
			fmt.Println("Keyword: ", keyword, "\nVisiting: ", request.URL.String())
		})

		collector.Visit("https://timesofindia.indiatimes.com/topic/"+keyword+"/"+p)
	}
	return News
}