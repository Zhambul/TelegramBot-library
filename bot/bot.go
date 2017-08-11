package bot

import (
	"bot/comm"
	"time"
	"log"
)

var defaultHandlers map[Matcher]Handler
var inlineHandler InlineHandler
var contexts map[int]*Context

func Init(token string) {
	comm.Init(token)
	defaultHandlers = make(map[Matcher]Handler)
	contexts = make(map[int]*Context)
	go deleteOldContexts()
}

func deleteOldContexts() {
	ticker := time.NewTicker(1 * time.Hour)

	for {
		<-ticker.C
		log.Println("Looking for old contexts")
		for chatId, c := range contexts {
			log.Println("Deleting old context")
			if time.Since(c.lastModified) > 6*time.Hour {
				delete(contexts, chatId)
			}
		}
	}
}

func RegisterInlineHandler(h InlineHandler) {
	inlineHandler = h
}

func RegisterMatchedHandler(m Matcher, h Handler) {
	defaultHandlers[m] = h
}

func RegisterHandler(text string, h Handler) {
	RegisterMatchedHandler(&simpleMatcher{text}, h)
}

func GetContext(acc *BotAccount) *Context {
	if c, has := contexts[acc.ChatId]; has {
		return c
	}
	c := newContext(acc)
	for m, h := range defaultHandlers {
		c.RegisterHandler(m, h)
	}
	contexts[acc.ChatId] = c
	return c
}
