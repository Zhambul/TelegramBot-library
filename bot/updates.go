package bot

import (
	"log"
	"bot/comm"
	"math/rand"
)

func (c *Context) onInline(i *Inline) {
	c.log.info("Context::onInline START")
	c.log.info("Context::lock")
	defer func() {
		c.log.info("Context::unlock")
		c.lock.Unlock()
	}()

	c.Inline = i
	if inlineHandler != nil {
		r := inlineHandler.Handle(c)
		if r != nil {
			c.sendInlineAnswer(r)
		}
	}
	c.Inline = nil

	c.log.info("Context::onInline END")
}

func (c *Context) onCallback(r *Response) {
	c.log.info("Context::onCallback START")
	c.log.info("Context::lock")
	c.lock.Lock()
	defer func() {
		c.log.info("Context::unlock")
		c.lock.Unlock()
	}()

	c.CurrentResponse = r
	h := r.ClickedButton.Handler
	if h != nil {
		c.log.info("Context::onCallback. Calling handler START")
		newR := h.Handle(c)
		c.log.info("Context::onCallback. Calling handler END")
		c.SendReply(newR)
		r.ClickedButton = nil
		c.CurrentResponse = nil
	} else {
		c.log.info("Context::unlock")
		c.log.err("Context::onCallback END. Handler is nil")
	}

	c.log.info("Context::onCallback END")
}

func (c *Context) onReply(m *Message, repliedToId int) {
	c.log.info("Context::onReply START")
	c.log.info("Context::lock")
	c.lock.Lock()
	defer func() {
		c.log.info("Context::unlock")
		c.lock.Unlock()
	}()

	canHandle := make([]Handler, 0)
	for _, resp := range c.responses {
		if resp.ReplyHandler == nil {
			continue
		}

		if resp.messageId != repliedToId {
			continue
		}

		canHandle = append(canHandle, resp.ReplyHandler)
	}

	if len(canHandle) < 1 {
		c.log.info("Context::onReply END. Could not find response to reply")
		return
	}

	if len(canHandle) > 1 {
		c.log.err("Context::onReply END. Too many handlers")
		return
	}

	c.Message = m
	r := canHandle[0].Handle(c)
	c.SendReply(r)
	c.Message = nil
	c.log.info("Context::onReply END")
}

func (c *Context) onMessage(m *Message) {
	c.log.info("Context::onMessage START")
	c.log.info("Context::lock")
	c.lock.Lock()
	defer func() {
		c.log.info("Context::unlock")
		c.lock.Unlock()
	}()

	c.Message = m
	c.CurrentResponse = &Response{}

	canHandle := make([]Handler, 0)
	for m, h := range c.handlers {
		if m.Match(c) {
			canHandle = append(canHandle, h)
		}
	}

	if c.NextHandler != nil {
		canHandle = append(canHandle, c.NextHandler)
		c.NextHandler = nil
	}

	if len(canHandle) < 1 {
		c.log.info("Context::onMessage END. No handler to handle")
		return
	}

	if len(canHandle) > 1 {
		c.log.info("Context::onMessage END. Too much handler to handle")
		return
	}

	r := canHandle[0].Handle(c)
	c.SendReply(r)

	c.log.info("Context::onMessage END")
}

func (c *Context) sendInlineAnswer(a *InlineAnswer) {
	c.log.info("Context::sendInlineAnswer START")
	res := &comm.InlineQueryResult{
		Type:  "article",
		Id:    "qweqew",
		Title: a.Title,
		InputMessageContent: &comm.InputMessageContent{
			MessageText: a.MessageText,
			ParseMode:   "Markdown",
		},
		Description: a.Description,
	}

	if a.Button != nil {
		m := make([][]*comm.InlineKeyboardButton, 0)
		m1 := make([]*comm.InlineKeyboardButton, 0)
		m1 = append(m1, &comm.InlineKeyboardButton{
			Text: a.Button.Text,
			Url:  a.Button.URL,
		})
		m = append(m, m1)
		res.ReplyMarkup = &comm.InlineKeyboardMarkup{
			InlineKeyboard: m,
		}
	} else {
		//todo set ReplyMarkup
	}

	inlineResult := make([]*comm.InlineQueryResult, 0)
	inlineResult = append(inlineResult, res)

	answer := &comm.InlineQueryAnswer{
		InlineQueryId: a.InlineId,
		Results:       inlineResult,
		CacheTime:     0,
	}
	err := comm.AnswerInlineQuery(answer)
	if err != nil {
		c.log.err("Context::sendInlineAnswer END, %v", err)
		return
	}
	c.log.info("Context::sendInlineAnswer END")
}

func (c *Context) SendReply(r *Response) {
	c.log.info("Context::SendReply START")
	if r == nil {
		c.log.info("Context::SendReply END, response is nil")
		return
	}
	reply := &comm.Reply{
		ChatId:      c.BotAccount.ChatId,
		Text:        r.Text,
		ReplyMarkup: c.parseKeyboard(r),
		ParseMode:   "Markdown",
	}

	if r.messageId != 0 {
		reply.MessageId = r.messageId
		err := comm.UpdateMessage(reply)
		if err != nil {
			c.log.err("Context::SendReply END%v", err)
		}
		return
	}

	id, err := comm.SendMessage(reply)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return
	}

	r.messageId = id

	if c.getResponseByMessageId(r.messageId) == nil {
		c.responses = append(c.responses, r)
	}
	c.log.info("Context::SendReply END")
}

func (c *Context) getResponseByMessageId(msgId int) *Response {
	c.log.info("Context::getResponseByMessageId START")
	for _, resp := range c.responses {
		if resp.messageId == msgId {
			c.log.info("Context::getResponseByMessageId END. Found resp %+v", resp)
			return resp
		}
	}
	c.log.info("Context::getResponseByMessageId END, could not find")
	return nil
}

func (c *Context) parseKeyboard(r *Response) *comm.InlineKeyboardMarkup {
	c.log.info("Context::parseKeyboard START")
	if len(r.Buttons) == 0 {
		return nil
	}

	ik1 := make([][]*comm.InlineKeyboardButton, 0)

	for _, buttonRow := range r.Buttons {
		ik2 := make([]*comm.InlineKeyboardButton, 0)
		for _, b := range buttonRow {
			if b.SwitchInlineQuery == "" {
				b.callbackData = randString()
			}

			ik2 = append(ik2, &comm.InlineKeyboardButton{
				Text:              b.Text,
				CallbackData:      b.callbackData,
				Url:               b.URL,
				SwitchInlineQuery: b.SwitchInlineQuery,
			})
		}
		ik1 = append(ik1, ik2)
	}

	c.log.info("Context::parseKeyboard END")
	return &comm.InlineKeyboardMarkup{
		InlineKeyboard: ik1,
	}
}

func (c *Context) deleteResponseByMessageId(messageId int) error {
	c.log.info("Context::deleteResponseByMessageId")
	return comm.DeleteMessage(&comm.DeleteMsg{
		ChatId:    c.BotAccount.ChatId,
		MessageId: messageId,
	})
}

func (c *Context) DeleteResponse(response *Response) error {
	//todo delete from c.responses
	c.log.info("Context::DeleteResponse")
	return c.deleteResponseByMessageId(response.messageId)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString() string {
	n := 12
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
