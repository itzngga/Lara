package general

import (
	"fmt"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/valyala/fasttemplate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"strings"
	"time"
)

var helpMenuTemplate *fasttemplate.Template

func init() {
	embed.Commands.Add(help)

	helpMenuTemplate = fasttemplate.New(constant.HELP_MENU_FORMAT, "{{ ", " }}")
}

func generateHelp(prefix, category string, cmdMap map[string][]*command.Command) string {
	var helpStr = fmt.Sprintf("*── 「 %s 」 ──*\n", strings.ToUpper(category))
	for index, cmd := range cmdMap[category] {
		var parsedAliases string
		for _, alias := range cmd.Aliases {
			parsedAliases += "*" + prefix + alias + "*, "
		}
		if len(cmd.Aliases) == 0 {
			parsedAliases = "-"
		}
		helpStr += fmt.Sprintf("%d. *%s%s*\n%s\nAliases: %s\n\n", index+1, prefix, cmd.Name, cmd.Description, parsedAliases)
	}
	return helpStr
}

var help = &command.Command{
	Name:        "help",
	Aliases:     []string{"menu", "tulung", "list"},
	Description: "Bot Menu",
	Category:    constant.GENERAL_CATEGORY,
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var cmdMap = make(map[string][]*command.Command, 0)
		for _, cmd := range embed.Commands.Get() {
			for _, category := range embed.Categories.Get() {
				if category == cmd.Category {
					if _, ok := cmdMap[category]; ok {
						cmdMap[category] = append(cmdMap[category], cmd)
					} else {
						cmdMap[category] = make([]*command.Command, 0)
						cmdMap[category] = append(cmdMap[category], cmd)
					}
					break
				}

			}
		}

		if len(ctx.Arguments) > 0 {
			args := strings.ToLower(ctx.Arguments[0])
			if args == "1" || args == "general" {
				result := generateHelp(ctx.Prefix, constant.GENERAL_CATEGORY, cmdMap)
				return ctx.GenerateReplyMessage(result)
			} else if args == "2" || args == "media" {
				result := generateHelp(ctx.Prefix, constant.MEDIA_CATEGORY, cmdMap)
				return ctx.GenerateReplyMessage(result)
			} else if args == "3" || args == "utility" {
				result := generateHelp(ctx.Prefix, constant.UTILITY_CATEGORY, cmdMap)
				return ctx.GenerateReplyMessage(result)
			} else {
				var listMenu string
				var index = 1
				for category := range cmdMap {
					listMenu += fmt.Sprintf("*[%d]* %s\n", index, strings.Title(category))
					index += 1
				}
				result := helpMenuTemplate.ExecuteString(map[string]interface{}{
					"pushname": []byte(ctx.MessageInfo.PushName),
					"date":     []byte(time.Now().Format("2006-01-02 15:04:05")),
					"menu":     []byte(listMenu),
					"prefix":   []byte(ctx.Prefix),
				})

				return ctx.GenerateReplyMessage(result)
			}
		} else {
			var listMenu string
			var index = 1
			for category := range cmdMap {
				listMenu += fmt.Sprintf("*[%d]* %s\n", index, strings.Title(category))
				index += 1
			}
			result := helpMenuTemplate.ExecuteString(map[string]interface{}{
				"pushname": []byte(ctx.MessageInfo.PushName),
				"date":     []byte(time.Now().Format("2006-01-02 15:04:05")),
				"menu":     []byte(listMenu),
				"prefix":   []byte(ctx.Prefix),
			})

			return ctx.GenerateReplyMessage(result)
		}
	},
}
