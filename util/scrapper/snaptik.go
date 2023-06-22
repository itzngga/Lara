package scrapper

import (
	"bytes"
	_ "embed"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/go-resty/resty/v2"
	"strings"
)

//go:embed snap.js
var snapJS []byte

type SnaptikResponse struct {
	Username    string   `json:"username"`
	Description string   `json:"description"`
	VideoUrl    []string `json:"video_url"`
	ImageUrl    []string `json:"image_url"`
}

func (scrapper *Scrapper) GetSnaptik(tiktok string) (response SnaptikResponse, err error) {
	defer TimeElapsed("Scrap Snaptik")()

	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", browser.Firefox()).
		Get("https://snaptik.app/ID")
	if err != nil {
		return response, err
	}

	defer resp.RawBody().Close()

	document, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return response, err
	}

	var param = map[string]string{
		"url": tiktok,
	}

	document.Find("form input").Each(func(index int, selector *goquery.Selection) {
		name, ok := selector.Attr("name")
		if ok && name == "lang" {
			value, _ := selector.Attr("value")
			param["lang"] = value
		}

		if ok && name == "token" {
			value, _ := selector.Attr("value")
			param["token"] = value
		}

		selector.Next()
	})

	resp, err = client.R().
		SetFormData(param).
		SetHeader("Origin", "https://snaptik.app").
		SetHeader("Referer", "https://snaptik.app/ID").
		SetHeader("User-Agent", browser.Firefox()).
		Post("https://snaptik.app/abc2.php")
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

	document, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(parsedHtml)))
	if err != nil {
		return response, err
	}

	response = SnaptikResponse{}
	response.Username = document.Find("div.info > span").Text()
	response.Description = document.Find("div.info > div").Text()
	response.VideoUrl = make([]string, 0)
	response.ImageUrl = make([]string, 0)

	document.Find("div.video-links > a").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if ok {
			if strings.Contains(href, "https://cdn.snaptik.app") {
				response.VideoUrl = append(response.VideoUrl, href)
			}
		}
		selection.Next()
	})

	document.Find("div.dl-footer").Each(func(i int, selection *goquery.Selection) {
		src, ok := selection.Find("a").Attr("href")
		if ok && strings.Contains(src, "https://tikcdn.net") {
			response.ImageUrl = append(response.ImageUrl, src)
		}
	})

	return response, nil
}
