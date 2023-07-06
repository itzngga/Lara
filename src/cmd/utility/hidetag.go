package cmd

import (
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/types"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"strings"
)

func init() {
	embed.Commands.Add(hidetag)
}

var hidetag = &command.Command{
	Name:        "hidetag",
	Category:    constant.UTILITY_CATEGORY,
	Description: "Hidetag all user in group.",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if ctx.MessageInfo.IsGroup {
			groupInfo, err := ctx.GetGroupInfo(ctx.MessageInfo.Chat.ToNonAD().String())
			if err != nil {
				return ctx.GenerateReplyMessage("error: " + err.Error())
			}
			var mentionedJids []string

			for _, participant := range groupInfo.Participants {
				if participant.AddRequest == nil {
					mentionedJids = append(mentionedJids, participant.JID.ToNonAD().String())
				}
			}

			return &waProto.Message{
				ExtendedTextMessage: &waProto.ExtendedTextMessage{
					Text: types.String(strings.Join(ctx.Arguments, " ")),
					ContextInfo: &waProto.ContextInfo{
						MentionedJid: mentionedJids,
					},
				},
			}
		}
		return nil
	},
}
