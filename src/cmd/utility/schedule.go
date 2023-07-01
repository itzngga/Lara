package cmd

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Lara/util"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/xhit/go-str2duration/v2"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

func init() {
	embed.Commands.Add(schedule)

	cron = gocron.NewScheduler(time.Local)
	cron.TagsUnique()

	cron.StartAsync()
}

var cron *gocron.Scheduler

var schedule = &command.Command{
	Name:        "schedule",
	Category:    constant.UTILITY_CATEGORY,
	Description: "Schedule a message.",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if len(ctx.Arguments) == 1 && ctx.Arguments[0] == "stop" {
			cron.RemoveByTag(ctx.Number)
			return ctx.GenerateReplyMessage("Schedule stopped")
		}

		var duration string
		var captured *waProto.Message
		command.NewUserQuestion(ctx).
			SetQuestion("Every Duration?", &duration).
			CaptureMediaQuestion("Please send/reply a media message", &captured).
			WithOkEmoji().
			ExecWithParser()

		timeDuration, err := str2duration.ParseDuration(duration)
		if err != nil {
			return ctx.GenerateReplyMessage(fmt.Sprintf("error: invalid %s duration", duration))
		}

		if timeDuration < time.Second*5 {
			return ctx.GenerateReplyMessage("error: to avoid spamming, scheduling is allowed greater than 5 seconds")
		}

		ctx.SendReplyMessage("Success scheduling message every " + util.HumanizeDuration(timeDuration) + fmt.Sprintf("\n\nTo stop, do *%s%s stop*", ctx.Prefix, ctx.CurrentCommand.Name))
		job, _ := cron.Every(timeDuration).Tag(ctx.Number).Do(func() {
			ctx.SendReplyMessage(captured)
		})

		return ctx.GenerateReplyMessage(job.Tag)
	},
}
