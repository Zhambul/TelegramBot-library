package bot

import (
	"bot/comm"
	"errors"
)

func (c *Context) toMessage(msg *comm.Message) *Message {
	c.log.info("Bot::toMessage %v", msg)
	return &Message{
		Text: msg.Text,
	}
}

func (c *Context) toResponse(callback *comm.Callback) (*Response, error) {
	c.log.info("Bot::toResponse START")
	if len(c.responses) == 0 {
		c.log.err("Bot::toResponse. No responses in context")
	}
	for _, r := range c.responses {
		for _, btnRow := range r.Buttons {
			for _, btn := range btnRow {
				c.log.info("Bot::toResponse. Comparing callbackData %v - %v", btn.callbackData, callback.CallbackData)
				if btn.callbackData == callback.CallbackData {
					r.ClickedButton = btn
					c.log.info("Bot::toResponse END. Found clicked button in response %+v", r)
					return r, nil
				}
			}
		}
	}
	c.log.err("Bot::toResponse END. could no find response by callback")
	return nil, errors.New("Not found")
}

func (c *Context) toInline(i *comm.Inline) *Inline {
	c.log.info("Bot::toInline")
	return &Inline{
		Id:    i.Id,
		Query: i.Query,
	}
}
