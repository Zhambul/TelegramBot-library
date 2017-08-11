package bot

import "bot/comm"

type BotAccount struct {
	FirstName string
	LastName  string
	ChatId    int
}

func BotAccountFrom(from *comm.From) *BotAccount {
	return &BotAccount{
		ChatId:    from.Id,
		FirstName: from.FirstName,
		LastName:  from.LastName,
	}
}
