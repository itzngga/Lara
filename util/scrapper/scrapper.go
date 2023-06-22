package scrapper

import (
	"fmt"
	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/go-resty/resty/v2"
	"github.com/itzngga/Lara/util"
	"net/http"
	"regexp"
	"time"
)

type Scrapper struct {
}

func (scrapper *Scrapper) NewCloudflareBypass() *resty.Client {
	httpTransport := &http.Client{}
	httpTransport.Transport = cloudflarebp.AddCloudFlareByPass(httpTransport.Transport)

	client := resty.New()
	client.SetTransport(httpTransport.Transport)

	return client
}

//type ScrapperResponse interface {
//	SnapInstaResponse | SnaptikResponse
//}
//
//func Scrape[T ScrapperResponse](link string) (T, error) {
//	scrapper := Scrapper{}
//	if strings.Contains(link, "instagram") {
//		result, err := scrapper.GetSnapInsta(link)
//
//		return result, err
//	}
//}

var innerHtml = regexp.MustCompile("\\.innerHTML = \"(.*?)\";")
var innerToken = regexp.MustCompile("get_progressApi\\('/render\\.php\\?token=(.*?)'\\);")

func TimeElapsed(name string) func() {
	start := time.Now().In(time.Local)
	return func() {
		fmt.Printf("[%s] took %s\n", name, util.HumanizeDuration(time.Since(start)))
	}
}
