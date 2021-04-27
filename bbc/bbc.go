package bbc

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


func getheadline(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("bbc.co.uk", "www.bbc.co.uk"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("id"), "main-heading") {
			content=element.Text
		}
	})
	collector.Visit(s)

	return content
}

func getcontent(s string) string{
	var content string

	collector := colly.NewCollector(
		colly.AllowedDomains("bbc.co.uk", "www.bbc.co.uk"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "ArticleWrapper") {
			content=element.ChildText("p")
		}
	})
	collector.Visit(s)

	return content
}

func howmany() int{
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")
	collector:=colly.NewCollector(
		colly.AllowedDomains("bbc.co.uk", "www.bbc.co.uk"),
	)
	collector.OnHTML(`*`, func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("class"), "PageButtonContainer-StyledNumberedPageButton"){
			ch, err := strconv.Atoi(element.ChildText("a"))
			if err != nil {
				panic(err)
			}
			fmt.Printf("Scanned: %d\n", ch)
			pages=ch
		}
	})

	collector.Visit("https://www.bbc.co.uk/search?q="+keyword)

	fmt.Println(pages)

	return pages
}

func GetNews(wg *sync.WaitGroup) []model.News{
	defer wg.Done()
	defer fmt.Println("done: bbc")
	var (
		i int
		p string
	)

	n:=model.News{}
	howmany()
	keyword:=strings.ReplaceAll(key.Keyword, " ", "+")

	//for i=1; i < pages; i++ {
	for i=1; i < 4; i++ {

		p=strconv.Itoa(i)

		collector:=colly.NewCollector(
			colly.AllowedDomains("bbc.co.uk", "www.bbc.co.uk"),
		)
		collector.OnHTML(`*`, func(element *colly.HTMLElement) {

			if strings.Contains(element.Attr("class"),"PromoContent") {
				n.URL = element.ChildAttr("a", "href")
				if strings.Contains(n.URL, "news"){
					n.Headline = getheadline(n.URL)
					n.Content = strings.ReplaceAll(getcontent(n.URL), "\"", "'")
					if len(n.Content)>0{
						News = append(News, n)
					}
				}
			}
		})

		collector.OnRequest(func(request *colly.Request) {
			fmt.Println("Keyword: ", keyword,"\nVisiting: ", request.URL.String())
		})

		collector.Visit("https://www.bbc.co.uk/search?q=" + keyword + "&page=" + p)
	}
	return News
}