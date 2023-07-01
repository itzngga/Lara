package media

import (
	"fmt"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Lara/util"
	"github.com/itzngga/Lara/util/scrapper"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/valyala/fasttemplate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"strconv"
	"time"

	"strings"
)

func init() {
	embed.Commands.Add(downloader)
}

var downloader = &command.Command{
	Name:        "downloader",
	Aliases:     []string{"down", "get", "unduh", "g"},
	Description: "Download media based on Url",
	Category:    "media",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var link string
		command.NewUserQuestion(ctx).
			SetQuestion("Please send media url link", &link).
			WithLikeEmoji().
			ExecWithParser()

		if link != "" {
			if !util.ParseURL(link) {
				return ctx.GenerateReplyMessage("errors: invalid url scheme")
			}
		}

		if strings.Contains(link, "youtu.be") || strings.Contains(link, "youtube") {
			response, err := scrapper.GetY2Mate(link)
			if err != nil {
				return ctx.GenerateReplyMessage("error: " + err.Error())
			}

			var (
				downloads string
				i         int
			)

			downloads += "*── 「 OTHER 」 ──*\n"
			for key, data := range response.Links.Other {
				i++
				downloads += fmt.Sprintf("%d. *MediaID*: %s\n*Format*: %s (%s)\n*Size*: %s\n\n", i, key, data.Query, data.Format, data.Size)
			}

			i = 0
			downloads += "*── 「 MP4 」 ──*\n"
			for key, data := range response.Links.Mp4 {
				i++
				downloads += fmt.Sprintf("%d. *MediaID*: %s\n*Format*: %s (%s)\n*Size*: %s\n\n", i, key, data.Query, data.Format, data.Size)
			}

			i = 0
			downloads += "*── 「 MP3 」 ──*\n"
			for key, data := range response.Links.Mp3 {
				i++
				downloads += fmt.Sprintf("%d. *MediaID*: %s\n*Format*: %s (%s)\n*Size*: %s\n\n", i, key, data.Query, data.Format, data.Size)
			}

			template := fasttemplate.New(constant.Y2MATE_DESC, "[", "]")
			result := template.ExecuteString(map[string]interface{}{
				"url":       []byte("https://youtube.com/" + response.VideoId),
				"title":     []byte(response.Title),
				"duration":  []byte(util.HumanizeDuration(time.Second * time.Duration(response.Second))),
				"channel":   []byte(response.Channel),
				"downloads": []byte(downloads),
			})

			ctx.SendReplyMessage(result)

			var selector string
			command.NewUserQuestion(ctx).
				SetNoAskQuestions(&selector).
				WithOkEmoji().
				Exec()

			var selected scrapper.Y2MateVideoData
			var typeMedia int
			if val, ok := response.Links.Other[selector]; ok {
				typeMedia = 1
				selected = val
			}
			if val, ok := response.Links.Mp3[selector]; ok {
				typeMedia = 0
				selected = val
			}
			if val, ok := response.Links.Mp4[selector]; ok {
				typeMedia = 1
				selected = val
			}

			if selected.Token == "" {
				return ctx.GenerateReplyMessage("error: invalid given media id")
			} else {
				ctx.SendReplyMessage("Downloading...")
			}

			downloaded, err := scrapper.GetY2MateFromToken(response.VideoId, selected)
			if err != nil {
				return ctx.GenerateReplyMessage("error: " + err.Error())
			}

			template = fasttemplate.New(constant.Y2MATE_RESULT_CAPTION, "[", "]")
			result = template.ExecuteString(map[string]interface{}{
				"url":      []byte("https://youtube.com/" + response.VideoId),
				"title":    []byte(response.Title),
				"duration": []byte(util.HumanizeDuration(time.Second * time.Duration(response.Second))),
				"channel":  []byte(response.Channel),
				"size":     []byte(selected.Size),
				"format":   []byte(selected.Query + "(" + selected.Format + ")"),
			})

			if typeMedia == 1 {
				uploaded, err := ctx.UploadVideoFromUrl(downloaded.MediaUrl, result)
				if err != nil {
					return ctx.GenerateReplyMessage("error: " + err.Error())
				}
				return ctx.GenerateReplyMessage(uploaded)
			} else {
				uploaded, err := ctx.UploadAudioFromUrl(downloaded.MediaUrl)
				if err != nil {
					return ctx.GenerateReplyMessage("error: " + err.Error())
				}
				return ctx.GenerateReplyMessage(uploaded)
			}
		} else if strings.Contains(link, "tiktok") {
			result, err := scrapper.GetSnaptik(link)
			if err != nil {
				return ctx.GenerateReplyMessage("error: " + err.Error())
			}

			var slides = make([]struct {
				index int
				media string
				url   string
			}, 0)

			var index = 1
			for _, imageUrl := range result.ImageUrl {
				slides = append(slides, struct {
					index int
					media string
					url   string
				}{index: index, media: "image", url: imageUrl})
				index++
			}
			for _, videoUrl := range result.VideoUrl {
				slides = append(slides, struct {
					index int
					media string
					url   string
				}{index: index, media: "video", url: videoUrl})
				index++
			}

			if len(slides) == 1 {
				template := fasttemplate.New(constant.SNAPTIK_RESULT, "[", "]")
				result := template.ExecuteString(map[string]interface{}{
					"username":    []byte(result.Username),
					"description": []byte(result.Description),
				})

				if slides[0].media == "image" {
					uploaded, err := ctx.UploadImageFromUrl(slides[0].url, result)
					if err != nil {
						return ctx.GenerateReplyMessage("error: " + err.Error())
					}
					return ctx.GenerateReplyMessage(uploaded)
				} else {
					uploaded, err := ctx.UploadVideoFromUrl(slides[0].url, result)
					if err != nil {
						return ctx.GenerateReplyMessage("error: " + err.Error())
					}
					return ctx.GenerateReplyMessage(uploaded)
				}
			} else {
				var parsedSlides string
				for _, data := range slides {
					parsedSlides += fmt.Sprintf("%d. Slide ke-%d\n\n", data.index, data.index)
				}
				template := fasttemplate.New(constant.SNAPTIK_LIST, "[", "]")
				templateResult := template.ExecuteString(map[string]interface{}{
					"username":    []byte(result.Username),
					"description": []byte(result.Description),
					"slides":      []byte(parsedSlides),
				})

				ctx.SendReplyMessage(templateResult)

				var selector string
				command.NewUserQuestion(ctx).
					SetNoAskQuestions(&selector).
					WithOkEmoji().
					Exec()

				selected, err := strconv.Atoi(selector)
				if err != nil {
					return ctx.GenerateReplyMessage("error: invalid slide")
				}

				var selectedStruct struct {
					index int
					media string
					url   string
				}

				var ok bool
				for _, slide := range slides {
					if slide.index == selected {
						selectedStruct = slide
						ok = true
						break
					}
				}

				if !ok {
					return ctx.GenerateReplyMessage("error: invalid slide")
				}

				template = fasttemplate.New(constant.SNAPTIK_RESULT, "[", "]")
				result := template.ExecuteString(map[string]interface{}{
					"username":    []byte(result.Username),
					"description": []byte(result.Description),
				})

				if selectedStruct.media == "image" {
					uploaded, err := ctx.UploadImageFromUrl(selectedStruct.url, result)
					if err != nil {
						return ctx.GenerateReplyMessage("error: " + err.Error())
					}
					return ctx.GenerateReplyMessage(uploaded)
				} else {
					uploaded, err := ctx.UploadVideoFromUrl(selectedStruct.url, result)
					if err != nil {
						return ctx.GenerateReplyMessage("error: " + err.Error())
					}
					return ctx.GenerateReplyMessage(uploaded)
				}
			}
		} else if strings.Contains(link, "instagram") {
			result, err := scrapper.GetSnapInsta(link)
			if err != nil {
				return ctx.GenerateReplyMessage("error: " + err.Error())
			}

			if len(result.ResultMedia) == 1 {
				template := fasttemplate.New(constant.SNAPINSTA_RESULT, "[", "]")
				templateRes := template.ExecuteString(map[string]interface{}{
					"username": []byte(result.Username),
				})
				return ctx.GenerateReplyMessage(map[string]string{result.ResultMedia[0]: templateRes})
			} else {
				var parsedSlides string
				for i := range result.ResultMedia {
					parsedSlides += fmt.Sprintf("%d. Slide ke-%d\n\n", i+1, i+1)
				}
				template := fasttemplate.New(constant.SNAPINSTA_LIST, "[", "]")
				templateResult := template.ExecuteString(map[string]interface{}{
					"username": []byte(result.Username),
					"slides":   []byte(parsedSlides),
				})

				ctx.SendReplyMessage(templateResult)

				var selector string
				command.NewUserQuestion(ctx).
					SetNoAskReplyQuestion(&selector).
					WithOkEmoji().
					Exec()

				selected, err := strconv.Atoi(selector)
				if err != nil {
					return ctx.GenerateReplyMessage("error: invalid slide")
				}
				var selectedSlide string
				var ok bool
				for i, data := range result.ResultMedia {
					if i+1 == selected {
						selectedSlide = data
						ok = true
						break
					}
				}

				if !ok {
					return ctx.GenerateReplyMessage("error: invalid slide")
				}

				template = fasttemplate.New(constant.SNAPINSTA_RESULT, "[", "]")
				templateRes := template.ExecuteString(map[string]interface{}{
					"username": []byte(result.Username),
				})
				return ctx.GenerateReplyMessage(map[string]string{selectedSlide: templateRes})
			}
		} else if strings.Contains(link, "twitter") {
			result, err := scrapper.GetSnapTwitter(link)
			if err != nil {
				return ctx.GenerateReplyMessage("error: " + err.Error())
			}
			template := fasttemplate.New(constant.SNAPTWIT_RESULT, "[", "]")
			resultText := template.ExecuteString(map[string]interface{}{
				"username":    []byte(result.Username),
				"description": []byte(result.Description),
			})
			return ctx.GenerateReplyMessage(map[string]string{result.MediaUrl: resultText})
		} else if strings.Contains(link, "facebook") || strings.Contains(link, "fb.watch") {
			result, err := scrapper.GetSnapSave(link)
			if err != nil {
				return ctx.GenerateReplyMessage("error: " + err.Error())
			}
			var parsedQualities string
			for i, data := range result {
				parsedQualities += fmt.Sprintf("%d. *Format*: %s\n\n", i+1, data.Quality)
			}
			template := fasttemplate.New(constant.SNAPSAVE_RESULT, "[", "]")
			templateResult := template.ExecuteString(map[string]interface{}{
				"quality": []byte(parsedQualities),
			})

			ctx.SendReplyMessage(templateResult)

			var selector string
			command.NewUserQuestion(ctx).
				SetNoAskReplyQuestion(&selector).
				WithOkEmoji().
				Exec()

			selected, err := strconv.Atoi(selector)
			if err != nil {
				return ctx.GenerateReplyMessage("error: invalid slide")
			}
			var selectedSlide scrapper.SnapSaveResponse
			var ok bool
			for i, data := range result {
				if i+1 == selected {
					selectedSlide = data
					ok = true
					break
				}
			}

			if !ok {
				return ctx.GenerateReplyMessage("error: invalid slide")
			}

			if selectedSlide.Render {
				url, err := scrapper.GetRenderSnapSave(selectedSlide.Link)
				if err != nil {
					return ctx.GenerateReplyMessage("error: " + err.Error())
				}
				selectedSlide.Link = url
			}

			return ctx.GenerateReplyMessage(map[string]string{selectedSlide.Link: "*── 「 FACEBOOK 」 ──*"})
		}

		return ctx.GenerateReplyMessage("error: download url not valid")
	},
}
