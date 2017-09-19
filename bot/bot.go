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
	log.Println("Bot::Init START")
	comm.Init(token)
	defaultHandlers = make(map[Matcher]Handler)
	contexts = make(map[int]*Context)
	//go deleteOldContexts()
	log.Println("Bot::Init END")
}

func RegisterInlineHandler(h InlineHandler) {
	log.Println("Bot::RegisterInlineHandler")
	inlineHandler = h
}

func RegisterMatchedHandler(m Matcher, h Handler) {
	log.Println("Bot::RegisterMatchedHandler")
	defaultHandlers[m] = h
}

func RegisterHandler(text string, h Handler) {
	log.Println("Bot::RegisterHandler")
	RegisterMatchedHandler(&simpleMatcher{text}, h)
}

func EnableWebhook(host string) error {
	log.Println("Bot::EnableWebhook")
	return comm.EnableWebhook(host)
}

func GetContext(acc *BotAccount) *Context {
	log.Println("Bot::GetContext START")
	if contexts == nil {
		contexts = make(map[int]*Context)
	}
	if c, has := contexts[acc.ChatId]; has {
		log.Printf("Bot::GetContext END. Found old context for %+v\n", acc)
		c.lastModified = time.Now()
		return c
	}
	c := newContext(acc)
	for m, h := range defaultHandlers {
		c.RegisterHandler(m, h)
	}
	contexts[acc.ChatId] = c
	log.Printf("Bot::GetContext END. Created new context for %+v\n", acc)
	return c
}

func deleteOldContexts() {
	log.Println("Bot::deleteOldContexts START")
	ticker := time.NewTicker(1 * time.Hour)

	for {
		<-ticker.C
		log.Println("Bot::deleteOldContexts. Looking for old contexts")
		for chatId, c := range contexts {
			if time.Since(c.lastModified) > 24 *time.Hour {
				log.Println("Bot::deleteOldContexts. Deleting old context")
				delete(contexts, chatId)
			}
		}
	}

	log.Println("Bot::deleteOldContexts END")
}