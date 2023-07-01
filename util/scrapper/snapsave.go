package scrapper

import (
	"bytes"
	"encoding/json"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type SnapSaveResponse struct {
	Quality string `json:"quality"`
	Render  bool   `json:"render"`
	Link    string `json:"link"`
}

type SnapSaveRenderResponse struct {
	Status int `json:"status"`
	Data   struct {
		Progress int    `json:"progress"`
		FilePath string `json:"file_path"`
	} `json:"data"`
}

func GetRenderSnapSave(token string) (string, error) {
	client := NewCloudflareBypass()
	resp, err := client.R().
		SetQueryParam("token", token).
		SetHeader("User-Agent", browser.Firefox()).
		Get("https://snapsave.app/render.php")

	if err != nil {
		return "", err
	}

	var result SnapSaveRenderResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	defer resp.RawBody().Close()

	if result.Status == 1 {
		if result.Data.Progress != 100 {
			return GetRenderSnapSave(token)
		} else {
			return result.Data.FilePath, nil
		}
	} else {
		return GetRenderSnapSave(token)
	}
}

func GetSnapSave(requiredUrl string) (response []SnapSaveResponse, err error) {
	defer TimeElapsed("Scrap SnapSave")()

	client := NewCloudflareBypass()
	resp, err := client.R().
		SetFormData(map[string]string{
			"url":    requiredUrl,
			"action": "post",
			"lang":   "id",
		}).
		SetHeader("Origin", "https://snapinsta.app").
		SetHeader("Referer", "https://snapinsta.app/id").
		SetHeader("User-Agent", browser.Firefox()).
		Post("https://snapsave.app/action.php?lang=id")
	if err != nil {
		return response, err
	}

	defer resp.RawBody().Close()

	result, err := DecodeSnap(resp.String())
	if err != nil {
		return response, err
	}

	html := innerHtml.FindStringSubmatch(result)[1]
	parsedHtml := strings.ReplaceAll(html, `\"`, "")

	document, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(parsedHtml)))
	if err != nil {
		return response, err
	}

	response = make([]SnapSaveResponse, 0)
	document.Find("table > tbody > tr").Each(func(i int, selection *goquery.Selection) {
		result := SnapSaveResponse{
			Quality: selection.Find("td.video-quality").Text(),
		}

		href, ok := selection.Find("td > a").Attr("href")
		if ok {
			result.Link = href
		} else {
			onClick, ok := selection.Find("td > button").Attr("onclick")
			if ok {
				token := strings.ReplaceAll(onClick, `\'`, "'")
				token = innerToken.FindStringSubmatch(token)[1]
				result.Link = token
				result.Render = true
			}
		}

		response = append(response, result)
	})

	return response, nil
}
