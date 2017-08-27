package bot

import (
	"bot/comm"
	"log"
)

/**
 1. wait for a update
 2. find a context for update
 3. map update structs to bot known structs
 4. find a handler to handle
 5. prepare for the handler call
 6. call the handler
 7. get a response from the handler
 8. map response to telegram reply
 9. serialize and send
 10. go to #1.
 */
func Run() {
	log.Println("Bot::Run START")
	for {
		updates, err := comm.GetUpdates()
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			continue
		}
		for _, msg := range updates.Messages {
			log.Println("Bot::Run Message")
			c := GetContext(BotAccountFrom(msg.From))
			if msg.ReplyToMessage != nil {
				go c.onReply(c.toMessage(msg), msg.ReplyToMessage.MessageId)
			} else {
				go c.onMessage(c.toMessage(msg))
			}
		}

		for _, callback := range updates.Callbacks {
			log.Println("Bot::Run Callback")
			c := GetContext(BotAccountFrom(callback.From))
			r, err := c.toResponse(callback)
			if err != nil {
				//go c.deleteResponseByMessageId(callback.Message.MessageId)
				continue
			}

			go c.onCallback(r)
		}

		for _, inline := range updates.Inlines {
			log.Println("Bot::Run Inline")
			c := GetContext(BotAccountFrom(inline.From))
			go c.onInline(c.toInline(inline))
		}
	}
}
