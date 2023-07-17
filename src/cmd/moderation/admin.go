package moderation

import (
	"fmt"
	"github.com/itzngga/Lara/entity"
	"github.com/itzngga/Lara/repo"
	"github.com/itzngga/Lara/src/cmd/constant"
	"github.com/itzngga/Lara/src/mid"
	util2 "github.com/itzngga/Lara/util"
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util"
	"github.com/valyala/fasttemplate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"math"
	"strconv"
	"strings"
	"time"
)

func init() {
	embed.Commands.Add(admin)
}

var admin = &command.Command{
	Name:        "admin",
	Aliases:     []string{"adm", "min"},
	Description: "Admin console",
	Category:    constant.MODERATION_CATEGORY,
	Middleware:  mid.AdminMiddleware,
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if len(ctx.Arguments) >= 1 {
			switch ctx.Arguments[0] {
			case "list":
				var admins []entity.AdminEntity
				var totalData int

				if len(ctx.Arguments) == 1 || ctx.Arguments[1] != "all" {
					var q, page string
					command.NewUserQuestion(ctx).SetQuestion("Search Query?", &q).SetQuestion("Select Page", &page).ExecWithParser()

					if q != "" {
						q = strings.TrimPrefix(q, "list ")
					}
					pPage, err := strconv.Atoi(page)
					if err != nil {
						return ctx.GenerateReplyMessage("error: invalid page number")
					}

					admins, totalData, err = repo.AdminRepository.GetAllAdmins(q, pPage, 25)
					if err != nil {
						fmt.Printf("error: %v\n", err)
						return nil
					}
				} else {
					var err error
					admins, totalData, err = repo.AdminRepository.GetAllAdmins("", 1, 0)
					if err != nil {
						fmt.Printf("error: %v\n", err)
						return nil
					}
				}

				var parsedAdmins string
				for i, admin := range admins {
					parsed, _ := util.JIDToString(admin.Username)
					toNow := util2.HumanizeDuration(time.Now().In(time.Local).Sub(admin.CreatedAt))
					parsedAdmins += fmt.Sprintf("%d. \n*Username*: %s\n*WaMe*: %s%s\n*CreatedAt*: %s\n\n", i+1, parsed, "wa.me/", parsed, toNow)
				}

				totalPage := int32(math.Ceil(float64(totalData) / float64(25)))
				if totalPage == 0 {
					totalPage = 1
				}

				template := fasttemplate.New(constant.ADMIN_LIST_ALL_TEMPLATE, "[", "]")
				result := template.ExecuteString(map[string]interface{}{
					"total":      []byte(strconv.Itoa(totalData)),
					"page":       []byte("1"),
					"content":    []byte(parsedAdmins),
					"total_page": []byte(strconv.Itoa(int(totalPage))),
				})

				return ctx.GenerateReplyMessage(result)
			case "add":
				var target string

				if len(ctx.Arguments) >= 2 {
					target = ctx.Arguments[1]
				} else {
					return ctx.GenerateReplyMessage("error: no tag user or number")
				}
				mentioned := util.ParseMentionedJid(ctx.Message)
				if len(mentioned) != 0 {
					target = mentioned[0]
				}

				jid := util.ParseUserJid(constant.NumberRegex.FindString(target))
				if jid.String() == "" {
					return ctx.GenerateReplyMessage("error: invalid number")
				}

				result := jid.ToNonAD().String()
				ok, err := repo.AdminRepository.IsValidAdmin(result)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return nil
				}
				if ok {
					return ctx.GenerateReplyMessage("error: target already admin")
				}

				now := time.Now().In(time.Local)
				err = repo.AdminRepository.InsertNewAdmin(entity.AdminEntity{
					Username:  result,
					Password:  "123",
					CreatedAt: now,
				})

				if err != nil {
					fmt.Printf("error: %v\n", err)
					return nil
				}

				return ctx.GenerateReplyMessage("success: add target to admin")
			case "del":
				var target string

				if len(ctx.Arguments) >= 2 {
					target = ctx.Arguments[1]
				} else {
					return ctx.GenerateReplyMessage("error: no tag user or number")
				}
				mentioned := util.ParseMentionedJid(ctx.Message)
				if len(mentioned) != 0 {
					target = mentioned[0]
				}

				jid := util.ParseUserJid(constant.NumberRegex.FindString(target))
				if jid.String() == "" {
					return ctx.GenerateReplyMessage("error: invalid number")
				}

				result := jid.ToNonAD().String()
				ok, err := repo.AdminRepository.IsValidAdmin(result)
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return nil
				}
				if !ok {
					return ctx.GenerateReplyMessage("error: not a admin user")
				}

				err = repo.AdminRepository.DeleteAdmin(jid.ToNonAD().String())
				if err != nil {
					fmt.Printf("error: %v\n", err)
					return nil
				}

				return ctx.GenerateReplyMessage("success: remove target from admin")
			}
		}
		return nil
	},
}
