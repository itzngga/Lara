package general

import (
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Lara/util"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func init() {
	embed.Commands.Add(translate)
}

var translate = &command.Command{
	Name:        "translate",
	Aliases:     []string{"tl"},
	Description: "Translate text from source to target",
	Category:    constant.GENERAL_CATEGORY,
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var text, source, target string
		command.NewUserQuestion(ctx).
			SetQuestion("Text to Translate?", &text).
			SetQuestion("Source Language?", &source).
			SetQuestion("Target Language?", &target).
			ExecWithParser()

		text, err := util.Translate(text, source, text)
		if err != nil {
			return ctx.GenerateReplyMessage("error: " + err.Error())
		}

		return ctx.GenerateReplyMessage(text)
	},
}
