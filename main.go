package main

import (
	"bytes"
	"encoding/json"
	"example.com/parser/bbc"
	"example.com/parser/brazilian"
	"example.com/parser/canada"
	"example.com/parser/chinapost"
	"example.com/parser/euronews"
	"example.com/parser/fr"
	"example.com/parser/india"
	"example.com/parser/japantoday"
	"example.com/parser/key"
	"example.com/parser/model"
	"example.com/parser/nytimes"
	"example.com/parser/rt"
	"example.com/parser/tass"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {

	//Создание json файла, в который будут записаны ссылки, заголовки, тела новостей
	encfile, _ := os.OpenFile("big_encode.json", os.O_CREATE, os.ModePerm)
	defer encfile.Close()

	wg:=new(sync.WaitGroup)
	key.Init()
	t:=time.Now()

	AllNews:=make([]model.News, 0)

	wg.Add(1)
	go nytimes.GetNews(wg)
	wg.Add(1)
	go brazilian.GetNews(wg)
	wg.Add(1)
	go india.GetNews(wg)
	wg.Add(1)
	go japantoday.GetNews(wg)
	wg.Add(1)
	go bbc.GetNews(wg)
	wg.Add(1)
	go canada.GetNews(wg)
	wg.Add(1)
	go fr.GetNews(wg)
	wg.Add(1)
	go euronews.GetNews(wg)
	wg.Add(1)
	go chinapost.GetNews(wg)
	wg.Add(1)
	go tass.GetNews(wg)
	wg.Add(1)
	go rt.GetNews(wg)

	wg.Wait()

	//Working:
	AllNews = append(AllNews, 	nytimes.News...)
	//AllNews = append(AllNews, 	brazilian.News...)
	AllNews = append(AllNews, 	india.News...)
	AllNews = append(AllNews, 	japantoday.News...)
	AllNews = append(AllNews, 	bbc.News...)
	AllNews = append(AllNews,   canada.News...)
	AllNews = append(AllNews,   fr.News...)
	AllNews = append(AllNews,	euronews.News...)
	AllNews = append(AllNews,	chinapost.News...)
	//AllNews = append(AllNews,	tass.News...)
	AllNews = append(AllNews,	rt.News...)

	//In progress:
	//AllNews = append(AllNews, 	thelocal.GetNews()...)
	//AllNews = append(AllNews, 	ecns.GetNews()...)
	//AllNews = append(AllNews, 	scmp.GetNews()...)
	//AllNews = append(AllNews,     shine.GetNews()...)

	enc:=json.NewEncoder(os.Stdout)

	enc.SetIndent(""," ")
	enc.Encode(AllNews)

	encoded:=json.NewEncoder(encfile)
	encoded.SetIndent(""," ")
	encoded.Encode(AllNews)

	fmt.Println("\nВремя поиска: ", time.Since(t))

	//Создание клиента и отправка запроса на сервер (Flask)
	url := "http://127.0.0.1:5000/"

	content, err := ioutil.ReadFile("big_encode.json")
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//Вывод результатов запроса
	fmt.Println("\n\n\n\n\n\nParsed data")
	fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	//Создание и запись файла, содержащего результат запроса
	f, err := os.Create("parsednews.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(string(body))

	if err2 != nil {
		log.Fatal(err2)
	}

}
