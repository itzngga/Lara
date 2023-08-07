package media

import (
	"bytes"
	"github.com/itzngga/Lara/repo"
	"github.com/itzngga/Lara/src/cmd/constant"
	util2 "github.com/itzngga/Lara/util"
	"github.com/itzngga/Lara/util/metadata"
	"github.com/itzngga/Lara/util/scrapper"
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
	defer scrapper.TimeElapsed("Sticker Video")()

	ctx.SendEmoji("👌")

	data, err := ctx.Client.Download(video)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	exifFile := "temp/" + util2.MakeMD5UUID() + ".exif"
	webpFile := "temp/" + util2.MakeMD5UUID() + ".webp"

	var exif []byte
	if wm, err := repo.WMRepository.GetWMByJid(ctx.MessageInfo.Sender.ToNonAD().String()); wm.JID != "" || err != nil {
		exif = metadata.CreateMetadata(metadata.StickerMetadata{
			Name:      wm.StickerName,
			Publisher: wm.StickerPublisher,
		})
	} else {
		exif = metadata.CreateMetadata(metadata.StickerMetadata{
			Name:      "Sticker",
			Publisher: "Roxy",
		})
	}

	err = os.WriteFile(exifFile, exif, os.ModePerm)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	_, err = cli.ExecPipeline("ffmpeg", data,
		"-y", "-hide_banner", "-loglevel", "error",
		"-i", "pipe:0", "-f", "mp4",
		"-ss", "00:00:00", "-t", "00:00:15",
		"-vf", "fps=10,scale=720:-1:flags=lanczos:force_original_aspect_ratio=increase,crop=512:512,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=#00000000,setsar=1",
		"-compression_level", "6",
		"-q:v", "60", "-loop", "0",
		"-preset", "picture", "-an", "-fps_mode", "auto",
		"-f", "webp",
		webpFile,
	)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	err = cmdchain.Builder().
		Join("webpmux", "-set", "exif", exifFile, webpFile, "-o", webpFile).
		Finalize().Run()
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	webp, err := ctx.UploadStickerMessageFromPath(webpFile)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	defer func() {
		data = nil

		os.Remove(exifFile)
		os.Remove(webpFile)
		exif = nil
		webp = nil
	}()

	return ctx.GenerateReplyMessage(webp)
}
func StickerImage(ctx *command.RunFuncContext, img *waProto.ImageMessage) *waProto.Message {
	defer scrapper.TimeElapsed("Sticker Image")()

	ctx.SendEmoji("👌")

	data, err := ctx.Client.Download(img)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	reader := bytes.NewReader(data)
	exifFile := "temp/" + util2.MakeMD5UUID() + ".exif"
	webpFile := "temp/" + util2.MakeMD5UUID() + ".webp"

	var exif []byte
	if wm, err := repo.WMRepository.GetWMByJid(ctx.MessageInfo.Sender.ToNonAD().String()); wm.JID != "" || err != nil {
		exif = metadata.CreateMetadata(metadata.StickerMetadata{
			Name:      wm.StickerName,
			Publisher: wm.StickerPublisher,
		})
	} else {
		exif = metadata.CreateMetadata(metadata.StickerMetadata{
			Name:      "Sticker",
			Publisher: "Roxy",
		})
	}

	err = os.WriteFile(exifFile, exif, os.ModePerm)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	err = cmdchain.Builder().
		Join("convert", "-", "-resize", "512x512", "-background", "none", "-compose", "Copy", "-gravity", "center", "-extent", "512x512", "-quality", "100", "png:-").
		WithInjections(reader).ForwardError().
		Join("cwebp", "-quiet", "-mt", "-exact", "-q", "100", "-m", "6", "-alpha_q", "100", "-o", webpFile, "--", "-").DiscardStdOut().
		Finalize().WithError(os.Stderr).Run()
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	err = cmdchain.Builder().
		Join("webpmux", "-set", "exif", exifFile, webpFile, "-o", webpFile).
		Finalize().Run()
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	webp, err := ctx.UploadStickerMessageFromPath(webpFile)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	defer func() {
		reader = nil
		data = nil

		os.Remove(exifFile)
		os.Remove(webpFile)
		exif = nil
		webp = nil
	}()

	return ctx.GenerateReplyMessage(webp)
}

func StickerWM(ctx *command.RunFuncContext) *waProto.Message {
	defer scrapper.TimeElapsed("Sticker WM")()

	ctx.SendEmoji("👌")

	exifFile := "temp/" + util2.MakeMD5UUID() + ".exif"
	webpFile := "temp/" + util2.MakeMD5UUID() + ".webp"

	_, err := ctx.DownloadToFile(true, webpFile)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	var exif []byte
	if wm, err := repo.WMRepository.GetWMByJid(ctx.MessageInfo.Sender.ToNonAD().String()); wm.JID != "" || err != nil {
		exif = metadata.CreateMetadata(metadata.StickerMetadata{
			Name:      wm.StickerName,
			Publisher: wm.StickerPublisher,
		})
	} else {
		exif = metadata.CreateMetadata(metadata.StickerMetadata{
			Name:      "Sticker",
			Publisher: "Roxy",
		})
	}

	err = os.WriteFile(exifFile, exif, os.ModePerm)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	err = cmdchain.Builder().
		Join("webpmux", "-set", "exif", exifFile, webpFile, "-o", webpFile).
		Finalize().Run()
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	webp, err := ctx.UploadStickerMessageFromPath(webpFile)
	if err != nil {
		return ctx.GenerateReplyMessage("error: " + err.Error())
	}

	defer func() {
		os.Remove(exifFile)
		os.Remove(webpFile)
		exif = nil
		webp = nil
	}()

	return ctx.GenerateReplyMessage(webp)
}
