package general

import (
	"fmt"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

func init() {
	embed.Commands.Add(speed)
}

var speed = &command.Command{
	Name:        "speed",
	Aliases:     []string{"ping", "sp"},
	Description: "Check Bot Speed",
	Category:    constant.GENERAL_CATEGORY,
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		t := time.Now()
		ctx.SendReplyMessage("wait...")
		return ctx.GenerateReplyMessage(fmt.Sprintf("Latency: %f seconds", time.Now().Sub(t).Seconds()))
	},
}
