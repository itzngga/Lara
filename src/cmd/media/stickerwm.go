package media

import (
	"github.com/itzngga/Lara/entity"
	"github.com/itzngga/Lara/repo"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func init() {
	embed.Commands.Add(stickerwm)
}

var stickerwm = &command.Command{
	Name:        "stickerwm",
	Aliases:     []string{"swm", "stkwm"},
	Category:    constant.MEDIA_CATEGORY,
	Description: "Create sticker from image or video with watermark",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if len(ctx.Arguments) != 0 {
			var name, publisher string

			command.NewUserQuestion(ctx).
				SetNoAskReplyQuestion(&name).
				SetNoAskReplyQuestion(&publisher).
				ExecWithParser()

			err := repo.WMRepository.PutWM(entity.WMEntity{
				JID:              ctx.MessageInfo.Sender.ToNonAD().String(),
				StickerName:      name,
				StickerPublisher: publisher,
			})

			if err != nil {
				return ctx.GenerateReplyMessage("error: " + err.Error())
			}
		}

		if ctx.Message.GetImageMessage() != nil {
			return StickerImage(ctx, ctx.Message.GetImageMessage())
		} else if util.ParseQuotedMessage(ctx.Message).GetImageMessage() != nil {
			return StickerImage(ctx, util.ParseQuotedMessage(ctx.Message).GetImageMessage())
		} else if ctx.Message.GetVideoMessage() != nil {
			return StickerVideo(ctx, ctx.Message.GetVideoMessage())
		} else if util.ParseQuotedMessage(ctx.Message).GetVideoMessage() != nil {
			return StickerVideo(ctx, util.ParseQuotedMessage(ctx.Message).GetVideoMessage())
		} else if util.ParseQuotedMessage(ctx.Message).GetStickerMessage() != nil {
			return StickerWM(ctx)
		}
		return ctx.GenerateReplyMessage("error: unexpected action")
	},
}
