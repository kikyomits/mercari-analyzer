package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)


func main() {
	var url = "https://www.mercari.com/jp/search/"
	var headers = map[string]string {
		"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
		"accept-encoding": "gzip, deflate, br",
		"accept-language": "en-US,en;q=0.9,ja;q=0.8,pt;q=0.7",
		"cache-control": "no-cache",
		"pragma": "no-cache",
		"upgrade-insecure-requests": "1",
		"user-agent": "Mozilla/5.0",
	}



	req, _ := http.NewRequest("GET", url, nil)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	q := req.URL.Query()
    q.Add("keyword", "ルンバ")
    req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, _ := http.DefaultClient.Do(req)
	fmt.Println(resp)
	doc := goquery.NewDocumentFromNode(resp.Body)
	selection := doc.Find("main")
	innerSeceltion := selection.Find("a")

	innerSeceltion.Each(func(index int, s *goquery.Selection) {
		fmt.Println(s)
	})

}