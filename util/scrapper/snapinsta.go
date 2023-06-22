package scrapper

import (
	"bytes"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"strings"
)

type SnapInstaResponse struct {
	Username    string   `json:"username"`
	Avatar      string   `json:"avatar"`
	ResultMedia []string `json:"result_media"`
}

func (scrapper *Scrapper) GetSnapInsta(instagram string) (response SnapInstaResponse, err error) {
	defer TimeElapsed("Scrap SnapInsta")()

	client := scrapper.NewCloudflareBypass()
	resp, err := client.R().
		SetFormData(map[string]string{
			"url":    instagram,
			"action": "post",
			"lang":   "id",
		}).
		SetHeader("Origin", "https://snapinsta.app").
		SetHeader("Referer", "https://snapinsta.app/id").
		SetHeader("User-Agent", browser.Firefox()).
		Post("https://snapinsta.app/action2.php")

	if err != nil {
		return response, err
	}

	defer resp.RawBody().Close()

	vm := goja.New()
	_, err = vm.RunString(string(snapJS))
	if err != nil {
		return response, err
	}

	var fn func(string) string
	err = vm.ExportTo(vm.Get("Decode"), &fn)
	if err != nil {
		return response, err
	}

	result := fn(resp.String())
	html := innerHtml.FindStringSubmatch(result)[1]
	parsedHtml := strings.ReplaceAll(html, `\"`, "")

	document, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(parsedHtml)))
	if err != nil {
		return response, err
	}

	response = SnapInstaResponse{
		ResultMedia: make([]string, 0),
	}
	response.Username = document.Find("div.download-top > div").First().Text()
	avatar, _ := document.Find("div.download-top > div > img").First().Attr("src")
	response.Avatar = avatar

	document.Find("div.download-bottom > a").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if ok {
			response.ResultMedia = append(response.ResultMedia, href)
		}
	})

	return response, nil
}
