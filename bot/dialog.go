package bot
//
//import (
//	"time"
//	"sync"
//	"log"
//	"bot/comm"
//)
//
//type Dialog struct {
//	//next handler to handle
//	BotAccount         *BotAccount
//	NextMessageHandler MessageHandler
//
//	responses []*Response
//
//	messageHandlers map[MessageMatcher]MessageHandler
//	buttonHandlers  map[string]MessageHandler
//	inlineHandler   InlineHandler
//
//	//to track and delete old contexts
//	lastModified time.Time
//
//	lock sync.Mutex
//
//	log *contextLogger
//}
//
//type MessageMatcher interface {
//	Match(Message) bool
//}
//func (d *Dialog) onInline(i *Inline) {
//	d.log.info("Context::onInline START")
//	d.log.info("Context::lock")
//	d.lock.Lock()
//	defer func() {
//		d.log.info("Context::unlock")
//		d.lock.Unlock()
//	}()
//
//	if inlineHandler != nil {
//		d.sendInlineAnswer(inlineHandler.Handle(d))
//	}
//	d.Inline = nil
//
//	d.log.info("Context::onInline END")
//}
//
//func (d *Dialog) onCallback(r *Response) {
//	d.log.info("Context::onCallback START")
//	d.log.info("Context::lock")
//	d.lock.Lock()
//	defer func() {
//		d.log.info("Context::unlock")
//		d.lock.Unlock()
//	}()
//
//	d.CurrentResponse = r
//	if r.ClickedButton == nil {
//		d.log.err("Context::onCallback END. ClickedButton is nil")
//		return
//	}
//	h := r.ClickedButton.Handler
//	if h != nil {
//		d.log.info("Context::onCallback. Calling handler START")
//		newR := h.Handle(d)
//		d.log.info("Context::onCallback. Calling handler END")
//		d.Send(newR)
//		r.ClickedButton = nil
//		d.CurrentResponse = nil
//	} else {
//		d.log.info("Context::unlock")
//		d.log.err("Context::onCallback END. Handler is nil")
//	}
//
//	d.log.info("Context::onCallback END")
//}
//
//func (d *Dialog) onReply(m *Message, repliedToId int) {
//	d.log.info("Context::onReply START")
//	d.log.info("Context::lock")
//	d.lock.Lock()
//	defer func() {
//		d.log.info("Context::unlock")
//		d.lock.Unlock()
//	}()
//
//	canHandle := make([]Handler, 0)
//	for _, resp := range d.responses {
//		if resp.ReplyHandler == nil {
//			continue
//		}
//
//		if resp.messageId != repliedToId {
//			continue
//		}
//
//		canHandle = append(canHandle, resp.ReplyHandler)
//	}
//
//	if len(canHandle) < 1 {
//		d.log.info("Context::onReply END. Could not find response to reply")
//		return
//	}
//
//	if len(canHandle) > 1 {
//		d.log.err("Context::onReply END. Too many handlers")
//		return
//	}
//
//	d.Message = m
//	r := canHandle[0].Handle(d)
//	d.Send(r)
//	d.Message = nil
//	d.log.info("Context::onReply END")
//}
//
//func (d *Dialog) onMessage(m *Message) {
//	d.log.info("Dialog::onMessage START")
//	d.log.info("Dialog::lock")
//	d.lock.Lock()
//	defer func() {
//		d.log.info("Dialog::unlock")
//		d.lock.Unlock()
//	}()
//
//	canHandle := make([]MessageHandler, 0)
//	for matcher, h := range d.messageHandlers {
//		if matcher.Match(m) {
//			canHandle = append(canHandle, h)
//		}
//	}
//
//	if d.NextMessageHandler != nil {
//		canHandle = append(canHandle, d.NextMessageHandler)
//		d.NextMessageHandler = nil
//	}
//
//	if len(canHandle) < 1 {
//		d.log.info("Context::onMessage END. No handler to handle")
//		return
//	}
//
//	if len(canHandle) > 1 {
//		d.log.info("Context::onMessage END. Too much handler to handle")
//		return
//	}
//
//	r := canHandle[0].Handle(d, m)
//	d.Send(r)
//}
//
//func (d *Dialog) Send(r *Response) {
//	d.log.info("Context::SendReply START")
//	if r == nil {
//		d.log.info("Context::SendReply END, response is nil")
//		return
//	}
//
//	defer func() {
//		d.addResponse(r)
//	}()
//
//	reply := &comm.Reply{
//		ChatId:      d.BotAccount.ChatId,
//		Text:        r.Text,
//		ReplyMarkup: parseKeyboard(r.Buttons),
//		ParseMode:   "Markdown",
//	}
//
//	if r.messageId != 0 {
//		reply.MessageId = r.messageId
//		err := comm.UpdateMessage(reply)
//		if err != nil {
//			d.log.err("Context::SendReply END%v", err)
//		}
//		return
//	}
//
//	id, err := comm.SendMessage(reply)
//	if err != nil {
//		log.Printf("ERROR: %v\n", err)
//		return
//	}
//
//	r.messageId = id
//	d.log.info("Context::SendReply END")
//}
//
//func (d *Dialog) addResponse(r *Response) {
//	d.log.info("Context::addResponse START")
//	if d.getResponseByMessageId(r.messageId) == nil {
//		d.responses = append(d.responses, r)
//		d.log.info("Context::addResponse END. New response added")
//		return
//	}
//	d.log.info("Context::addResponse END. No need to add")
//}
//
//func (d *Dialog) getResponseByMessageId(msgId int) *Response {
//	d.log.info("Context::getResponseByMessageId START")
//	for _, resp := range d.responses {
//		if resp.messageId == msgId {
//			d.log.info("Context::getResponseByMessageId END. Found resp %+v", resp)
//			return resp
//		}
//	}
//	d.log.info("Context::getResponseByMessageId END, could not find")
//	return nil
//}
//func parseKeyboard(buttons [][]*Button) *comm.InlineKeyboardMarkup {
//	if len(buttons) == 0 {
//		return nil
//	}
//
//	ik1 := make([][]*comm.InlineKeyboardButton, 0)
//
//	for _, buttonRow := range buttons {
//		ik2 := make([]*comm.InlineKeyboardButton, 0)
//		for _, b := range buttonRow {
//			if b.SwitchInlineQuery == "" {
//				b.callbackData = randString()
//			}
//
//			ik2 = append(ik2, &comm.InlineKeyboardButton{
//				Text:              b.Text,
//				CallbackData:      b.callbackData,
//				Url:               b.URL,
//				SwitchInlineQuery: b.SwitchInlineQuery,
//			})
//		}
//		ik1 = append(ik1, ik2)
//	}
//
//	return &comm.InlineKeyboardMarkup{
//		InlineKeyboard: ik1,
//	}
//}
//
////-------------------------------
////type MessageContext struct {
////	Message Message
////}
////
////type ButtonContext struct {
////	Button          *Button
////	CurrentResponse Response
////}
////
////type InlineContext struct {
////	Inline Inline
////}
//
////-------------------------------
//
//type MessageHandler interface {
//	Handle(d *Dialog, c *Message) *Response
//}
//
//type ButtonHandler interface {
//	Handle(d *Dialog, b *Button, currentResponse *Response) *Response
//}
//
//type InlineHandler1 interface {
//	Handle(d *Dialog, i Inline)
//}
