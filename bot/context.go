package bot

import (
	"time"
	"sync"
	"log"
)

type Context struct {
	//state for handler to inspect
	Message         *Message
	CurrentResponse *Response
	Inline          *Inline

	//telegram account info
	BotAccount *BotAccount

	//next handler to handle
	NextHandler Handler

	//inner state to choose handler
	responses []*Response
	handlers  map[Matcher]Handler

	//to track and delete old contexts
	lastModified time.Time
	lock         sync.Mutex

	log *contextLogger
}

//construct a context for an account
func newContext(acc *BotAccount) *Context {
	log.Println("Bot::newContext")
	c := &Context{
		BotAccount:      acc,
		handlers:        make(map[Matcher]Handler),
		CurrentResponse: &Response{},
	}
	c.log = newContextLogger(acc, c)
	return c
}

func (c *Context) RegisterHandler(m Matcher, h Handler) {
	log.Println("Context::RegisterHandler")
	c.handlers[m] = h
}
