package media

import (
	"fmt"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util"
	"github.com/itzngga/Roxy/util/cli"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func init() {
	embed.Commands.Add(sticker)
}

var sticker = &command.Command{
	Name:        "sticker",
	Aliases:     []string{"s", "stiker"},
	Category:    constant.MEDIA_CATEGORY,
	Description: "Create sticker from image or video",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if ctx.Message.GetImageMessage() != nil {
			return StickerImage(ctx, ctx.Message.GetImageMessage())
		} else if util.ParseQuotedMessage(ctx.Message).GetImageMessage() != nil {
			return StickerImage(ctx, util.ParseQuotedMessage(ctx.Message).GetImageMessage())
		} else if ctx.Message.GetVideoMessage() != nil {
			return StickerVideo(ctx, ctx.Message.GetVideoMessage())
		} else if util.ParseQuotedMessage(ctx.Message).GetVideoMessage() != nil {
			return StickerVideo(ctx, util.ParseQuotedMessage(ctx.Message).GetVideoMessage())
		}
		return ctx.GenerateReplyMessage("Invalid Sticker Action")
	},
}

func StickerVideo(ctx *command.RunFuncContext, video *waProto.VideoMessage) *waProto.Message {
	data, err := ctx.Client.Download(video)
	if err != nil {
		fmt.Printf("Failed to download video: %v\n", err)
		return nil
	}

	var qValue int
	switch dataLen := len(data); {
	case dataLen < 300000:
		qValue = 25
	case dataLen < 400000:
		qValue = 15
	default:
		qValue = 8
	}

	resultData, err := cli.ExecPipeline("ffmpeg", data,
		"-y", "-hide_banner", "-loglevel", "panic",
		"-i", "pipe:0",
		"-filter:v", "fps=fps=15",
		"-compression_level", "0",
		"-q:v", fmt.Sprintf("%d", qValue),
		"-loop", "0",
		"-preset", "picture",
		"-an", "-vsync", "0",
		"-s", "512:512",
		"-f", "webp",
		"pipe:1",
	)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	ctx.SendReplyMessage(resultData)
	return nil
}
func StickerImage(ctx *command.RunFuncContext, img *waProto.ImageMessage) *waProto.Message {
	data, err := ctx.Client.Download(img)
	if err != nil {
		fmt.Printf("Failed to download image: %v\n", err)
		return nil
	}

	resultData, err := cli.ExecPipeline("ffmpeg", data,
		"-y", "-hide_banner", "-loglevel", "panic",
		"-i", "pipe:0",
		"-f", "webp",
		"-s", "512:512",
		"-preset", "picture",
		"pipe:1",
	)

	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	ctx.SendReplyMessage(resultData)
	return nil
}
