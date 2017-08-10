package bot

type Update struct {
	Messages []*Message
	Buttons  []*Button
	Inlines  []*Inline
}

type Button struct {
	Text              string
	Handler           Handler
	URL               string
	SwitchInlineQuery string
	callbackData      string
}

type Message struct {
	Text string
}

type Inline struct {
	Id    string
	Query string
}

type Response struct {
	Text          string
	Buttons       [][]*Button
	ClickedButton *Button
	ReplyHandler  Handler

	messageId int
}

type InlineAnswer struct {
	InlineId    string
	Title       string
	MessageText string
	Description string
	Button      *Button
}

func (r *Response) AddButtonString(text string, handler Handler) {
	b := make([]*Button, 0)
	b = append(b, &Button{Text: text, Handler: handler})
	r.AddButtonRow(b...)
}

func (r *Response) AddButton(b *Button) {
	buttonRow := make([]*Button, 0)
	buttonRow = append(buttonRow, b)
	r.AddButtonRow(buttonRow...)
}

func (r *Response) AddButtonRow(b ...*Button) {
	r.Buttons = append(r.Buttons, b)
}

func (r *Response) ClearButtons() {
	r.Buttons = make([][]*Button, 0)
}
