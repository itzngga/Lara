package mid

import (
	"fmt"
	"github.com/itzngga/Lara/repo"
	"github.com/itzngga/Roxy/command"
)

func AdminMiddleware(ctx *command.RunFuncContext) bool {
	ok, err := repo.AdminRepository.IsValidAdmin(ctx.MessageInfo.Sender.ToNonAD().String())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return false
	}

	return ok
}

func PremiumMiddleware(ctx *command.RunFuncContext) bool {
	if AdminMiddleware(ctx) {
		return true
	}

	ok, err := repo.PremiumRepository.IsValidPremium(ctx.MessageInfo.Sender.ToNonAD().String())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return false
	}

	return ok
}
