package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yoonhero/ohpotatocoin/utils"
)

func Scrape() {
	var baseURL string = "https://search.shopping.naver.com/catalog/12542966302"
	res, err := http.Get(baseURL)
	utils.HandleErr(err)
	defer res.Body.Close()
	checkCode(res)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	utils.HandleErr(err)

	fmt.Println(doc.Find(".style_container__3iYev").Text())

	circles := doc.Find("#_line_chart_container > svg > g > g.bb-chart > g.bb-chart-lines > g > g.bb-shapes.bb-shapes-y.bb-circles.bb-circles-y")
	// text := doc.Find("text.bb-text-6").Text()

	circles.Each(func(idx int, sel *goquery.Selection) {
		title, ok := sel.Find(".bb-circle").Attr("cy")

		fmt.Println(title, goquery.NodeName(sel), idx, "aaa", ok)
	})

	fmt.Println(doc.Find("#_line_chart_container > svg > g > g.bb-chart > g.bb-chart-lines > g > g.bb-shapes.bb-shapes-y.bb-circles.bb-circles-y > circle.bb-shape.bb-shape-0.bb-circle.bb-circle-0"))
	//bb-circles
	// id, _ := card.Attr("data-jk")
	// title := cleanString(card.Find(".title>a").Text())
	// location := cleanString(card.Find(".sjcl").Text())
	// salary := cleanString(card.Find(".salaryText").Text())
	// summary := cleanString(card.Find(".summary").Text())

}

// CleanString cleans a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
