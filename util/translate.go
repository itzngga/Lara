package util

import (
	"encoding/json"
	"errors"
	"github.com/itzngga/Roxy/util"
	"net/url"
	"strings"
)

func Translate(source, sourceLang, targetLang string) (string, error) {
	params := url.Values{}
	params.Add("client", "gtx")
	params.Add("sl", sourceLang)
	params.Add("tl", targetLang)
	params.Add("dt", "t")
	params.Add("q", source)

	body, err := util.DoHTTPRequest("GET", "https://translate.googleapis.com/translate_a/single?"+params.Encode())
	if err != nil {
		return "", err
	}

	if ok := strings.Contains(string(body), `<title>Error 400 (Bad Request)`); ok {
		return "err", errors.New("error: 400 (Bad Request)")
	}

	var text []string
	var result []interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "err", errors.New("error: unmarshaling data")
	}

	defer func() {
		body = nil
		text = nil
		result = nil
	}()

	if len(result) > 0 {
		inner := result[0]
		for _, slice := range inner.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				text = append(text, translatedText.(string))
				break
			}
		}
		cText := strings.Join(text, "")

		return cText, nil
	} else {
		return "", errors.New("error: no translated data in response")
	}
}
