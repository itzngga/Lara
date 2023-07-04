package media

import (
	"bytes"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util"
	"github.com/itzngga/Roxy/util/cli"
	cmdchain "github.com/rainu/go-command-chain"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"os"
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
		return ctx.GenerateReplyMessage("error: unexpected action")
	},
}

func StickerVideo(ctx *command.RunFuncContext, video *waProto.VideoMessage) *waProto.Message {
	data, err := ctx.Client.Download(video)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	resultData, err := cli.ExecPipeline("ffmpeg", data,
		"-y", "-hide_banner", "-loglevel", "error",
		"-i", "pipe:0", "-f", "mp4",
		"-ss", "00:00:00", "-t", "00:00:15",
		"-vf", "fps=10,scale=720:-1:flags=lanczos:force_original_aspect_ratio=increase,crop=512:512,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=#00000000,setsar=1",
		"-compression_level", "6",
		"-q:v", "60", "-loop", "0",
		"-preset", "picture", "-an", "-fps_mode", "auto",
		"-f", "webp",
		"pipe:1",
	)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	defer func() {
		data = nil
		resultData = nil
	}()

	return ctx.GenerateReplyMessage(resultData)
}
func StickerImage(ctx *command.RunFuncContext, img *waProto.ImageMessage) *waProto.Message {
	data, err := ctx.Client.Download(img)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	reader := bytes.NewReader(data)
	writer := bytes.NewBuffer(nil)

	defer func() {
		reader = nil
		writer = nil
		data = nil
	}()

	err = cmdchain.Builder().
		Join("convert", "-", "-resize", "512x512", "-background", "black", "-compose", "Copy", "-gravity", "center", "-extent", "512x512", "-quality", "100", "-").
		WithInjections(reader).ForwardError().
		Join("cwebp", "-quiet", "-mt", "-exact", "-q", "100", "-m", "6", "-alpha_q", "100", "-o", "-", "--", "-").
		WithOutputForks(writer).WithErrorForks(os.Stderr).
		Finalize().WithError(os.Stderr).Run()
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	return ctx.GenerateReplyMessage(writer.Bytes())
}
