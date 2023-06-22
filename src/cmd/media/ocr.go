package media

import (
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util/cli"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func init() {
	embed.Commands.Add(ocr)
}

var ocr = &command.Command{
	Name:        "ocr",
	Category:    constant.MEDIA_CATEGORY,
	Description: "Scan text on images",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var lang, psm string
		var msg *waProto.Message
		if msg = ctx.GetDownloadable(true); msg == nil {
			command.NewUserQuestion(ctx).
				SetQuestion("Please select a language", &lang).
				SetQuestion("Please select a psm", &psm).
				CaptureMediaQuestion("Please send/reply a media message", &msg).
				Exec()
		}

		if lang == "" || psm == "" {
			command.NewUserQuestion(ctx).
				SetQuestion("Please select a language", &lang).
				SetQuestion("Please select a psm", &psm).
				Exec()
		}

		result, err := ctx.DownloadMessage(msg, false)
		if err != nil {
			return ctx.GenerateReplyMessage("error: no downloadable message")
		}
		res := cli.ExecPipeline("tesseract", result, "stdin", "stdout", "-l", lang, "--psm", psm)
		if res == nil {
			return ctx.GenerateReplyMessage("error: no response")
		}

		return ctx.GenerateReplyMessage(string(res))
	},
}
