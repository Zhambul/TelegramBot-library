package bot

import (
	"time"
	"sync"
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
}

//construct a context for an account
func newContext(acc *BotAccount) *Context {
	return &Context{
		BotAccount:      acc,
		handlers:        make(map[Matcher]Handler),
		CurrentResponse: &Response{},
	}
}

func (c *Context) RegisterHandler(m Matcher, h Handler) {
	c.handlers[m] = h
}
