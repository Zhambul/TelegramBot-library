package bot

import (
	"bot/comm"
	"errors"
)

func (c *Context) toMessage(msg *comm.Message) *Message {
	return &Message{
		Text: msg.Text,
	}
}

func (c *Context) toResponse(callback *comm.Callback) (*Response, error) {
	for _, r := range c.responses {
		for _, btnRow := range r.Buttons {
			for _, btn := range btnRow {
				if btn.callbackData == callback.CallbackData {
					r.ClickedButton = btn
					return r, nil
				}
			}
		}
	}
	return nil, errors.New("Not found")
}

func (c *Context) toInline(i *comm.Inline) *Inline {
	return &Inline{
		Id:    i.Id,
		Query: i.Query,
	}
}
