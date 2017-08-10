package bot

import (
	"log"
	"bot/comm"
	"math/rand"
)

func (c *Context) onInline(i *Inline) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Inline = i
	log.Println("ON INLINE")
	sendInlineAnswer(inlineHandler.Handle(c))
	c.Inline = nil
}

func (c *Context) onCallback(r *Response) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.CurrentResponse = r
	newR := r.ClickedButton.Handler.Handle(c)
	c.SendReply(newR)
	r.ClickedButton = nil
	c.CurrentResponse = nil
}

func (c *Context) onReply(m *Message, repliedToId int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if repliedToId == 0 {
		panic("qweqweqeweqw")
	}
	log.Printf("Checking replyed to %v\n", repliedToId)
	canHandle := make([]Handler, 0)
	for _, resp := range c.responses {
		if resp.ReplyHandler == nil {
			continue
		}

		if resp.messageId != repliedToId {
			continue
		}

		log.Println("Calling Reply Handler")
		canHandle = append(canHandle, resp.ReplyHandler)

	}
	if len(canHandle) < 1 {
		log.Println("Could not find response to reply")
		return
	}

	if len(canHandle) > 1 {
		panic("Too many handlers")
	}

	c.Message = m
	r := canHandle[0].Handle(c)
	c.SendReply(r)
	c.Message = nil
}

func (c *Context) onMessage(m *Message) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Message = m
	if c.CurrentResponse == nil {
		c.CurrentResponse = &Response{}
	}

	if c.NextHandler != nil {
		r := c.NextHandler.Handle(c)
		c.NextHandler = nil
		c.SendReply(r)
		return
	}

	canHandle := make([]Handler, 0)
	for m, h := range c.handlers {
		if m.Match(c) {
			canHandle = append(canHandle, h)
		}
	}

	if len(canHandle) < 1 {
		log.Println("WARNING: NO HANLDER TO HANDLE")
		return
	}

	if len(canHandle) > 1 {
		log.Println("WARNING: TOO MORE HANLDER TO HANDLE")
		return
	}
	r := canHandle[0].Handle(c)
	c.SendReply(r)
}

func sendInlineAnswer(a *InlineAnswer) {

	log.Println("sendInlineAnswer START")
	inlineResult := make([]*comm.InlineQueryResult, 0)
	m := make([][]*comm.InlineKeyboardButton, 0)
	m1 := make([]*comm.InlineKeyboardButton, 0)
	m1 = append(m1, &comm.InlineKeyboardButton{
		Text: a.Button.Text,
		Url:  a.Button.URL,
	})

	m = append(m, m1)
	inlineResult = append(inlineResult, &comm.InlineQueryResult{
		Type:  "article",
		Id:    "qweqew",
		Title: a.Title,
		InputMessageContent: &comm.InputMessageContent{
			MessageText: a.MessageText,
			ParseMode:   "Markdown",
		},
		Description: a.Description,
		ReplyMarkup: &comm.InlineKeyboardMarkup{
			InlineKeyboard: m,
		},
	})
	answer := &comm.InlineQueryAnswer{
		InlineQueryId: a.InlineId,
		Results:       inlineResult,
		CacheTime:     0,
	}
	err := AnswerInlineQuery(answer)
	if err != nil {
		panic(err)
	}
	log.Println("sendInlineAnswer END")
}

func (c *Context) SendReply(r *Response) {
	log.Println("SEND REPLY START")
	if r == nil {
		return
	}
	//
	//c.responses = append(c.responses, r)

	reply := &comm.Reply{
		ChatId:      c.BotAccount.ChatId,
		Text:        r.Text,
		ReplyMarkup: c.parseKeyboard(r),
		ParseMode:   "Markdown",
	}

	if r.messageId != 0 {
		log.Println("Context::updating message")
		reply.MessageId = r.messageId
		err := UpdateMessage(reply)
		if err != nil {

			/*TODO
    2017/08/09 22:42:16 HTTP POST - https://api.telegram.org/bot366621722:AAH5scmfkscK8_Es0dNIJj8gZ-lxluCYD1o/editMessageText
	2017/08/09 22:42:16 {"chat_id":104563894,"text":"*Gleb* из группы *Тюлени* вызывает через *10* минут\n\nZhambyl - Нет\nMaria - Нет\nGleb - Да\n","message_id":3427,"reply_markup":{"inline_keyboard":[[{"text":"Да","callback_data":"WAijDehyBvKT"},{"text":"Нет","callback_data":"MwOQaYLVcMfr"}],[{"text":"Отменить","callback_data":"mFurHTuhzgZF"}]]},"parse_mode":"Markdown"}
	panic: HTTP ERROR: url - https://api.telegram.org/bot366621722:AAH5scmfkscK8_Es0dNIJj8gZ-lxluCYD1o/editMessageText
	, status code - 400 body - {"ok":false,"error_code":400,"description":"Bad Request: message to edit not found"}
			  */

			panic(err)
		}
		return
	}

	log.Println("Context::sending message")
	id, err := SendMessage(reply)
	if err != nil {
		panic(err)
	}

	r.messageId = id

	if c.getResponseByMessageId(r.messageId) == nil {
		log.Println("APPENDING RESPONSE")
		c.responses = append(c.responses, r)
	} else {
		log.Println("NOT APPENDING RESPONSE")
	}

	log.Println("SEND REPLY END")
}

func (c *Context) getResponseByMessageId(msgId int) *Response {
	for _, resp := range c.responses {
		if resp.messageId == msgId {
			return resp
		}
	}
	return nil
}

func (c *Context) parseKeyboard(r *Response) *comm.InlineKeyboardMarkup {
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

	return &comm.InlineKeyboardMarkup{
		InlineKeyboard: ik1,
	}
}

func (c *Context) DeleteResponse(response *Response) error {
	return deleteMessage(&comm.DeleteMessage{
		ChatId:    c.BotAccount.ChatId,
		MessageId: response.messageId,
	})
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
