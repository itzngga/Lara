package state

import (
	"github.com/itzngga/Roxy/command"
)

var question = &command.StateCommand{
	Name: "question",
	RunFunc: func(c *command.StateFuncContext) {
		if question, ok := c.Locals["curr_question"]; ok {
			c.SendMessage(question)
			return
		}
		answers, ok := c.Locals["answers"].([]string)
		if ok {
			c.Locals["answers"] = answers
		}
	},
}
