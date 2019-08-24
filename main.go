package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type Output struct {
	Item struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		IsSold      bool   `json:"isSold"`
		Price       string `json:"price"`
		Tax         string `json:"tax"`
		ShippingFee string `json:"shippingFee"`
		Size        string `json:"size"`
		Brand       struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"brand"`
		Category struct {
			Large  string `json:"large"`
			Medium string `json:"medium"`
			Small  string `json:"small"`
		} `json:"category"`
		Condition string `json:"condition"`
	} `json:"item"`
	Seller struct {
		URL       string `json:"url"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		ImageURL  string ``
		Evaluates struct {
			Good   string `json:"good"`
			Normal string `json:"normal"`
			Bad    string `json:"bad"`
		} `json:"evaluates"`
		Description string `json:"description"`
		Items       struct {
			Selling int `json:"selling"`
			Sold    int `json:"sold"`
		} `json:"items"`
	} `json:"seller"`
	Shipping struct {
		Payer          string `json:"payer"`
		OriginLocation string `json:"originLocation"`
		Type           string `json:"type"`
		LeadTime       string `json:"leadTime"`
	} `json:"shipping"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func getItemData(keyword string) {
	url := "https://www.mercari.com/jp/search/?keyword=" + keyword
	doc, _ := goquery.NewDocument(url)

	// sections of items on the mercari
	selector := "body > div.default-container > main > div.l-content > section > div.items-box-content.clearfix > section.items-box"
	doc.Find(selector).EachWithBreak(func(i int, s *goquery.Selection) bool {
		// Get item detail page url
		inner := s.Find("a")
		itemURL, isExists := inner.Attr("href")

		if isExists {
			getItemDetail(itemURL)
			return false
		}
		return true
	})
}

func getItemDetail(url string) {
	out := Output{}

	// Get Item ID from URL
	out.Item.ID = strings.Split(strings.Split(url, "jp/")[1], "/?")[0]
	out.Item.URL = url

	// Analyse item detail page
	doc, _ := goquery.NewDocument(url)
	body := doc.Find("body > div.default-container > section")

	out.Item.Name = body.Find("h1.item-name").Text()
	out.Item.Description = body.Find("div.item-description > p").Text()

	priceBox := body.Find("div.item-price-box")
	out.Item.Price = priceBox.Find("span.item-price").First().Text()
	tmp := priceBox.Find("span.item-tax").First().Text()
	out.Item.Tax = strings.Replace(strings.Replace(tmp, " (", "", -1), ")", "", -1)
	out.Item.ShippingFee = priceBox.Find("span.item-shipping-fee").First().Text()

	// selector of item detail table in the page
	selector := "div.item-main-content > table > tbody > tr"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// table key
		key := s.Find("th").Text()
		switch key {
		case "出品者":
			// get seller name
			out.Seller.Name = s.Find("td > a").Text()
			out.Seller.URL, _ = s.Find("td > a").Attr("href")
			tmp := strings.Split(out.Seller.URL, "/")
			out.Seller.ID = tmp[len(tmp)-2]

			s.Find("td > div > div").Each(func(j int, inner *goquery.Selection) {
				v := inner.Find("span").Text()
				switch j {
				case 0:
					out.Seller.Evaluates.Good = v
				case 1:
					out.Seller.Evaluates.Normal = v
				case 2:
					out.Seller.Evaluates.Bad = v
				default:
					log.Warn("Scraping field warning. Field 'Seller.Evaluates' Got unexpected key: " + v)
				}
			})

		case "カテゴリー":
			s.Find("td > a").Each(func(j int, inner *goquery.Selection) {
				switch j {
				case 0:
					out.Item.Category.Large = inner.Text()
				case 1:
					out.Item.Category.Medium = strings.Replace(inner.Find("div").Text(), " ", "", 1)
				case 2:
					out.Item.Category.Small = strings.Replace(inner.Find("div").Text(), " ", "", 1)
				default:
					log.Warn("Scraping field warning. Field 'Item.Category' Got unexpected key: " + inner.Text())
				}
			})
		case "ブランド":
			out.Item.Brand.Name = strings.Replace(strings.Replace(s.Find("td > a > div").Text(), " ", "", -1), "\n", "", -1)
			out.Item.Brand.URL, _ = s.Find("td > a").Attr("href")
			tmp := strings.Split(out.Item.Brand.URL, "/")
			out.Item.Brand.ID = tmp[len(tmp)-2]

		case "商品の状態":
			out.Item.Condition = s.Find("td").Text()

		case "商品のサイズ":
			out.Item.Size = s.Find("td").Text()

		case "配送料の負担":
			out.Shipping.Payer = s.Find("td").Text()

		case "配送の方法":
			out.Shipping.Type = s.Find("td").Text()

		case "発送日の目安":
			out.Shipping.LeadTime = s.Find("td").Text()

		case "配送元地域":
			out.Shipping.OriginLocation = s.Find("td > a").Text()

		default:
			log.Info("Skipped: " + s.Text())
		}
	})

	outputJson, _ := json.Marshal(out)

	fmt.Println(string(outputJson))

}

func main() {
	getItemData("ルンバ")
}
