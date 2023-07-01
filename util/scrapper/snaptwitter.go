package scrapper

import (
	"bytes"
	"encoding/json"
	"errors"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"strings"
)

type GetSnapTwitterResponse struct {
	Username    string `json:"username"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	MediaUrl    string `json:"media_url"`
}

func GetSnapTwitter(requiredUrl string) (response GetSnapTwitterResponse, err error) {
	defer TimeElapsed("SnapTwitter")()

	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", browser.Firefox()).
		Get("https://snaptwitter.com/id")
	if err != nil {
		return response, err
	}

	defer resp.RawBody().Close()

	document, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return response, err
	}

	var param = map[string]string{
		"url": requiredUrl,
	}

	document.Find("form input").Each(func(index int, selector *goquery.Selection) {
		name, ok := selector.Attr("name")
		if ok && name == "token" {
			value, _ := selector.Attr("value")
			param["token"] = value
		}

		selector.Next()
	})

	resp, err = client.R().
		SetFormData(param).
		SetHeader("Origin", "https://snaptwitter.com").
		SetHeader("Referer", "https://snaptwitter.com/id").
		SetHeader("User-Agent", browser.Firefox()).
		Post("https://snaptwitter.com/action.php")
	if err != nil {
		return response, err
	}

	defer resp.RawBody().Close()

	var result = map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return response, err
	}

	if result["error"] != "false" {
		parsedHtml := strings.ReplaceAll(result["data"].(string), `\"`, "\"")
		document, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(parsedHtml)))
		if err != nil {
			return response, err
		}

		src, _ := document.Find("div.videotikmate-left > img").Attr("src")
		username := document.Find("div.videotikmate-middle > div > h1 > a").Text()
		description := document.Find("div.videotikmate-middle > p > span").Text()
		token, _ := document.Find("#download-block > div > a").Attr("href")

		return GetSnapTwitterResponse{
			Username:    username,
			Avatar:      src,
			Description: description,
			MediaUrl:    "https://snaptwitter.com" + token,
		}, nil
	} else {
		return response, errors.New("error: error == true")
	}
}
