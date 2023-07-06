package media

import (
	"github.com/itzngga/Lara/entity"
	"github.com/itzngga/Lara/repo"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func init() {
	embed.Commands.Add(wm)
}

var wm = &command.Command{
	Name:        "wm",
	Aliases:     []string{"exif", "watermark"},
	Category:    constant.MEDIA_CATEGORY,
	Description: "Set watermark for creating sticker",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var name, publisher string

		command.NewUserQuestion(ctx).
			SetQuestion("Sticker Name?", &name).
			SetQuestion("Sticker Publisher?", &publisher).
			ExecWithParser()

		err := repo.WMRepository.PutWM(entity.WMEntity{
			JID:              ctx.MessageInfo.Sender.ToNonAD().String(),
			StickerName:      name,
			StickerPublisher: publisher,
		})
		if err != nil {
			return ctx.GenerateReplyMessage("error: " + err.Error())
		}

		return ctx.GenerateReplyMessage("success set watermark for sticker")
	},
}
