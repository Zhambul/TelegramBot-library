package bot

type Dialog struct {
	//next handler to handle
	BotAccount *BotAccount
	NextHandler Handler
	responses []*Response
}

//-------------------------------

type MessageContext struct {
	Text string
}

type ButtonContext struct {
	Button *Button
}

type InlineContext struct {

}

//-------------------------------

type MessageHandler interface {
	Handle(d *Dialog)
}

type ButtonHandler interface {
	Handle(d *Dialog)
}

type InlineHandler1 interface {
	Handle(d *Dialog)
}
