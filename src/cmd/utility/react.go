package cmd

import (
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
	"time"
)

func init() {
	embed.Commands.Add(react)
}

var react = &command.Command{
	Name:        "react",
	Category:    constant.UTILITY_CATEGORY,
	Description: "React a message.",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var textReact string
		if len(ctx.Arguments) >= 1 {
			textReact = ctx.Arguments[0]
		}

		id := ctx.MessageInfo.ID
		chat := ctx.MessageInfo.Chat
		sender := ctx.MessageInfo.Sender
		key := &waProto.MessageKey{
			FromMe:    proto.Bool(true),
			Id:        proto.String(id),
			RemoteJid: proto.String(chat.String()),
		}

		if !sender.IsEmpty() && sender.User != ctx.Client.Store.ID.String() {
			key.FromMe = proto.Bool(false)
			key.Participant = proto.String(sender.ToNonAD().String())
		}

		return &waProto.Message{
			ReactionMessage: &waProto.ReactionMessage{
				Key:               key,
				Text:              proto.String(textReact),
				SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
			},
		}
	},
}
