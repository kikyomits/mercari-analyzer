package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"io/ioutil"
)


func convrtToUTF8(str string, origin string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(origin, byteReader)
	strBytes, _ = ioutil.ReadAll(reader)
	return string(strBytes)
}

func getItemData(keyword string) {
	url := "https://www.mercari.com/jp/search/?keyword=" + keyword
	doc, _ := goquery.NewDocument(url)
	// sections of items on the mercari
	selector := "body > div.default-container > main > div.l-content > section > div.items-box-content.clearfix > section.items-box"
	doc.Find(selector).EachWithBreak(func(i int, s *goquery.Selection) bool {
		// Get item detail page url
		inner := s.Find("a")
		itemUrl, isExists := inner.Attr("href")

		if isExists {
			getItemDetail(itemUrl)
			return true
		}
		return true
	})
}

func getItemDetail(url string) {
	selector := "body > div.default-container > section > div.item-main-content.clearfix > table > tbody > tr"
	doc, _ := goquery.NewDocument(url)
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// Get item detail page url
		//body > div.default-container > section > div.item-main-content.clearfix > table > tbody > tr:nth-child(2) > td > a:nth-child(2) > div > i
		//body > div.default-container > section > div.item-main-content.clearfix > table > tbody > tr:nth-child(2) > td > a:nth-child(1) > div
		//body > div.default-container > section > div.item-main-content.clearfix > table > tbody > tr:nth-child(2) > td > a:nth-child(2) > div > i

		key := s.Find("th").Text()

		value := ""
		if key == "出品者" {

		} else if key == "カテゴリー" {
			s.Find("td > a").Each(func(j int, inner *goquery.Selection) {
				if j == 0 {
					value = inner.Text()
				} else {
					value = value + "," + inner.Find("div").Text()
				}
			})
		} else {
			value = s.Find("td").Text()
		}
		fmt.Println(key, value)
		//itemUrl, isExists := inner.Attr("href")
		//
		//if isExists {
		//	fmt.Println(itemUrl)
		//}
	})
}


func main() {
	getItemData("ルンバ")

	//fileInfos, _ := ioutil.ReadFile("./test.html")
	//stringReader := strings.NewReader(string(fileInfos))
	//doc, _ := goquery.NewDocumentFromReader(stringReader)



	//innerSeceltion := selection.Find("a")
	//
	//innerSeceltion.Each(func(index int, s *goquery.Selection) {
	//	fmt.Println(s)
	//})

}