package cmd

import (
	"github.com/itzngga/Lara/src/cmd/constant"
	_ "github.com/itzngga/Lara/src/cmd/general"
	_ "github.com/itzngga/Lara/src/cmd/media"
	_ "github.com/itzngga/Lara/src/cmd/utility"
	"github.com/itzngga/Roxy/embed"
)

func init() {
	embed.Categories.Add(constant.GENERAL_CATEGORY)
	embed.Categories.Add(constant.MEDIA_CATEGORY)
	embed.Categories.Add(constant.UTILITY_CATEGORY)
}
