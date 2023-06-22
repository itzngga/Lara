package scrapper

import (
	"bytes"
	"encoding/json"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
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

func (scrapper *Scrapper) GetRenderSnapSave(token string) (string, error) {
	client := scrapper.NewCloudflareBypass()
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
			return scrapper.GetRenderSnapSave(token)
		} else {
			return result.Data.FilePath, nil
		}
	} else {
		return scrapper.GetRenderSnapSave(token)
	}
}

func (scrapper *Scrapper) GetSnapSave(requiredUrl string) (response []SnapSaveResponse, err error) {
	defer TimeElapsed("Scrap SnapSave")()

	client := scrapper.NewCloudflareBypass()
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
