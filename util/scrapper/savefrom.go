package scrapper

import (
	"encoding/json"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/go-resty/resty/v2"
)

type Y2MateVideoData struct {
	Size   string `json:"size"`
	Format string `json:"f"`
	Query  string `json:"q"`
	Token  string `json:"k"`
}
type Y2MateResponse struct {
	Status  string `json:"status"`
	VideoId string `json:"vid"`
	Title   string `json:"title"`
	Second  int    `json:"t"`
	Channel string `json:"a"`
	Links   struct {
		Mp4   map[string]Y2MateVideoData `json:"mp4"`
		Mp3   map[string]Y2MateVideoData `json:"mp3"`
		Other map[string]Y2MateVideoData `json:"other"`
	}
}

type DownloadY2MateResponse struct {
	Status   string `json:"status"`
	VideoID  string `json:"vid"`
	Title    string `json:"title"`
	MediaUrl string `json:"dlink"`
}

func GetY2MateFromToken(vidId string, data Y2MateVideoData) (response DownloadY2MateResponse, err error) {
	client := resty.New()
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Origin":       "https://www.y2mate.com",
			"Referer":      "https://www.y2mate.com/en625",
			"User-Agent":   browser.Firefox(),
		}).
		SetFormData(map[string]string{
			"vid": vidId,
			"k":   data.Token,
		}).
		Post("https://www.y2mate.com/mates/convertV2/index")
	if err != nil {
		return response, err
	}

	var result DownloadY2MateResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetY2Mate(urlTarget string) (response Y2MateResponse, err error) {
	client := resty.New()
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Origin":       "https://www.y2mate.com",
			"Referer":      "https://www.y2mate.com/id625",
			"User-Agent":   browser.Firefox(),
		}).
		SetFormData(map[string]string{
			"k_query": urlTarget,
			"k_page":  "home",
			"hl":      "id",
			"q_auto":  "0",
		}).
		Post("https://www.y2mate.com/mates/analyzeV2/ajax")
	if err != nil {
		return response, err
	}

	var result Y2MateResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
