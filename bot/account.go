package bot

import (
	"bot/comm"
	"log"
)

type BotAccount struct {
	FirstName string
	LastName  string
	ChatId    int
}

func BotAccountFrom(from *comm.From) *BotAccount {
	log.Println("Bot::BotAccountFrom")
	return &BotAccount{
		ChatId:    from.Id,
		FirstName: from.FirstName,
		LastName:  from.LastName,
	}
}
