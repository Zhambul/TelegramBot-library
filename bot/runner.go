package bot

import (
	"log"
)

// 0. wait for a update
// 1. find a context for update
// 2. map to bot structs
// 3. prepare for a handler call
// 4. get a response from handler
// 5. map response to telegram reply
// 6. serialize and send
func Run() {
	var offset int
	for {
		updates := GetUpdates(offset)
		log.Println("Bot RUNNER START")

		for _, msg := range updates.Messages {
			c := RegisterContext(NewBotAccount(msg.From))
			if msg.ReplyToMessage != nil {
				log.Println("ON REPLY")
				c.onReply(c.toMessage(msg), msg.ReplyToMessage.MessageId)
			} else {
				c.onMessage(c.toMessage(msg))
			}
		}

		for _, callback := range updates.Callbacks {
			c := RegisterContext(NewBotAccount(callback.From))
			r, err := c.toResponse(callback)
			if err != nil {
				log.Println("could not find Callback Response")
				continue
			}

			c.onCallback(r)
		}

		for _, inline := range updates.Inlines {
			c := RegisterContext(NewBotAccount(inline.From))
			c.onInline(c.toInline(inline))
		}

		offset = updates.NextUpdateId
		log.Println("Bot RUNNER END")
	}
}
